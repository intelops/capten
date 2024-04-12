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
		// TODO: Add test cases.
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
	type args struct {
		captenConfig config.CaptenConfig
		hc           *helm.Client
		appConfigs   []types.AppConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
