package terralu

type TerraluProviderInfo struct {
	Alias     string `validate:"required"`
	Region    string `validate:"required"`
	ApiKey    string `validate:"required"`
	KeyID     string
	KeySecret string
}

// VirtualMachineInstance represents the VM instance with required and optional fields
type VirtualMachineInstance struct {
	RequiredFields VirtualMachineRequiredFields
	OptionalFields VirtualMachineOptionalFields
}

// VirtualMachineRequiredFields defines fields that are mandatory for creating a virtual machine
type VirtualMachineRequiredFields struct {
	Name        string             `validate:"required"`
	MachineType *MachineTypeSchema `validate:"required"`
	Image       *ImageSchema       `validate:"required"`
	SSHKeyName  string             `validate:"required"`
}

// VirtualMachineOptionalFields holds optional fields for VM creation
type VirtualMachineOptionalFields struct {
	NameIsPrefix bool
	Network      NetworkSchema
}

// ImageSchema represents the nested schema for image configuration
type ImageSchema struct {
	Name string `validate:"required"`
}

// MachineTypeSchema represents the nested schema for machine type configuration
type MachineTypeSchema struct {
	Name string `validate:"required"`
}

// NetworkSchema holds the network configuration details
type NetworkSchema struct {
	AssociatePublicIP bool
	DeletePublicIP    bool
	Interface         *NetworkInterface
	VPC               *VPCSchema
}

// NetworkInterface represents the configuration of network interface
type NetworkInterface struct {
	SecurityGroups []SecurityGroup
}

// SecurityGroup represents a security group associated with a network interface
type SecurityGroup struct {
	ID string
}

// VPCSchema represents the VPC configuration for the network
type VPCSchema struct {
	ID   string
	Name string
}
