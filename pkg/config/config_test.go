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

func TestUpdateLBEndpointFile(t *testing.T) {
	mockOS := new(MockOS)
	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("error while fetching current dir", err)
	}
	dirpath, err := getRelativePathUpTo(currentdir)
	if err != nil {
		log.Println("error while fetching relative dir", err)
	}
	tmp := "/" + dirpath + "/config/"
	cfg := &CaptenConfig{
		CaptenHostValuesFileName: "capten-lb-endpoint.yaml",
		ConfigDirPath:            tmp,
	}

	yamlContent := `loadBalancerHost: oldhost`
	mockOS.On("ReadFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName)).Return([]byte(yamlContent), nil)

	mockOS.On("WriteFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName), mock.Anything, os.FileMode(0644)).Return(nil)

	err = UpdateLBEndpointFile(cfg, "newhost", "")
	require.NoError(t, err)
	assert.Equal(t, "newhost", cfg.LoadBalancerHost)

	mockOS.On("ReadFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName)).Return(nil, errors.New("failed to read file"))

	err = UpdateLBEndpointFile(cfg, "newhost", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")

	mockOS.On("WriteFile", cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenHostValuesFileName), mock.Anything, os.FileMode(0644)).Return(errors.New("failed to write file"))

	err = UpdateLBEndpointFile(cfg, "newhost", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write file")
}

func getRelativePathUpTo(currentPath string) (string, error) {
	targetDir := "capten"
	parts := strings.Split(currentPath, string(filepath.Separator))

	for i, part := range parts {
		if part == targetDir {
			return filepath.Join(parts[:i+1]...), nil
		}
	}

	return "", fmt.Errorf("directory %s not found in path %s", targetDir, currentPath)
}
