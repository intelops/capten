package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func Test_readAndValidResourceIdentfierFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name             string
		args             args
		wantResourceType string
		wantId           string
		wantErr          bool
	}{
		{
			name: "Valid Resource",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("resource-type", "test-type", "resource type")
					cmd.Flags().String("id", "test-id", "resource id")
					return cmd
				}(),
			},
			wantResourceType: "test-type",
			wantId:           "test-id",
			wantErr:          false,
		},
		{
			name: "Invalid Resource",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("resource-type", "", "resource type")
					cmd.Flags().String("id", "", "resource id")
					return cmd
				}(),
			},
			wantResourceType: "",
			wantId:           "",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResourceType, gotId, err := readAndValidResourceIdentfierFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidResourceIdentfierFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResourceType != tt.wantResourceType {
				t.Errorf("readAndValidResourceIdentfierFlags() gotResourceType = %v, want %v", gotResourceType, tt.wantResourceType)
			}
			if gotId != tt.wantId {
				t.Errorf("readAndValidResourceIdentfierFlags() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func Test_readCloudTypeAttributesFlags(t *testing.T) {
	type args struct {
		cmd       *cobra.Command
		cloudType string
	}

	tests := []struct {
		name           string
		args           args
		wantAttributes map[string]string
		wantErr        bool
	}{
		{
			name: "valid AWS flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("access-key", "test-access-key", "access key")
					cmd.Flags().String("secret-key", "test-secret-key", "secret key")
					return cmd
				}(),
				cloudType: "aws",
			},
			wantAttributes: map[string]string{
				"access-key": "test-access-key",
				"secret-key": "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "valid Azure flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("client-id", "test-client-id", "client id")
					cmd.Flags().String("client-secret", "test-client-secret", "client secret")
					return cmd
				}(),
				cloudType: "azure",
			},
			wantAttributes: map[string]string{
				"client-id":     "test-client-id",
				"client-secret": "test-client-secret",
			},
			wantErr: false,
		},
		{
			name: "invalid cloud type",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					return cmd
				}(),
				cloudType: "invalid",
			},
			wantAttributes: nil,
			wantErr:        true,
		},
		{
			name: "empty AWS flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					return cmd
				}(),
				cloudType: "aws",
			},
			wantAttributes: nil,
			wantErr:        true,
		},
		{
			name: "empty Azure flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					return cmd
				}(),
				cloudType: "azure",
			},
			wantAttributes: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAttributes, err := readCloudTypeAttributesFlags(tt.args.cmd, tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCloudTypeAttributesFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readCloudTypeAttributesFlags() = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAttributes, err := readCloudTypeAttributesFlags(tt.args.cmd, tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCloudTypeAttributesFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readCloudTypeAttributesFlags() = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}

func Test_readAndValidResourceDataFlags(t *testing.T) {
	type args struct {
		cmd          *cobra.Command
		resourceType string
	}

	tests := []struct {
		name           string
		args           args
		wantAttributes map[string]string
		wantErr        bool
	}{
		{
			name: "ResourceType is empty",
			args: args{
				cmd: &cobra.Command{},
			},
			wantAttributes: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAttributes, err := readAndValidResourceDataFlags(tt.args.cmd, tt.args.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidResourceDataFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readAndValidResourceDataFlags() = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAttributes, err := readAndValidResourceDataFlags(tt.args.cmd, tt.args.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidResourceDataFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readAndValidResourceDataFlags() = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}

func Test_readAndValidCreateResourceFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name             string
		args             args
		wantResourceType string
		wantAttributes   map[string]string
		wantErr          bool
	}{
		{
			name: "missing resource type",
			args: args{
				cmd: &cobra.Command{},
			},
			wantResourceType: "",
			wantAttributes:   nil,
			wantErr:          true,
		},
		{
			name: "no attributes",
			args: args{
				cmd: &cobra.Command{},
			},
			wantResourceType: "test-type",
			wantAttributes:   map[string]string{},
			wantErr:          false,
		},
		{
			name: "with attributes",
			args: args{
				cmd: &cobra.Command{},
			},
			wantResourceType: "test-type",
			wantAttributes: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResourceType, gotAttributes, err := readAndValidCreateResourceFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidCreateResourceFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResourceType != tt.wantResourceType {
				t.Errorf("readAndValidCreateResourceFlags() gotResourceType = %v, want %v", gotResourceType, tt.wantResourceType)
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readAndValidCreateResourceFlags() gotAttributes = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}

func Test_readAndValidUpdateResourceFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name             string
		args             args
		wantResourceType string
		wantId           string
		wantAttributes   map[string]string
		wantErr          bool
	}{
		{
			name: "valid flags",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("resource-type", "test-type", "resource type")
					cmd.Flags().String("id", "test-id", "resource id")
					cmd.Flags().String("key1", "value1", "attribute key1")
					cmd.Flags().String("key2", "value2", "attribute key2")
					return cmd
				}(),
			},
			wantResourceType: "test-type",
			wantId:           "test-id",
			wantAttributes: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "missing resource type",
			args: args{
				cmd: &cobra.Command{},
			},
			wantResourceType: "",
			wantId:           "",
			wantAttributes:   nil,
			wantErr:          true,
		},
		{
			name: "missing id",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("resource-type", "test-type", "resource type")
					return cmd
				}(),
			},
			wantResourceType: "test-type",
			wantId:           "",
			wantAttributes:   nil,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResourceType, gotId, gotAttributes, err := readAndValidUpdateResourceFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndValidUpdateResourceFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResourceType != tt.wantResourceType {
				t.Errorf("readAndValidUpdateResourceFlags() gotResourceType = %v, want %v", gotResourceType, tt.wantResourceType)
			}
			if gotId != tt.wantId {
				t.Errorf("readAndValidUpdateResourceFlags() gotId = %v, want %v", gotId, tt.wantId)
			}
			if !reflect.DeepEqual(gotAttributes, tt.wantAttributes) {
				t.Errorf("readAndValidUpdateResourceFlags() gotAttributes = %v, want %v", gotAttributes, tt.wantAttributes)
			}
		})
	}
}
