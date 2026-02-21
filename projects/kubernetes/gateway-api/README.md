k3d cluster create test \
--port "80:80@loadbalancer" \
--port "443:443@loadbalancer" \
--k3s-arg "--disable=traefik@server:*"

# https://github.com/traefik/traefik-helm-chart/tree/master/traefik
helm repo add traefik https://traefik.github.io/charts
helm repo update

helm install traefik traefik/traefik \
--version=39.0.2 \
--create-namespace \
--wait \
--namespace=traefik \
--set providers.kubernetesIngress.enabled=false \
--set providers.kubernetesGateway.enabled=true \
--set gatewayClass.enabled=false \
--set gateway.enabled=false \
--set ports.web.port=80 \
--set ports.web.exposedPort=80 \
--set ports.web.protocol=TCP \
--set ports.websecure.port=443 \
--set ports.websecure.exposedPort=443 \
--set ports.websecure.protocol=TCP

kubectl create deployment nginx --image=nginx --namespace=default --port=80
kubectl expose deployment nginx --namespace=default --port=80 --target-port=80

kubectl apply --filename - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: traefik
spec:
  controllerName: traefik.io/gateway-controller
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway-api
  namespace: default
spec:
  gatewayClassName: traefik
  listeners:
    - name: http
      protocol: HTTP
      port: 80
      allowedRoutes:
        namespaces:
          from: All
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: nginx-route
  namespace: default
spec:
  parentRefs:
  - name: gateway-api
    sectionName: http
    kind: Gateway
  hostnames:
  - <REDACTED>.nip.io
  rules:
  - matches:
    - path: 
        type: PathPrefix
        value: /nginx
    filters:
      - type: URLRewrite
        urlRewrite:
          path:
            replacePrefixMatch: /
            type: ReplacePrefixMatch
    backendRefs:
    - name: nginx
      port: 80
EOF

curl http://<REDACTED>.nip.io/nginx