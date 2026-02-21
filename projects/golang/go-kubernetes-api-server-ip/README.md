# go-kubernetes-api-server-ip

This project uses `client-go` to connect to a Kubernetes cluster and print the IP address of the Kubernetes API Server.

## How it works

- Builds Kubernetes client config from your local kubeconfig (default: `~/.kube/config`)
- Queries `discovery.k8s.io/v1` EndpointSlices in `default`
- Filters by label `kubernetes.io/service-name=kubernetes`
- Prints the endpoint IP address(es), which correspond to the Kubernetes API Server

## Prerequisites

- Go installed
- Access to a Kubernetes cluster
- Valid kubeconfig file (default path: `~/.kube/config`)

## Run

```bash
go mod tidy
go run main.go
```

Or with an explicit kubeconfig path:

```bash
go run main.go -kubeconfig /path/to/kubeconfig
```

## Example output

```text
[10.96.0.1]
```