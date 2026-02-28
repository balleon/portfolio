# Kubernetes API Server IP

## Overview
This project uses `client-go` to discover and print the Kubernetes API server endpoint IP address from `EndpointSlice` resources.

## Goals
- Build Kubernetes client configuration from kubeconfig.
- Query and filter `EndpointSlice` objects for the `kubernetes` service.
- Print API server endpoint IPs for troubleshooting and diagnostics.

## Repository Structure
- `main.go`: CLI entrypoint and API server IP lookup logic.
- `go.mod`: module dependencies.

## Prerequisites
- Go installed
- Access to a Kubernetes cluster
- Valid kubeconfig (default `~/.kube/config`)

## Usage
### 1) Run with default kubeconfig
```bash
go mod tidy
go run main.go
```

### 2) Run with explicit kubeconfig
```bash
go run main.go -kubeconfig /path/to/kubeconfig
```

## Validation
Expected output format:
```text
[10.96.0.1]
```

## Cleanup
No cluster resources are created by this tool.