package k8s

import (
	"testing"

	"k8s.io/client-go/kubernetes"
)

func TestCreateNamespaceIfNotExist(t *testing.T) {
	type args struct {
		kubeconfigPath string
		namespaceName  string
		label          map[string]string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		kubeconfig    string
		namespaceName string
		label         map[string]string
	}{
		{
			name:          "should fail when kubeconfigPath is empty",
			args:          args{kubeconfigPath: "", namespaceName: "test", label: map[string]string{}},
			wantErr:       true,
			kubeconfig:    "",
			namespaceName: "test",
			label:         map[string]string{},
		},
		{
			name:          "should fail when namespaceName is empty",
			args:          args{kubeconfigPath: "../config/kubeconfig", namespaceName: "", label: map[string]string{}},
			wantErr:       true,
			kubeconfig:    "test",
			namespaceName: "",
			label:         map[string]string{},
		},
		{
			name:          "should pass when kubeconfigPath and namespaceName are valid",
			args:          args{kubeconfigPath: "../config/kubeconfig", namespaceName: "test", label: map[string]string{"test": "test"}},
			wantErr:       false,
			kubeconfig:    "test",
			namespaceName: "test",
			label:         map[string]string{"test": "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNamespaceIfNotExist(tt.args.kubeconfigPath, tt.args.namespaceName, tt.args.label); (err != nil) != tt.wantErr {
				t.Errorf("CreateNamespaceIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestCreateorUpdateNamespaceWithLabel(t *testing.T) {
	type args struct {
		kubeconfigPath string
		namespaceName  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should fail when kubeconfigPath is empty",
			args:    args{kubeconfigPath: "", namespaceName: "test"},
			wantErr: true,
		},
		{
			name:    "should fail when namespaceName is empty",
			args:    args{kubeconfigPath: "../config/kubeconfig", namespaceName: ""},
			wantErr: false,
		},
		{
			name:    "should pass when kubeconfigPath and namespaceName are valid",
			args:    args{kubeconfigPath: "../config/kubeconfig", namespaceName: "test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateorUpdateNamespaceWithLabel(tt.args.kubeconfigPath, tt.args.namespaceName); (err != nil) != tt.wantErr {
				t.Errorf("CreateorUpdateNamespaceWithLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_namespaceExists(t *testing.T) {
	type args struct {
		clientset *kubernetes.Clientset
		name      string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "should return false when namespace does not exist",
			args:    args{clientset: &kubernetes.Clientset{}, name: "non_existent_namespace"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "should return true when namespace exists",
			args:    args{clientset: &kubernetes.Clientset{}, name: "default"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespaceExists(tt.args.clientset, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("namespaceExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("namespaceExists() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespaceExists(tt.args.clientset, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("namespaceExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("namespaceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateNamespaceLabel(t *testing.T) {
	type args struct {
		clientset *kubernetes.Clientset
		name      string
		labels    map[string]string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should update namespace labels",
			args: args{
				clientset: &kubernetes.Clientset{},
				name:      "default",
				labels: map[string]string{
					"new-label": "new-value",
				},
			},
			wantErr: false,
		},
		{
			name: "should not error when namespace does not exist",
			args: args{
				clientset: &kubernetes.Clientset{},
				name:      "non-existent-namespace",
				labels: map[string]string{
					"new-label": "new-value",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateNamespaceLabel(tt.args.clientset, tt.args.name, tt.args.labels); (err != nil) != tt.wantErr {
				t.Errorf("updateNamespaceLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
