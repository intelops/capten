package k3s

import (
	"capten/pkg/config"
	"reflect"
	"testing"
)

func Test_getClusterInfo(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getClusterInfo(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClusterInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClusterInfo() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	type args struct {
		captenConfig     config.CaptenConfig
		clusterInfo      interface{}
		templateFileName string
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
			if err := generateTemplateVarFile(tt.args.captenConfig, tt.args.clusterInfo, tt.args.templateFileName); (err != nil) != tt.wantErr {
				t.Errorf("generateTemplateVarFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
