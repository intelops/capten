package k8s

import (
	"capten/pkg/config"
	"errors"
	"reflect"
	"testing"
	"time"

	clientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
)

func Test_getOpenEBSClient(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		want    *clientset.Clientset
		wantErr bool
	}{
		{
			name: "Empty config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty kubeconfig",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "",
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOpenEBSClient(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOpenEBSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOpenEBSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOpenEBSBlockDevices(t *testing.T) {
	type args struct {
		openebsClientset *clientset.Clientset
		captenConfig     config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				openebsClientset: nil,
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "some-pool-cluster",
					PoolClusterNamespace: "some-pool-cluster-ns",
				},
			},
			want: []map[string]string{
				{
					"blockDevice": "bd1",
					"nodeName":    "node1",
				},
				{
					"blockDevice": "bd2",
					"nodeName":    "node2",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOpenEBSBlockDevices(tt.args.openebsClientset, tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOpenEBSBlockDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOpenEBSBlockDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCStorPoolClusters(t *testing.T) {
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
			if err := CreateCStorPoolClusters(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateCStorPoolClusters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_retry(t *testing.T) {
	type args struct {
		retries  int
		interval time.Duration
		f        func() error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful after one retry",
			args: args{
				retries:  2,
				interval: 10 * time.Millisecond,
				f: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Unsuccessful after max retries",
			args: args{
				retries:  2,
				interval: 10 * time.Millisecond,
				f: func() error {
					return errors.New("some error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.retries, tt.args.interval, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCStorPoolClusterWithRetries(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
			},
			wantErr: false,
		},
		{
			name: "Unsuccessful",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName:   "kubeconfig",
					PoolClusterName:      "pool-cluster",
					PoolClusterNamespace: "pool-cluster-ns",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateCStorPoolClusterWithRetries(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateCStorPoolClusterWithRetries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
