package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,path=vpcs,singular=vpc,shortName=linvpc,categories=linode
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type=integer,priority=0,JSONPath=`.spec.id`,description="VPC ID"
// +kubebuilder:printcolumn:name="Label",type=string,priority=0,JSONPath=`.spec.label`,description="VPC Label"
// +kubebuilder:printcolumn:name="Region",type=string,priority=0,JSONPath=`.spec.region`,description="Linode Region"
// +kubebuilder:selectablefield:JSONPath=`.spec.id`
// +kubebuilder:selectablefield:JSONPath=`.spec.label`
// +kubebuilder:selectablefield:JSONPath=`.spec.region`

// VPC is the Schema for the vpcs API.
// CREATE: https://techdocs.akamai.com/linode-api/reference/post-vpc
// UPDATE: https://techdocs.akamai.com/linode-api/reference/put-vpc
// READ:   https://techdocs.akamai.com/linode-api/reference/get-vpc
type VPC struct {
	metav1.TypeMeta `json:",inline"`

	// +k8s:optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// +k8s:required
	Spec VPCSpec `json:"spec"`

	// +k8s:optional
	Status VPCStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// VPCList contains a list of VPC.
type VPCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []VPC `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VPC{}, &VPCList{})
}

// VPCSpec defines the desired state of a VPC.
type VPCSpec struct {
	// +kubebuilder:title:="ID"
	// +k8s:optional
	ID int64 `json:"id,omitzero" linode:"id" api:"readonly,filterable"`

	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=64
	// +k8s:required
	Label string `json:"label" linode:"label" api:"required,filterable"`

	// +k8s:optional
	Description string `json:"description,omitzero" linode:"description"`

	// +k8s:required
	Region string `json:"region" linode:"region" api:"required,immutable,filterable"`

	// +listType=atomic
	// +k8s:optional
	Subnets []VPCSubnet `json:"subnets,omitzero" linode:"subnets"`
}

// VPCStatus defines the observed state of a VPC.
type VPCStatus struct {
	// +k8s:optional
	Created metav1.Time `json:"created,omitzero" linode:"created" api:"readonly"`

	// +k8s:optional
	Updated metav1.Time `json:"updated,omitzero" linode:"updated" api:"readonly"`

	// +listType=map
	// +listMapKey=type
	// +k8s:optional
	Conditions []metav1.Condition `json:"conditions,omitzero"`
}
