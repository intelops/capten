package helm

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				captenConfig: config.CaptenConfig{
					KubeConfigFileName: "kubeconfig",
				},
			},
			wantErr: false,
		},
		{
			name: "empty config",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_Install(t *testing.T) {

	type args struct {
		ctx       context.Context
		appConfig *types.AppConfig
	}
	type fields struct {
		Settings       *cli.EnvSettings
		defaultTimeout time.Duration
		captenConfig   config.CaptenConfig
	}

	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantAlreadyInstalled bool
		wantErr              bool
	}{
		{
			name: "valid app config",
			args: args{
				ctx: context.TODO(),
				appConfig: &types.AppConfig{
					ChartName: "mysql",
					Name:      "mysql-test",
					Namespace: "default",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid app config",
			args: args{
				ctx: context.TODO(),
				appConfig: &types.AppConfig{
					ChartName: "unknown-chart",
					Name:      "mysql-test",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
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
		{
			name: "valid action config",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				ctx:          context.TODO(),
				settings:     &cli.EnvSettings{},
				actionConfig: &action.Configuration{},
				appConfig:    &types.AppConfig{},
			},
			wantErr: false,
		},
		{
			name: "invalid action config",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				ctx:          context.TODO(),
				settings:     &cli.EnvSettings{},
				actionConfig: &action.Configuration{},
				appConfig:    &types.AppConfig{},
			},
			wantErr: true,
		},
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
		{
			name: "valid action config",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				ctx:          context.TODO(),
				settings:     &cli.EnvSettings{},
				actionConfig: &action.Configuration{},
				appConfig:    &types.AppConfig{},
			},
			wantErr: false,
		},
		{
			name: "invalid action config",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				ctx:          context.TODO(),
				settings:     &cli.EnvSettings{},
				actionConfig: &action.Configuration{},
				appConfig:    &types.AppConfig{},
			},
			wantErr: true,
		},
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
		{
			name: "app installed",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				actionConfig: &action.Configuration{},
				releaseName:  "installed-app",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "app not installed",
			fields: fields{
				Settings:       &cli.EnvSettings{},
				defaultTimeout: time.Second * 10,
				captenConfig:   config.CaptenConfig{},
			},
			args: args{
				actionConfig: &action.Configuration{},
				releaseName:  "not-installed-app",
			},
			want:    false,
			wantErr: false,
		},
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
		{
			name: "Valid template with single value",
			args: args{
				data: []byte("Hello {{ .Name }}!"),
				values: map[string]interface{}{
					"Name": "Alice",
				},
			},
			wantTransformedData: []byte("Hello Alice!"),
			wantErr:             false,
		},
		{
			name: "Valid template with multiple values",
			args: args{
				data: []byte("{{ .Name }} is {{ .Age }} years old."),
				values: map[string]interface{}{
					"Name": "Bob",
					"Age":  25,
				},
			},
			wantTransformedData: []byte("Bob is 25 years old."),
			wantErr:             false,
		},
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

var mockWriteFile = func(filename string, data []byte, perm os.FileMode) error {
	return nil
}

func TestPrepareAppValues_Success(t *testing.T) {
	client := &Client{ /* initialize fields if necessary */ }
	appConfig := &types.AppConfig{
		Name:           "test-app",
		TemplateValues: []byte("valid-template"),
		OverrideValues: map[string]interface{}{"key": "value"},
	}

	mockWriteFile = func(filename string, data []byte, perm os.FileMode) error {
		return nil
	}

	tmpValuesPath, err := client.prepareAppValues(appConfig)

	assert.NoError(t, err)
	assert.Contains(t, tmpValuesPath, "test-app-values.yaml")
}

func TestPrepareAppValues_WriteFileError(t *testing.T) {
	client := &Client{}
	appConfig := &types.AppConfig{
		Name:           "test-app",
		TemplateValues: []byte("valid-template"),
		OverrideValues: map[string]interface{}{"key": "value"},
	}

	tmpValuesPath, err := client.prepareAppValues(appConfig)

	assert.Error(t, err)
	assert.Equal(t, "", tmpValuesPath)
	assert.Contains(t, err.Error(), "failed to write app values to file")
}
