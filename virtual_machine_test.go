package terralu

import (
	"bytes"
	"strings"
	"testing"
)

func TestTerraluImpl_GenerateTerraformVirtualMachineConfig(t *testing.T) {
	type fields struct {
		credentials *TerraluProviderInfo
		buffer      bytes.Buffer
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
			name: "Test GenerateTerraformVirtualMachineConfig",
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
						}},
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
          image         = "ami-name-123456"
        }`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TerraluImpl{
				credentials: tt.fields.credentials,
				buffer:      tt.fields.buffer,
			}
			got, err := tr.GenerateTerraformVirtualMachineConfig(tt.args.vm, tt.args.pInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("TerraluImpl.GenerateTerraformVirtualMachineConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			normalizedGot := strings.ReplaceAll(strings.ReplaceAll(got, "\n", ""), " ", "")
			normalizedWant := strings.ReplaceAll(strings.ReplaceAll(tt.want, "\n", ""), " ", "")
			if normalizedGot != normalizedWant {
				t.Errorf("TerraluImpl.GenerateTerraformVirtualMachineConfig() = %v, want %v", got, tt.want)
			}

		})
	}
}
