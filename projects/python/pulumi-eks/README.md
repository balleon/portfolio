# Pulumi Amazon EKS Cluster with Python

## Overview
This project provisions AWS networking and an Amazon EKS cluster using Pulumi and Python. It demonstrates infrastructure-as-code practices for reproducible, cloud-native platform deployment.

## Security Warning
This example provisions EKS with public endpoint access for demonstration purposes. For production, restrict public access and use VPN or private connectivity. Additionally, EKS nodes and control plane communication should be encrypted and monitored in production environments.

## Goals
- Provision VPC with public and private subnets across multiple availability zones.
- Deploy EKS cluster with auto-scaling node pools.
- Demonstrate Pulumi infrastructure-as-code patterns in Python.
- Use configuration management for environment-specific settings.

## Repository Structure
- `__main__.py`: Pulumi program defining AWS infrastructure.
- `Pulumi.yaml`: Project metadata and runtime configuration.
- `requirements.txt`: Python dependencies (Pulumi, AWS provider).

## Prerequisites
- AWS account with permissions for VPC, EKS, and IAM resources.
- Pulumi CLI installed and authenticated.
- AWS credentials configured (via AWS CLI or environment variables).
- Python 3.8 or later.

## Usage

### 1) Initialize stack
```bash
export AWS_ACCESS_KEY_ID="<REDACTED>"
export AWS_SECRET_ACCESS_KEY="<REDACTED>"
export PULUMI_CONFIG_PASSPHRASE=""

pulumi login --local
pulumi stack init dev
```

### 2) Set stack configuration
```bash
pulumi config set pulumi-eks:eks_version "1.35"
pulumi config set aws:region eu-west-3
pulumi config set pulumi-eks:vpc_cidr "10.0.0.0/16"
pulumi config set --path pulumi-eks:private_subnet_cidrs[0] "10.0.0.0/24"
pulumi config set --path pulumi-eks:private_subnet_cidrs[1] "10.0.1.0/24"
pulumi config set --path pulumi-eks:private_subnet_cidrs[2] "10.0.2.0/24"
pulumi config set --path pulumi-eks:public_subnet_cidrs[0] "10.0.3.0/24"
pulumi config set --path pulumi-eks:public_subnet_cidrs[1] "10.0.4.0/24"
pulumi config set --path pulumi-eks:public_subnet_cidrs[2] "10.0.5.0/24"
pulumi config set --path pulumi-eks:availability_zones[0] "eu-west-3a"
pulumi config set --path pulumi-eks:availability_zones[1] "eu-west-3b"
pulumi config set --path pulumi-eks:availability_zones[2] "eu-west-3c"
```

### 4) Install dependencies
```bash
pulumi install
```

### 5) Preview infrastructure
```bash
pulumi preview
```

### 6) Deploy infrastructure
```bash
pulumi up
```

## Validation
After deployment completes, configure kubectl to access the cluster:
```bash
aws eks update-kubeconfig --name eks-pulumi-eks --region eu-west-3
kubectl get nodes
kubectl get pods --all-namespaces
```

## Cleanup
To destroy all resources:
```bash
pulumi destroy
pulumi stack rm dev
```