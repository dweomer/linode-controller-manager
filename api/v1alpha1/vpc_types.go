package v1alpha1

// VPCSubnet defines a subnet within a VPC.
type VPCSubnet struct {
	// +kubebuilder:title:="ID"
	// +k8s:optional
	ID int64 `json:"id,omitzero" linode:"id" api:"readonly,filterable"`

	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=64
	// +k8s:required
	Label string `json:"label" linode:"label" api:"required,filterable"`

	// +kubebuilder:title:="IPv4"
	// RFC1918 range in CIDR notation. Prefix lengths 1-29.
	// +k8s:required
	IPv4 string `json:"ipv4" linode:"ipv4" api:"required,immutable"`
}
