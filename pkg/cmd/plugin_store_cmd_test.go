package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_readAndValidatePluginStoreTypeFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name          string
		args          args
		wantStoreType string
		wantErr       bool
	}{
		{
			name: "valid store type flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "test-store-type", "store type")
					return cmd
				}(),
			},
			wantStoreType: "test-store-type",
			wantErr:       false,
		},
		{
			name: "invalid store type flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "", "store type")
					return cmd
				}(),
			},
			wantStoreType: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStoreType, err := readAndValidatePluginStoreTypeFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidatePluginStoreTypeFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStoreType != tt.wantStoreType {
				t.Errorf("readAndValidatePluginStoreTypeFlags() = %v, want %v", gotStoreType, tt.wantStoreType)
			}
		})
	}
}

func Test_readAndValidatePluginStoreShowFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name           string
		args           args
		wantPluginName string
		wantStoreType  string
		wantErr        bool
	}{
		{
			name: "valid flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("plugin", "test-plugin", "plugin name")
					cmd.Flags().String("store-type", "test-store-type", "store type")
					return cmd
				}(),
			},
			wantPluginName: "test-plugin",
			wantStoreType:  "test-store-type",
			wantErr:        false,
		},
		{
			name: "missing store type flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("plugin", "test-plugin", "plugin name")
					return cmd
				}(),
			},
			wantPluginName: "",
			wantStoreType:  "",
			wantErr:        true,
		},
		{
			name: "missing plugin flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "test-store-type", "store type")
					return cmd
				}(),
			},
			wantPluginName: "",
			wantStoreType:  "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPluginName, gotStoreType, err := readAndValidatePluginStoreShowFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidatePluginStoreShowFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPluginName != tt.wantPluginName {
				t.Errorf("readAndValidatePluginStoreShowFlags() gotPluginName = %v, want %v", gotPluginName, tt.wantPluginName)
			}
			if gotStoreType != tt.wantStoreType {
				t.Errorf("readAndValidatePluginStoreShowFlags() gotStoreType = %v, want %v", gotStoreType, tt.wantStoreType)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPluginName, gotStoreType, err := readAndValidatePluginStoreShowFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidatePluginStoreShowFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPluginName != tt.wantPluginName {
				t.Errorf("readAndValidatePluginStoreShowFlags() gotPluginName = %v, want %v", gotPluginName, tt.wantPluginName)
			}
			if gotStoreType != tt.wantStoreType {
				t.Errorf("readAndValidatePluginStoreShowFlags() gotStoreType = %v, want %v", gotStoreType, tt.wantStoreType)
			}
		})
	}
}

func Test_readAndValidatePluginStoreConfigFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name             string
		args             args
		wantStoreType    string
		wantGitProjectId string
		wantErr          bool
	}{
		{
			name: "valid flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "test-store-type", "store type")
					cmd.Flags().String("git-project-id", "test-git-project-id", "git project id")
					return cmd
				}(),
			},
			wantStoreType:    "test-store-type",
			wantGitProjectId: "test-git-project-id",
			wantErr:          false,
		},
		{
			name: "invalid store type flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "", "store type")
					cmd.Flags().String("git-project-id", "test-git-project-id", "git project id")
					return cmd
				}(),
			},
			wantStoreType:    "",
			wantGitProjectId: "",
			wantErr:          true,
		},
		{
			name: "invalid git project id flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("store-type", "test-store-type", "store type")
					cmd.Flags().String("git-project-id", "", "git project id")
					return cmd
				}(),
			},
			wantStoreType:    "",
			wantGitProjectId: "",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStoreType, gotGitProjectId, err := readAndValidatePluginStoreConfigFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidatePluginStoreConfigFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStoreType != tt.wantStoreType {
				t.Errorf("readAndValidatePluginStoreConfigFlags() gotStoreType = %v, want %v", gotStoreType, tt.wantStoreType)
			}
			if gotGitProjectId != tt.wantGitProjectId {
				t.Errorf("readAndValidatePluginStoreConfigFlags() gotGitProjectId = %v, want %v", gotGitProjectId, tt.wantGitProjectId)
			}
		})
	}
}
