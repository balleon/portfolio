# GitLab on EKS with RDS, ElastiCache and S3

## Overview
This project deploys [GitLab](https://gitlab.com/gitlab-org/charts/gitlab) (Helm chart `10.2.0`) onto an existing EKS cluster via Terraform, backed entirely by managed AWS services instead of the chart's bundled, in-cluster PostgreSQL/Redis/MinIO: an RDS PostgreSQL instance, an ElastiCache Redis cluster, and a set of S3 buckets for object storage (artifacts, uploads, LFS, packages, backups, etc.).

It assumes an EKS cluster already exists (e.g. `kubernetes/eks`) with the AWS Load Balancer Controller, Traefik, and the EBS CSI driver already installed. This project does not create the cluster, install an ingress controller, or provision storage classes — it only adds GitLab and its data-plane dependencies on top.

This project uses HTTP for demonstration purposes only and is not intended for production use.

## Goals
- Deploy GitLab CE via the official Helm chart, routed through the cluster's existing Traefik ingress controller.
- Replace the chart's bundled PostgreSQL and Redis with an RDS instance and an ElastiCache cluster.
- Use S3 for GitLab's object storage (LFS, artifacts, uploads, packages, external diffs, Terraform state, dependency proxy, backups, container registry) via IRSA — no static AWS credentials.

## Architecture
```
EKS cluster (existing)
  ├── Traefik (existing, cluster-wide ingress) ◄── Ingress: gitlab.<domain>, registry.<domain>
  │
  └── namespace: gitlab
        ├── webservice / sidekiq / gitaly / toolbox / registry
        │     ├── ServiceAccount (IRSA) ──► IAM role ──► S3 buckets
        │     ├── ── psql (global.psql) ──────────────► RDS PostgreSQL
        │     └── ── redis (global.redis) ────────────► ElastiCache Redis
        └── Secrets: gitlab-postgresql-password, gitlab-rails-storage

AWS
  ├── RDS PostgreSQL      (private subnets of the EKS VPC)
  ├── ElastiCache Redis   (private subnets of the EKS VPC)
  ├── S3 buckets          (artifacts, uploads, lfs, packages, external-diffs,
  │                        terraform-state, dependency-proxy, backups, registry)
  └── IAM role (IRSA)     assumable by system:serviceaccount:gitlab:gitlab
```

Security groups for RDS and ElastiCache only allow inbound traffic from the EKS cluster's security group (`data.aws_eks_cluster.this.vpc_config[0].cluster_security_group_id`), and both are placed in the same private subnets as the cluster.

## Repository Structure
- `data.tf`: lookups against the existing EKS cluster (VPC, subnets, security group, OIDC provider).
- `rds.tf`: RDS PostgreSQL instance, subnet group and security group.
- `elasticache.tf`: ElastiCache Redis cluster, subnet group and security group.
- `s3.tf`: S3 buckets and the IAM role/policy (IRSA) GitLab uses to reach them.
- `kubernetes.tf`: namespace and the two secrets GitLab reads (DB password, object storage connection).
- `main.tf`: the `helm_release` for GitLab, with its values built as an HCL map (`local.gitlab_values`) and passed via `yamlencode()` — external psql/redis/object storage, ingress, disabled bundled subcharts.
- `outputs.tf` / `variables.tf` / `providers.tf` / `versions.tf`: Terraform plumbing.

## Prerequisites
- An existing EKS cluster with the AWS Load Balancer Controller, Traefik and the EBS CSI driver deployed, and a default `StorageClass` backed by it (Gitaly's repository storage needs a PVC).
- `terraform`, `aws` CLI (used by the `exec`-based Kubernes/Helm provider auth via `aws eks get-token`), `kubectl`.
- A DNS name you control, to point `gitlab.<domain>` and `registry.<domain>` at Traefik's Service/LoadBalancer once deployed.
- IAM permissions to create RDS, ElastiCache, S3, IAM and Kubernetes resources.

## Deployment
### 1) Initialize
```bash
export TF_VAR_cluster_name=<eks-cluster-name>
export TF_VAR_domain=<your-domain>
terraform init
```

### 2) Deploy
```bash
terraform apply
```
The GitLab chart deploys `webservice`, `sidekiq`, `gitaly`, `toolbox` and the container registry; first boot (migrations, initial admin bootstrap) can take several minutes.

### 3) Point DNS
Get Traefik's external address and create `gitlab.<domain>` / `registry.<domain>` records (or entries in `/etc/hosts` for local testing) pointing to it:
```bash
kubectl get svc -n traefik traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

## Validation
```bash
kubectl get pods -n gitlab
```

Retrieve the initial `root` password:
```bash
kubectl get secret -n gitlab gitlab-gitlab-initial-root-password -o jsonpath='{.data.password}' | base64 --decode
```

Browse to `http://gitlab.<domain>` and sign in as `root`. Confirm object storage is wired up by uploading an attachment/avatar and checking the corresponding S3 bucket:
```bash
terraform output s3_buckets
aws s3 ls "s3://$(terraform output -json s3_buckets | jq -r .uploads)"
```

## Security Warning
- Ingress is plain HTTP (`global.hosts.gitlab.https: false`) and ElastiCache has no AUTH token / in-transit encryption — both are demo simplifications, not suitable for production.
- RDS and ElastiCache are single-AZ/single-node with `skip_final_snapshot = true` and no deletion protection, to keep `terraform destroy` clean for a demo.

## Cleanup
```bash
terraform destroy
unset TF_VAR_cluster_name TF_VAR_domain
```
