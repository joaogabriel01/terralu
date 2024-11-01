package terralu

import "bytes"

// TerraluImpl is the concrete implementation of the Terralu and TerraformGenerator interfaces
type TerraluImpl struct {
	credentials *TerraluProviderInfo
	buffer      bytes.Buffer
	filename    string
}

// Set stores the credentials and region
func (t *TerraluImpl) Set(info *TerraluProviderInfo) {
	t.credentials = info
}

// Get returns the credentials and region
func (t *TerraluImpl) Get() *TerraluProviderInfo {
	return t.credentials
}

func NewTerralu(credentials *TerraluProviderInfo, filename string) Terralu {
	return &TerraluImpl{
		credentials: credentials,
		filename:    filename,
		buffer:      bytes.Buffer{},
	}
}

func (t *TerraluImpl) GetTerraluVirtualMachine() TerraformVirtualMachineGenerator{
	return t
}