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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppSpec defines the desired state of App
// type AppSpec struct {
// 	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
// 	// Important: Run "make" to regenerate code after modifying this file
//
// 	// Foo is an example field of App. Edit app_types.go to remove/update
// 	Foo string `json:"foo,omitempty"`
// }

type AppSpec struct {
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas *int32 `json:"replicas,omitempty"`

	// Container port to expose
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// Optional environment variables
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// AppStatus defines the observed state of App
// type AppStatus struct {
// 	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
// 	// Important: Run "make" to regenerate code after modifying this file
// }

type AppStatus struct {
	// Conditions of the app
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ReadyReplicas shows how many pods are ready
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Phase e.g. Pending, Running, Failed
	Phase string `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// App is the Schema for the apps API
type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec,omitempty"`
	Status AppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppList contains a list of App
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

func init() {
	SchemeBuilder.Register(&App{}, &AppList{})
}
