package helm

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"context"
	"reflect"
	"testing"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

func TestNewClient(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Install(t *testing.T) {
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}
	type args struct {
		ctx       context.Context
		appConfig *types.AppConfig
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantAlreadyInstalled bool
		wantErr              bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Client{
				Settings:       tt.fields.Settings,
				defaultTimeout: tt.fields.defaultTimeout,
				captenConfig:   tt.fields.captenConfig,
			}
			gotAlreadyInstalled, err := h.Install(tt.args.ctx, tt.args.appConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAlreadyInstalled != tt.wantAlreadyInstalled {
				t.Errorf("Client.Install() = %v, want %v", gotAlreadyInstalled, tt.wantAlreadyInstalled)
			}
		})
	}
}

func TestClient_installApp(t *testing.T) {
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}
	type args struct {
		ctx          context.Context
		settings     *cli.EnvSettings
		actionConfig *action.Configuration
		appConfig    *types.AppConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Client{
				Settings:       tt.fields.Settings,
				defaultTimeout: tt.fields.defaultTimeout,
				captenConfig:   tt.fields.captenConfig,
			}
			if err := h.installApp(tt.args.ctx, tt.args.settings, tt.args.actionConfig, tt.args.appConfig); (err != nil) != tt.wantErr {
				t.Errorf("Client.installApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_upgradeApp(t *testing.T) {
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}
	type args struct {
		ctx          context.Context
		settings     *cli.EnvSettings
		actionConfig *action.Configuration
		appConfig    *types.AppConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Client{
				Settings:       tt.fields.Settings,
				defaultTimeout: tt.fields.defaultTimeout,
				captenConfig:   tt.fields.captenConfig,
			}
			if err := h.upgradeApp(tt.args.ctx, tt.args.settings, tt.args.actionConfig, tt.args.appConfig); (err != nil) != tt.wantErr {
				t.Errorf("Client.upgradeApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_IsAppInstalled(t *testing.T) {
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}
	type args struct {
		actionConfig *action.Configuration
		releaseName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Client{
				Settings:       tt.fields.Settings,
				defaultTimeout: tt.fields.defaultTimeout,
				captenConfig:   tt.fields.captenConfig,
			}
			got, err := h.IsAppInstalled(tt.args.actionConfig, tt.args.releaseName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.IsAppInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.IsAppInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executeAppConfigTemplate(t *testing.T) {
	type args struct {
		data   []byte
		values map[string]interface{}
	}
	tests := []struct {
		name                string
		args                args
		wantTransformedData []byte
		wantErr             bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTransformedData, err := executeAppConfigTemplate(tt.args.data, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeAppConfigTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTransformedData, tt.wantTransformedData) {
				t.Errorf("executeAppConfigTemplate() = %v, want %v", gotTransformedData, tt.wantTransformedData)
			}
		})
	}
}

func TestClient_prepareAppValues(t *testing.T) {
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}
	type args struct {
		appConfig *types.AppConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Client{
				Settings:       tt.fields.Settings,
				defaultTimeout: tt.fields.defaultTimeout,
				captenConfig:   tt.fields.captenConfig,
			}
			got, err := h.prepareAppValues(tt.args.appConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.prepareAppValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.prepareAppValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
