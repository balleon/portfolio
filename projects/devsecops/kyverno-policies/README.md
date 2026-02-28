# Kyverno with Policy Enforcement

## Overview
This project deploys Kyverno and provides sample policies for baseline Kubernetes policy enforcement.

## Goals
- Install Kyverno in a dedicated namespace.
- Enforce non-privileged container execution.
- Enforce required labels on resources.

## Repository Structure
- `policies/disallow-privileged-containers.yaml`: blocks privileged containers.
- `policies/require-labels.yaml`: requires label key `test`.

## Prerequisites
- A Kubernetes cluster
- `kubectl`
- `helm`

## Usage
### 1) Install Kyverno
```bash
helm repo add kyverno https://kyverno.github.io/kyverno/
helm repo update

helm install kyverno kyverno/kyverno \
  --version=3.4.4 \
  --namespace=kyverno \
  --create-namespace \
  --wait \
  --wait-for-jobs
```

### 2) Apply Policies
```bash
kubectl apply --filename=policies/disallow-privileged-containers.yaml
kubectl apply --filename=policies/require-labels.yaml
```

## Validation
```bash
kubectl run nginx --namespace=default --image=nginx --privileged=true
kubectl run nginx --namespace=default --image=nginx --labels="test=true"
```

## Cleanup
```bash
kubectl delete --filename=policies/disallow-privileged-containers.yaml
kubectl delete --filename=policies/require-labels.yaml
helm uninstall kyverno --namespace=kyverno
```