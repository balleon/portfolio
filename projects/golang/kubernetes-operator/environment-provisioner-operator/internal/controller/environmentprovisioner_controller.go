/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	balleonlocalv1alpha1 "github.com/balleon/portfolio/environment-provisioner-operator/api/v1alpha1"
)

const (
	resourceQuotaName = "environment-provisioner-quota"
	provisionerLabel  = "balleon.local/environment-provisioner"
)

// EnvironmentProvisionerReconciler reconciles a EnvironmentProvisioner object
type EnvironmentProvisionerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=balleon.local,resources=environmentprovisioners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=balleon.local,resources=environmentprovisioners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=balleon.local,resources=environmentprovisioners/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete

func (r *EnvironmentProvisionerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	ep := &balleonlocalv1alpha1.EnvironmentProvisioner{}
	if err := r.Get(ctx, req.NamespacedName, ep); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !ep.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	if err := r.reconcileNamespace(ctx, ep); err != nil {
		log.Error(err, "failed to reconcile Namespace")
		return ctrl.Result{}, err
	}

	if err := r.reconcileResourceQuota(ctx, ep); err != nil {
		log.Error(err, "failed to reconcile ResourceQuota")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EnvironmentProvisionerReconciler) reconcileNamespace(ctx context.Context, ep *balleonlocalv1alpha1.EnvironmentProvisioner) error {
	desired := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ep.Spec.NamespaceName,
		},
	}
	if err := ctrl.SetControllerReference(ep, desired, r.Scheme); err != nil {
		return err
	}

	existing := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{Name: ep.Spec.NamespaceName}, existing)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, desired)
	}
	return err
}

func (r *EnvironmentProvisionerReconciler) reconcileResourceQuota(ctx context.Context, ep *balleonlocalv1alpha1.EnvironmentProvisioner) error {
	desired := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceQuotaName,
			Namespace: ep.Spec.NamespaceName,
			Labels: map[string]string{
				provisionerLabel: ep.Name,
			},
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceLimitsCPU:      ep.Spec.ResourcesQuota.Limits.CPU,
				corev1.ResourceLimitsMemory:   ep.Spec.ResourcesQuota.Limits.Memory,
				corev1.ResourceRequestsCPU:    ep.Spec.ResourcesQuota.Request.CPU,
				corev1.ResourceRequestsMemory: ep.Spec.ResourcesQuota.Request.Memory,
			},
		},
	}

	existing := &corev1.ResourceQuota{}
	err := r.Get(ctx, types.NamespacedName{Name: resourceQuotaName, Namespace: ep.Spec.NamespaceName}, existing)
	if apierrors.IsNotFound(err) {
		return r.Create(ctx, desired)
	}
	if err != nil {
		return err
	}

	existing.Spec.Hard = desired.Spec.Hard
	return r.Update(ctx, existing)
}

func (r *EnvironmentProvisionerReconciler) mapResourceQuotaToProvisioner(_ context.Context, obj client.Object) []reconcile.Request {
	if name, ok := obj.GetLabels()[provisionerLabel]; ok {
		return []reconcile.Request{{NamespacedName: types.NamespacedName{Name: name}}}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EnvironmentProvisionerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&balleonlocalv1alpha1.EnvironmentProvisioner{}).
		Owns(&corev1.Namespace{}).
		Watches(
			&corev1.ResourceQuota{},
			handler.EnqueueRequestsFromMapFunc(r.mapResourceQuotaToProvisioner),
		).
		Named("environmentprovisioner").
		Complete(r)
}
