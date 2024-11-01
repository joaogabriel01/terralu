package terralu

import (
	"bytes"

	"github.com/google/uuid"
)

// TerraluImpl is the concrete implementation of the Terralu and TerraformGenerator interfaces
type TerraluImpl struct {
	credentials *TerraluProviderInfo
	buffer      bytes.Buffer
	dir         string
	mainPath    string
}

// Get returns the credentials and region
func (t *TerraluImpl) GetTerraluProviderInfo() *TerraluProviderInfo {
	return t.credentials
}

// NewTerralu creates a new Terralu instance
func NewTerralu(credentials *TerraluProviderInfo) Terralu {
	impl := &TerraluImpl{
		credentials: credentials,
		dir:         uuid.New().String(),
		buffer:      bytes.Buffer{},
	}
	err := impl.CreateDirectory()
	if err != nil {
		panic(err)
	}
	return impl
}
