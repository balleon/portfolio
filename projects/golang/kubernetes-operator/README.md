# Kubernetes Operator

## Overview
This folder is the entrypoint documentation for the Kubernetes operator project, implemented in the `app-operator/` subdirectory.

## Goals
- Describe the purpose and capabilities of the operator.
- Point to the operator implementation location.
- Provide standard local and cluster workflows.

## Repository Structure
- `app-operator/`: Kubebuilder-based operator project.

## Prerequisites
- Go `1.22+`
- Docker (or compatible OCI runtime)
- `kubectl` configured to a Kubernetes cluster
- `make`

## Usage
### 1) Move into operator project
```bash
cd app-operator
```

### 2) Run locally
```bash
make install
make run
```

### 3) Deploy in cluster
```bash
make docker-build docker-push IMG=<registry>/app-operator:<tag>
make deploy IMG=<registry>/app-operator:<tag>
```

## Validation
```bash
kubectl get crd apps.apps.test.local
kubectl get pods --namespace=default
```

## Cleanup
```bash
make undeploy
make uninstall
```