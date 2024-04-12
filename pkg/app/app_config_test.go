package app

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"reflect"
	"testing"
)

func TestGetClusterGlobalValues(t *testing.T) {
	type args struct {
		valuesFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClusterGlobalValues(tt.args.valuesFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterGlobalValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClusterGlobalValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApps(t *testing.T) {
	type args struct {
		appListFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApps(tt.args.appListFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppConfig(t *testing.T) {
	type args struct {
		appConfigFilePath string
		globalValues      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    types.AppConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppConfig(tt.args.appConfigFilePath, tt.args.globalValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppValuesTemplate(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		appName      string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppValuesTemplate(tt.args.captenConfig, tt.args.appName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppValuesTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteAppConfig(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		appConfig    types.AppConfig
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
			if err := WriteAppConfig(tt.args.captenConfig, tt.args.appConfig); (err != nil) != tt.wantErr {
				t.Errorf("WriteAppConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrepareGlobalVaules(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareGlobalVaules(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareGlobalVaules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareGlobalVaules() = %v, want %v", got, tt.want)
			}
		})
	}
}
