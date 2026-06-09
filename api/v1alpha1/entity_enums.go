package v1alpha1

const (
	EntityTypeAccount        EntityType = "account"
	EntityTypeBackups        EntityType = "backups"
	EntityTypeCommunity      EntityType = "community"
	EntityTypeDatabase       EntityType = "database"
	EntityTypeDisks          EntityType = "disks"
	EntityTypeDomain         EntityType = "domain"
	EntityTypeEntityTransfer EntityType = "entity_transfer"
	EntityTypeFirewall       EntityType = "firewall"
	EntityTypeImage          EntityType = "image"
	EntityTypeIPAddress      EntityType = "ipaddress"
	EntityTypeLinode         EntityType = "linode"
	EntityTypeLKECluster     EntityType = "lkecluster"
	EntityTypeLoadBalancer   EntityType = "loadbalancer"
	EntityTypeLongview       EntityType = "longview"
	EntityTypeManagedService EntityType = "managed_service"
	EntityTypeNodeBalancer   EntityType = "nodebalancer"
	EntityTypeOAuthClient    EntityType = "oauth_client"
	EntityTypePlacementGroup EntityType = "placement_group"
	EntityTypeProfile        EntityType = "profile"
	EntityTypeStackscript    EntityType = "stackscript"
	EntityTypeSubnet         EntityType = "subnet"
	EntityTypeTag            EntityType = "tag"
	EntityTypeTicket         EntityType = "ticket"
	EntityTypeToken          EntityType = "token"
	EntityTypeUser           EntityType = "user"
	EntityTypeUserSSHKey     EntityType = "user_ssh_key"
	EntityTypeVolume         EntityType = "volume"
	EntityTypeVPC            EntityType = "vpc"
)

type (
	// EntityType represents the type of a Linode API entity.
	// +k8s:enum
	EntityType string
)
