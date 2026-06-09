package v1alpha1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,path=instances,singular=instance,shortName=li;lin;linodes,categories=linode
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type=integer,priority=0,JSONPath=`.spec.id`,description="Instance ID"
// +kubebuilder:printcolumn:name="Label",type=string,priority=0,JSONPath=`.spec.label`,description="Instance Label"
// +kubebuilder:printcolumn:name="Region",type=string,priority=0,JSONPath=`.spec.region`,description="Linode Region"
// +kubebuilder:printcolumn:name="Group",type=string,priority=9,JSONPath=`.spec.group`,description="Linode Group"
// +kubebuilder:printcolumn:name="Generation",type=string,priority=1,JSONPath=`.spec.interfaceGeneration`,description="Interface Generation"
// +kubebuilder:printcolumn:name="Type",type=string,priority=1,JSONPath=`.spec.type`,description="Linode Type"
// +kubebuilder:printcolumn:name="Placement",type=string,priority=1,JSONPath=`.spec.placementGroup.label`,description="Placement Group Label"
// +kubebuilder:printcolumn:name="Status",type=string,priority=0,JSONPath=`.status.status`,description="Linode Status"
// +kubebuilder:selectablefield:JSONPath=`.spec.id`
// +kubebuilder:selectablefield:JSONPath=`.spec.label`
// +kubebuilder:selectablefield:JSONPath=`.spec.region`
// +kubebuilder:selectablefield:JSONPath=`.spec.group`
// +kubebuilder:selectablefield:JSONPath=`.spec.interfaceGeneration`
// +kubebuilder:selectablefield:JSONPath=`.spec.placementGroup.id`
// +kubebuilder:selectablefield:JSONPath=`.spec.placementGroup.label`
// +kubebuilder:selectablefield:JSONPath=`.spec.placementGroup.type`

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +k8s:optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of Instance
	// +k8s:required
	Spec InstanceSpec `json:"spec"`

	// status defines the observed state of Instance
	// +k8s:optional
	Status InstanceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}

// InstanceSpec defines the desired state of Instance
// CREATE: https://techdocs.akamai.com/linode-api/reference/post-linode-instance
// UPDATE: https://techdocs.akamai.com/linode-api/reference/put-linode-instance
// READ:   https://techdocs.akamai.com/linode-api/reference/get-linode-instance
type InstanceSpec struct {
	// +k8s:optional
	ID int64 `json:"id,omitzero" linode:"id" api:"readonly,filterable,status"`

	// +kubebuilder:validation:MinLength:=3
	// +kubebuilder:validation:MaxLength:=64
	// +k8s:optional
	Label string `json:"label,omitzero" linode:"label" api:"required,filterable,status"`

	// +k8s:optional
	Image string `json:"image,omitzero" linode:"image" api:"immutable,status"`

	// +kubebuilder:validation:MinLength:=7
	// +kubebuilder:validation:MaxLength:=128
	// +k8s:optional
	RootPassword string `json:"rootPassword,omitzero" linode:"root_password" datapolicy:"password" api:"required,writeonly"`

	// +k8s:optional
	Region string `json:"region,omitzero" linode:"region" api:"required,filterable,immutable,status"`

	// +k8s:optional
	Type string `json:"type,omitzero" linode:"type" api:"required,immutable,status"`

	// +listType=set
	// +k8s:optional
	Tags []string `json:"tags,omitzero" linode:"tags" api:"filterable,status"`

	// +listType=set
	// +k8s:optional
	AuthorizedKeys []string `json:"authorizedKeys,omitzero" linode:"authorized_keys" api:"writeonly"`

	// +listType=set
	// +k8s:optional
	AuthorizedUsers []string `json:"authorizedUsers,omitzero" linode:"authorized_users" api:"writeonly"`

	// +k8s:optional
	// +default:value=true
	Booted bool `json:"booted" linode:"booted" default:"true" api:"writeonly"`

	// +k8s:optional
	// +default:value="enabled"
	DiskEncryption InstanceDiskEncryption `json:"diskEncryption,omitzero" linode:"disk_encryption" default:"enabled" api:"immutable,status"`

	// +k8s:optional
	InterfaceGeneration InstanceInterfaceGeneration `json:"interfaceGeneration,omitzero" linode:"interface_generation" api:"immutable,filterable,status"`

	// +k8s:optional
	Metadata InstanceMetadata `json:"metadata,omitzero" linode:"metadata" api:"writeonly"`

	// +k8s:optional
	// +default:value="linode/power_off_on"
	MaintenancePolicy InstanceMaintenancePolicy `json:"maintenancePolicy,omitzero" linode:"maintenance_policy" default:"linode/power_off_on" api:"status"`

	// +k8s:optional
	StackscriptID int64 `json:"stackscriptID,omitzero" linode:"stackscript_id" api:"writeonly"`
	// +k8s:optional
	StackscriptData json.RawMessage `json:"stackscriptData,omitzero" linode:"stackscript_data" api:"writeonly"`

	// +k8s:optional
	BackupID int64 `json:"backupID,omitzero" linode:"backup_id" api:"writeonly"`
	// +k8s:optional
	BackupsEnabled *bool `json:"backupsEnabled,omitzero" linode:"backups_enabled" api:"writeonly"`

	// +k8s:optional
	FirewallID int64 `json:"firewallID,omitzero" linode:"firewall_id" api:"writeonly"`

	// +k8s:optional
	Group string `json:"group,omitzero" linode:"group" api:"deprecated,filterable,status"`

	// +k8s:optional
	Interfaces InstanceInterfaces `json:"interfaces,omitzero" linode:"interfaces" api:"writeonly"`

	// +k8s:optional
	NetworkHelper *bool `json:"networkHelper,omitzero" linode:"network_helper" api:"writeonly"`

	// +k8s:optional
	PlacementGroup InstancePlacementGroup `json:"placementGroup,omitzero" linode:"placement_group" api:"status"`

	// +k8s:optional
	// +default:value:=false
	PrivateIP bool `json:"privateIP,omitzero" linode:"private_ip" api:"writeonly"`

	// +k8s:optional
	// +default:value:=512
	SwapSize *int64 `json:"swapSize,omitzero" linode:"swap_size" api:"writeonly"`

	// +k8s:optional
	Kernel *string `json:"kernel,omitzero" linode:"kernel" api:"writeonly"`

	// +k8s:optional
	BootSize *int `json:"bootSize,omitzero" linode:"boot_size" api:"writeonly"`

	// +listType=set
	// +k8s:optional
	IPv4 []string `json:"ipv4,omitzero" linode:"ipv4" api:"writeonly"`

	// +k8s:optional
	WatchdogEnabled *bool `json:"watchdogEnabled,omitzero" linode:"watchdog_enabled" api:"status"`

	// +structType=atomic
	// +k8s:optional
	Alerts *InstanceAlerts `json:"alerts,omitzero" linode:"alerts" api:"status"`

	// +k8s:optional
	BackupSchedule *InstanceBackupSchedule `json:"backupSchedule,omitzero" linode:"backups" api:"status"`

	// +listType=set
	// +k8s:optional
	Locks []InstanceLockType `json:"locks,omitzero" linode:"locks" api:"status,v4beta"`
}

