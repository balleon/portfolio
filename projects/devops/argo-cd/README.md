# GitOps with Argo CD

## Overview
This project bootstraps Argo CD with Helmfile and deploys an NGINX workload through an Argo CD `Application` resource.

## Goals
- Install Argo CD in a Kubernetes cluster using declarative Helmfile configuration.
- Deploy an NGINX application from GitOps manifests.
- Validate Argo CD synchronization flow end to end.

## Repository Structure
- `helmfile.yaml`: Argo CD installation definition.
- `application.yaml`: Argo CD `Application` for NGINX deployment.

## Prerequisites
- A Kubernetes cluster (Minikube, EKS, GKE, or equivalent)
- `kubectl`
- `helm`
- `helmfile`

## Usage
### 1) Deploy Argo CD
```bash
helmfile --file=helmfile.yaml sync
```

### 2) Deploy NGINX with Argo CD
```bash
kubectl apply --filename=application.yaml
```

## Validation
```bash
kubectl get application --namespace=argo
kubectl describe application nginx --namespace=argo
kubectl get pods --namespace=nginx
```

## Cleanup
```bash
kubectl delete --filename=application.yaml
helmfile --file=helmfile.yaml destroy
```