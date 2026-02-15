# app-operator

A lightweight Kubernetes Operator (built with kubebuilder) in Go that manages `App` custom resources
(group `apps.test.local`, version `v1`) to deploy and reconcile containerized applications on a cluster.

# Usage

```
mkdir app-operator && cd app-operator
```

```
kubebuilder init --domain=test.local --repo=github.com/balleon/app-operator
```

```
kubebuilder create api --group=apps --version=v1 --kind=App --resource=true --controller=true
```

```
make manifests
```

```
kubectl apply --filename=config/crd/bases/apps.test.local_apps.yaml

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
    value: "hello-operator"
EOF

kubectl get pods --namespace=default

kubectl delete App nginx-demo --namespace=default
```