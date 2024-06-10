package agent

import (
	"capten/pkg/agent/pb/agentpb"
	"capten/pkg/agent/pb/vaultcredpb"
	"capten/pkg/config"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials"
)

func TestGetAgentClient(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	type args struct {
		config      config.CaptenConfig
		clusterconf config.CaptenClusterValues
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
					AgentSecure:        true,
					AgentHostName:      "captenagent",
					CertDirPath:        "/" + presentdir + "/cert/",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
				},
				clusterconf: config.CaptenClusterValues{
					DomainName: "aws.optimizor.app",
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
				clusterconf: config.CaptenClusterValues{
					DomainName: "aws.optimizor.app",
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

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	dir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

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
					AgentSecure:        true,
					AgentHostName:      "captenagent",
					CertDirPath:        "/" + dir + "/cert/",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
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

// func TestLoadTLSCredentials(t *testing.T) {
// 	// Test case 1: LoadX509KeyPair fails

// 	currentdir, err := os.Getwd()
// 	if err != nil {
// 		log.Println("Error while getting cuerent dir", err)
// 	}
// 	dir, err := getRelativePathUpTo(currentdir)

// 	if err != nil {
// 		log.Println("Error while getting working dir", err)
// 	}

// 	captenConfig := config.CaptenConfig{
// 		CertDirPath:        "/" + dir + "/cert",
// 		ClientCertFileName: "client.crt",
// 		ClientKeyFileName:  "client.key",
// 		CAFileName:         "ca.crt",
// 	}
// 	os.MkdirAll(captenConfig.CertDirPath, os.ModePerm)
// 	defer os.RemoveAll(captenConfig.CertDirPath)

// 	_, err = loadTLSCredentials(captenConfig)
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}

// 	// Test case 2: AppendCertsFromPEM fails
// 	certFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName)
// 	keyFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName)
// 	caFile := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName)

// 	err = os.WriteFile(certFile, []byte("dummy cert"), os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write cert file: %v", err)
// 	}

// 	err = os.WriteFile(keyFile, []byte("dummy key"), os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write key file: %v", err)
// 	}

// 	err = os.WriteFile(caFile, []byte("dummy ca"), os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write ca file: %v", err)
// 	}

// 	_, err = loadTLSCredentials(captenConfig)
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}

// 	// Test case 3: Successful load
// 	certPEM, err := os.ReadFile("cert/cert.pem")
// 	if err != nil {
// 		t.Fatalf("Failed to read cert file: %v", err)
// 	}

// 	keyPEM, err := os.ReadFile("cert/key.pem")
// 	if err != nil {
// 		t.Fatalf("Failed to read key file: %v", err)
// 	}

// 	caPEM, err := os.ReadFile("cert/ca.pem")
// 	if err != nil {
// 		t.Fatalf("Failed to read ca file: %v", err)
// 	}

// 	err = os.WriteFile(certFile, certPEM, os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write cert file: %v", err)
// 	}

// 	err = os.WriteFile(keyFile, keyPEM, os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write key file: %v", err)
// 	}

// 	err = os.WriteFile(caFile, caPEM, os.ModePerm)
// 	if err != nil {
// 		t.Fatalf("Failed to write ca file: %v", err)
// 	}

// 	_, err = loadTLSCredentials(captenConfig)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// }

func getRelativePathUpTo(currentPath string) (string, error) {
	targetDir := "capten"
	// Split the path into parts
	parts := strings.Split(currentPath, string(filepath.Separator))

	// Traverse the path parts and look for the target directory
	for i, part := range parts {
		if part == targetDir {
			// Join the parts up to the target directory
			return filepath.Join(parts[:i+1]...), nil
		}
	}

	return "", fmt.Errorf("directory %s not found in path %s", targetDir, currentPath)
}

func Test_loadTLSCredentials(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	dir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		want    credentials.TransportCredentials
		wantErr bool
	}{
		{
			name: "valid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{
					AgentCertFileName: "agent.crt",
					AgentKeyFileName:  "agent.key",
					CAFileName:        "ca.crt",
					//CertDirPath: ,
					CertDirPath:        "/" + dir + "/cert/",
					ClientCertFileName: "client.crt",
					ClientKeyFileName:  "client.key",
				},
			},
			want: credentials.NewTLS(&tls.Config{
				// Certificates: []tls.Certificate{cert},
				// RootCAs:      caCertPool,
				ClientAuth: tls.RequireAnyClientCert,
				// ClientCAs:    caCertPool,
			}),

			wantErr: true,
		},
		{
			name: "invalid cert file",
			args: args{
				captenConfig: config.CaptenConfig{
					AgentCertFileName:  "client.crt",
					AgentKeyFileName:   "client.key",
					CAFileName:         "ca.crt",
					CertDirPath:        "/" + dir + "/cert/",
					ClientCertFileName: "dbcjd.key",
					ClientKeyFileName:  "client.key",
				},
			},
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name: "invalid key file",
		// 	args: args{
		// 		captenConfig: config.CaptenConfig{
		// 			AgentCertFileName:  "agent.crt",
		// 			AgentKeyFileName:   "client.key",
		// 			CAFileName:         "ca.crt",
		// 			CertDirPath:        "/" + dir + "/cert/",
		// 			ClientCertFileName: "client.crt",
		// 			ClientKeyFileName:  "agent.key",
		// 		},
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
		// {
		// 	name: "invalid ca file",
		// 	args: args{
		// 		captenConfig: config.CaptenConfig{
		// 			AgentCertFileName:  "agent.crt",
		// 			AgentKeyFileName:   "agent.key",
		// 			CAFileName:         "root.crt",
		// 			CertDirPath:        "/" + dir + "/cert/",
		// 			ClientCertFileName: "client.crt",
		// 			ClientKeyFileName:  "client.key",
		// 		},
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadTLSCredentials(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadTLSCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("loadTLSCredentials() = %v, want %v", got, tt.want)
			// }
		})
	}
}
