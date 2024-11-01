package terralu

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

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
        }`,
			wantErr: false,
		},
		{
			name:  "Nil Provider Info",
			input: nil,
			want: `terraform {
        	  required_providers {
        		mgc = {
        		  source = "magalucloud/mgc"
        		}
        	}
        }`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TerraluImpl{
				credentials: tt.input,
				buffer:      bytes.Buffer{},
			}
			got, err := tr.GenerateTerraformGenericProviderConfig(nil, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTerraformGenericProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Normalize whitespace for comparison
			normalizedGot := strings.ReplaceAll(strings.ReplaceAll(got, "\n", ""), " ", "")
			normalizedWant := strings.ReplaceAll(strings.ReplaceAll(tt.want, "\n", ""), " ", "")
			if normalizedGot != normalizedWant {
				t.Errorf("GenerateTerraformGenericProviderConfig() = %v, want %v", got, tt.want)
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

// // TestGenerateTerraformVirtualMachineConfig tests the GenerateTerraformVirtualMachineConfig method
// func TestTerraluImpl_GenerateTerraformVirtualMachineConfig(t *testing.T) {
// 	type fields struct {
// 		credentials *TerraluProviderInfo
// 		buffer      bytes.Buffer
// 	}
// 	type args struct {
// 		vm    *VirtualMachineInstance
// 		pInfo *TerraluProviderInfo
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "Basic VM Config",
// 			fields: fields{
// 				credentials: &TerraluProviderInfo{
// 					Region:    "us-east-1",
// 					Alias:     "test",
// 					ApiKey:    "access",
// 					KeyID:     "key",
// 					KeySecret: "secret",
// 				},
// 			},
// 			args: args{
// 				vm: &VirtualMachineInstance{
// 					RequiredFields: VirtualMachineRequiredFields{
// 						Name: "test",
// 						MachineType: &MachineTypeSchema{
// 							Name: "t2.micro",
// 						},
// 						Image: &ImageSchema{
// 							Name: "ami-name-123456",
// 						},
// 					},
// 					OptionalFields: VirtualMachineOptionalFields{},
// 				},
// 				pInfo: &TerraluProviderInfo{
// 					Alias:     "test",
// 					Region:    "us-east-1",
// 					ApiKey:    "access",
// 					KeyID:     "key",
// 					KeySecret: "secret",
// 				},
// 			},
// 			want: `resource "mgc_virtual_machine_instances" "test" {
//   provider      = mgc.test
//   name          = "test"
//   machine_type  = {
//     name  = "t2.micro"
//     image         = "ami-name-123456"
//   }
// }`,
// 			wantErr: false,
// 		},
// 		{
// 			name: "VM Config with Optional Fields",
// 			fields: fields{
// 				credentials: &TerraluProviderInfo{
// 					Region:    "eu-central-1",
// 					Alias:     "prod",
// 					ApiKey:    "prod-access",
// 					KeyID:     "prod-key",
// 					KeySecret: "prod-secret",
// 				},
// 			},
// 			args: args{
// 				vm: &VirtualMachineInstance{
// 					RequiredFields: VirtualMachineRequiredFields{
// 						Name: "prod-vm",
// 						MachineType: &MachineTypeSchema{
// 							Name: "m5.large",
// 						},
// 						Image: &ImageSchema{
// 							Name: "ami-prod-789012",
// 						},
// 						SSHKeyName: "prod-keypair",
// 					},
// 					OptionalFields: VirtualMachineOptionalFields{
// 						NameIsPrefix: true,
// 						Network: &NetworkSchema{
// 							AssociatePublicIP: true,
// 							DeletePublicIP:    false,
// 							Interface: &NetworkInterface{
// 								SecurityGroups: []SecurityGroup{
// 									{ID: "sg-12345"},
// 									{ID: "sg-67890"},
// 								},
// 							},
// 							VPC: &VPCSchema{
// 								ID: "vpc-abcdef",
// 							},
// 						},
// 					},
// 				},
// 				pInfo: &TerraluProviderInfo{
// 					Alias:     "prod",
// 					Region:    "eu-central-1",
// 					ApiKey:    "prod-access",
// 					KeyID:     "prod-key",
// 					KeySecret: "prod-secret",
// 				},
// 			},
// 			want: `resource "mgc_virtual_machine_instances" "prod-vm" {
//   provider      = mgc.prod
//   name          = "prod-vm"
//   machine_type  = {
//     name  = "m5.large"
//     image         = "ami-prod-789012"
//     name_is_prefix = true
//     network {
//       associate_public_ip = true
//       interface {
//         security_group_ids = ["sg-12345"]
//         security_group_ids = ["sg-67890"]
//       }
//       vpc_id = "vpc-abcdef"
//       private_address = "10.0.0.10"
//       public_address  = "54.210.123.45"
//       ipv6            = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
//     }
//     ssh_key_name = "prod-keypair"
//   }
// }`,
// 			wantErr: false,
// 		},
// 		{
// 			name: "Template Parsing Error",
// 			fields: fields{
// 				credentials: &TerraluProviderInfo{
// 					Region: "us-east-1",
// 					Alias:  "test",
// 				},
// 			},
// 			args: args{
// 				vm: &VirtualMachineInstance{
// 					RequiredFields: VirtualMachineRequiredFields{
// 						Name: "test",
// 						MachineType: &MachineTypeSchema{
// 							Name: "t2.micro",
// 						},
// 						Image: &ImageSchema{
// 							Name: "ami-name-123456",
// 						},
// 					},
// 					OptionalFields: VirtualMachineOptionalFields{},
// 				},
// 				pInfo: &TerraluProviderInfo{
// 					Alias:     "test",
// 					Region:    "us-east-1",
// 					ApiKey:    "access",
// 					KeyID:     "key",
// 					KeySecret: "secret",
// 				},
// 			},
// 			// Introduce an error by modifying the template to include invalid syntax
// 			want:    "",
// 			wantErr: true,
// 		},
// 		{
// 			name: "Template Execution Error",
// 			fields: fields{
// 				credentials: &TerraluProviderInfo{
// 					Region: "us-east-1",
// 					Alias:  "test",
// 				},
// 			},
// 			args: args{
// 				vm: &VirtualMachineInstance{
// 					RequiredFields: VirtualMachineRequiredFields{
// 						Name: "test",
// 						MachineType: &MachineTypeSchema{
// 							Name: "t2.micro",
// 						},
// 						Image: &ImageSchema{
// 							Name: "ami-name-123456",
// 						},
// 					},
// 					// Inject nil to cause execution error
// 					OptionalFields: VirtualMachineOptionalFields{
// 						Network: nil,
// 					},
// 				},
// 				pInfo: &TerraluProviderInfo{
// 					Alias:     "test",
// 					Region:    "us-east-1",
// 					ApiKey:    "access",
// 					KeyID:     "key",
// 					KeySecret: "secret",
// 				},
// 			},
// 			// Since the template handles nils gracefully, expect no error
// 			want: `resource "mgc_virtual_machine_instances" "test" {
//   provider      = mgc.test
//   name          = "test"
//   machine_type  = {
//     name  = "t2.micro"
//     image         = "ami-name-123456"
//   }
// }`,
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tr := &TerraluImpl{
// 				credentials: tt.fields.credentials,
// 				buffer:      tt.fields.buffer,
// 			}

// 			// If the test is "Template Parsing Error", modify the template to cause a parsing error
// 			if tt.name == "Template Parsing Error" {
// 				originalTemplate := `
// resource "mgc_virtual_machine_instances" "{{ .RequiredFields.Name }}" {
//   provider      = mgc.{{ .Alias }}
//   name          = "{{ .RequiredFields.Name }}"
//   machine_type  = {
//     name  = "{{ .RequiredFields.MachineType.Name }}"
//     image         = "{{ .RequiredFields.Image.Name }}"
//     // Missing closing brace to induce an error
// `
// 				constTerraformTemplate := originalTemplate
// 				// Replace the GenerateTerraformVirtualMachineConfig method's template temporarily
// 				tr.GenerateTerraformVirtualMachineConfig = func(vm *VirtualMachineInstance, pInfo *TerraluProviderInfo) (string, error) {
// 					tmpl, err := template.New("terraform").Parse(constTerraformTemplate)
// 					if err != nil {
// 						return "", fmt.Errorf("error parsing the template: %w", err)
// 					}
// 					err = tmpl.Execute(&tr.buffer, struct {
// 						VirtualMachineInstance
// 						TerraluProviderInfo
// 					}{
// 						VirtualMachineInstance: *vm,
// 						TerraluProviderInfo:    *pInfo,
// 					})
// 					if err != nil {
// 						return "", fmt.Errorf("error executing the template: %w", err)
// 					}
// 					return tr.buffer.String(), nil
// 				}
// 			}

// 			got, err := tr.GenerateTerraformVirtualMachineConfig(tt.args.vm, tt.args.pInfo)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GenerateTerraformVirtualMachineConfig() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if tt.wantErr {
// 				// If an error is expected, no need to check the output
// 				return
// 			}

// 			// Normalize whitespace for comparison
// 			normalizedGot := strings.ReplaceAll(strings.ReplaceAll(got, "\n", ""), " ", "")
// 			normalizedWant := strings.ReplaceAll(strings.ReplaceAll(tt.want, "\n", ""), " ", "")
// 			if normalizedGot != normalizedWant {
// 				t.Errorf("GenerateTerraformVirtualMachineConfig() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
