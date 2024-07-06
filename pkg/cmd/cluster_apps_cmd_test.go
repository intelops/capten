package cmd

import (
	"capten/pkg/config"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func Test_readAppsNameFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}

	tests := []struct {
		name         string
		args         args
		wantAppsName string
		wantErr      bool
	}{
		{
			name: "Valid Apps Name Flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("apps-name", "test-app", "apps name")
					return cmd
				}(),
			},
			wantAppsName: "test-app",
			wantErr:      false,
		},
		{
			name: "Invalid Apps Name Flag",
			args: args{
				cmd: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().String("apps-name", "", "apps name")
					return cmd
				}(),
			},
			wantAppsName: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAppsName, err := readAppsNameFlags(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAppsNameFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAppsName != tt.wantAppsName {
				t.Errorf("readAppsNameFlags() = %v, want %v", gotAppsName, tt.wantAppsName)
			}
		})
	}
}

func Test_loadSetupAppsActions(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		want    *SetupAppsActionList
		wantErr bool
	}{
		{
			name: "Valid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
				},
			},
			want:    &SetupAppsActionList{},
			wantErr: false,
		},
		{
			name: "Invalid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadSetupAppsActions(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadSetupAppsActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadSetupAppsActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEnabled(t *testing.T) {
	type args struct {
		actionConfig map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "True",
			args: args{
				actionConfig: map[string]interface{}{
					"enabled": true,
				},
			},
			want: true,
		},
		{
			name: "False",
			args: args{
				actionConfig: map[string]interface{}{
					"enabled": false,
				},
			},
			want: false,
		},
		{
			name: "Nil",
			args: args{
				actionConfig: nil,
			},
			want: false,
		},
		{
			name: "Empty",
			args: args{
				actionConfig: map[string]interface{}{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEnabled(tt.args.actionConfig); got != tt.want {
				t.Errorf("isEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_execActionIfEnabled(t *testing.T) {
	type args struct {
		actionConfig map[string]interface{}
		f            func() error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Enabled",
			args: args{
				actionConfig: map[string]interface{}{
					"enabled": true,
				},
				f: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Disabled",
			args: args{
				actionConfig: map[string]interface{}{
					"enabled": false,
				},
				f: func() error {
					return fmt.Errorf("should not be called")
				},
			},
			wantErr: false,
		},
		{
			name: "Nil",
			args: args{
				actionConfig: nil,
				f:            func() error { return fmt.Errorf("should not be called") },
			},
			wantErr: false,
		},
		{
			name: "Empty",
			args: args{
				actionConfig: map[string]interface{}{},
				f:            func() error { return fmt.Errorf("should not be called") },
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				actionConfig: map[string]interface{}{
					"enabled": true,
				},
				f: func() error {
					return fmt.Errorf("some error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := execActionIfEnabled(tt.args.actionConfig, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("execActionIfEnabled() error = %v, wantErr %v", err, tt.wantErr)
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
					return fmt.Errorf("some error")
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
