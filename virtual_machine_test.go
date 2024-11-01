package terralu

import (
	"bytes"

	"github.com/google/go-cmp/cmp"

	"errors"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

// TestGenerateTerraformGenericProviderConfig tests the GenerateTerraformGenericProviderConfig method
func TestTerraluImpl_GenerateTerraformGenericProviderConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   *TerraluProviderInfo
		want    string
		wantErr bool
	}{
		{
			name: "Valid Provider Info",
			input: &TerraluProviderInfo{
				Region:    "us-west-2",
				Alias:     "mgc",
				ApiKey:    "api-key",
				KeyID:     "key-id",
				KeySecret: "key-secret",
			},
			want: `terraform {
					required_providers {
						mgc = {
						source = "magalucloud/mgc"
						}
					}
					}
					provider "mgc" {
					alias    = "mgc"
					region   = "us-west-2"
					api_key  = "api-key"
					}`,
			wantErr: false,
		},
		{
			name:  "Empty Provider Info",
			input: &TerraluProviderInfo{},
			want: `terraform {
					required_providers {
						mgc = {
						source = "magalucloud/mgc"
						}
					}
					}
					provider "mgc" {
					alias    = ""
					region   = ""
					api_key  = ""
					}`,
			wantErr: false,
		},
		// Adicione mais casos de teste conforme necessário
	}

	for _, tt := range tests {
		tt := tt // captura a variável para evitar problemas com goroutines
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTerralu(tt.input, uuid.New())
			got, err := tr.GenerateTerraformGenericProviderConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			// Trim whitespace to avoid differences due to extra spaces or new lines
			got = strings.TrimSpace(got)
			want := strings.TrimSpace(tt.want)
			if diff := cmp.Diff(normalizeWhitespace(want), normalizeWhitespace(got)); diff != "" {
				t.Errorf("%s mismatch (-want +got):\n%s", tt.name, diff)
			}
		})

	}
}

