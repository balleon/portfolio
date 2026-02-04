# unused-secret

A Kubernetes utility tool written in Go that identifies and lists unused Secrets across a cluster.

## Overview

This tool scans a Kubernetes cluster to find all Secrets that are not referenced by any workloads. It checks for Secret usage across:

- **Deployments** - mounted volumes and environment variables
- **StatefulSets** - mounted volumes and environment variables
- **DaemonSets** - mounted volumes and environment variables
- **Jobs** - mounted volumes and environment variables
- **CronJobs** - mounted volumes and environment variables
- **Ingresses** - TLS certificates
- **Pod specifications** - image pull secrets, volume mounts, and environment variable references

## Usage

### Prerequisites

- Go 1.x or later
- Kubernetes cluster access configured via kubeconfig
- kubectl CLI (optional, for managing test environments)

### Setup and Run

```bash
# Deploy test Kubernetes resources
kubectl apply -f ./kubernetes/

# Initialize Go module dependencies
go mod tidy

# Run the utility (uses default kubeconfig at ~/.kube/config)
go run main.go

# Or specify a custom kubeconfig
go run main.go -kubeconfig=/path/to/kubeconfig
```

## Output

The tool prints each unused Secret in the format:
```
Secret namespace/secret-name is unused.
```

## Kubernetes Manifests

The `kubernetes/` directory contains sample manifests for testing:
- `namespace.yaml` - Test namespace
- `secrets.yaml` - Test Secrets
- `deployment.yaml` - Test Deployment using a Secret
- `statefulset.yaml` - Test StatefulSet
- `daemonset.yaml` - Test DaemonSet