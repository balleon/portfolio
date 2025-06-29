# Purpose of this Go application

This Go program runs a web server that listens for HTTP GET requests at `/version`. When accessed, it connects to a Kubernetes cluster (using your `kubeconfig`) and returns the cluster's Kubernetes version in JSON format.

## Prerequisites

- Access to a Linux or Unix-like terminal
- A Kubernetes cluster (Minikube, EKS, GKE)
- [Go](https://go.dev/doc/install)
- [golangci-lint](https://golangci-lint.run/welcome/install/)

## Usage

### 1. Initializes Go package

Initializes, downloads dependencies and runs a set of checks:

```bash
go mod init github.com/balleon/kubernetes-version
go get .

golangci-lint fmt
golangci-lint run
```

### 2. Run application

Run the application and retrieve the Kubernetes version:

```bash
go run main.go

curl http://<server-private-ip>:8080/version
```