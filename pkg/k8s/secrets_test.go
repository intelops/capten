package k8s

import (
	"capten/pkg/config"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func Test_createOrUpdateSecret(t *testing.T) {
	type args struct {
		k8sClient *kubernetes.Clientset
		secret    *corev1.Secret
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test empty client",
			args: args{
				k8sClient: nil,
				secret:    &corev1.Secret{},
			},
			wantErr: true,
		},
		{
			name: "Test empty secret",
			args: args{
				k8sClient: &kubernetes.Clientset{},
				secret:    nil,
			},
			wantErr: true,
		},
		{
			name: "Test valid input",
			args: args{
				k8sClient: &kubernetes.Clientset{},
				secret:    &corev1.Secret{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateSecret(tt.args.k8sClient, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOrUpdateCertSecrets(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test nil config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
		{
			name: "Test empty config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Test valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateOrUpdateCertSecrets(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateCertSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateAgentCertSecret(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		k8sClient    *kubernetes.Clientset
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test nil config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
		{
			name: "Test empty config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Test valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
				k8sClient: &kubernetes.Clientset{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateAgentCertSecret(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateAgentCertSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateAgentCertSecret(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateAgentCertSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateAgentCACert(t *testing.T) {

	type args struct {
		captenConfig config.CaptenConfig
		k8sClient    *kubernetes.Clientset
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test nil config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
		{
			name: "Test empty config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Test valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
				k8sClient: &kubernetes.Clientset{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateAgentCACert(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateAgentCACert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateClusterCAIssuerSecret(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		k8sClient    *kubernetes.Clientset
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test nil config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
		{
			name: "Test empty config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "",
					PoolClusterNamespace: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Test valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
				k8sClient: &kubernetes.Clientset{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateClusterCAIssuerSecret(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateClusterCAIssuerSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateClusterCAIssuerSecret(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateClusterCAIssuerSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
