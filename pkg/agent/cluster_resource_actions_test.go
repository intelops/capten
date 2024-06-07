package agent

import (
	"capten/pkg/config"
	"reflect"
	"testing"
)

func TestListClusterResources(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		resourceType string
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
			if err := ListClusterResources(tt.args.captenConfig, tt.args.resourceType); (err != nil) != tt.wantErr {
				t.Errorf("ListClusterResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddClusterResource(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		resourceType string
		attributes   map[string]string
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
			if err := AddClusterResource(tt.args.captenConfig, tt.args.resourceType, tt.args.attributes); (err != nil) != tt.wantErr {
				t.Errorf("AddClusterResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateClusterResource(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		resourceType string
		id           string
		attributes   map[string]string
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
			if err := UpdateClusterResource(tt.args.captenConfig, tt.args.resourceType, tt.args.id, tt.args.attributes); (err != nil) != tt.wantErr {
				t.Errorf("UpdateClusterResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_prepareCloudAttributes(t *testing.T) {
	type args struct {
		attributes map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareCloudAttributes(tt.args.attributes)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareCloudAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareCloudAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveClusterResource(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		resourceType string
		id           string
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
			if err := RemoveClusterResource(tt.args.captenConfig, tt.args.resourceType, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveClusterResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
