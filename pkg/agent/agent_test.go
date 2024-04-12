package agent

import (
	"capten/pkg/agent/agentpb"
	"capten/pkg/agent/vaultcredpb"
	"capten/pkg/config"
	"reflect"
	"testing"

	"google.golang.org/grpc/credentials"
)

func TestGetAgentClient(t *testing.T) {
	type args struct {
		config config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    agentpb.AgentClient
		wantErr bool
	}{

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAgentClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAgentClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAgentClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadTLSCredentials(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    credentials.TransportCredentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadTLSCredentials(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadTLSCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadTLSCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVaultClient(t *testing.T) {
	type args struct {
		config config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    vaultcredpb.VaultCredClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVaultClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVaultClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVaultClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
