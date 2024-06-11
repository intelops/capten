package app

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClusterGlobalValues(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testValues := `key1: value1
key2:
  subkey1: subvalue1
  subkey2: subvalue2`
		valuesFilePath := writeTempFile(t, testValues)

		values, err := GetClusterGlobalValues(valuesFilePath)

		assert.NoError(t, err)

		assert.Equal(t, "value1", values["key1"])
		assert.Equal(t, map[interface{}]interface{}{"subkey1": "subvalue1", "subkey2": "subvalue2"}, values["key2"])
	})

	t.Run("FileNotFound", func(t *testing.T) {
		values, err := GetClusterGlobalValues("/non/existent/file.yaml")

		assert.Error(t, err)
		assert.Nil(t, values)
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		invalidValues := `key1: value1
key2:
  subkey1: subvalue1
  subkey2: subvalue2`
		valuesFilePath := writeTempFile(t, invalidValues)
		corruptValues := `key1: value1
key2:
  subkey1: subvalue1
  subkey2:: subvalue2`
		ioutil.WriteFile(valuesFilePath, []byte(corruptValues), 0644)

		// Invoke the function
		values, err := GetClusterGlobalValues(valuesFilePath)

		// Assert an error occurred
		assert.Error(t, err)
		assert.Nil(t, values)
	})
}

// Helper function to write temporary file with test data
func writeTempFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "test-values-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Sync(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func TestGetApps(t *testing.T) {
	type args struct {
		appListFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Valid app list file",
			args: args{
				appListFilePath: "testdata/valid-apps.yaml",
			},
			want:    []string{"app1", "app2"},
			wantErr: false,
		},
		{
			name: "Invalid app list file",
			args: args{
				appListFilePath: "testdata/invalid-apps.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Non-existent app list file",
			args: args{
				appListFilePath: "testdata/non-existent-file.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApps(tt.args.appListFilePath, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppConfig(t *testing.T) {
	type args struct {
		appConfigFilePath string
		globalValues      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    types.AppConfig
		wantErr bool
	}{
		{
			name: "Valid app config file",
			args: args{
				appConfigFilePath: "testdata/valid_app_config.yaml",
				globalValues:      map[string]interface{}{},
			},
			want: types.AppConfig{
				Name:        "test-app",
				Description: "Test app",
				//	Chart:       "test-chart",
				Version: "1.0.0",
				//	Values:      map[string]interface{}{"key1": "value1", "key2": map[string]interface{}{"subkey1": "subvalue1", "subkey2": "subvalue2"}},
			},
			wantErr: false,
		},
		{
			name: "Invalid app config file",
			args: args{
				appConfigFilePath: "testdata/invalid_app_config.yaml",
				globalValues:      map[string]interface{}{},
			},
			want:    types.AppConfig{},
			wantErr: true,
		},
		{
			name: "Non-existent app config file",
			args: args{
				appConfigFilePath: "testdata/nonexistent_app_config.yaml",
				globalValues:      map[string]interface{}{},
			},
			want:    types.AppConfig{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppConfig(tt.args.appConfigFilePath, tt.args.globalValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetAppValuesTemplate(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		appName      string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Valid template",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsValuesDirPath: "testdata/valid-values",
				},
				appName: "app1",
			},
			want: []byte("image: foo/bar:v1\n"),
		},
		{
			name: "Non-existent template",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsValuesDirPath: "testdata/non-existent-values",
				},
				appName: "app1",
			},
			want: nil,
		},
		{
			name: "Invalid template",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsValuesDirPath: "testdata/invalid-values",
				},
				appName: "app1",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAppValuesTemplate(tt.args.captenConfig, tt.args.appName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppValuesTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestWriteAppConfig(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		appConfig    types.AppConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid app config",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsTempDirPath: "/tmp",
				},
				appConfig: types.AppConfig{
					Name: "test-app",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid app config - no name",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsTempDirPath: "/tmp",
				},
				appConfig: types.AppConfig{
					Name: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid app config - no temp dir path",
			args: args{
				captenConfig: config.CaptenConfig{
					AppsTempDirPath: "",
				},
				appConfig: types.AppConfig{
					Name: "test-app",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteAppConfig(tt.args.captenConfig, tt.args.appConfig); (err != nil) != tt.wantErr {
				t.Errorf("WriteAppConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestPrepareGlobalVaules(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Valid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			want:    map[string]interface{}{},
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
		{
			name: "Empty captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareGlobalVaules(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareGlobalVaules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareGlobalVaules() = %v, want %v", got, tt.want)
			}
		})
	}
}
