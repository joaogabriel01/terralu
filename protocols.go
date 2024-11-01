package terralu

// Terralu is an interface for managing virtual machine instances
type Terralu interface {
	TerraformGenerator
	Save() error
}

// TerraformGenerator defines the contract for generating Terraform code
type TerraformGenerator interface {
	TerraluCredentialsAndRegion
	GenerateTerraformGenericProviderConfig(vm *VirtualMachineInstance, pInfo *TerraluProviderInfo) (string, error)
	TerraformVirtualMachineGenerator
}

// TerraformVirtualMachineGenerator defines the contract for generating Terraform configuration for virtual machines
type TerraformVirtualMachineGenerator interface {
	GenerateTerraformVirtualMachineConfig(vm *VirtualMachineInstance, pInfo *TerraluProviderInfo) (string, error)
}

// TerraluCredentialsAndRegion manages personal info for API authorization and region setting
type TerraluCredentialsAndRegion interface {
	Set(*TerraluProviderInfo)
	Get() *TerraluProviderInfo
}
