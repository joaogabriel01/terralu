package terralu

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-playground/validator/v10"
)

// GenerateTerraformConfig generates the Terraform generic configuration
func (t *TerraluImpl) GenerateTerraformGenericProviderConfig() (string, error) {
	if t.credentials == nil {
		return "", fmt.Errorf("credentials are not set")
	}
	const terraformTemplate = `terraform {
	required_providers {
		mgc = {
			source = "magalucloud/mgc"
		}
	}
}
provider "mgc" {
	alias    = "{{ .Alias }}"
	region   = "{{ .Region }}"
	api_key  = "{{ .ApiKey }}"
}`

	tmpl, err := template.New("terraform").Parse(terraformTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing the template: %w", err)
	}

	// Execute the template with the provided data
	err = tmpl.Execute(&t.buffer, *t.credentials)
	if err != nil {
		return "", fmt.Errorf("error executing the template: %w", err)
	}
	t.buffer.WriteString("\n")
	return t.buffer.String(), nil
}

// GenerateTerraformConfig generates the Terraform configuration based on the VM and personal information
func (t *TerraluImpl) GenerateTerraformVirtualMachineConfig(vm *VirtualMachineInstance) (string, error) {
	validate := validator.New()
	err := validate.Struct(vm)
	if err != nil {
		return "", fmt.Errorf("error validating the virtual machine instance: %w", err)
	}

	const terraformTemplate = `
resource "mgc_virtual_machine_instances" "{{ .RequiredFields.Name }}" {
  provider      = mgc.{{ .Alias }}
  name          = "{{ .RequiredFields.Name }}"
  machine_type  = {
	name  = "{{ .RequiredFields.MachineType.Name }}"
  }
  image         = {
	name  = "{{ .RequiredFields.Image.Name }}"
  }
  
  {{- if .OptionalFields.NameIsPrefix }}
  name_is_prefix = true
  {{- end }}
  network = {
    associate_public_ip = {{ .OptionalFields.Network.AssociatePublicIP }}
    {{- if .OptionalFields.Network.DeletePublicIP }}
    delete_public_ip    = {{ .OptionalFields.Network.DeletePublicIP }}
    {{- end }}
    {{- if .OptionalFields.Network.Interface }}
    interface {
      {{- range .OptionalFields.Network.Interface.SecurityGroups }}
      security_group_ids = ["{{ .ID }}"]
      {{- end }}
    }
    {{- end }}
    {{- if .OptionalFields.Network.VPC }}
    vpc_id = "{{ .OptionalFields.Network.VPC.ID }}"
    {{- end }}
  }

  ssh_key_name = "{{ .RequiredFields.SSHKeyName }}"
}
`

	tmpl, err := template.New("terraform").Parse(terraformTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing the template: %w", err)
	}

	// Execute the template with the provided data
	err = tmpl.Execute(&t.buffer, struct {
		VirtualMachineInstance
		TerraluProviderInfo
	}{
		VirtualMachineInstance: *vm,
		TerraluProviderInfo:    *t.credentials,
	})
	if err != nil {
		return "", fmt.Errorf("error executing the template: %w", err)
	}
	t.buffer.WriteString("\n")
	return t.buffer.String(), nil
}

// Save saves the buffer content to a file
func (t *TerraluImpl) Save() error {
	if t.buffer.Len() == 0 {
		return fmt.Errorf("buffer is empty, nothing to save")
	}

	// Create or overwrite the file
	file, err := os.Create(t.filename)
	if err != nil {
		return fmt.Errorf("error creating the file: %w", err)
	}
	defer file.Close()

	// Write the buffer to the file
	_, err = file.Write(t.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}

	return nil
}
