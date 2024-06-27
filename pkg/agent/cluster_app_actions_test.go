package agent

import (
	"capten/pkg/config"
	"log"
	"os"
	"testing"
)

func TestListClusterApplications(t *testing.T) {

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
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					CertDirPath: "/" + presentdir + "/cert/",

					AgentHostName:      "captenagent",
					KubeConfigFileName: "kubeconfig",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "awsagent.optimizor.app",
					},
					CaptenClusterHost: config.CaptenClusterHost{
						LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListClusterApplications(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("ListClusterApplications() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShowClusterAppData(t *testing.T) {
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
		appName      string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config and app name",
			args: args{
				captenConfig: config.CaptenConfig{
					CertDirPath:        "/" + presentdir + "/cert/",
					ConfigDirPath:      "/" + presentdir + "/config/",
					KubeConfigFileName: "kubeconfig",
					AgentHostName:      "captenagent",
					ClientKeyFileName:  "client.key",
					ClientCertFileName: "client.crt",
					CAFileName:         "ca.crt",
					CaptenClusterValues: config.CaptenClusterValues{
						DomainName: "awsagent.optimizor.app",
					},
					CaptenClusterHost: config.CaptenClusterHost{
						LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
					},
				},

				appName: "external-secrets",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowClusterAppData(tt.args.captenConfig, tt.args.appName); (err != nil) != tt.wantErr {
				t.Errorf("ShowClusterAppData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeployDefaultApps(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeployDefaultApps(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("DeployDefaultApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeployDefaultApps(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("DeployDefaultApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitAndTrackDefaultAppsDeploymentStatus(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
			},
		},
		{
			name: "Invalid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WaitAndTrackDefaultAppsDeploymentStatus(tt.args.captenConfig)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WaitAndTrackDefaultAppsDeploymentStatus(tt.args.captenConfig)
		})
	}
}

func TestGetDefaultAppsDeploymentStatus(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
			},
			want:    true,
			want1:   "status",
			wantErr: false,
		},
		{
			name: "Invalid config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			want:    false,
			want1:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetDefaultAppsDeploymentStatus(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultAppsDeploymentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultAppsDeploymentStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetDefaultAppsDeploymentStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
