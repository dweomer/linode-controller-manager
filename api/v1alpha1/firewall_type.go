package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,path=firewalls,singular=firewall,shortName=linfw,categories=linode
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type=integer,priority=0,JSONPath=`.spec.id`,description="Firewall ID"
// +kubebuilder:printcolumn:name="Label",type=string,priority=0,JSONPath=`.spec.label`,description="Firewall Label"
// +kubebuilder:printcolumn:name="Status",type=string,priority=0,JSONPath=`.status.status`,description="Firewall Status"
// +kubebuilder:selectablefield:JSONPath=`.spec.id`
// +kubebuilder:selectablefield:JSONPath=`.spec.label`

// Firewall is the Schema for the firewalls API.
// CREATE: https://techdocs.akamai.com/linode-api/reference/post-firewalls
// UPDATE: https://techdocs.akamai.com/linode-api/reference/put-firewall
// READ:   https://techdocs.akamai.com/linode-api/reference/get-firewall
type Firewall struct {
	metav1.TypeMeta `json:",inline"`

	// +k8s:optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// +k8s:required
	Spec FirewallSpec `json:"spec"`

	// +k8s:optional
	Status FirewallStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// FirewallList contains a list of Firewall.
type FirewallList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Firewall `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Firewall{}, &FirewallList{})
}

// FirewallSpec defines the desired state of a Firewall.
type FirewallSpec struct {
	// +kubebuilder:title:="ID"
	// +k8s:optional
	ID int64 `json:"id,omitzero" linode:"id" api:"readonly,filterable"`

	// +kubebuilder:validation:MinLength:=3
	// +kubebuilder:validation:MaxLength:=32
	// +k8s:required
	Label string `json:"label" linode:"label" api:"required,filterable"`

	// +listType=set
	// +k8s:optional
	Tags []string `json:"tags,omitzero" linode:"tags" api:"filterable"`

	// +k8s:required
	Rules FirewallRules `json:"rules" linode:"rules"`
}

// FirewallStatus defines the observed state of a Firewall.
type FirewallStatus struct {
	// +kubebuilder:validation:Enum=enabled;disabled;deleted
	// +k8s:optional
	Status FirewallStatusType `json:"status,omitzero" linode:"status" api:"readonly"`

	// +k8s:optional
	Created metav1.Time `json:"created,omitzero" linode:"created" api:"readonly"`

	// +k8s:optional
	Updated metav1.Time `json:"updated,omitzero" linode:"updated" api:"readonly"`

	// +listType=atomic
	// +k8s:optional
	Entities []Entity `json:"entities,omitzero" linode:"entities" api:"readonly"`

	// +listType=map
	// +listMapKey=type
	// +k8s:optional
	Conditions []metav1.Condition `json:"conditions,omitzero"`
}
