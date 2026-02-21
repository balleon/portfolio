# app-operator

A lightweight Kubernetes Operator (Kubebuilder + controller-runtime) written in Go.

The operator manages a custom resource named `App` in API group `apps.test.local/v1`, and reconciles:

- a `Deployment` (container image, replicas, env vars)
- a `Service` (ClusterIP exposing the configured port)
- status fields (`phase`, `readyReplicas`)

## Project structure

This folder documents the operator project located at:

```sh
app-operator/
```

## Prerequisites

- Go `1.22+`
- Docker (or compatible OCI runtime)
- `kubectl` configured for a reachable Kubernetes cluster
- `make`

## How this operator was created

The project was scaffolded with Kubebuilder using the following commands:

```sh
mkdir app-operator && cd app-operator
kubebuilder init --domain=test.local --repo=github.com/balleon/app-operator
kubebuilder create api --group=apps --version=v1 --kind=App --resource=true --controller=true
```

After scaffolding, the following files were customized:

- `app-operator/api/v1/app_types.go`
- `app-operator/internal/controller/app_controller.go`

Then manifests were regenerated:

```sh
make manifests
```

## Quickstart (local development)

From this directory:

```sh
cd app-operator
```

Install CRDs:

```sh
make install
```

Run the controller locally (outside the cluster):

```sh
make run
```

In another terminal, create an `App` resource:

```sh
cat <<EOF | kubectl apply -f -
apiVersion: apps.test.local/v1
kind: App
metadata:
  name: nginx-demo
  namespace: default
spec:
  image: nginx:1.25
  replicas: 1
  port: 80
  env:
    - name: MY_ENV
      value: hello-operator
EOF
```

Verify reconciliation:

```sh
kubectl get app nginx-demo --namespace=default --output=yaml
kubectl get deployment nginx-demo-app --namespace=default
kubectl get service nginx-demo-svc --namespace=default
kubectl get pods --namespace=default --selector=app=nginx-demo
```

## Deploy controller in cluster

Build and push image:

```sh
make docker-build docker-push IMG=<registry>/app-operator:<tag>
```

Deploy manager:

```sh
make deploy IMG=<registry>/app-operator:<tag>
```

## API (`App` spec)

- `spec.image` (string, required)
- `spec.port` (int, required, `1..65535`)
- `spec.replicas` (int, optional, `1..10`)
- `spec.env` (optional list of Kubernetes env vars)

## Cleanup

```sh
kubectl delete app nginx-demo --namespace=default
make undeploy
make uninstall
```