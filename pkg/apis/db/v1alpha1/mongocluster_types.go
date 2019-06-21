package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
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

	Mongo    MongoSettings `json:"mongo,omitempty"`
}

// MongoSettings define the specification of the mongo cluster
type MongoSettings struct {
	Image           string              `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy   `json:"imagePullPolicy,omitempty"`
	Replicas        int32               `json:"replicas,omitempty"`
	Resources       MongoResources      `json:"resources,omitempty"`
	Command         []string            `json:"command,omitempty"`
	Storage         MongoStorage        `json:"storage,omitempty"`
	Tolerations     []corev1.Toleration `json:"tolerations,omitempty"`
}

// MongoResources sets the limits and requests for a container
type MongoResources struct {
	Requests CPUAndMem `json:"requests,omitempty"`
	Limits   CPUAndMem `json:"limits,omitempty"`
}

// CPUAndMem defines how many cpu and ram the container will request/limit
type CPUAndMem struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type MongoStorage struct {
	StorageClassName string                      `json:"storageClassName,omitempty"`
	Resources        corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
}

// MongoClusterStatus defines the observed state of MongoCluster
// +k8s:openapi-gen=true
type MongoClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`
	Replicas           int32  `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
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
