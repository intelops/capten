package k8s

import (
	//"reflect"
	"testing"

	"k8s.io/client-go/kubernetes"
)

func TestMakeNamespacePrivilege(t *testing.T) {
	type args struct {
		kubeconfigPath string
		ns             string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test empty kubeconfig path",
			args:    args{"", "testns"},
			wantErr: true,
		},
		{
			name:    "Test empty namespace",
			args:    args{"../config/kubeconfig", ""},
			wantErr: true,
		},
		{
			name:    "Test valid case",
			args:    args{"../config/kubeconfig", "testns"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeNamespacePrivilege(tt.args.kubeconfigPath, tt.args.ns); (err != nil) != tt.wantErr {
				t.Errorf("MakeNamespacePrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeNamespacePrivilege(tt.args.kubeconfigPath, tt.args.ns); (err != nil) != tt.wantErr {
				t.Errorf("MakeNamespacePrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetK8SClient(t *testing.T) {
	type args struct {
		kubeconfigPath string
	}

	tests := []struct {
		name    string
		args    args
		want    *kubernetes.Clientset
		wantErr bool
	}{
		{
			name:    "Test empty kubeconfig path",
			args:    args{""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test valid kubeconfig path",
			args:    args{"../config/kubeconfig"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetK8SClient(tt.args.kubeconfigPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetK8SClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != (tt.want == nil) {
				t.Errorf("GetK8SClient() = %v, want %v", got, tt.want)
			}
		})
	}

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, err := GetK8SClient(tt.args.kubeconfigPath)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("GetK8SClient() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("GetK8SClient() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

func TestCreateNamespaceIfNotExists(t *testing.T) {
	type args struct {
		kubeconfigPath string
		namespace      string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test empty kubeconfig path",
			args:    args{"", "testns"},
			wantErr: true,
		},
		{
			name:    "Test empty namespace",
			args:    args{"../config/kubeconfig", ""},
			wantErr: true,
		},
		{
			name:    "Test valid case",
			args:    args{"../config/kubeconfig", "testns"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNamespaceIfNotExists(tt.args.kubeconfigPath, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("CreateNamespaceIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNamespaceIfNotExists(tt.args.kubeconfigPath, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("CreateNamespaceIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
