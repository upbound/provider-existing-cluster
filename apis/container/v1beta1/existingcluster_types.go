/*
Copyright 2019 The Crossplane Authors.

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

package v1beta1

import (
	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cluster states.
const (
	ClusterStateRunning = "RUNNING"
)

// Defaults for Existing Cluster resources.
const (
	DefaultReclaimPolicy = runtimev1alpha1.ReclaimRetain
)

// ExistingClusterObservation is used to show the observed state of the existing cluster cluster resource.
type ExistingClusterObservation struct {
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Endpoint      string `json:"endpoint,omitempty"`
}

// ExistingClusterParameters define the desired state of an existing cluster.
type ExistingClusterParameters struct {
}

// A ExistingClusterSpec defines the desired state of a ExistingCluster.
type ExistingClusterSpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  ExistingClusterParameters `json:"forProvider"`
}

// A ExistingClusterStatus represents the observed state of a ExistingCluster.
type ExistingClusterStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     ExistingClusterObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ExistingCluster is a managed resource that represents a Google Kubernetes Engine
// cluster.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.bindingPhase"
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.atProvider.status"
// +kubebuilder:printcolumn:name="ENDPOINT",type="string",JSONPath=".status.atProvider.endpoint"
// +kubebuilder:printcolumn:name="CLUSTER-CLASS",type="string",JSONPath=".spec.classRef.name"
// +kubebuilder:printcolumn:name="RECLAIM-POLICY",type="string",JSONPath=".spec.reclaimPolicy"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster
type ExistingCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExistingClusterSpec   `json:"spec"`
	Status ExistingClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExistingClusterList contains a list of ExistingCluster items
type ExistingClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExistingCluster `json:"items"`
}
