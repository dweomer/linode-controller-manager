package v1alpha1

// FirewallRules defines inbound and outbound rules with default policies.
type FirewallRules struct {
	// +kubebuilder:validation:Enum=ACCEPT;DROP
	// +k8s:required
	InboundPolicy FirewallPolicy `json:"inboundPolicy" linode:"inbound_policy" api:"required"`

	// +listType=atomic
	// +k8s:optional
	Inbound []FirewallRule `json:"inbound,omitzero" linode:"inbound"`

	// +kubebuilder:validation:Enum=ACCEPT;DROP
	// +k8s:required
	OutboundPolicy FirewallPolicy `json:"outboundPolicy" linode:"outbound_policy" api:"required"`

	// +listType=atomic
	// +k8s:optional
	Outbound []FirewallRule `json:"outbound,omitzero" linode:"outbound"`

	// +k8s:optional
	Fingerprint string `json:"fingerprint,omitzero" linode:"fingerprint" api:"readonly"`
}

// FirewallRule defines a single inbound or outbound firewall rule.
type FirewallRule struct {
	// +kubebuilder:validation:Enum=ACCEPT;DROP
	// +k8s:required
	Action FirewallPolicy `json:"action" linode:"action" api:"required"`

	// +kubebuilder:validation:Enum=TCP;UDP;ICMP;IPENCAP
	// +k8s:required
	Protocol FirewallProtocol `json:"protocol" linode:"protocol" api:"required"`

	// Comma-separated list of ports or port ranges (e.g. "22,80,443,8000-9000").
	// +k8s:optional
	Ports string `json:"ports,omitzero" linode:"ports"`

	// +k8s:required
	Addresses FirewallAddresses `json:"addresses" linode:"addresses" api:"required"`

	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=32
	// +k8s:optional
	Label string `json:"label,omitzero" linode:"label"`

	// +k8s:optional
	Description string `json:"description,omitzero" linode:"description"`
}

// FirewallAddresses
// +structType=atomic
type FirewallAddresses struct {
	// +kubebuilder:title:="IPv4"
	// +listType=atomic
	// +k8s:optional
	IPv4 []string `json:"ipv4,omitzero" linode:"ipv4"`

	// +kubebuilder:title:="IPv6"
	// +listType=atomic
	// +k8s:optional
	IPv6 []string `json:"ipv6,omitzero" linode:"ipv6"`
}

