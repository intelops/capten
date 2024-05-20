package agent

import (
	"capten/pkg/agent/pb/agentpb"
	"capten/pkg/agent/pb/vaultcredpb"
	"capten/pkg/config"
	"os"

	//"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{
			name: "Secure connection",
			args: args{
				config: config.CaptenConfig{
					AgentSecure:   true,
					AgentHostName: "captenagent",
					// DomainName:      "com",
					// CaptenAgentPort: "50051",
				},
			},
			wantErr: false,
		},
		{
			name: "Insecure connection",
			args: args{
				config: config.CaptenConfig{
					AgentSecure:   false,
					AgentHostName: "captenagent",
					// DomainName:      "com",
					// CaptenAgentPort: "50051",
				},
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAgentClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAgentClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
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
		{
			name: "Secure connection",
			args: args{
				config: config.CaptenConfig{
					AgentSecure:   true,
					AgentHostName: "captenagent",
				},
			},
			wantErr: false,
		},
		{
			name: "Insecure connection",
			args: args{
				config: config.CaptenConfig{
					AgentSecure:   false,
					AgentHostName: "captenagent",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing host name",
			args: args{
				config: config.CaptenConfig{
					AgentSecure: true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetVaultClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVaultClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}

}

func TestLoadTLSCredentials(t *testing.T) {
	// Test case 1: LoadX509KeyPair fails
	captenConfig := config.CaptenConfig{
		//CertDirPath:    "/path/to/certs",
		ClientCertFileName: "client.crt",
		ClientKeyFileName:  "client.key",
		CAFileName:         "ca.crt",
	}
	os.MkdirAll(captenConfig.CertDirPath, os.ModePerm)
	defer os.RemoveAll(captenConfig.CertDirPath)

	_, err := loadTLSCredentials(captenConfig)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case 2: AppendCertsFromPEM fails
	certFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName)
	keyFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName)
	caFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName)

	err = os.WriteFile(certFile, []byte("dummy cert"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write cert file: %v", err)
	}

	err = os.WriteFile(keyFile, []byte("dummy key"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	err = os.WriteFile(caFile, []byte("dummy ca"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write ca file: %v", err)
	}

	_, err = loadTLSCredentials(captenConfig)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case 3: Successful load
	certPEM, err := os.ReadFile("testdata/cert.pem")
	if err != nil {
		t.Fatalf("Failed to read cert file: %v", err)
	}

	keyPEM, err := os.ReadFile("testdata/key.pem")
	if err != nil {
		t.Fatalf("Failed to read key file: %v", err)
	}

	caPEM, err := os.ReadFile("testdata/ca.pem")
	if err != nil {
		t.Fatalf("Failed to read ca file: %v", err)
	}

	err = os.WriteFile(certFile, certPEM, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write cert file: %v", err)
	}

	err = os.WriteFile(keyFile, keyPEM, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	err = os.WriteFile(caFile, caPEM, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write ca file: %v", err)
	}

	_, err = loadTLSCredentials(captenConfig)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
