package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_readAndValidClusterFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	cloudService := "aws"
	clusterType := "k3s"

	tests := []struct {
		name             string
		args             args
		wantCloudService string
		wantClusterType  string
		wantErr          bool
	}{
		{
			name: "Valid Cluster Flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("cloud-service", cloudService, "cloud service")
					cmd.Flags().String("cluster-type", clusterType, "cluster type")
					return cmd
				}(),
			},
			wantCloudService: cloudService,
			wantClusterType:  clusterType,
			wantErr:          false,
		},
		{
			name: "Invalid Cluster Flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("cloud-service", "", "cloud service")
					cmd.Flags().String("cluster-type", "", "cluster type")
					return cmd
				}(),
			},
			wantCloudService: "",
			wantClusterType:  "",
			wantErr:          true,
		},
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
		{
			name: "supported cloud service and cluster type",
			args: args{
				cloudService: "aws",
				clusterType:  "talos",
			},
			wantErr: false,
		},
		{
			name: "unsupported cloud service",
			args: args{
				cloudService: "not-supported",
				clusterType:  "talos",
			},
			wantErr: true,
		},
		{
			name: "unsupported cluster type",
			args: args{
				cloudService: "aws",
				clusterType:  "not-supported",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateClusterFlags(tt.args.cloudService, tt.args.clusterType); (err != nil) != tt.wantErr {
				t.Errorf("validateClusterFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
