package v1alpha1

import "k8s.io/apimachinery/pkg/apis/meta/v1"

type InstanceAlerts struct {
	CPU           int64 `json:"cpu,omitzero" linode:"cpu"`
	IO            int64 `json:"io,omitzero" linode:"io"`
	NetworkIn     int64 `json:"networkIn,omitzero" linode:"network_in"`
	NetworkOut    int64 `json:"networkOut,omitzero" linode:"network_out"`
	TransferQuota int64 `json:"transferQuota,omitzero" linode:"transfer_quota"`
}

type InstanceBackups struct {
	Available      bool                   `json:"available,omitzero" linode:"available" api:"readonly"`
	Enabled        bool                   `json:"enabled,omitzero" linode:"enabled" api:"readonly"`
	Schedule       InstanceBackupSchedule `json:"schedule,omitzero" linode:"schedule" api:"readonly"`
	LastSuccessful v1.Time                `json:"lastSuccessful,omitzero" linode:"last_successful" api:"readonly"`
}

type InstanceBackupSchedule struct {
	Day    string `json:"day,omitzero" linode:"day"`
	Window string `json:"window,omitzero" linode:"window"`
}

type InstanceMetadata struct {
	UserData string `json:"userData,omitempty" linode:"user_data"`
}

type InstancePlacementGroup struct {
	ID int64 `json:"id,omitzero" linode:"id" api:"filterable,readonly"`
	// +kubebuilder:validation:MinLength:=1
	Label  string                       `json:"label,omitzero" linode:"label" api:"filterable"`
	Policy InstancePlacementGroupPolicy `json:"policy,omitzero" linode:"placement_group_policy"`
	Type   InstancePlacementGroupType   `json:"type,omitzero" linode:"placement_group_type" api:"filterable"`

	// +k8s:optional
	MigratingTo *int64 `json:"migratingTo,omitzero" linode:"migrating_to" api:"readonly"`
	// +k8s:optional
	CompliantOnly *bool `json:"compliantOnly,omitzero" linode:"compliant_only" api:"writeonly"`
}

type InstanceSpecs struct {
	Disk               int32 `json:"disk,omitzero" linode:"disk" api:"readonly"`
	GPUs               int32 `json:"gpus,omitzero" linode:"gpus" api:"readonly"`
	Memory             int32 `json:"memory,omitzero" linode:"memory" api:"readonly"`
	Transfer           int32 `json:"transfer,omitzero" linode:"transfer" api:"readonly"`
	VCPUs              int32 `json:"vcpus,omitzero" linode:"vcpus" api:"readonly"`
	AcceleratedDevices int32 `json:"acceleratedDevices,omitzero" linode:"accelerated_devices" api:"readonly"`
}
