# Strict mTLS with Istio

## Overview
This project deploys Istio and enforces strict mutual TLS (mTLS) for inbound traffic to a namespace, using the sidecar proxy pattern rather than any application code changes.

Every Pod in the `secure` namespace runs alongside an Istio sidecar proxy that intercepts its traffic. A `PeerAuthentication` resource set to `STRICT` mode instructs those sidecars to reject any inbound connection that isn't mTLS, so only clients inside the mesh — carrying a sidecar-issued certificate — can reach the service.

This project uses HTTP application traffic and namespace-wide `STRICT` mode for demonstration purposes only and is not intended for production use.

## Goals
- Install Istio (`base` + `istiod`) via Helm.
- Enable automatic sidecar injection for a namespace.
- Enforce `STRICT` mTLS on inbound traffic to that namespace with `PeerAuthentication`.
- Demonstrate that an in-mesh client is admitted while an out-of-mesh client is rejected.

## Architecture
```
curl (secure ns, sidecar injected)  ──mTLS──►     httpbin (secure ns)
                                                    ▲
curl (default ns, no sidecar)       ──plaintext──┘   PeerAuthentication: STRICT
                                                    └──► rejected (exit code 56)
```

## Repository Structure
- `manifests/httpbin.yaml`: ServiceAccount, Service, and Deployment for the httpbin backend (source: [istio/istio samples](https://github.com/istio/istio/tree/1.30.4/samples)).
- `manifests/curl.yaml`: ServiceAccount, Service, and Deployment for the curl test client (source: [istio/istio samples](https://github.com/istio/istio/tree/1.30.4/samples)).
- `manifests/peer-authentication.yaml`: `PeerAuthentication` enforcing `STRICT` mTLS on inbound traffic to the `secure` namespace.

## Prerequisites
- A Kubernetes cluster
- `kubectl`
- `helm`

## Usage
### 1) Install Istio
```bash
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
```

### 2) Create the `secure` namespace with automatic sidecar injection
```bash
kubectl create namespace secure
kubectl label namespace secure istio-injection=enabled
```

### 3) Deploy the httpbin server and curl clients
```bash
kubectl apply --namespace=secure --filename=./manifests/httpbin.yaml
kubectl apply --namespace=secure --filename=./manifests/curl.yaml
kubectl apply --namespace=default --filename=./manifests/curl.yaml
```

### 4) Enforce mTLS for inbound traffic to the `secure` namespace
```bash
kubectl apply --namespace=secure --filename=./manifests/peer-authentication.yaml
```

## Validation
In-mesh client (sidecar injected, `secure` namespace) reaches httpbin over mTLS:
```bash
kubectl exec --namespace=secure deployment/curl -- curl --silent --output /dev/null --write-out "%{http_code}\n" http://httpbin.secure.svc.cluster.local:8000/headers
# HTTP code 200
```

Out-of-mesh client (no sidecar, `default` namespace) is rejected — plaintext traffic cannot reach a `STRICT` mTLS namespace:
```bash
kubectl exec --namespace=default deployment/curl -- curl --silent --output /dev/null --write-out "%{http_code}\n" http://httpbin.secure.svc.cluster.local:8000/headers
# Exit code 56 (connection reset by peer)
```

## Cleanup
```bash
kubectl delete --namespace=secure --filename=./manifests/peer-authentication.yaml
kubectl delete --namespace=secure --filename=./manifests/curl.yaml
kubectl delete --namespace=secure --filename=./manifests/httpbin.yaml
kubectl delete --namespace=default --filename=./manifests/curl.yaml
kubectl delete namespace secure

helm uninstall istiod --namespace=istio-system
helm uninstall istio-base --namespace=istio-system
kubectl delete namespace istio-system
```
