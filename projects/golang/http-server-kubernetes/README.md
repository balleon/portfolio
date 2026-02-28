# Go HTTP Server on Kubernetes

## Overview
This project exposes an HTTP endpoint (`/version`) returning Kubernetes server version information from a Go application deployed to Kubernetes.

## Goals
- Build and lint a Go HTTP service.
- Package the application as a container image.
- Deploy with Kubernetes manifests (Namespace, Deployment, Service, Ingress).

## Repository Structure
- `source/`: Go source code and module.
- `deploy/kubernetes/`: deployment manifests.
- `Dockerfile`: container build definition.

## Prerequisites
- A Kubernetes cluster
- Docker
- Go
- `golangci-lint`
- `kubectl`

## Usage
### 1) Build and lint
```bash
docker build --tag kube-version:latest .

cd source
go get .
golangci-lint fmt
golangci-lint run
cd ..
```

### 2) Deploy manifests
```bash
kubectl apply --filename=./deploy/kubernetes/{namespace.yaml,deployment.yaml,service.yaml,ingress.yaml}
```

## Validation
```bash
curl http://<hostname>/version
```

## Cleanup
```bash
kubectl delete --filename=./deploy/kubernetes/{namespace.yaml,deployment.yaml,service.yaml,ingress.yaml}
```