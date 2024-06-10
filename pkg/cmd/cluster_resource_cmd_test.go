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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
