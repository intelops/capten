package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_readAndValidClusterFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name             string
		args             args
		wantCloudService string
		wantClusterType  string
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCloudService, gotClusterType, err := readAndValidClusterFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidClusterFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCloudService != tt.wantCloudService {
				t.Errorf("readAndValidClusterFlags() gotCloudService = %v, want %v", gotCloudService, tt.wantCloudService)
			}
			if gotClusterType != tt.wantClusterType {
				t.Errorf("readAndValidClusterFlags() gotClusterType = %v, want %v", gotClusterType, tt.wantClusterType)
			}
		})
	}
}

func Test_validateClusterFlags(t *testing.T) {
	type args struct {
		cloudService string
		clusterType  string
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
			if err := validateClusterFlags(tt.args.cloudService, tt.args.clusterType); (err != nil) != tt.wantErr {
				t.Errorf("validateClusterFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
