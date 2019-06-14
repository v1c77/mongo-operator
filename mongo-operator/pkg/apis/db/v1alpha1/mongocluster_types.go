package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MongoClusterSpec defines the desired state of MongoCluster
// +k8s:openapi-gen=true
type MongoClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas int32 `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
}

// MongoClusterStatus defines the observed state of MongoCluster
// +k8s:openapi-gen=true
type MongoClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`
	Replicas int32 `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoCluster is the Schema for the mongoclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type MongoCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoClusterSpec   `json:"spec,omitempty"`
	Status MongoClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoClusterList contains a list of MongoCluster
type MongoClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoCluster{}, &MongoClusterList{})
}
