# Go HTTP server on Kubernetes

This Go program runs a web server that listens for HTTP GET requests at `/version`. When accessed, it connects to a Kubernetes cluster (using your `kubeconfig` or Kubernetes `ServiceAccount`) and returns the cluster's Kubernetes version in JSON format. This guide describes how to build the application's Docker image from the Go source code, then deploy it on Kubernetes.

## Prerequisites

- Access to a Linux or Unix-like terminal
- A Kubernetes cluster (Minikube, EKS, GKE)
- A configured `~/.kube/config` file
- A Docker registry
- [Docker](https://docs.docker.com/engine/install/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Go](https://go.dev/doc/install)
- [golangci-lint](https://golangci-lint.run/welcome/install/)

## Usage

### 1. Initializes Go package

Initializes, downloads dependencies, runs a set of checks and build the Docker image:

```bash
docker build --tag kube-version:latest .

cd ./source

go mod init github.com/balleon/kubernetes-version
go get .

golangci-lint fmt
golangci-lint run

cd ..
```

### 2. Deploy application in Kubernetes cluster

Before creating Kubernetes resources, the `kube-version` Docker image must be available in a Docker repository.  
Create Kubernetes resources at once, including a Namespace, Deployment, Service, and Ingress (`host` need to adapted in `ingress.yaml` file), using their respective YAML configuration files:

```bash
kubectl apply --filename=./deploy/kubernetes/{namespace.yaml,deployment.yaml,service.yaml,ingress.yaml}

curl http://<hostname>/version
```

### 3. Cleanup

Tear down everything created in this guide:

```bash
kubectl delete --filename={namespace.yaml,deployment.yaml,service.yaml,ingress.yaml}
```