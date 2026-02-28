# Terraform Amazon EKS Cluster with Ingress NGINX Controller

## Overview
This project provisions AWS networking, an Amazon EKS cluster, and the NGINX Ingress Controller using Terraform.

## Goals
- Create a VPC with public and private subnets.
- Provision EKS and configure access.
- Deploy NGINX Ingress Controller through Helm.

## Repository Structure
- `main.tf`: core infrastructure and Helm resources.
- `variables.tf`: project input variables.
- `outputs.tf`: exported values.
- `providers.tf`, `terraform.tf`: provider and backend settings.

## Prerequisites
- AWS account with permissions for VPC/EKS/NLB
- S3 bucket for Terraform state
- `terraform`
- `aws`
- `kubectl`

## Usage
### 1) Initialize Terraform backend
```bash
export AWS_ACCESS_KEY_ID=<REDACTED>
export AWS_SECRET_ACCESS_KEY=<REDACTED>
export AWS_REGION=<REDACTED>

terraform init \
    -backend-config="bucket=<REDACTED>" \
    -backend-config="key=state/terraform.tfstate" \
    -backend-config="region=${AWS_REGION}"
```

### 2) Deploy infrastructure
```bash
terraform apply
```

### 3) Configure kubeconfig
```bash
aws eks update-kubeconfig --name <cluster_name>
```

## Validation
```bash
kubectl get service --namespace=ingress-nginx
curl http://<nlb-address>
curl https://<nlb-address> --insecure
```

## Cleanup
```bash
helm uninstall ingress-nginx --namespace=ingress-nginx
terraform destroy
```