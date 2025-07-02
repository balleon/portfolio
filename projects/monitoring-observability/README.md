# Monitoring and observability with Prometheus, Grafana and Loki

This guide shows how to deploy and use `Monitoring` and `Observability` tools.

This Terraform configuration deploys the [Kube-Prometheus-Stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack), [Loki](https://github.com/grafana/loki/tree/main/production/helm/loki) and [Promtail](https://github.com/grafana/helm-charts/tree/main/charts/promtail) using the [Helm provider](https://registry.terraform.io/providers/hashicorp/helm/latest/docs) on a Kubernetes cluster and configures:

- Configures the `Helm provider` to use your local Kubernetes config.
- Installs the `Grafana` Helm chart with pre-configured dashboards and datasources.
- Installs the `Alertmanager` Helm chart with pre-configured alerts.
- Installs the `Prometheus Operator` Helm chart with `Node Exporter`.
- Installs the `Loki` Helm chart.
- Installs the `Promtail` Helm chart.

All Terraform resources are defined in a single file: `main.tf`.

## Components

- **Namespace**
    - `prometheus` (created automatically)
    - `loki` (created automatically)
    - `promtail` (created automatically)
- **Helm Releases** 
    - `kube-prometheus-stack`
    - `loki`
    - `promtail`

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
Grafana should be available at `https://<hostname>/grafana` with `Prometheus` and `Loki` datasources  

### 3. Cleanup

To delete all resources created by Terraform:

```bash
terraform destroy

unset TF_VAR_ingress_hostname
```