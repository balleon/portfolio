# GitLab on EKS with RDS, ElastiCache and S3

## Overview
This project deploys [GitLab](https://gitlab.com/gitlab-org/charts/gitlab) (Helm chart `10.2.0`) onto an existing EKS cluster using Terraform. It assumes the cluster (e.g. `kubernetes/eks`) already has the AWS Load Balancer Controller, Traefik and the EBS CSI driver installed, and reuses Traefik's existing NLB hostname as GitLab's Ingress host — no owned domain required. PostgreSQL runs on a managed RDS instance, Redis runs on a managed ElastiCache replication group, and object storage uses a set of S3 buckets reached via IRSA. A GitLab Runner (Kubernetes executor) is deployed alongside GitLab, with its runner authentication token provisioned automatically by a bootstrap Job.

GitLab is exposed over plain HTTP for testing purposes only; this is not recommended for production use.

## Security Warning
- Ingress is plain HTTP (`global.hosts.gitlab.https: false`) — a demo simplification, not suitable for production. Without an owned domain there's no name to issue a publicly-trusted TLS certificate for, so HTTPS isn't available here.
- RDS is single-AZ with `skip_final_snapshot = true` and no deletion protection, and ElastiCache runs as a single node (`num_cache_clusters = 1`, no automatic failover) — sized for demo traffic, not production.

## Goals
- Deploy GitLab CE via the official Helm chart, routed through the cluster's existing Traefik ingress controller.
- Replace the chart's bundled PostgreSQL with an RDS instance, and Redis with a managed ElastiCache replication group (the chart dropped its own bundled Redis subchart in v10.0.0), both reached with encryption in transit.
- Use S3 for GitLab's object storage (LFS, artifacts, uploads, packages, external diffs, Terraform state, dependency proxy, backups) via IRSA — no static AWS credentials.
- Deploy a GitLab Runner (Kubernetes executor) via Helm, with its runner authentication token minted automatically by a bootstrap Job that execs into the GitLab toolbox pod and registers the runner via `gitlab-rails runner` — no manual registration step.
- Set explicit CPU/memory requests and limits for every deployed component.

The container registry is disabled/out of scope.

## Repository Structure
- `data.tf`: lookups against the existing EKS cluster (VPC, subnets, security group, OIDC provider) and Traefik's Service (NLB hostname).
- `rds.tf`: RDS PostgreSQL instance, subnet group and security group.
- `redis.tf`: ElastiCache Redis replication group, subnet group and security group.
- `runner.tf`: GitLab Runner Helm release, plus the RBAC and bootstrap Job that provision its runner authentication token from the GitLab toolbox pod.
- `s3.tf`: S3 buckets and the IAM role/policy (IRSA) GitLab uses to reach them.
- `kubernetes.tf`: namespace, the shared ServiceAccount, and the secrets GitLab reads (DB password, Redis password, object storage connection).
- `main.tf`: the `helm_release` for GitLab, with its values built as an HCL map and passed via `yamlencode()`.
- `outputs.tf` / `variables.tf` / `providers.tf` / `versions.tf`: Terraform plumbing.

## Prerequisites
- An existing EKS cluster with the AWS Load Balancer Controller, Traefik and the EBS CSI driver deployed, and a default `StorageClass` backed by it (Gitaly needs a PVC).
- AWS account with permissions to create RDS, S3, IAM and Kubernetes resources
- S3 bucket for Terraform state
- `terraform`
- `aws`
- `kubectl`
- No domain required — GitLab is reached at the hostname of Traefik's existing NLB Service, read via `data.kubernetes_service_v1.traefik`.

## Usage
### 1) Initialize Terraform backend
```bash
export AWS_ACCESS_KEY_ID=<REDACTED>
export AWS_SECRET_ACCESS_KEY=<REDACTED>
export AWS_REGION=<REDACTED>
export TF_VAR_cluster_name=<eks-cluster-name>

terraform init \
-backend-config="bucket=<REDACTED>" \
-backend-config="key=state/$(basename $(pwd))/terraform.tfstate" \
-backend-config="region=${AWS_REGION}"
```

### 2) Deploy GitLab
```bash
terraform apply
```
The GitLab chart deploys `webservice`, `sidekiq`, `gitaly` and `toolbox`; first boot (migrations, initial admin bootstrap) can take several minutes. Once a `toolbox` pod is running, a bootstrap Job execs into it to mint a runner token and installs the GitLab Runner Helm release — the first apply can take a while.

### 3) Get the URL
```bash
terraform output gitlab_url
```

## Validation
```bash
kubectl get pods --namespace=gitlab
kubectl get pods --namespace=gitlab-runner
```

Retrieve the initial `root` password:
```bash
kubectl get secret --namespace=gitlab gitlab-gitlab-initial-root-password -o jsonpath='{.data.password}' | base64 --decode
```

Browse to `$(terraform output -raw gitlab_url)` and sign in as `root`. Confirm object storage is wired up by uploading an attachment/avatar and checking the corresponding S3 bucket:
```bash
terraform output s3_buckets
aws s3 ls "s3://$(terraform output -json s3_buckets | jq -r .uploads)"
```

Confirm the runner registered under **Admin > CI/CD > Runners**, then trigger a pipeline and check it picks up the job.

## Cleanup
```bash
terraform destroy
unset TF_VAR_cluster_name
```
