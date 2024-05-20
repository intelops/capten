package agent

import (
	"capten/pkg/agent/pb/vaultcredpb"
	"capten/pkg/config"
	"capten/pkg/types"
	"reflect"
	"testing"

	"github.com/pkg/errors"
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
		{
			name: "Test with empty config and values",
			args: args{
				captenConfig:    config.CaptenConfig{},
				appGlobalVaules: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "Test with valid config and values",
			args: args{
				captenConfig: config.CaptenConfig{
					VaultCredHostName: "vault-cred",

					//	VaultAddress: "http://localhost:8200",
					//	VaultToken:   "s.1234567890",
				},
				appGlobalVaules: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantErr: false,
		},
		{
			name: "Test with invalid config",
			args: args{
				captenConfig: config.CaptenConfig{
					VaultCredHostName: "vault-cred",
					//	VaultToken:   "s.1234567890",
				},
				appGlobalVaules: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantErr: true,
		},
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
		{
			name: "Test with empty config and values",
			args: args{
				captenConfig:    config.CaptenConfig{},
				appGlobalVaules: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "Test with valid config and values",
			args: args{
				captenConfig: config.CaptenConfig{
					VaultCredHostName: "",
				},
				appGlobalVaules: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantErr: false,
		},
		{
			name: "Test with invalid config",
			args: args{
				captenConfig: config.CaptenConfig{
					VaultCredHostName: "",
					//	VaultToken:   "s.1234567890",
				},
				appGlobalVaules: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantErr: true,
		},
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
		{
			name: "Test with empty config and vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{},
				//vaultClient:  vaultcredpb.NewVaultCredClient()
			},
			wantErr: true,
		},
		{
			name: "Test with empty config and valid vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{},
				//vaultClient:  vaultcredpb.VaultCredClient{},
			},
			wantErr: true,
		},
		{
			name: "Test with valid config and empty vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
					ConfigDirPath:      "./testdata/kubeconfig",
				},
				//	vaultClient: vaultcredpb.VaultCredClient{},
			},
			wantErr: true,
		},
		{
			name: "Test with valid config and valid vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
					ConfigDirPath:      "./testdata/kubeconfig",
				},
				//	vaultClient: vaultcredpb.VaultCredClient{},
			},
			wantErr: true,
		},
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
		{
			name: "Test with empty config and empty vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{},
				//vaultClient:  vaultcredpb.V
			},
			wantErr: true,
		},
		{
			name: "Test with empty config and valid vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{},
				vaultClient:  vaultcredpb.NewVaultCredClient(nil),
			},
			wantErr: true,
		},
		{
			name: "Test with valid config and empty vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{
					CaptenGlobalValuesFileName: "global_values.yaml",
					ConfigDirPath:              "./testdata",
				},
				//	vaultClient: vaultcredpb.VaultCredClient{},
			},
			wantErr: true,
		},
		{
			name: "Test with valid config and valid vaultclient",
			args: args{
				captenConfig: config.CaptenConfig{
					CaptenGlobalValuesFileName: "global_values.yaml",
					ConfigDirPath:              "./testdata",
				},
				vaultClient: vaultcredpb.NewVaultCredClient(nil),
			},
			wantErr: false,
		},
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
		{
			name:    "Generate token of length 32",
			want:    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
			wantErr: false,
		},
		{
			name:    "Generate token of length 16",
			want:    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
			wantErr: false,
		},
		{
			name:    "Generate token of length 64",
			want:    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
			wantErr: false,
		},
		{
			name:    "Generate token of length 0",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randomTokenGeneration()
			if (err != nil) != tt.wantErr {
				t.Errorf("randomTokenGeneration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("randomTokenGeneration() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_generateCosignKeyPair_Success(t *testing.T) {
	// Mock the generation of private and public keys
	generateCosignKeyPair := func() ([]byte, []byte, error) {
		return []byte("MockPrivateKey"), []byte("MockPublicKey"), nil
	}

	tests := []struct {
		name    string
		want    []byte
		want1   []byte
		wantErr bool
	}{
		{
			name:    "Test successful key generation",
			want:    []byte("MockPrivateKey"),
			want1:   []byte("MockPublicKey"),
			wantErr: false,
		},
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

func Test_generateCosignKeyPair_Failure(t *testing.T) {
	// Mock the failure in key generation
	generateCosignKeyPair := func() ([]byte, []byte, error) {
		return nil, nil, errors.New("Key generation failed")
	}

	tests := []struct {
		name    string
		want    []byte
		want1   []byte
		wantErr bool
	}{
		{
			name:    "Test key generation failure",
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
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
			if err := configureCosignKeysSecret(tt.args.captenConfig, tt.args.vaultClient, types.CredentialAppConfig{}); (err != nil) != tt.wantErr {
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
			if err := storeCredentials(tt.args.captenConfig, tt.args.appGlobalVaules, tt.args.vaultClient, types.CredentialAppConfig{}); (err != nil) != tt.wantErr {
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
