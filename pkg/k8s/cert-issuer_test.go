package k8s

import (
	"capten/pkg/config"
	"testing"
)

func TestCreateOrUpdateClusterIssuer(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "aws.intelops.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config: missing domain",
			args: args{
				captenConfig: config.CaptenConfig{
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config: missing email",
			args: args{
				captenConfig: config.CaptenConfig{
					CaptenClusterValues: config.CaptenClusterValues{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateOrUpdateClusterIssuer(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateClusterIssuer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
