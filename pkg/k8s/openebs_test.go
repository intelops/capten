package k8s

import (
	"capten/pkg/config"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateCStorPoolClusterWithRetries(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateCStorPoolClusterWithRetries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
