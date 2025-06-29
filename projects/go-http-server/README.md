# Purpose of this Go application

This Go program runs a web server that listens for HTTP GET requests at `/version`. When accessed, it connects to a Kubernetes cluster (using your `kubeconfig`) and returns the cluster's Kubernetes version in JSON format.

## Prerequisites

- Access to a Linux or Unix-like terminal
- A Kubernetes cluster (Minikube, EKS, GKE)
- A configured `~/.kube/config` file
- [Docker](https://docs.docker.com/engine/install/)
- [Go](https://go.dev/doc/install)
- [golangci-lint](https://golangci-lint.run/welcome/install/)

## Usage

### 1. Initializes Go package

Initializes, downloads dependencies and runs a set of checks:

```bash
cd ./source

go mod init github.com/balleon/kubernetes-version
go get .

golangci-lint fmt
golangci-lint run
```

### 2. Run the application

Run the application locally and retrieve the Kubernetes version:

```bash
cd ./source

go run main.go

curl http://<server-private-ip>:8080/version
```

Run the application inside a container and retrieve the Kubernetes version:

```bash
docker build --tag kube-version:latest .

docker run \
--name kube-version \
--detach \
--publish 8080:8080 \
--volume ${HOME}/.kube/config:/root/.kube/config:ro \
kube-version:latest

curl http://<server-private-ip>:8080/version

docker rm --force kube-version
```