package agent

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
	//"gopkg.in/yaml.v2"
)

func TestSyncInstalledAppConfigsOnAgent(t *testing.T) {
	writeTestAppConfig := func(appConfig map[string]interface{}, dir string) string {
		appConfigFilePath := filepath.Join(dir, "test-app.yaml")
		appConfigFile, err := os.Create(appConfigFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer appConfigFile.Close()

		appConfigBytes, err := yaml.Marshal(appConfig)
		if err != nil {
			t.Fatal(err)
		}

		_, err = appConfigFile.Write(appConfigBytes)
		if err != nil {
			t.Fatal(err)
		}

		return dir
	}

	type args struct {
		config config.CaptenConfig
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty directory",
			args: args{
				config: config.CaptenConfig{
					AppsDirPath: t.TempDir(),
				},
			},
			wantErr: false,
		},
		{
			name: "Not existing directory",
			args: args{
				config: config.CaptenConfig{
					AppsDirPath: "/not/existing/dir",
				},
			},
			wantErr: true,
		},
		{
			name: "Valid app config",
			args: args{
				config: config.CaptenConfig{
					AppsDirPath: writeTestAppConfig(map[string]interface{}{
						"name":    "valid-app",
						"version": "1.2.3",
					}, t.TempDir()),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid app config - no version",
			args: args{
				config: config.CaptenConfig{
					AppsDirPath: writeTestAppConfig(map[string]interface{}{
						"name": "valid-app",
						// No version provided
					}, t.TempDir()),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncInstalledAppConfigsOnAgent(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("SyncInstalledAppConfigsOnAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteTestAppConfig(t *testing.T) {
	tests := []struct {
		name    string
		values  map[string]interface{}
		dir     string
		wantErr bool
	}{
		{
			name:    "Valid app config",
			values:  map[string]interface{}{"name": "test-app", "version": "1.0.0"},
			dir:     t.TempDir(),
			wantErr: false,
		},
		{
			name:    "Invalid app config - no name",
			values:  map[string]interface{}{"version": "1.0.0"},
			dir:     t.TempDir(),
			wantErr: true,
		},
		{
			name:    "Invalid app config - no version",
			values:  map[string]interface{}{"name": "invalid-app"},
			dir:     t.TempDir(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeTestAppConfig(t, tt.values, tt.dir)
			_, err := os.Stat(filepath.Join(dir, "test-app.yaml"))
			if (err != nil) != tt.wantErr {
				t.Errorf("writeTestAppConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func Test_readInstalledAppConfigs(t *testing.T) {
	type args struct {
		config config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantRet []types.AppConfig
		wantErr bool
	}{
		{
			name: "Empty directory",
			args: args{
				config: config.CaptenConfig{
					AppsTempDirPath: t.TempDir(),
				},
			},
			wantRet: []types.AppConfig{},
			wantErr: false,
		},
		{
			name: "Not existing directory",
			args: args{
				config: config.CaptenConfig{
					AppsTempDirPath: "/not/existing/dir",
				},
			},
			wantRet: []types.AppConfig{},
			wantErr: true,
		},
		{
			name: "Valid app config",
			args: args{
				config: config.CaptenConfig{
					AppsTempDirPath: writeTestAppConfig(t, map[string]interface{}{
						"name":    "valid-app",
						"version": "1.2.3",
					}, t.TempDir()),
				},
			},
			wantRet: []types.AppConfig{
				{
					Name:    "valid-app",
					Version: "1.2.3",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := readInstalledAppConfigs(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("readInstalledAppConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("readInstalledAppConfigs() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func writeTestAppConfig(t *testing.T, appConfig map[string]interface{}, dir string) string {
	appConfigFilePath := filepath.Join(dir, "test-app.yaml")
	appConfigFile, err := os.Create(appConfigFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer appConfigFile.Close()

	appConfigBytes, err := yaml.Marshal(appConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = appConfigFile.Write(appConfigBytes)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}
