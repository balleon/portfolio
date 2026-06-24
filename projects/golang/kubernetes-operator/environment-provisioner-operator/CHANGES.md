# Implementation Changes

## Overview

Two files were updated from Kubebuilder scaffolding to a working reconciler that provisions a `Namespace` and a `ResourceQuota` from a single `EnvironmentProvisioner` custom resource.

---

## `api/v1alpha1/environmentprovisioner_types.go`

### What changed

Replaced the placeholder `Foo *string` field with the real domain types.

**Added types:**

```go
type ResourceList struct {
    CPU    resource.Quantity `json:"cpu"`
    Memory resource.Quantity `json:"memory"`
}

type ResourcesQuota struct {
    Limits  ResourceList `json:"limits"`
    Request ResourceList `json:"request"`
}
```

**Updated `EnvironmentProvisionerSpec`:**

```go
type EnvironmentProvisionerSpec struct {
    NamespaceName  string         `json:"namespaceName"`
    ResourcesQuota ResourcesQuota `json:"resourcesQuota"`
}
```

`api/v1alpha1/zz_generated.deepcopy.go` was regenerated automatically via `make generate` to add `DeepCopyInto` / `DeepCopy` for `ResourceList` and `ResourcesQuota`.

---

## `internal/controller/environmentprovisioner_controller.go`

### What changed

Replaced the empty `Reconcile` stub with a full implementation.

### New constants

| Constant | Value | Purpose |
|---|---|---|
| `resourceQuotaName` | `"environment-provisioner-quota"` | Fixed name for the `ResourceQuota` in the provisioned namespace |
| `provisionerLabel` | `"balleon.local/environment-provisioner"` | Label key used to link a `ResourceQuota` back to its owning CR |

### New RBAC markers

```go
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
```

### Reconcile logic

```
Reconcile(req)
  ├── Fetch EnvironmentProvisioner CR → return if NotFound or being deleted
  ├── reconcileNamespace
  │     ├── Build desired Namespace with SetControllerReference (owner ref → CR)
  │     └── Create if missing, skip if already exists
  └── reconcileResourceQuota
        ├── Build desired ResourceQuota with label provisionerLabel=<cr-name>
        ├── Create if missing
        └── Update spec.Hard if already exists
```

### Drift detection

| Manually deleted | Trigger | Effect |
|---|---|---|
| `Namespace` | `Owns(&corev1.Namespace{})` watch enqueues the owning CR | Namespace (and its ResourceQuota) are recreated |
| `ResourceQuota` | `Watches(&corev1.ResourceQuota{}, mapResourceQuotaToProvisioner)` enqueues the CR by label | ResourceQuota is recreated |
| `Namespace` spec edited | same `Owns` watch | No-op on namespace spec (immutable), ResourceQuota is re-synced |
| `ResourceQuota` spec edited | same label watch | `spec.Hard` is overwritten to match the CR |

### Deletion

`SetControllerReference` sets an owner reference from the `Namespace` to the `EnvironmentProvisioner` CR. When the CR is deleted, the Kubernetes garbage collector removes the `Namespace`, which cascade-deletes the `ResourceQuota` inside it. No finalizer is required.

### `SetupWithManager`

```go
ctrl.NewControllerManagedBy(mgr).
    For(&EnvironmentProvisioner{}).
    Owns(&corev1.Namespace{}).                                       // watches owned Namespaces
    Watches(&corev1.ResourceQuota{}, EnqueueRequestsFromMapFunc(…)). // watches labeled ResourceQuotas
    Complete(r)
```

---

## `config/samples/v1alpha1_environmentprovisioner.yaml`

Updated the sample to match the CR described in the README:

```yaml
spec:
  namespaceName: test
  resourcesQuota:
    limits:
      cpu: "2"
      memory: 1Gi
    request:
      cpu: "1"
      memory: 1Gi
```
