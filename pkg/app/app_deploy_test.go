package app

import (
	"capten/pkg/config"
	"capten/pkg/helm"
	"capten/pkg/types"

	"reflect"
	"testing"
)

func TestDeployApps(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		globalValues map[string]interface{}
		groupFile    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid arguments",
			args: args{
				captenConfig: config.CaptenConfig{},
				globalValues: map[string]interface{}{},
				groupFile:    "testdata/app_group.yaml",
			},
			wantErr: false,
		},
		{
			name: "Invalid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{},
				globalValues: map[string]interface{}{},
				groupFile:    "testdata/app_group.yaml",
			},
			wantErr: true,
		},
		{
			name: "Invalid groupFile",
			args: args{
				captenConfig: config.CaptenConfig{},
				globalValues: map[string]interface{}{},
				groupFile:    "invalid_file.yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeployApps(tt.args.captenConfig, tt.args.globalValues, tt.args.groupFile); (err != nil) != tt.wantErr {
				t.Errorf("DeployApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_installAppGroup(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			captenConfig config.CaptenConfig
			hc           *helm.Client
			appConfigs   []types.AppConfig
		}
		want bool
	}{
		{
			name: "Empty appConfigs",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs:   []types.AppConfig{},
			},
			want: true,
		},
		{
			name: "AppConfig with privileged namespace",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs: []types.AppConfig{
					{
						PrivilegedNamespace: true,
						Namespace:           "test",
					},
				},
			},
			want: true,
		},
		{
			name: "AppConfig with error on creating namespace",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs: []types.AppConfig{
					{
						PrivilegedNamespace: true,
						Namespace:           "test",
					},
				},
			},
			want: false,
		},
		{
			name: "AppConfig with successful installation",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs: []types.AppConfig{
					{
						Name: "test",
					},
				},
			},
			want: true,
		},
		{
			name: "AppConfig with failed installation",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs: []types.AppConfig{
					{
						Name: "test",
					},
				},
			},
			want: false,
		},
		{
			name: "AppConfig with error on writing appConfig",
			args: struct {
				captenConfig config.CaptenConfig
				hc           *helm.Client
				appConfigs   []types.AppConfig
			}{
				captenConfig: config.CaptenConfig{},
				hc:           &helm.Client{},
				appConfigs: []types.AppConfig{
					{
						Name: "test",
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := installAppGroup(tt.args.captenConfig, tt.args.hc, tt.args.appConfigs); got != tt.want {
				t.Errorf("installAppGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareAppGroupConfigs(t *testing.T) {
	type args struct {
		captenConfig     config.CaptenConfig
		globalValues     map[string]interface{}
		appGroupNameFile string
	}
	tests := []struct {
		name           string
		args           args
		wantAppConfigs []types.AppConfig
		wantErr        bool
	}{
		{
			name: "Empty captenConfig, empty globalValues, empty appGroupNameFile",
			args: args{
				captenConfig:     config.CaptenConfig{},
				globalValues:     map[string]interface{}{},
				appGroupNameFile: "",
			},
			wantAppConfigs: []types.AppConfig{},
			wantErr:        true,
		},
		{
			name: "Empty captenConfig, empty globalValues, non-empty appGroupNameFile",
			args: args{
				captenConfig:     config.CaptenConfig{},
				globalValues:     map[string]interface{}{},
				appGroupNameFile: "test",
			},
			wantAppConfigs: []types.AppConfig{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAppConfigs, err := prepareAppGroupConfigs(tt.args.captenConfig, tt.args.globalValues, tt.args.appGroupNameFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareAppGroupConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAppConfigs, tt.wantAppConfigs) {
				t.Errorf("prepareAppGroupConfigs() = %v, want %v", gotAppConfigs, tt.wantAppConfigs)
			}
		})
	}
}
func Test_replaceOverrideTemplateValues(t *testing.T) {

	type args struct {
		templateData map[string]interface{}
		values       map[string]interface{}
	}

	tests := []struct {
		name                string
		args                args
		wantTransformedData map[string]interface{}
		wantErr             bool
	}{
		{
			name: "valid template with single value",
			args: args{
				templateData: map[string]interface{}{
					"name": "{{ .Name }}",
					"age":  10,
				},
				values: map[string]interface{}{
					"Name": "Alice",
				},
			},
			wantTransformedData: map[string]interface{}{
				"name": "Alice",
				"age":  10,
			},
			wantErr: false,
		},
		{
			name: "valid template with multiple values",
			args: args{
				templateData: map[string]interface{}{
					"name": "{{ .Name }}",
					"age":  "{{ .Age }}",
				},
				values: map[string]interface{}{
					"Name": "Bob",
					"Age":  20,
				},
			},
			wantTransformedData: map[string]interface{}{
				"name": "Bob",
				"age":  20,
			},
			wantErr: false,
		},
		{
			name: "invalid template",
			args: args{
				templateData: map[string]interface{}{
					"name": "{{ .Name }}",
					"age":  10,
				},
				values: map[string]interface{}{
					"Name": "Charlie",
				},
			},
			wantTransformedData: nil,
			wantErr:             true,
		},
		{
			name: "template data is empty",
			args: args{
				templateData: map[string]interface{}{},
				values:       map[string]interface{}{},
			},
			wantTransformedData: map[string]interface{}{},
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTransformedData, err := replaceOverrideTemplateValues(tt.args.templateData, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceOverrideTemplateValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTransformedData, tt.wantTransformedData) {
				t.Errorf("replaceOverrideTemplateValues() = %v, want %v", gotTransformedData, tt.wantTransformedData)
			}
		})
	}
}

func Test_replaceTemplateStringValues(t *testing.T) {
	type args struct {
		templateStringData string
		values             map[string]interface{}
	}
	tests := []struct {
		name                      string
		args                      args
		wantTransformedStringData string
		wantErr                   bool
	}{
		{
			name: "Valid template with single value",
			args: args{
				templateStringData: "Hello {{ .Name }}!",
				values: map[string]interface{}{
					"Name": "Alice",
				},
			},
			wantTransformedStringData: "Hello Alice!",
			wantErr:                   false,
		},
		{
			name: "Valid template with multiple values",
			args: args{
				templateStringData: "{{ .Name }} is {{ .Age }} years old.",
				values: map[string]interface{}{
					"Name": "Bob",
					"Age":  25,
				},
			},
			wantTransformedStringData: "Bob is 25 years old.",
			wantErr:                   false,
		},
		{
			name: "Invalid template",
			args: args{
				templateStringData: "{{ .Name }} is {{ .Age }} years old.",
				values: map[string]interface{}{
					"Name": "Charlie",
				},
			},
			wantTransformedStringData: "",
			wantErr:                   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTransformedStringData, err := replaceTemplateStringValues(tt.args.templateStringData, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceTemplateStringValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTransformedStringData != tt.wantTransformedStringData {
				t.Errorf("replaceTemplateStringValues() = %v, want %v", gotTransformedStringData, tt.wantTransformedStringData)
			}
		})
	}
}
