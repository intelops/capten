package agent

import (
	"capten/pkg/config"
	"testing"
)

func TestListClusterPlugins(t *testing.T) {
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
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
			},
			wantErr: false,
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListClusterPlugins(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("ListClusterPlugins() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeployPlugin(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
		pluginName   string
		version      string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config and valid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				storeType: "helm",
			},
			wantErr: false,
		},
		{
			name: "Valid config and invalid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				storeType: "invalid-store",
			},
			wantErr: true,
		},
		{
			name: "Invalid config and valid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
				storeType: "helm",
			},
			wantErr: true,
		},
		{
			name: "Invalid config and invalid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
				storeType: "invalid-store",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeployPlugin(tt.args.captenConfig, tt.args.storeType, tt.args.pluginName, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("DeployPlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnDeployPlugin(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
		pluginName   string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config and valid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				storeType:  "helm",
				pluginName: "some-plugin",
			},
			wantErr: false,
		},
		{
			name: "Valid config and invalid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				storeType:  "invalid-store",
				pluginName: "some-plugin",
			},
			wantErr: true,
		},
		{
			name: "Invalid config and valid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
				storeType:  "helm",
				pluginName: "some-plugin",
			},
			wantErr: true,
		},
		{
			name: "Invalid config and invalid store type",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
				storeType:  "invalid-store",
				pluginName: "some-plugin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnDeployPlugin(tt.args.captenConfig, tt.args.storeType, tt.args.pluginName); (err != nil) != tt.wantErr {
				t.Errorf("UnDeployPlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShowClusterPluginData(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		pluginName   string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config and valid plugin name",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				pluginName: "some-plugin",
			},
			wantErr: false,
		},
		{
			name: "Valid config and invalid plugin name",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
				pluginName: "invalid-plugin",
			},
			wantErr: true,
		},
		{
			name: "Invalid config and valid plugin name",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
				pluginName: "some-plugin",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowClusterPluginData(tt.args.captenConfig, tt.args.pluginName); (err != nil) != tt.wantErr {
				t.Errorf("ShowClusterPluginData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
