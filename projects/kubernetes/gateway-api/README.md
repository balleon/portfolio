# Gateway API with Traefik

## Overview
This project deploys NGINX and exposes it through Kubernetes Gateway API resources managed by Traefik.

## Security Warning
This setup exposes traffic over HTTP on port 80 for example purposes. Plain HTTP is not encrypted; use HTTPS/TLS in production and enforce HTTP-to-HTTPS redirects.

## Goals
- Install Traefik with Kubernetes Gateway provider enabled.
- Deploy NGINX application resources.
- Route traffic with `GatewayClass`, `Gateway`, and `HTTPRoute`.

## Repository Structure
- `manifests/nginx/`: Namespace, Deployment, and Service for NGINX.
- `manifests/gateway-class.yaml`: GatewayClass definition.
- `manifests/gateway.yaml`: Gateway listener.
- `manifests/http-route.yaml`: path routing to NGINX.

## Prerequisites
- Kubernetes cluster access
- `kubectl`
- `helm`
- DNS entry matching host configured in `http-route.yaml`

## Usage
### 1) Install Traefik
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

### 2) Deploy application and Gateway API manifests
```bash
kubectl apply --filename=./manifests/nginx/{namespace.yaml,deployment.yaml,service.yaml}
kubectl apply --filename=./manifests/{gateway-class.yaml,gateway.yaml,http-route.yaml}
```

## Validation
```bash
kubectl get gatewayclass
kubectl get gateway --namespace=nginx
kubectl get httproute --namespace=nginx
curl http://<REDACTED>/nginx
```

## Cleanup
```bash
kubectl delete --filename=./manifests/{gateway-class.yaml,gateway.yaml,http-route.yaml}
kubectl delete --filename=./manifests/nginx/{service.yaml,deployment.yaml,namespace.yaml}
helm uninstall traefik --namespace=traefik
kubectl delete namespace traefik
```