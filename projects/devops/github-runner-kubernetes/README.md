# Terraform GitHub Actions Runner Controller Deployment

## Overview
This project deploys GitHub Actions Runner Controller (ARC) on Kubernetes with Terraform and Helm.

## Goals
- Install `gha-runner-scale-set-controller`.
- Install `gha-runner-scale-set` for self-hosted runners.
- Manage the full lifecycle with Terraform.

## Repository Structure
- `main.tf`: Terraform resources for providers and Helm releases.
- `variables.tf`: Input variables (including GitHub settings).
- `outputs.tf`: Terraform outputs.
- `versions.tf`: Provider and Terraform version constraints.

## Prerequisites
- A Kubernetes cluster with access from your workstation
- `~/.kube/config` configured
- GitHub token with required permissions
- `terraform`
- `kubectl`

## Usage
### 1) Initialize
```bash
export TF_VAR_github_url=<REDACTED>
export TF_VAR_github_token=<REDACTED>

terraform init
```

### 2) Deploy
```bash
terraform apply
```

## Validation
```bash
kubectl get pods --namespace=arc-systems
kubectl get deployment --namespace=arc-systems
```

GitHub runners are visible in your repository settings under Actions runners.

## Cleanup
```bash
terraform destroy
unset TF_VAR_github_url TF_VAR_github_token
```