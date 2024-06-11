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
		{
			name: "list git-project",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "git-project",
			},
			wantErr: false,
		},
		{
			name: "list cloud-provider",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "cloud-provider",
			},
			wantErr: false,
		},
		{
			name: "list container-registry",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "container-registry",
			},
			wantErr: false,
		},
		{
			name: "list unknown-resource",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "unknown-resource",
			},
			wantErr: true,
		},
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
		{
			name: "add git-project",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "git-project",
				attributes: map[string]string{
					"git-project-url": "https://github.com/example/example-project.git",
					"labels":          "label1,label2",
					"access-token":    "testAccessToken",
					"user-id":         "testUserId",
				},
			},
			wantErr: false,
		},
		{
			name: "add cloud-provider",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "cloud-provider",
				attributes: map[string]string{
					"cloud-type":  "aws",
					"labels":      "label1,label2",
					"cloud-key":   "testCloudKey",
					"cloud-token": "testCloudToken",
				},
			},
			wantErr: false,
		},
		{
			name: "add container-registry",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "container-registry",
				attributes: map[string]string{
					"registry-url":      "https://example.com",
					"labels":            "label1,label2",
					"registry-type":     "harbor",
					"registry-username": "testUsername",
					"registry-password": "testPassword",
				},
			},
			wantErr: false,
		},
		{
			name: "add unknown-resource",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "unknown-resource",
				attributes:   map[string]string{},
			},
			wantErr: true,
		},
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
		{
			name: "update git-project",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "git-project",
				id:           "test-id",
				attributes: map[string]string{
					"git-project-url": "https://github.com/example/example-project.git",
					"labels":          "label1,label2",
					"access-token":    "updatedAccessToken",
					"user-id":         "updatedUserId",
				},
			},
			wantErr: false,
		},
		{
			name: "update cloud-provider",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "cloud-provider",
				id:           "test-id",
				attributes: map[string]string{
					"cloud-type":  "aws",
					"labels":      "label1,label2",
					"cloud-key":   "updatedCloudKey",
					"cloud-token": "updatedCloudToken",
				},
			},
			wantErr: false,
		},
		{
			name: "update container-registry",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "container-registry",
				id:           "test-id",
				attributes: map[string]string{
					"registry-url":      "https://example.com",
					"labels":            "label1,label2",
					"registry-type":     "harbor",
					"registry-username": "updatedUsername",
					"registry-password": "updatedPassword",
				},
			},
			wantErr: false,
		},
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
		{
			name: "valid AWS cloud attributes",
			args: args{
				attributes: map[string]string{
					"access-key": "test-access-key",
					"secret-key": "test-secret-key",
				},
			},
			want: map[string]string{
				"accessKey": "test-access-key",
				"secretKey": "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "valid Azure cloud attributes",
			args: args{
				attributes: map[string]string{
					"client-id":     "test-client-id",
					"client-secret": "test-client-secret",
				},
			},
			want: map[string]string{
				"clientID":     "test-client-id",
				"clientSecret": "test-client-secret",
			},
			wantErr: false,
		},
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
		{
			name: "valid AWS resource",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "aws",
				id:           "test-id",
			},
			wantErr: false,
		},
		{
			name: "valid Azure resource",
			args: args{
				captenConfig: config.CaptenConfig{},
				resourceType: "azure",
				id:           "test-id",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveClusterResource(tt.args.captenConfig, tt.args.resourceType, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveClusterResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveClusterResource(tt.args.captenConfig, tt.args.resourceType, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveClusterResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
