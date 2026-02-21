# Terraform Amazon EKS Cluster with Ingress NGINX Controller

This repository contains Terraform code to deploy a Kubernetes environment on AWS. It provisions:

- A **VPC** with public and private subnets
- An **Amazon EKS Auto Mode** cluster
- An **Helm Release** for **Ingress NGINX Controller**

All Terraform resources are defined in a single file: `main.tf`.

## Components

- **Amazon VPC** (via `terraform-aws-modules/vpc`)
- **Amazon EKS** (via `terraform-aws-modules/eks`)
- **Helm provider** configured with EKS outputs
- **Helm Releases** 
    - `Ingress NGINX Controller for Kubernetes`

## Prerequisites

- [Terraform](https://www.terraform.io/downloads)
- [AWS CLI](https://aws.amazon.com/cli/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- An Amazon account with permissions to provision VPC, EKS, and Network Load Balancer
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

### 4. Accessing NGINX Network Load Balancer

Get the Amazon Network Load Balancer address then test it:

```bash
kubectl get service --namespace=ingress-nginx

curl http://<nlb-address>
curl https://<nlb-address> --insecure
```

### 5. Cleanup

To delete all resources created by Terraform:

```bash
helm uninstall ingress-nginx --namespace=ingress-nginx

terraform destroy
```