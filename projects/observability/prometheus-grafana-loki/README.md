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

## Grafana Dashboard: SLI/SLO (Availability)
This project includes a dashboard to illustrate SRE concepts using the NGINX HTTP server metrics exposed by `nginx-prometheus-exporter`.

- Dashboard file (embedded): `manifests/grafana-dashboard-configmap.yaml`
- SLI: availability computed from `nginx_up`
- SLO target: shown on dashboard as fixed `99.9%`
- Includes: 1d availability, error budget remaining, 1h burn rate, availability trend, and request rate

### Import via ConfigMap (auto-provisioning)
If you prefer GitOps/Kubernetes-native dashboard provisioning, apply:

```bash
kubectl apply -f manifests/grafana-dashboard-configmap.yaml
```

This ConfigMap is created in namespace `prometheus` with label `grafana_dashboard: "1"`, which is watched by Grafana sidecar in `kube-prometheus-stack`.

### PromQL used for SLI/SLO
- SLI (1d availability): `100 * avg_over_time((avg(nginx_up{namespace="nginx"}) or on() vector(0))[1d:1m])`
- Error budget remaining (%): `100 * clamp_min((avg_over_time((avg(nginx_up{namespace="nginx"}) or on() vector(0))[1d:1m]) - 0.999) / (1 - 0.999), 0)`
- Burn rate (1h): `(1 - avg_over_time((avg(nginx_up{namespace="nginx"}) or on() vector(0))[1h:1m])) / (1 - 0.999)`

## Cleanup
```bash
terraform destroy
unset TF_VAR_ingress_hostname
```