package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mocking os package functions
type MockOS struct {
	mock.Mock
}

func (m *MockOS) Getwd() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockOS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

var originalOsGetwd = os.Getwd
var originalOsReadFile = os.ReadFile
var originalOsWriteFile = os.WriteFile

// func resetMocks() {
// 	osGetwd = originalOsGetwd
// 	osReadFile = originalOsReadFile
// 	osWriteFile = originalOsWriteFile
// }

// func TestLoadConfig(t *testing.T) {
// 	// Mock environment variables
// 	os.Setenv("VAULTCRED_HOST_NAME", "vault-cred")
// 	os.Setenv("AGENT_HOST_NAME", "captenagent")
// 	// Add more environment variable settings as needed

// 	cfg, err := LoadConfig()
// 	require.NoError(t, err)
// 	assert.Equal(t, "vault-cred", cfg.VaultCredHostName)
// 	assert.Equal(t, "captenagent", cfg.AgentHostName)
// 	// Add more assertions as needed

// 	// Test error handling in Getwd
// 	mockOS := new(MockOS)
// 	mockOS.On("Getwd").Return("", errors.New("error getting current directory"))
// 	osGetwd = mockOS.Getwd

// 	_, err = LoadConfig()
// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "error getting current directory")
// }

func TestUpdateLBEndpointFile(t *testing.T) {
	mockOS := new(MockOS)
	// osReadFile = mockOS.ReadFile
	// osWriteFile = mockOS.WriteFile
	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("error while fetching current dir", err)
	}
	dirpath, err := getRelativePathUpTo(currentdir)
	if err != nil {
		log.Println("error while fetching relative dir", err)
	}
	tmp := "/" + dirpath + "/config/"
	log.Println("Dir path ", tmp)
	cfg := &CaptenConfig{
		CaptenHostValuesFileName: "capten-lb-endpoint.yaml",
		ConfigDirPath:            tmp,
		//ConfigDirPath: "/home/shifnazarnaz/go/src/github.com/intelops/capten/config/",
	}

	// Mock reading YAML file
	yamlContent := `loadBalancerHost: oldhost`
	mockOS.On("ReadFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName)).Return([]byte(yamlContent), nil)

	// Mock writing YAML file
	mockOS.On("WriteFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName), mock.Anything, os.FileMode(0644)).Return(nil)

	err = UpdateLBEndpointFile(cfg, "newhost")
	require.NoError(t, err)
	assert.Equal(t, "newhost", cfg.LoadBalancerHost)

	// Test error handling in ReadFile
	mockOS.On("ReadFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName)).Return(nil, errors.New("failed to read file"))

	err = UpdateLBEndpointFile(cfg, "newhost")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")

	// Test error handling in WriteFile
	mockOS.On("WriteFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName), mock.Anything, os.FileMode(0644)).Return(errors.New("failed to write file"))

	err = UpdateLBEndpointFile(cfg, "newhost")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write file")
}

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

