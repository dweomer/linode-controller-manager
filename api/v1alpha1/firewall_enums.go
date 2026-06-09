package v1alpha1

const (
	FirewallPolicyAccept FirewallPolicy = "ACCEPT"
	FirewallPolicyDrop   FirewallPolicy = "DROP"

	FirewallProtocolTCP     FirewallProtocol = "TCP"
	FirewallProtocolUDP     FirewallProtocol = "UDP"
	FirewallProtocolICMP    FirewallProtocol = "ICMP"
	FirewallProtocolIPENCAP FirewallProtocol = "IPENCAP"

	FirewallStatusEnabled  FirewallStatusType = "enabled"
	FirewallStatusDisabled FirewallStatusType = "disabled"
	FirewallStatusDeleted  FirewallStatusType = "deleted"
)

type (
	// FirewallPolicy represents the allowed values for firewall rule actions and default policies.
	// +k8s:enum
	FirewallPolicy string
	// FirewallProtocol represents the allowed values for firewall rule protocols.
	// +k8s:enum
	FirewallProtocol string
	// FirewallStatusType represents the allowed status values for a Firewall.
	// +k8s:enum
	FirewallStatusType string
)
