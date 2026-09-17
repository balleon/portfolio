helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

# Install Istio Custom Resource Definitions
helm install istio-base istio/base \
--version=1.30.4 \
--namespace=istio-system \
--create-namespace \
--wait

# Install Istio
helm install istiod istio/istiod \
--version=1.30.4 \
--namespace=istio-system \
--create-namespace \
--set=pilot.resources.requests.cpu=100m \
--set=pilot.resources.requests.memory=256Mi \
--set=global.proxy.resources.requests.cpu=10m \
--set=global.proxy.resources.requests.memory=40Mi \
--wait

# Create `secure` Namespace with automatic Istio sidecar injection
kubectl create namespace secure
kubectl label namespace secure istio-injection=enabled

# Deployment HTTP server and curl Pods (source https://github.com/istio/istio/tree/1.30.4/samples)
kubectl apply --namespace=secure --filename=./manifests/httpbin.yaml
kubectl apply --namespace=secure --filename=./manifests/curl.yaml
kubectl apply --namespace=default --filename=./manifests/curl.yaml

# Enforce mTLS for inbound traffic to the `secure` Namespace
kubectl apply --namespace=secure --filename=./manifests/peer-authentication.yaml

# Test
kubectl exec --namespace=secure deployment/curl -- curl --silent --output /dev/null --write-out "%{http_code}\n" http://httpbin.secure.svc.cluster.local:8000/headers  # HTTP code 200
kubectl exec --namespace=default deployment/curl -- curl --silent --output /dev/null --write-out "%{http_code}\n" http://httpbin.secure.svc.cluster.local:8000/headers # Exit code 56