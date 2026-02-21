# GitOps with Argo CD

This guide demonstrates bootstrapping Argo CD using Helmfile and defining an Argo CD `Application` that deploys NGINX from GitHub. It provisions:

- A **Helm Release** for **Argo CD** with Helm
- A **Kubernetes Deployment** for **NGINX** with Argo CD

## Prerequisites

- Access to a Linux or Unix-like terminal
- A Kubernetes cluster (Minikube, EKS, GKE)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)
- [Helmfile](https://helmfile.readthedocs.io/en/latest/#installation)

## Usage

### 1. Deploy Argo CD

Use Helmfile to bootstrap and manage Argo CD installation in a reproducible, Git-controlled way:

```bash
helmfile --file=helmfile.yaml sync
```

### 2. Deploy NGINX from an Argo CD Application

Create Argo CD Application to deploy NGINX HTTP server from a manifest stored in a GitHub repository:

```bash
kubectl apply --filename=application.yaml
```

### 3. Argo CD Application synchronization

Check Argo CD Application status and NGINX deployment:

```bash
kubectl get application --namespace=argo

kubectl describe application nginx --namespace=argo

kubectl get pods --namespace=nginx
```

### 4. Cleanup

Tear down everything created in this guide:

```bash
kubectl delete --filename=application.yaml

helmfile --file=helmfile.yaml destroy
```