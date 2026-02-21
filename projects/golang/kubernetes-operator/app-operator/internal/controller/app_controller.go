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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "github.com/balleon/app-operator/api/v1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.test.local,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.test.local,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.test.local,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
// func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
// 	_ = log.FromContext(ctx)
//
// 	// TODO(user): your logic here
//
// 	return ctrl.Result{}, nil
// }

func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. Fetch the App CR
	app := &appv1.App{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Reconcile Deployment
	dep := r.desiredDeployment(app)
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, func() error {
		// Mutate: set desired spec (idempotent)
		dep.Spec.Replicas = app.Spec.Replicas // assumes non-nil or you add default logic
		dep.Spec.Template.Spec.Containers[0].Image = app.Spec.Image
		dep.Spec.Template.Spec.Containers[0].Env = app.Spec.Env
		// You can add more (resources, probes, etc.) later
		return nil
	})
	if err != nil {
		log.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}
	log.Info("Deployment reconciled", "operation", op, "name", dep.Name)

	// 3. Reconcile Service
	svc := r.desiredService(app)
	op, err = controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		svc.Spec.Ports[0].Port = app.Spec.Port
		svc.Spec.Ports[0].TargetPort = intstr.FromInt32(app.Spec.Port)
		// Add more ports or type change if needed
		return nil
	})
	if err != nil {
		log.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}
	log.Info("Service reconciled", "operation", op, "name", svc.Name)

	// 4. Update status (simple version; improve with conditions later)
	// For better status, re-fetch dep to get latest .Status
	if err := r.Get(ctx, client.ObjectKeyFromObject(dep), dep); err != nil {
		log.Error(err, "Failed to refresh Deployment status")
		// continue anyway
	}

	app.Status.Phase = "Running" // or compute from conditions
	app.Status.ReadyReplicas = dep.Status.ReadyReplicas

	if err := r.Status().Update(ctx, app); err != nil {
		log.Error(err, "Failed to update App status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
// func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
// 	return ctrl.NewControllerManagedBy(mgr).
// 		For(&appsv1.App{}).
// 		Complete(r)
// }

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.App{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *AppReconciler) desiredDeployment(app *appv1.App) *appsv1.Deployment {
	labels := map[string]string{"app": app.Name}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-app",
			Namespace: app.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "app",
						Image: app.Spec.Image, // will be overridden in mutate if changed
						Ports: []corev1.ContainerPort{{
							ContainerPort: app.Spec.Port,
						}},
						Env: app.Spec.Env, // will be overridden in mutate
					}},
				},
			},
		},
	}

	// Always set owner reference (garbage collection + event trigger)
	ctrl.SetControllerReference(app, dep, r.Scheme)
	return dep
}

func (r *AppReconciler) desiredService(app *appv1.App) *corev1.Service {
	labels := map[string]string{"app": app.Name}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-svc",
			Namespace: app.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       80,                   // default; mutate will set to spec.port
				TargetPort: intstr.FromInt32(80), // default; mutate overrides
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	ctrl.SetControllerReference(app, svc, r.Scheme)
	return svc
}
