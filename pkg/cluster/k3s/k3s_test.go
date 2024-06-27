package k3s

import (
	"capten/pkg/config"
	"strings"
	"testing"
)

func Test_getClusterInfo(t *testing.T) {
	validConfig := config.CaptenConfig{}
	type args struct {
		captenConfig config.CaptenConfig
	}
	invalidConfig := config.CaptenConfig{}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid cluster configuration",
			args: args{
				captenConfig: validConfig,
			},
			wantErr: false,
		},
		{
			name: "Invalid cluster configuration",
			args: args{
				captenConfig: invalidConfig,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getClusterInfo(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClusterInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func Test_createOrDestroyCluster(t *testing.T) {

	type args struct {
		captenConfig config.CaptenConfig
		action       string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful creation of a cluster",
			args: args{
				captenConfig: config.CaptenConfig{},
				action:       "create",
			},
			wantErr: false,
		},
		{
			name: "Successful destruction of a cluster",
			args: args{
				captenConfig: config.CaptenConfig{},
				action:       "destroy",
			},
			wantErr: false,
		},
		{
			name: "Error handling when creating a cluster with invalid config",
			args: args{
				captenConfig: config.CaptenConfig{},
				action:       "create",
			},
			wantErr: true,
		},
		{
			name: "Error handling when destroying a cluster with invalid config",
			args: args{
				captenConfig: config.CaptenConfig{},
				action:       "destroy",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrDestroyCluster(tt.args.captenConfig, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("createOrDestroyCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrDestroyCluster(tt.args.captenConfig, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("createOrDestroyCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestCreate(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful creation with valid config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: false,
		},
		{
			name: "Error handling with invalid config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Create(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDestroy(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful destruction of a cluster",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: false,
		},
		{
			name: "Error handling when trying to destroy a non-existing cluster",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Destroy(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateTemplateVarFile(t *testing.T) {
	tests := []struct {
		name           string
		captenConfig   config.CaptenConfig
		clusterInfo    interface{}
		templateFile   string
		expectedErrMsg string
	}{
		{
			name:         "Successful Generation",
			captenConfig: config.CaptenConfig{},
			clusterInfo:  "testClusterInfo",
			templateFile: "testTemplateFile",
		},
		{
			name:           "Error Reading Template File",
			captenConfig:   config.CaptenConfig{},
			clusterInfo:    "testClusterInfo",
			templateFile:   "invalidTemplateFile",
			expectedErrMsg: "failed to read template file",
		},
		{
			name:           "Error Creating Template File",
			captenConfig:   config.CaptenConfig{},
			clusterInfo:    "testClusterInfo",
			templateFile:   "testTemplateFile",
			expectedErrMsg: "failed to create template file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateTemplateVarFile(tt.captenConfig, tt.clusterInfo, tt.templateFile)
			if tt.expectedErrMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error containing '%s', but got: %v", tt.expectedErrMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
