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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createOrUpdateClusterCAIssuerSecret(tt.args.captenConfig, tt.args.k8sClient); (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateClusterCAIssuerSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
