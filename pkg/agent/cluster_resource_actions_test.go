package agent

import (
	"capten/pkg/config"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestListClusterResources(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}
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
				captenConfig: config.CaptenConfig{
					AgentSecure:        true,
					AgentHostName:      "captenagent",
					CertDirPath:        "/" + presentdir + "/cert/",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
					VaultCredHostName:  "vault-cred",
					AgentHostPort:      ":443",
					CaptenClusterHost: config.CaptenClusterHost{
						LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
					},
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "awsdemo.optimizor.app",
					},
				},
				resourceType: "git-project",
			},
			wantErr: false,
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

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

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
				captenConfig: config.CaptenConfig{

					AgentSecure:        true,
					AgentHostName:      "captenagent",
					CertDirPath:        "/" + presentdir + "/cert/",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
					VaultCredHostName:  "vault-cred",
					AgentHostPort:      ":443",
					CaptenClusterHost: config.CaptenClusterHost{
						LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
					},
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "awsdemo.optimizor.app",
					},
				},
				resourceType: "git-project",
				attributes: map[string]string{
					"git-project-url": "https://github.com/example/example-project.git",
					//"labels":          "label1,label2",
					"access-token": "testAccessToken",
					"user-id":      "testUserId",
				},
			},
			wantErr: false,
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

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}
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
				captenConfig: config.CaptenConfig{
					AgentSecure:        true,
					AgentHostName:      "captenagent",
					CertDirPath:        "/" + presentdir + "/cert/",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
					VaultCredHostName:  "vault-cred",
					AgentHostPort:      ":443",
					CaptenClusterHost: config.CaptenClusterHost{
						LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
					},
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "awsdemo.optimizor.app",
					},
				},
				resourceType: "git-project",
				id:           "9e9cb5a1-49bc-46d6-b773-5423672438cc",
				attributes: map[string]string{
					"git-project-url": "https://github.com/example/example-project.git",
					"labels":          "IntelopsCi",
					"access-token":    "ksjdksjdk",
					"user-id":         "updatedUserId",
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
					"cloud-type": "aws",
					"accessKey":  "dkdndnfdnf",
					"secretKey":  "nSlkdnnns",
				},
			},
			want: map[string]string{
				//"cloud-type": "aws",
				"accessKey": "dkdndnfdnf",
				"secretKey": "nSlkdnnns",
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

}
