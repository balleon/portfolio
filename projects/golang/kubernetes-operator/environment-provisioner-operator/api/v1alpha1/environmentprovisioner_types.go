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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceList holds CPU and memory quantities for a resource constraint.
type ResourceList struct {
	CPU    resource.Quantity `json:"cpu"`
	Memory resource.Quantity `json:"memory"`
}

// ResourcesQuota defines the limits and requests for the ResourceQuota.
type ResourcesQuota struct {
	Limits  ResourceList `json:"limits"`
	Request ResourceList `json:"request"`
}

// EnvironmentProvisionerSpec defines the desired state of EnvironmentProvisioner
type EnvironmentProvisionerSpec struct {
	// namespaceName is the name of the Namespace to provision.
	// +kubebuilder:validation:MinLength=1
	NamespaceName string `json:"namespaceName"`

	// resourcesQuota defines the resource limits and requests applied to the provisioned Namespace.
	ResourcesQuota ResourcesQuota `json:"resourcesQuota"`
}

// EnvironmentProvisionerStatus defines the observed state of EnvironmentProvisioner.
type EnvironmentProvisionerStatus struct {
	// conditions represent the current state of the EnvironmentProvisioner resource.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// EnvironmentProvisioner is the Schema for the environmentprovisioners API
type EnvironmentProvisioner struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of EnvironmentProvisioner
	// +required
	Spec EnvironmentProvisionerSpec `json:"spec"`

	// status defines the observed state of EnvironmentProvisioner
	// +optional
	Status EnvironmentProvisionerStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// EnvironmentProvisionerList contains a list of EnvironmentProvisioner
type EnvironmentProvisionerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []EnvironmentProvisioner `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &EnvironmentProvisioner{}, &EnvironmentProvisionerList{})
		return nil
	})
}
