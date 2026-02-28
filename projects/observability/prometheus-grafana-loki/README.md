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

## Grafana Dashboard: SLI, SLO, SLA (demo app)

This repository includes a dashboard focused on SRE concepts for the app in `manifests/`:
- `dashboards/sre-sli-slo-sla-demo.json`

### 1) Validate dashboard JSON (with jq)
```bash
jq empty dashboards/sre-sli-slo-sla-demo.json
```

### 2) Import dashboard in Grafana
1. Open `https://<hostname>/grafana`
2. Go to **Dashboards** -> **New** -> **Import**
3. Upload `dashboards/sre-sli-slo-sla-demo.json`
4. Select the `Prometheus` datasource

### 3) Select service variable
In the dashboard top bar:
- Set **Service** to the Traefik exported service for the demo app: `demo-demo-80@kubernetes`
- Keep **SLO target (%)** at `99.9` (or adjust to show different error-budget behavior)

### 4) Generate traffic for the demo
Run requests continuously to populate SLI/SLO panels:
```bash
while true; do curl -sS http://<hostname>/demo > /dev/null; sleep 0.2; done
```

### 5) SLI / SLO / SLA used in this demo

This demo uses Traefik metrics for the demo backend `demo-demo-80@kubernetes`.

#### SLI (Service Level Indicator)
SLIs are measured values from live traffic:

- **Availability SLI** (success ratio):
	- `good = code !~ "5.."`
	- `total = all HTTP codes`
	- `availability = good / total`
	- Panels: **SLI: Availability (5m)** and **Availability SLI vs SLO Target**

- **Latency SLI**:
	- `p95` and `p99` from `traefik_service_request_duration_seconds_bucket` using `histogram_quantile`
	- Panel: **Latency SLI (p95/p99)**

#### SLO (Service Level Objective)
The SLO is the target objective for the chosen SLI.

- In this dashboard, availability objective is controlled by **SLO target (%)** (default `99.9`).
- Error budget is `1 - SLO` (for `99.9`, budget is `0.1%`).
- Burn rate is:
	- `error_rate / error_budget`
	- where `error_rate = 1 - availability`
	- Panels: **SLO Burn Rate (5m)** and **Multi-Window Burn Rate**

Interpretation:
- Burn rate `1.0` = consuming error budget at planned speed
- Burn rate `> 1.0` = budget is being consumed too fast
- Burn rate `< 1.0` = budget is being consumed slowly

#### SLA (Service Level Agreement)
SLA is the external commitment, usually on longer windows.

- This demo illustrates SLA view with 30-day observed availability:
	- Panel: **SLA Observed Availability (30d)**
- It also shows remaining budget over 30 days:
	- Panel: **Error Budget Remaining (30d)**

Note: this dashboard is for SRE education and visualization. A production SLA typically requires a formal contract definition (scope, exclusions, penalties, and reporting method).

## Cleanup
```bash
terraform destroy
unset TF_VAR_ingress_hostname
```