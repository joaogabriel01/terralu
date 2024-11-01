package terralu

import (
	"fmt"
	"os"
	"text/template"
)

// GenerateTerraformConfig generates the Terraform generic configuration
func (t *TerraluImpl) GenerateTerraformGenericProviderConfig(vm *VirtualMachineInstance, pInfo *TerraluProviderInfo) (string, error) {
	const terraformTemplate = `terraform {
	  required_providers {
		mgc = {
		  source = "magalucloud/mgc"
		}
	}
}`
	return terraformTemplate, nil
}

// GenerateTerraformConfig generates the Terraform configuration based on the VM and personal information
func (t *TerraluImpl) GenerateTerraformVirtualMachineConfig(vm *VirtualMachineInstance, pInfo *TerraluProviderInfo) (string, error) {
	const terraformTemplate = `
resource "mgc_virtual_machine_instances" "{{ .RequiredFields.Name }}" {
  provider      = mgc.{{ .Alias }}
  name          = "{{ .RequiredFields.Name }}"
  machine_type  = {
	name  = "{{ .RequiredFields.MachineType.Name }}"
  image         = "{{ .RequiredFields.Image.Name }}"
  
  {{- if .OptionalFields.NameIsPrefix }}
  name_is_prefix = true
  {{- end }}

  {{- if .OptionalFields.Network }}
  network {
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
    private_address = "{{ .OptionalFields.Network.PrivateAddress }}"
    public_address  = "{{ .OptionalFields.Network.PublicAddress }}"
    ipv6            = "{{ .OptionalFields.Network.IPv6 }}"
  }
  {{- end }}

  {{- if .RequiredFields.SSHKeyName }}
  ssh_key_name = "{{ .OptionalFields.SSHKeyName }}"
  {{- end }}
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
		TerraluProviderInfo:    *pInfo,
	})
	if err != nil {
		return "", fmt.Errorf("error executing the template: %w", err)
	}

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
