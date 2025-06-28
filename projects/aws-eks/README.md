# Terraform Amazon EKS Cluster with NGINX

This repository contains Terraform code to deploy a Kubernetes environment on AWS. It provisions:

- A **VPC** with public and private subnets
- An **Amazon EKS Auto Mode** cluster
- A **Kubernetes Deployment**, **Service**, and **Ingress** for **NGINX**
- An **Amazon Application Load Balancer** to expose NGINX publicly

All Terraform resources are defined in a single file: `main.tf`.

## Components

- **Amazon VPC** (via `terraform-aws-modules/vpc`)
- **Amazon EKS** (via `terraform-aws-modules/eks`)
- **Kubernetes provider** configured with EKS outputs
- **NGINX Deployment** in Kubernetes
- **Kubernetes Service** (port 8080 -> 80)
- **Kubernetes Ingress** backed by an **Amazon Application Load Balancer** (⚠️exposed over HTTP only for testing purposes⚠️)
- **IngressClass** for ALB with default behavior

## Prerequisites

- [Terraform](https://www.terraform.io/downloads)
- [AWS CLI](https://aws.amazon.com/cli/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- An Amazon account with permissions to provision VPC, EKS, and Application Load Balancer
- An Amazon S3 bucket to store Terraform state

## Usage

### 1. Initialize Terraform

Run this command to initialize Terraform and download necessary providers and modules:

```bash
export AWS_ACCESS_KEY_ID=<REDACTED>
export AWS_SECRET_ACCESS_KEY=<REDACTED>
export AWS_REGION=<REDACTED>

terraform init \
-backend-config="bucket=<REDACTED>" \
-backend-config="key=state/terraform.tfstate" \
-backend-config="region=${AWS_REGION}"
```

### 2. Apply Terraform

To deploy the infrastructure and Kubernetes resources, run:

```bash
terraform apply
```

### 3. Kubernetes context

To configure your local kubectl to connect to the cluster, run:

```bash
aws eks update-kubeconfig --name <cluster_name>
```

### 4. Accessing NGINX

Get the Amazon Application Load Balancer address then test it:

```bash
kubectl get ingress --namespace=nginx nginx

curl http://<alb-address>
```

### 5. Cleanup

To delete all resources created by Terraform:

```bash
terraform destroy
```