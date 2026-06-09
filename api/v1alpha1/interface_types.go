package v1alpha1

// InstanceInterfaces is a polymorphic list of network interfaces for an Instance.
// The valid shape depends on the parent Instance's interfaceGeneration:
//   - "linode": each entry uses exactly one of Public, VPC, or VLAN
//   - "legacy_config": each entry uses Purpose with the legacy flat fields
//
// +listType=atomic
type InstanceInterfaces []Interface

// Interface represents a single network interface.
// For linode-gen: set exactly one of Public, VPC, or VLAN.
// For legacy-gen: set Purpose and the corresponding flat fields.
type Interface struct {
	// --- linode generation fields ---

	// +kubebuilder:title:="Firewall ID"
	// +k8s:optional
	FirewallID *int64 `json:"firewallID,omitzero" linode:"firewall_id"`

	// +structType=atomic
	// +k8s:optional
	DefaultRoute *InterfaceDefaultRoute `json:"defaultRoute,omitzero" linode:"default_route"`

	// +k8s:optional
	Public *InterfacePublic `json:"public,omitzero" linode:"public"`

	// +kubebuilder:title:="VPC"
	// +k8s:optional
	VPC *InterfaceVPC `json:"vpc,omitzero" linode:"vpc"`

	// +kubebuilder:title:="VLAN"
	// +k8s:optional
	VLAN *InterfaceVLAN `json:"vlan,omitzero" linode:"vlan"`

	// --- legacy generation fields ---

	// +kubebuilder:validation:Enum=public;vlan;vpc
	// +k8s:optional
	Purpose InterfacePurpose `json:"purpose,omitzero" linode:"purpose" api:"immutable"`

	// +k8s:optional
	Primary bool `json:"primary,omitzero" linode:"primary"`

	// VLAN label. Required when purpose is "vlan".
	// +kubebuilder:validation:MaxLength:=64
	// +kubebuilder:validation:Pattern:="[a-zA-Z0-9-]+"
	// +k8s:optional
	Label string `json:"label,omitzero" linode:"label"`

	// VLAN IPAM address in CIDR notation. Only for purpose "vlan".
	// +kubebuilder:title:="IPAM Address"
	// +k8s:optional
	IPAMAddress string `json:"ipamAddress,omitzero" linode:"ipam_address"`

	// VPC subnet ID. Required when purpose is "vpc" (legacy).
	// +kubebuilder:title:="Subnet ID"
	// +k8s:optional
	SubnetID *int64 `json:"subnetID,omitzero" linode:"subnet_id"`

	// Legacy VPC IPv4 configuration.
	// +kubebuilder:title:="IPv4"
	// +structType=atomic
	// +k8s:optional
	IPv4 *InterfaceLegacyIPv4 `json:"ipv4,omitzero" linode:"ipv4"`

	// Legacy VPC IP ranges routed to this interface.
	// +kubebuilder:title:="IP Ranges"
	// +listType=set
	// +k8s:optional
	IPRanges []string `json:"ipRanges,omitzero" linode:"ip_ranges"`
}

// InterfaceDefaultRoute
// +structType=atomic
type InterfaceDefaultRoute struct {
	// +kubebuilder:title:="IPv4"
	// +k8s:optional
	IPv4 *bool `json:"ipv4,omitzero" linode:"ipv4"`
	// +kubebuilder:title:="IPv6"
	// +k8s:optional
	IPv6 *bool `json:"ipv6,omitzero" linode:"ipv6"`
}

type InterfacePublic struct {
	// +kubebuilder:title:="IPv4"
	// +k8s:optional
	IPv4 *InterfacePublicIPv4 `json:"ipv4,omitzero" linode:"ipv4"`
	// +kubebuilder:title:="IPv6"
	// +k8s:optional
	IPv6 *InterfacePublicIPv6 `json:"ipv6,omitzero" linode:"ipv6"`
}

type InterfacePublicIPv4 struct {
	// +listType=atomic
	// +k8s:optional
	Addresses []InterfacePublicIPv4Address `json:"addresses,omitzero" linode:"addresses"`
}

type InterfacePublicIPv4Address struct {
	// "auto" or a specific IPv4 address.
	Address string `json:"address" linode:"address"`
	// +k8s:optional
	Primary bool `json:"primary,omitzero" linode:"primary"`
}

type InterfacePublicIPv6 struct {
	// +listType=atomic
	// +k8s:optional
	Ranges []InterfaceIPv6Range `json:"ranges,omitzero" linode:"ranges"`
}

type InterfaceIPv6Range struct {
	Range string `json:"range" linode:"range"`
}

type InterfaceVPC struct {
	// +kubebuilder:title:="Subnet ID"
	// +k8s:required
	SubnetID int64 `json:"subnetID" linode:"subnet_id"`
	// +kubebuilder:title:="IPv4"
	// +k8s:optional
	IPv4 *InterfaceVPCIPv4 `json:"ipv4,omitzero" linode:"ipv4"`
}

type InterfaceVPCIPv4 struct {
	// +listType=atomic
	// +k8s:optional
	Addresses []InterfaceVPCIPv4Address `json:"addresses,omitzero" linode:"addresses"`
	// +listType=atomic
	// +k8s:optional
	Ranges []InterfaceVPCIPv4Range `json:"ranges,omitzero" linode:"ranges"`
}

type InterfaceVPCIPv4Address struct {
	// "auto" or a specific VPC subnet IPv4 address.
	Address string `json:"address" linode:"address"`
	// +k8s:optional
	Primary *bool `json:"primary,omitzero" linode:"primary"`
	// Bidirectional NAT mapping between this VPC address and a public IPv4.
	// "auto", a specific public IPv4, or omit for no NAT.
	// +kubebuilder:title:="NAT"
	// +k8s:optional
	NAT *string `json:"nat,omitzero" linode:"nat_1_1_address"`
}

type InterfaceVPCIPv4Range struct {
	// CIDR notation or prefix only (e.g. "10.0.0.0/28" or "/28").
	Range string `json:"range" linode:"range"`
}

type InterfaceVLAN struct {
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=64
	// +kubebuilder:validation:Pattern:="[a-zA-Z0-9-]+"
	// +k8s:required
	Label string `json:"label" linode:"vlan_label"`
	// +kubebuilder:title:="IPAM Address"
	// +k8s:optional
	IPAMAddress string `json:"ipamAddress,omitzero" linode:"ipam_address"`
}

// +structType=atomic
type InterfaceLegacyIPv4 struct {
	// VPC subnet IPv4 address.
	// +kubebuilder:title:="VPC"
	// +k8s:optional
	VPC string `json:"vpc,omitzero" linode:"vpc"`
	// Bidirectional NAT mapping between this VPC address and a public IPv4.
	// A specific address or "any".
	// +kubebuilder:title:="NAT"
	// +k8s:optional
	NAT *string `json:"nat,omitzero" linode:"nat_1_1"`
}
