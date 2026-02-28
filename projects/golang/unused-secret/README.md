# Unused Secret

## Overview
This Go utility scans a Kubernetes cluster and reports secrets that are not referenced by workloads.

## Goals
- Detect secrets unused by Deployments, StatefulSets, DaemonSets, Jobs, CronJobs, and Ingresses.
- Support default or custom kubeconfig paths.
- Provide a simple CLI workflow for cluster hygiene checks.

## Repository Structure
- `main.go`: unused secret detection logic.
- `kubernetes/`: sample manifests for test scenarios.

## Prerequisites
- Go installed
- Kubernetes access with a valid kubeconfig
- `kubectl` (optional for sample resources)

## Usage
### 1) (Optional) Deploy sample manifests
```bash
kubectl apply --filename=./kubernetes/
```

### 2) Run the scanner
```bash
go mod tidy
go run main.go
```

### 3) Run with a custom kubeconfig
```bash
go run main.go -kubeconfig=/path/to/kubeconfig
```

## Validation
Expected output format:
```text
Secret namespace/secret-name is unused.
```

## Cleanup
```bash
kubectl delete --filename=./kubernetes/
```