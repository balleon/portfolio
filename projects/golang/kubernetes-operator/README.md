# Environment Provisioner Operator

A Kubernetes operator built with [Kubebuilder](https://book.kubebuilder.io/) that provisions a `Namespace` and a `ResourceQuota` from a single custom resource.

## Initialization

```bash
mkdir environment-provisioner-operator && cd environment-provisioner-operator

kubebuilder init --domain=balleon.local --repo=github.com/balleon/portfolio/environment-provisioner-operator
kubebuilder create api --version=v1alpha1 --kind=EnvironmentProvisioner --namespaced=false --resource=true --controller=true
```

## Implementation

See [environment-provisioner-operator/CHANGES.md](environment-provisioner-operator/CHANGES.md) for a description of the types, reconcile logic, drift detection, and deletion behaviour.

## Apply
```bash
make manifests

kubectl apply -f config/crd/bases/balleon.local_environmentprovisioners.yaml

go run cmd/main.go

kubectl apply -f config/samples/v1alpha1_environmentprovisioner.yaml
```