package terralu

import (
	"bytes"

	"github.com/google/uuid"
)

// TerraluImpl is the concrete implementation of the Terralu and TerraformGenerator interfaces
type TerraluImpl struct {
	credentials *TerraluProviderInfo
	buffer      bytes.Buffer
	filename    string
}

// Get returns the credentials and region
func (t *TerraluImpl) GetTerraluProviderInfo() *TerraluProviderInfo {
	return t.credentials
}

func NewTerralu(credentials *TerraluProviderInfo, uuid uuid.UUID) Terralu {
	return &TerraluImpl{
		credentials: credentials,
		filename:    uuid.String(),
		buffer:      bytes.Buffer{},
	}
}

func (t *TerraluImpl) GetTerraluVirtualMachine() TerraformVirtualMachineGenerator {
	return t
}
