# Monitoring with Kube Prometheus Stack

This Terraform configuration deploys the [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs) on a Kubernetes cluster. This Terraform configuration:

- Configures the Helm provider to use your local Kubernetes config.
- Installs the `Grafana` Helm chart with pre-configured dashboards.
- Installs the `Alertmanager` Helm chart with pre-configured alerts.
- Installs the `Prometheus Operator` Helm chart with `Node Exporter`.

All Terraform resources are defined in a single file: `main.tf`.

## Components

- **Namespace** `prometheus` (created automatically)
- **Helm Release** `kube-prometheus-stack` (Kube Prometheus Stack)

## Prerequisites

- A Kubernetes cluster (Minikube, EKS, GKE) with an `Ingress Controller`
- A configured `~/.kube/config` file
- [Terraform](https://www.terraform.io/downloads)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Usage

### 1. Initialize Terraform

Run this command to initialize Terraform and download necessary providers and modules:

```bash
export TF_VAR_ingress_hostname=<hostname>

terraform init
```

### 2. Apply Terraform

To deploy the Kubernetes resources, run:

```bash
terraform apply
```

Prometheus should be available at `https://<hostname>/prometheus`  
Alertmanager should be available at `https://<hostname>/alertmanager`  
Grafana should be available at `https://<hostname>/grafana`  

### 3. Cleanup

To delete all resources created by Terraform:

```bash
terraform destroy

unset TF_VAR_ingress_hostname
```