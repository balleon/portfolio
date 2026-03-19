# Kagent + Ollama: OOMKilled Pod Remediation

## Overview
This project demonstrates AI-assisted Kubernetes operations using:
- Official `kagent` Helm charts: https://github.com/kagent-dev/kagent/tree/main/helm
- Official `ollama` Helm chart: https://github.com/otwld/ollama-helm

Scenario: a pod crashes with `OOMKilled`, and kagent (backed by a free local Ollama model) helps diagnose and remediate the issue.

## Architecture
- `ollama` runs in namespace `ollama` with a local model (`llama3.2`) pulled at startup.
- `kagent` runs in namespace `kagent` with provider default set to Ollama.
- A custom `ModelConfig` and `Agent` CR configure kagent for this use case.
- A demo workload in namespace `oom-demo` is intentionally under-provisioned to trigger OOM kills.

## Repository Structure
- `manifests/app/namespace.yaml`: namespace for the failing workload.
- `manifests/app/deployment.yaml`: intentionally memory-starved workload.
- `manifests/kagent/modelconfig-ollama.yaml`: kagent `ModelConfig` pointing to Ollama service.
- `manifests/kagent/agent-oom-remediator.yaml`: kagent `Agent` CR for OOM incident response.

## Prerequisites
- Kubernetes cluster access
- `kubectl`
- `helm`
- `kagent` CLI

## Deployment
### 1) Install Ollama (official chart)
```bash
helm repo add otwld https://helm.otwld.com/
helm repo update

helm upgrade --install ollama otwld/ollama \
--version=1.52.0 \
--namespace=ollama \
--create-namespace \
--set=ollama.models.pull[0]=llama2 \
--set=resources.requests.cpu=500m \
--set=resources.requests.memory=1Gi \
--set=resources.limits.cpu=2 \
--set=resources.limits.memory=4Gi
```

### 2) Install kagent CRDs and kagent (official charts)
```bash
helm upgrade --install kagent-crds \
oci://ghcr.io/kagent-dev/kagent/helm/kagent-crds \
--version=0.7.23 \
--namespace=kagent \
--create-namespace

helm upgrade --install kagent \
oci://ghcr.io/kagent-dev/kagent/helm/kagent \
--version=0.7.23 \
--namespace=kagent \
--create-namespace \
--set=providers.default=ollama
```

### 3) Configure kagent using Custom Resources
```bash
kubectl apply -f ./manifests/kagent/modelconfig-ollama.yaml
kubectl apply -f ./manifests/kagent/agent-oom-remediator.yaml
```

### 4) Deploy the failing workload
```bash
kubectl apply -f ./manifests/app/namespace.yaml
kubectl apply -f ./manifests/app/deployment.yaml
```

## Trigger and Observe the Incident
```bash
kubectl get pods -n oom-demo
kubectl describe pod -n oom-demo -l app=oom-demo
kubectl get events -n oom-demo --sort-by=.lastTimestamp
```

Expected signal: pod restarts with `Last State: Terminated` and `Reason: OOMKilled`.

## Use kagent CLI for Incident Automation
```bash
kagent get agent --namespace kagent

kagent invoke \
  --agent oom-remediator \
  --task "Investigate OOMKilled pods in namespace oom-demo, identify root cause, and provide kubectl commands to remediate safely."
```

## Apply Remediation (example)
If the agent recommends raising memory limits/requests for `oom-demo`, apply:
```bash
kubectl set resources deployment/oom-demo -n oom-demo \
  --requests=cpu=100m,memory=128Mi \
  --limits=cpu=200m,memory=256Mi

kubectl rollout status deployment/oom-demo -n oom-demo
kubectl get pods -n oom-demo
```

## Validation
```bash
kubectl get deployment oom-demo -n oom-demo -o wide
kubectl describe pod -n oom-demo -l app=oom-demo
```

Success criteria:
- Pod is `Running`
- No new `OOMKilled` events after rollout

## Cleanup
```bash
kubectl delete -f ./manifests/app/deployment.yaml
kubectl delete -f ./manifests/app/namespace.yaml

kubectl delete -f ./manifests/kagent/agent-oom-remediator.yaml
kubectl delete -f ./manifests/kagent/modelconfig-ollama.yaml

helm uninstall kagent -n kagent
helm uninstall kagent-crds -n kagent
helm uninstall ollama -n ollama
```

## Notes
- This project intentionally avoids Helm values files and uses `--set` only.
- kagent configuration is performed through Kubernetes Custom Resources.
- Because chart defaults can change, confirm final settings with:
```bash
helm show values oci://ghcr.io/kagent-dev/kagent/helm/kagent
helm show values otwld/ollama
```
