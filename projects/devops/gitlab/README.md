# GitLab on EKS with RDS and S3

## Overview
This project deploys [GitLab](https://gitlab.com/gitlab-org/charts/gitlab) (Helm chart `10.2.0`) onto an existing EKS cluster via Terraform. The database is backed by a managed RDS PostgreSQL instance, object storage uses a set of S3 buckets (artifacts, uploads, LFS, packages, backups, etc.), and Redis stays as the chart's bundled, in-cluster subchart. GitLab is served on the hostname of the cluster's existing Traefik NLB — no owned domain required.

It assumes an EKS cluster already exists (e.g. `kubernetes/eks`) with the AWS Load Balancer Controller, Traefik, and the EBS CSI driver already installed. This project does not create the cluster, install an ingress controller, or provision storage classes — it only adds GitLab and its data-plane dependencies on top.

This project uses HTTP for demonstration purposes only and is not intended for production use.

## Goals
- Deploy GitLab CE via the official Helm chart, routed through the cluster's existing Traefik ingress controller.
- Replace the chart's bundled PostgreSQL with an RDS instance; leave the chart's bundled Redis subchart in place.
- Use S3 for GitLab's object storage (LFS, artifacts, uploads, packages, external diffs, Terraform state, dependency proxy, backups) via IRSA — no static AWS credentials.

The container registry subchart is disabled (`registry.enabled: false`) — out of scope for this demo.

## Architecture
```
EKS cluster (existing)
  ├── Traefik (existing, cluster-wide ingress, fronted by an NLB)
  │     ◄── Ingress: <traefik-nlb-hostname> (read via data.kubernetes_service_v1.traefik)
  │
  └── namespace: gitlab
        ├── webservice / sidekiq / gitaly / toolbox
        │     ├── ServiceAccount (IRSA) ──► IAM role ──► S3 buckets
        │     └── ── psql (global.psql) ──────────────► RDS PostgreSQL
        ├── redis (bundled subchart, in-cluster, PVC via existing StorageClass)
        └── Secrets: gitlab-postgresql-password, gitlab-rails-storage

AWS
  ├── RDS PostgreSQL      (private subnets of the EKS VPC)
  ├── S3 buckets          (artifacts, uploads, lfs, packages, external-diffs,
  │                        terraform-state, dependency-proxy, backups)
  └── IAM role (IRSA)     assumable by system:serviceaccount:gitlab:gitlab
```

The RDS security group only allows inbound traffic from the EKS cluster's security group (`data.aws_eks_cluster.this.vpc_config[0].cluster_security_group_id`), and the instance is placed in the same private subnets as the cluster.

## Repository Structure
- `data.tf`: lookups against the existing EKS cluster (VPC, subnets, security group, OIDC provider) and Traefik's Service (NLB hostname).
- `rds.tf`: RDS PostgreSQL instance, subnet group and security group.
- `s3.tf`: S3 buckets and the IAM role/policy (IRSA) GitLab uses to reach them.
- `kubernetes.tf`: namespace and the two secrets GitLab reads (DB password, object storage connection).
- `main.tf`: the `helm_release` for GitLab, with its values built as an HCL map and passed via `yamlencode()` — external psql/object storage, ingress, disabled bundled subcharts (Redis stays enabled).
- `outputs.tf` / `variables.tf` / `providers.tf` / `versions.tf`: Terraform plumbing.

## Prerequisites
- An existing EKS cluster with the AWS Load Balancer Controller, Traefik and the EBS CSI driver deployed, and a default `StorageClass` backed by it (Gitaly and the bundled Redis subchart both need a PVC).
- `terraform`, `aws` CLI (used by the `exec`-based Kubernes/Helm provider auth via `aws eks get-token`), `kubectl`.
- No domain required — GitLab is reached at the hostname of Traefik's existing NLB Service, read via `data.kubernetes_service_v1.traefik`.
- IAM permissions to create RDS, S3, IAM and Kubernetes resources.

## Deployment
### 1) Initialize
```bash
export TF_VAR_cluster_name=<eks-cluster-name>
terraform init
```

### 2) Deploy
```bash
terraform apply
```
The GitLab chart deploys `webservice`, `sidekiq`, `gitaly` and `toolbox`; first boot (migrations, initial admin bootstrap) can take several minutes.

### 3) Get the URL
```bash
terraform output gitlab_url
```

## Validation
```bash
kubectl get pods -n gitlab
```

Retrieve the initial `root` password:
```bash
kubectl get secret -n gitlab gitlab-gitlab-initial-root-password -o jsonpath='{.data.password}' | base64 --decode
```

Browse to `$(terraform output -raw gitlab_url)` and sign in as `root`. Confirm object storage is wired up by uploading an attachment/avatar and checking the corresponding S3 bucket:
```bash
terraform output s3_buckets
aws s3 ls "s3://$(terraform output -json s3_buckets | jq -r .uploads)"
```

## Security Warning
- Ingress is plain HTTP (`global.hosts.gitlab.https: false`) — a demo simplification, not suitable for production. Without an owned domain there's no name to issue a publicly-trusted TLS certificate for, so HTTPS isn't available here.
- RDS is single-AZ with `skip_final_snapshot = true` and no deletion protection, to keep `terraform destroy` clean for a demo. The bundled Redis subchart also runs as a single, unclustered instance.

## Cleanup
```bash
terraform destroy
unset TF_VAR_cluster_name
```
