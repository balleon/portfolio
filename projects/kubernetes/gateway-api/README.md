k3d cluster create test \
--port "80:80@loadbalancer" \
--port "443:443@loadbalancer" \
--k3s-arg "--disable=traefik@server:*"

# https://github.com/traefik/traefik-helm-chart/tree/master/traefik
helm repo add traefik https://traefik.github.io/charts && helm repo update

helm install traefik traefik/traefik \
--version=39.0.2 \
--create-namespace \
--wait \
--namespace=traefik \
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

kubectl apply --filename=./manifests/nginx/{namespace.yaml,deployment.yaml,service.yaml}

# https://gateway-api.sigs.k8s.io/api-types/gatewayclass/
kubectl apply --filename=./manifests/gateway-class.yaml

# https://gateway-api.sigs.k8s.io/api-types/gateway/
kubectl apply --filename=./manifests/gateway.yaml

# https://gateway-api.sigs.k8s.io/api-types/httproute/
kubectl apply --filename=./manifests/http-route.yaml

curl http://<REDACTED>/nginx