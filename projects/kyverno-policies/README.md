# Kyverno with policy enforcement

This guide contains the Helm-based deployment of [Kyverno](https://kyverno.io/), a Kubernetes native policy engine, along with custom policies to enforce security and governance standards across your cluster. It provisions:

- A **Helm Release** for **Kyverno** with Helm
- Two policies are included in this setup:
  - **Disallow Privileged Containers**: Prevents the use of privileged containers to enforce Pods security.
  - **Require Label**: Ensures all resources have specific a label for proper organization and governance.

## Prerequisites

- Access to a Linux or Unix-like terminal
- A Kubernetes cluster (Minikube, EKS, GKE)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)

## Installation

### 1. Deploy Kyverno

Deploy Kyverno using Helm:

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

## Test

### 1. Disallow Privileged Containers

Blocks the creation of Pods that run containers with the `privileged: true` setting.

```bash
kubectl apply --filename=policies/disallow-privileged-containers.yaml

kubectl run nginx --namespace=default --image=nginx --privileged=true  # return a message denying the request.
kubectl run nginx --namespace=default --image=nginx --privileged=false # the request complies with the policy.

kubectl delete pod nginx --namespace=default
kubectl delete --filename=policies/disallow-privileged-containers.yaml
```

### 2. Require Label

Ensures Pods include mandatory label key `test`:

```bash
kubectl apply --filename=policies/require-labels.yaml

kubectl run nginx --namespace=default --image=nginx                      # return a message denying the request.
kubectl run nginx --namespace=default --image=nginx --labels="test=true" # request complies with the policy.

kubectl delete pod nginx --namespace=default
kubectl delete --filename=policies/require-labels.yaml
```

## Cleanup

Tear down everything created in this guide:

```bash
helm uninstall kyverno --namespace=kyverno
```