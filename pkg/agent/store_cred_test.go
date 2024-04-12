package agent

import (
	"capten/pkg/agent/vaultcredpb"
	"capten/pkg/config"
	"reflect"
	"testing"
)

func TestStoreCredentials(t *testing.T) {
	type args struct {
		captenConfig    config.CaptenConfig
		appGlobalVaules map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StoreCredentials(tt.args.captenConfig, tt.args.appGlobalVaules); (err != nil) != tt.wantErr {
				t.Errorf("StoreCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreClusterCredentials(t *testing.T) {
	type args struct {
		captenConfig    config.CaptenConfig
		appGlobalVaules map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StoreClusterCredentials(tt.args.captenConfig, tt.args.appGlobalVaules); (err != nil) != tt.wantErr {
				t.Errorf("StoreClusterCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storeKubeConfig(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		vaultClient  vaultcredpb.VaultCredClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeKubeConfig(tt.args.captenConfig, tt.args.vaultClient); (err != nil) != tt.wantErr {
				t.Errorf("storeKubeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storeClusterGlobalValues(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		vaultClient  vaultcredpb.VaultCredClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeClusterGlobalValues(tt.args.captenConfig, tt.args.vaultClient); (err != nil) != tt.wantErr {
				t.Errorf("storeClusterGlobalValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_randomTokenGeneration(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randomTokenGeneration()
			if (err != nil) != tt.wantErr {
				t.Errorf("randomTokenGeneration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("randomTokenGeneration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCosignKeyPair(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		want1   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := generateCosignKeyPair()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateCosignKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateCosignKeyPair() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("generateCosignKeyPair() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_configireCosignKeysSecret(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		vaultClient  vaultcredpb.VaultCredClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := configireCosignKeysSecret(tt.args.captenConfig, tt.args.vaultClient); (err != nil) != tt.wantErr {
				t.Errorf("configireCosignKeysSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storeCosignKeys(t *testing.T) {
	type args struct {
		captenConfig    config.CaptenConfig
		appGlobalVaules map[string]interface{}
		vaultClient     vaultcredpb.VaultCredClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeCosignKeys(tt.args.captenConfig, tt.args.appGlobalVaules, tt.args.vaultClient); (err != nil) != tt.wantErr {
				t.Errorf("storeCosignKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storeTerraformStateConfig(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		vaultClient  vaultcredpb.VaultCredClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storeTerraformStateConfig(tt.args.captenConfig, tt.args.vaultClient); (err != nil) != tt.wantErr {
				t.Errorf("storeTerraformStateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
