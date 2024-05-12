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
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateNamespaceLabel(tt.args.clientset, tt.args.name, tt.args.labels); (err != nil) != tt.wantErr {
				t.Errorf("updateNamespaceLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
