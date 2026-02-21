# Terraform GitHub Actions Runner Controller Deployment

This repository contains Terraform code to deploy the [GitHub Actions Runner Controller (ARC)](https://github.com/actions/actions-runner-controller) on a Kubernetes cluster using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs). This Terraform configuration:

- Configures the Helm provider to use your local Kubernetes config.
- Installs the `gha-runner-scale-set-controller` Helm chart (ARC Controller).
- Installs the `gha-runner-scale-set` Helm chart to manage self-hosted GitHub Actions runners.

All Terraform resources are defined in a single file: `main.tf`.

## Components

- **Namespace** `arc-systems` (created automatically)
- **Helm Releases** 
    - `gha-runner-scale-set-controller` (ARC controller)
    - `gha-runner-scale-set` (Runner scale set)

## Prerequisites

- A Kubernetes cluster (Minikube, EKS, GKE)
- A configured `~/.kube/config` file
- GitHub personal access token with appropriate permissions
- [Terraform](https://www.terraform.io/downloads)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Usage

### 1. Initialize Terraform

Run this command to initialize Terraform and download necessary providers and modules:

```bash
export TF_VAR_github_url=<REDACTED>
export TF_VAR_github_token=<REDACTED>

terraform init
```

### 2. Apply Terraform

To deploy the Kubernetes resources, run:

```bash
terraform apply
```

The GitHub runner should be available at `https://github.com/<repository>/settings/actions/runners`.

### 3. Cleanup

To delete all resources created by Terraform:

```bash
terraform destroy

unset TF_VAR_github_url TF_VAR_github_token
```