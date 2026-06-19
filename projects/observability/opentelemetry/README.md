# OpenTelemetry Auto-Instrumentation on Kubernetes

## Overview
This project deploys a full observability stack on Kubernetes and demonstrates zero-code auto-instrumentation of a Python application using the OpenTelemetry Operator.

Traces, metrics, and logs are collected by an `OpenTelemetryCollector` and routed to Tempo, Prometheus, and Loki respectively, then visualized in Grafana.

## Goals
- Deploy Prometheus, Grafana, Loki, and Tempo as the observability backend.
- Install the OpenTelemetry Operator and configure a Collector to receive OTLP signals.
- Auto-instrument a Python Flask application without modifying its source code.
- Route traces to Tempo, metrics to Prometheus, and logs to Loki through a single Collector pipeline.

## Architecture
```
Python Flask app
  │  (auto-instrumented via Annotation — no code changes)
  ▼
OpenTelemetry Collector  (OTLP gRPC :4317 / HTTP :4318)
  ├── traces   ──► Tempo
  ├── metrics  ──► Prometheus  (OTLP receiver)
  └── logs     ──► Loki        (OTLP endpoint)
                    │
                    └── Grafana (datasources: Prometheus, Loki, Tempo)
```

## Repository Structure
- `helmfile.yaml`: Traefik, Prometheus, Grafana, Loki, Tempo, cert-manager, and OpenTelemetry Operator releases.
- `manifests/collector.yaml`: `OpenTelemetryCollector` CR with OTLP receivers and backend exporters.
- `manifests/instrumentation.yaml`: `Instrumentation` CR for Python auto-instrumentation.
- `manifests/deployment.yaml`: demo Python Flask workload with auto-instrumentation annotation.

## Prerequisites
- `k3d`
- `kubectl`
- `helm`
- `helmfile`

## Deployment
### 1) Create the cluster
```bash
k3d cluster create test \
  --port "80:80@loadbalancer" \
  --port "443:443@loadbalancer" \
  --k3s-arg "--disable=traefik@server:*"
```

### 2) Deploy the observability stack
```bash
export INGRESS_HOSTNAME=<hostname>
export GRAFANA_ADMIN_PASSWORD=<password>
helmfile sync
```

### 3) Apply OpenTelemetry resources and the demo workload
```bash
kubectl apply --filename=manifests/collector.yaml
kubectl apply --filename=manifests/instrumentation.yaml
kubectl apply --filename=manifests/deployment.yaml
```

## Validation
Send a request to the instrumented endpoint:
```bash
curl http://${INGRESS_HOSTNAME}/test
```

Access the observability UIs:
- Grafana: `http://${INGRESS_HOSTNAME}/grafana`
- Prometheus: `http://${INGRESS_HOSTNAME}/prometheus` (`count by (__name__) ({job="test/test"})`)
  — counts the number of active time series per metric name for the instrumented app; confirms that OTLP metrics are reaching Prometheus.

Query traces, metrics, and logs in Grafana using the pre-configured Tempo, Prometheus, and Loki datasources.

## Cleanup
```bash
kubectl delete --filename=manifests/

helmfile destroy

kubectl delete pvc --all --namespace=loki

k3d cluster delete test

unset INGRESS_HOSTNAME
unset GRAFANA_ADMIN_PASSWORD
```
