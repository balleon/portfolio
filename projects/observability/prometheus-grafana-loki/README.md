# Monitoring and Observability with Prometheus, Grafana, and Loki

## Overview
This project deploys a Kubernetes observability stack using Terraform and Helm: kube-prometheus-stack, Loki, and Promtail.

## Goals
- Deploy Prometheus, Alertmanager, and Grafana.
- Deploy Loki and Promtail for log collection and querying.
- Expose observability tools behind a configured ingress hostname.

## Repository Structure
- `main.tf`: Helm releases and related resources.
- `variables.tf`: input variables such as ingress hostname.
- `outputs.tf`: Terraform outputs.
- `versions.tf`: Terraform/provider version constraints.

## Prerequisites
- Kubernetes cluster with an ingress controller
- Valid `~/.kube/config`
- `terraform`
- `kubectl`

## Usage
### 1) Initialize
```bash
export TF_VAR_ingress_hostname=<hostname>
terraform init
```

### 2) Deploy
```bash
terraform apply
```

## Validation
After deployment, verify endpoints:
- `https://<hostname>/prometheus`
- `https://<hostname>/alertmanager`
- `https://<hostname>/grafana`

## Cleanup
```bash
terraform destroy
unset TF_VAR_ingress_hostname
```