// InstanceStatus defines the observed state of Instance.
type InstanceStatus struct {
	// +k8s:optional
	Created metav1.Time `json:"created,omitzero" linode:"created" api:"readonly,filterable,status"`

	// +k8s:optional
	Updated metav1.Time `json:"updated,omitzero" linode:"updated" api:"readonly,filterable,status"`

	// +k8s:optional
	Backups InstanceBackups `json:"backups,omitzero" linode:"backups" api:"readonly,status"`

	// +listType=set
	// +k8s:optional
	Capabilities []string `json:"capabilities,omitzero" linode:"capabilities" api:"readonly,status"`

	// +k8s:optional
	HasUserData bool `json:"hasUserData,omitzero" linode:"has_user_data" api:"readonly,status"`

	// +k8s:optional
	HostUUID string `json:"hostUUID,omitzero" linode:"host_uuid" api:"readonly,status"`

	// +k8s:optional
	HyperVisor string `json:"hypervisor,omitzero" linode:"hypervisor" api:"readonly,status"`

	// +kubebuilder:title:="IPv4"
	// +k8s:optional
	IPv4 []string `json:"ipv4,omitzero" linode:"ipv4" api:"readonly,filterable,status"`

	// +kubebuilder:title:="IPv6"
	// +k8s:optional
	IPv6 string `json:"ipv6,omitzero" linode:"ipv6" api:"readonly,status"`

	// +kubebuilder:title:="LKE Cluster ID"
	// +k8s:optional
	LKEClusterID int64 `json:"lkeClusterID,omitzero" linode:"lke_cluster_id" api:"readonly,status"`

	// +k8s:optional
	Specs InstanceSpecs `json:"specs,omitzero" linode:"specs" api:"readonly,status"`

	// +k8s:optional
	Status InstanceStatusType `json:"status,omitzero" linode:"status" api:"readonly,status"`

	// +listType=map
	// +listMapKey=type
	// +k8s:optional
	Conditions []metav1.Condition `json:"conditions,omitzero"`
}
