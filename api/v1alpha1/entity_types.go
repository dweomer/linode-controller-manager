package v1alpha1

// Entity is a reference to a Linode API resource.
// +structType=atomic
type Entity struct {
	// +kubebuilder:title:="ID"
	ID int64 `json:"id,omitzero" linode:"id"`

	Type  EntityType `json:"type,omitzero" linode:"type"`
	Label string     `json:"label,omitzero" linode:"label"`

	// +kubebuilder:title:="URL"
	// +k8s:optional
	URL string `json:"url,omitzero" linode:"url"`

	// +k8s:optional
	ParentEntity *Entity `json:"parentEntity,omitzero" linode:"parent_entity"`
}