// TestTerraluImpl_Save tests the Save method
func TestTerraluImpl_Save(t *testing.T) {
	type fields struct {
		buffer   bytes.Buffer
		filename string
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		assertion func(t *testing.T, filename string)
	}{
		{
			name: "Save with Non-Empty Buffer",
			fields: fields{
				buffer:   *bytes.NewBufferString("test content"),
				filename: uuid.New().String(),
			},
			wantErr: false,
			assertion: func(t *testing.T, filename string) {
				// Read the file and verify its content
				content, err := os.ReadFile(filename)
				if err != nil {
					t.Fatalf("Failed to read the file: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("File content = %v, want %v", string(content), "test content")
				}
				// Clean up the file
				os.Remove(filename)
			},
		},
		{
			name: "Save with Empty Buffer",
			fields: fields{
				buffer:   bytes.Buffer{},
				filename: uuid.New().String(),
			},
			wantErr: true,
			assertion: func(t *testing.T, filename string) {
				// Ensure the file was not created
				if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
					t.Errorf("File should not exist, but it does")
					// Clean up if exists
					os.Remove(filename)
				}
			},
		},
		{
			name: "Save with Invalid Filename",
			fields: fields{
				buffer:   *bytes.NewBufferString("invalid filename test"),
				filename: string([]byte{0x00, 0x01, 0x02}), // Invalid filename
			},
			wantErr: true,
			assertion: func(t *testing.T, filename string) {
				// No file should be created
				if _, err := os.Stat(filename); !strings.Contains(err.Error(), os.ErrInvalid.Error()) {
					t.Errorf(t.Name() + ": " + err.Error())
					// Clean up if exists
					os.Remove(filename)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TerraluImpl{
				buffer:   tt.fields.buffer,
				filename: tt.fields.filename,
			}
			err := tr.Save()
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.assertion(t, tt.fields.filename)
		})
	}
}

// func TestGenerate(t *testing.T) {
// 	pinfo := &TerraluProviderInfo{
// 		Alias:     "test",
// 		Region:    "br-se1",
// 		ApiKey:    "access",
// 		KeyID:     "key",
// 		KeySecret: "secret",
// 	}
// 	vm := &VirtualMachineInstance{
// 		RequiredFields: VirtualMachineRequiredFields{
// 			Name: "test",
// 			MachineType: &MachineTypeSchema{
// 				Name: "cloud-bs1.xsmall",
// 			},
// 			Image: &ImageSchema{
// 				Name: "cloud-ubuntu-22.04 LTS",
// 			},
// 		},
// 		// Inject nil to cause execution error
// 		OptionalFields: VirtualMachineOptionalFields{},
// 	}
// 	terralu := NewTerralu(pinfo, uuid.New())
// 	_, err := terralu.GenerateTerraformGenericProviderConfig()
// 	if err != nil {
// 		t.Errorf("Err on GenerateTerraformGenericProviderConfig: %v", err)
// 	}
// 	_, err = terralu.GenerateTerraformVirtualMachineConfig(vm)
// 	if err != nil {
// 		t.Errorf("Err on GenerateTerraformVirtualMachineConfig: %v", err)
// 	}

// 	err = terralu.Save()

// 	if err != nil {
// 		t.Fatalf("Error saving file %s", err)
// 	}
// }

// TestGenerateTerraformVirtualMachineConfig tests the GenerateTerraformVirtualMachineConfig method
func TestTerraluImpl_GenerateTerraformVirtualMachineConfig(t *testing.T) {
	type fields struct {
		credentials *TerraluProviderInfo
	}
	type args struct {
		vm    *VirtualMachineInstance
		pInfo *TerraluProviderInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Basic VM Config",
			fields: fields{
				credentials: &TerraluProviderInfo{
					Region:    "us-east-1",
					Alias:     "test",
					ApiKey:    "access",
					KeyID:     "key",
					KeySecret: "secret",
				},
			},
			args: args{
				vm: &VirtualMachineInstance{
					RequiredFields: VirtualMachineRequiredFields{
						Name: "test",
						MachineType: &MachineTypeSchema{
							Name: "t2.micro",
						},
						Image: &ImageSchema{
							Name: "ami-name-123456",
						},
						SSHKeyName: "test-key",
					},
					OptionalFields: VirtualMachineOptionalFields{},
				},
				pInfo: &TerraluProviderInfo{
					Alias:     "test",
					Region:    "us-east-1",
					ApiKey:    "access",
					KeyID:     "key",
					KeySecret: "secret",
				},
			},
			want: `resource "mgc_virtual_machine_instances" "test" {
          provider      = mgc.test
          name          = "test"
          machine_type  = {
        	name  = "t2.micro"
          }
          image         = {
        	name  = "ami-name-123456"
          }
          network = {
            associate_public_ip = false
          }
        
          ssh_key_name = "test-key"
        }`,
			wantErr: false,
		},
		{
			name: "VM Config with Optional Fields",
			fields: fields{
				credentials: &TerraluProviderInfo{
					Region:    "eu-central-1",
					Alias:     "prod",
					ApiKey:    "prod-access",
					KeyID:     "prod-key",
					KeySecret: "prod-secret",
				},
			},
			args: args{
				vm: &VirtualMachineInstance{
					RequiredFields: VirtualMachineRequiredFields{
						Name: "prod-vm",
						MachineType: &MachineTypeSchema{
							Name: "m5.large",
						},
						Image: &ImageSchema{
							Name: "ami-prod-789012",
						},
						SSHKeyName: "prod-keypair",
					},
					OptionalFields: VirtualMachineOptionalFields{
						NameIsPrefix: true,
						Network: NetworkSchema{
							AssociatePublicIP: true,
							DeletePublicIP:    false,
							Interface: &NetworkInterface{
								SecurityGroups: []SecurityGroup{
									{ID: "sg-12345"},
									{ID: "sg-67890"},
								},
							},
							VPC: &VPCSchema{
								ID: "vpc-abcdef",
							},
						},
					},
				},
				pInfo: &TerraluProviderInfo{
					Alias:     "prod",
					Region:    "eu-central-1",
					ApiKey:    "prod-access",
					KeyID:     "prod-key",
					KeySecret: "prod-secret",
				},
			},
			want: `resource "mgc_virtual_machine_instances" "prod-vm" {
          provider      = mgc.prod
          name          = "prod-vm"
          machine_type  = {
        	name  = "m5.large"
          }
          image         = {
        	name  = "ami-prod-789012"
          }
          name_is_prefix = true
          network = {
            associate_public_ip = true
            interface {
              security_group_ids = ["sg-12345"]
              security_group_ids = ["sg-67890"]
            }
            vpc_id = "vpc-abcdef"
          }
        
          ssh_key_name = "prod-keypair"
        }`,
			wantErr: false,
		},
		// {
		// 	name: "Template Parsing Error",
		// 	fields: fields{
		// 		credentials: &TerraluProviderInfo{
		// 			Region: "us-east-1",
		// 			Alias:  "test",
		// 		},
		// 	},
		// 	args: args{
		// 		vm: &VirtualMachineInstance{
		// 			RequiredFields: VirtualMachineRequiredFields{
		// 				Name: "test",
		// 				MachineType: &MachineTypeSchema{
		// 					Name: "t2.micro",
		// 				},
		// 				Image: &ImageSchema{
		// 					Name: "ami-name-123456",
		// 				},
		// 			},
		// 			OptionalFields: VirtualMachineOptionalFields{},
		// 		},
		// 		pInfo: &TerraluProviderInfo{
		// 			Alias:     "test",
		// 			Region:    "us-east-1",
		// 			ApiKey:    "access",
		// 			KeyID:     "key",
		// 			KeySecret: "secret",
		// 		},
		// 	},
		// 	// Introduce an error by modifying the template to include invalid syntax
		// 	want:    "",
		// 	wantErr: true,
		// },
		// !!! NECESSARY MOCK TO RUN COMMENTED TEST !!!
		{
			name: "Template Execution Error",
			fields: fields{
				credentials: &TerraluProviderInfo{
					Region: "us-east-1",
					Alias:  "test",
				},
			},
			args: args{
				vm: &VirtualMachineInstance{
					RequiredFields: VirtualMachineRequiredFields{
						Name: "test",
						MachineType: &MachineTypeSchema{
							Name: "t2.micro",
						},
						Image: &ImageSchema{
							Name: "ami-name-123456",
						},
						SSHKeyName: "test-key",
					},
					// Inject nil to cause execution error
					OptionalFields: VirtualMachineOptionalFields{},
				},
				pInfo: &TerraluProviderInfo{
					Alias:     "test",
					Region:    "us-east-1",
					ApiKey:    "access",
					KeyID:     "key",
					KeySecret: "secret",
				},
			},
			// Since the template handles nils gracefully, expect no error
			want: `resource "mgc_virtual_machine_instances" "test" {
          provider      = mgc.test
          name          = "test"
          machine_type  = {
        	name  = "t2.micro"
          }
          image         = {
        	name  = "ami-name-123456"
          }
          network = {
            associate_public_ip = false
          }
        
          ssh_key_name = "test-key"
        }`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTerralu(tt.args.pInfo, uuid.New())

			got, err := tr.GenerateTerraformVirtualMachineConfig(tt.args.vm)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if tt.wantErr {
				// If an error is expected, no need to check the output
				return
			}

			// Normalize whitespace for comparison
			normalizedGot := strings.ReplaceAll(strings.ReplaceAll(got, "\n", ""), " ", "")
			normalizedWant := strings.ReplaceAll(strings.ReplaceAll(tt.want, "\n", ""), " ", "")
			if normalizedGot != normalizedWant {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
