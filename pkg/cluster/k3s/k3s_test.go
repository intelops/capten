package k3s

import (
	"capten/pkg/config"
	//	"reflect"
	"strings"
	"testing"
)

func Test_getClusterInfo(t *testing.T) {
	validConfig := config.CaptenConfig{
		// Add valid configuration fields here
	}
	type args struct {
		captenConfig config.CaptenConfig
	}
	invalidConfig := config.CaptenConfig{
		// Add invalid configuration fields here
	}

	tests := []struct {
		name string
		args args
		//	want    interface{}
		wantErr bool
	}{
		{
			name: "Valid cluster configuration",
			args: args{
				captenConfig: validConfig,
			},
			//		want:    // Add expected result for valid config,
			wantErr: false,
		},
		{
			name: "Invalid cluster configuration",
			args: args{
				captenConfig: invalidConfig,
			},
			//		want:    // Add expected result for invalid config,
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
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("getClusterInfo() = %v, want %v", got, tt.want)
			// }
		})
	}
}
func Test_createOrDestroyCluster(t *testing.T) {
	// Test create or destroy cluster using valid and invalid config
	// and both create and destroy actions
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
			captenConfig: config.CaptenConfig{}, // Fill in with appropriate config
			clusterInfo:  "testClusterInfo",     // Fill in with appropriate cluster info
			templateFile: "testTemplateFile",    // Fill in with appropriate template file
		},
		{
			name:           "Error Reading Template File",
			captenConfig:   config.CaptenConfig{}, // Fill in with appropriate config
			clusterInfo:    "testClusterInfo",     // Fill in with appropriate cluster info
			templateFile:   "invalidTemplateFile", // An invalid template file that should cause an error
			expectedErrMsg: "failed to read template file",
		},
		{
			name:           "Error Creating Template File",
			captenConfig:   config.CaptenConfig{}, // Fill in with appropriate config
			clusterInfo:    "testClusterInfo",     // Fill in with appropriate cluster info
			templateFile:   "testTemplateFile",    // Fill in with appropriate template file
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
