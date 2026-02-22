# Gateway API with Traefik

This project deploys a simple NGINX HTTP server and exposes it through [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/) resources using [Traefik](https://github.com/traefik/traefik) as the Gateway controller.

## ⚠️ Security warning

This example exposes traffic over HTTP on port 80, which is unencrypted. Do not use this configuration as-is for production workloads handling sensitive data. For production, terminate TLS on port 443, enforce HTTPS redirects, and apply appropriate network access controls.

## What gets deployed

- An NGINX HTTP server in namespace `nginx`
	- A dedicated `Namespace`
	- A `Deployment` running NGINX
	- A `Service` exposing NGINX on port 80
- A `GatewayClass` managed by the Traefik Gateway controller
- A `Gateway` in Namespace `nginx` listening on HTTP/80
- An `HTTPRoute` in Namespace `nginx`
	- Matches path prefix `/nginx`
	- Rewrites `/nginx` to `/`
	- Forwards traffic to `Service/nginx:80`

## Prerequisites

- A running Kubernetes cluster
- `kubectl` configured for the target cluster
- `helm` installed
- External access to Traefik on port 80 (for curl test)
- DNS for the hostname used in `manifests/http-route.yaml` (`<REDACTED>`)

## 1) Install Traefik (Gateway provider enabled)

Reference chart: https://github.com/traefik/traefik-helm-chart/tree/master/traefik

```bash
helm repo add traefik https://traefik.github.io/charts
helm repo update

helm install traefik traefik/traefik \
--version=39.0.2 \
--namespace=traefik \
--create-namespace \
--wait \
--set gateway.enabled=false \
--set gatewayClass.enabled=false \
--set providers.kubernetesIngress.enabled=false \
--set providers.kubernetesGateway.enabled=true \
--set ports.web.exposedPort=80 \
--set ports.web.port=80 \
--set ports.web.protocol=TCP \
--set ports.websecure.exposedPort=443 \
--set ports.websecure.port=443 \
--set ports.websecure.protocol=TCP
```

## 2) Deploy NGINX HTTP server

```bash
kubectl apply --filename=./manifests/nginx/{namespace.yaml,deployment.yaml,service.yaml}
```

## 3) Deploy Gateway API resources

```bash
kubectl apply --filename=./manifests/{gateway-class.yaml,gateway.yaml,http-route.yaml}
```

## 4) Verify resources

```bash
kubectl get gatewayclass
kubectl get gateway --namespace=nginx
kubectl get httproute --namespace=nginx
kubectl get pods --namespace=nginx
kubectl get service --namespace=nginx
```

Optional quick checks:

```bash
kubectl describe gateway gateway --namespace=nginx
kubectl describe httproute route --namespace=nginx
```

## 5) Test routing

```bash
curl http://<REDACTED>/nginx
```

Expected result: NGINX default welcome page HTML.

## Cleanup

```bash
kubectl delete --filename=./manifests/{gateway-class.yaml,gateway.yaml,http-route.yaml}
kubectl delete --filename=./manifests/nginx/{service.yaml,deployment.yaml,namespace.yaml}

helm uninstall traefik --namespace=traefik
kubectl delete namespace traefik
```