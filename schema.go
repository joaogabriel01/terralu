package terralu

type TerraluProviderInfo struct {
	Alias     string `json:"alias" validate:"required"`
	Region    string `json:"region" validate:"required"`
	ApiKey    string `json:"api_key" validate:"required"`
	KeyID     string `json:"key_id"`
	KeySecret string `json:"key_secret"`
}

// VirtualMachineInstance represents the VM instance with required and optional fields
type VirtualMachineInstance struct {
	RequiredFields VirtualMachineRequiredFields `json:",inline"`
	OptionalFields VirtualMachineOptionalFields `json:",inline"`
}

// VirtualMachineRequiredFields defines fields that are mandatory for creating a virtual machine
type VirtualMachineRequiredFields struct {
	Name        string             `json:"name" validate:"required"`
	MachineType *MachineTypeSchema `json:"machine_type" validate:"required"`
	Image       *ImageSchema       `json:"image" validate:"required"`
	SSHKeyName  *string            `json:"ssh_key_name,omitempty"`
}

// VirtualMachineOptionalFields holds optional fields for VM creation
type VirtualMachineOptionalFields struct {
	NameIsPrefix bool           `json:"name_is_prefix,omitempty"`
	Network      *NetworkSchema `json:"network,omitempty"`
}

// ImageSchema represents the nested schema for image configuration
type ImageSchema struct {
	Name string `json:"name" validate:"required"`
}

// MachineTypeSchema represents the nested schema for machine type configuration
type MachineTypeSchema struct {
	Name string `json:"name" validate:"required"`
}

// NetworkSchema holds the network configuration details
type NetworkSchema struct {
	AssociatePublicIP bool              `json:"associate_public_ip" validate:"required"`
	DeletePublicIP    *bool             `json:"delete_public_ip,omitempty"`
	Interface         *NetworkInterface `json:"interface,omitempty"`
	VPC               *VPCSchema        `json:"vpc,omitempty"`
}

// NetworkInterface represents the configuration of network interface
type NetworkInterface struct {
	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`
}

// SecurityGroup represents a security group associated with a network interface
type SecurityGroup struct {
	ID string `json:"id,omitempty"`
}

// VPCSchema represents the VPC configuration for the network
type VPCSchema struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
