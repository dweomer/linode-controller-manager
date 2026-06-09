package v1alpha1

const (
	InstanceLockCannotDelete                 InstanceLockType = "cannot_delete"
	InstanceLockCannotDeleteWithSubresources InstanceLockType = "cannot_delete_with_subresources"

	InstanceDiskEncryptionEnabled  InstanceDiskEncryption = "enabled"  // linodego.InstanceDiskEncryptionEnabled
	InstanceDiskEncryptionDisabled InstanceDiskEncryption = "disabled" // linodego.InstanceDiskEncryptionDisabled

	InstanceGenerationLinode InstanceInterfaceGeneration = "linode"        // linodego.GenerationLinode
	InstanceGenerationLegacy InstanceInterfaceGeneration = "legacy_config" // linodego.GenerationLegacyConfig

	InstanceMaintenancePolicyMigrate    InstanceMaintenancePolicy = "linode/migrate"
	InstanceMaintenancePolicyPowerOffOn InstanceMaintenancePolicy = "linode/power_off_on"

	InstancePlacementGroupAntiAffinityLocal InstancePlacementGroupType = "anti_affinity:local"

	InstancePlacementGroupPolicyStrict   InstancePlacementGroupPolicy = "strict"
	InstancePlacementGroupPolicyFlexible InstancePlacementGroupPolicy = "flexible"

	InstanceStatusBillingSuspension InstanceStatusType = "billing_suspension"
	InstanceStatusBooting           InstanceStatusType = "booting"
	InstanceStatusBusy              InstanceStatusType = "busy"
	InstanceStatusCloning           InstanceStatusType = "cloning"
	InstanceStatusDeleting          InstanceStatusType = "deleting"
	InstanceStatusMigrating         InstanceStatusType = "migrating"
	InstanceStatusOffline           InstanceStatusType = "offline"
	InstanceStatusProvisioning      InstanceStatusType = "provisioning"
	InstanceStatusRebooting         InstanceStatusType = "rebooting"
	InstanceStatusRebuilding        InstanceStatusType = "rebuilding"
	InstanceStatusRestoring         InstanceStatusType = "restoring"
	InstanceStatusRunning           InstanceStatusType = "running"
	InstanceStatusResizing          InstanceStatusType = "resizing"
	InstanceStatusShuttingDown      InstanceStatusType = "shutting_down"
	InstanceStatusStopped           InstanceStatusType = "stopped"

	InterfacePurposePublic InterfacePurpose = "public"
	InterfacePurposeVLAN   InterfacePurpose = "vlan"
	InterfacePurposeVPC    InterfacePurpose = "vpc"
)

type (
	// InstanceDiskEncryption represents the allowed values for disk_encryption
	// +k8s:enum
	InstanceDiskEncryption string
	// InstanceMaintenancePolicy represents the allowed values for maintenance_policy
	// +k8s:enum
	InstanceMaintenancePolicy string
	// InstanceInterfaceGeneration represents the allowed types for interface_generation
	// +k8s:enum
	InstanceInterfaceGeneration string
	// InterfacePurpose represents the allowed values for a legacy config interface purpose.
	// +k8s:enum
	InterfacePurpose string
	// InstancePlacementGroupPolicy represents the allowed types for placement_group.policy
	// +k8s:enum
	InstancePlacementGroupPolicy string
	// InstancePlacementGroupType represents the allowed types for placement_group.type
	// +k8s:enum
	InstancePlacementGroupType string
	// InstanceLockType represents the allowed lock types (v4beta only)
	// +k8s:enum
	InstanceLockType string
	// InstanceStatusType represents the allowed types for status
	// +k8s:enum
	InstanceStatusType string
)
