package agent

import (
	"capten/pkg/agent/pb/pluginstorepb"
	"capten/pkg/config"
	"reflect"
	"testing"
)

func Test_getStoreTypeEnum(t *testing.T) {
	type args struct {
		storeType string
	}

	tests := []struct {
		name    string
		args    args
		want    pluginstorepb.StoreType
		wantErr bool
	}{
		{
			name:    "central-store",
			args:    args{storeType: "central-store"},
			want:    pluginstorepb.StoreType_CENTRAL_STORE,
			wantErr: false,
		},
		{
			name:    "default",
			args:    args{storeType: "default"},
			want:    pluginstorepb.StoreType_DEFAULT_STORE,
			wantErr: false,
		},
		{
			name:    "local-store",
			args:    args{storeType: "local-store"},
			want:    pluginstorepb.StoreType_LOCAL_STORE,
			wantErr: true,
		},
	}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	want    pluginstorepb.StoreType
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStoreTypeEnum(tt.args.storeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStoreTypeEnum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStoreTypeEnum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListPluginStoreApps(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Central store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store"},
			wantErr: false,
		},
		{
			name:    "Default store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "default"},
			wantErr: false,
		},
		{
			name:    "Invalid store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "invalid-store"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListPluginStoreApps(tt.args.captenConfig, tt.args.storeType); (err != nil) != tt.wantErr {
				t.Errorf("ListPluginStoreApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	//}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListPluginStoreApps(tt.args.captenConfig, tt.args.storeType); (err != nil) != tt.wantErr {
				t.Errorf("ListPluginStoreApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigPluginStore(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
		gitProjectId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Central store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store", gitProjectId: "gpid"},
			wantErr: false,
		},
		{
			name:    "Default store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "default", gitProjectId: "gpid"},
			wantErr: false,
		},
		{
			name:    "Invalid store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "invalid-store", gitProjectId: "gpid"},
			wantErr: true,
		},
		{
			name:    "Empty gitProjectId",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store", gitProjectId: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigPluginStore(tt.args.captenConfig, tt.args.storeType, tt.args.gitProjectId); (err != nil) != tt.wantErr {
				t.Errorf("ConfigPluginStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	//}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigPluginStore(tt.args.captenConfig, tt.args.storeType, tt.args.gitProjectId); (err != nil) != tt.wantErr {
				t.Errorf("ConfigPluginStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSynchPluginStore(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Central store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store"},
			wantErr: false,
		},
		{
			name:    "Default store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "default"},
			wantErr: false,
		},
		{
			name:    "Invalid store",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "invalid-store"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SynchPluginStore(tt.args.captenConfig, tt.args.storeType); (err != nil) != tt.wantErr {
				t.Errorf("SynchPluginStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SynchPluginStore(tt.args.captenConfig, tt.args.storeType); (err != nil) != tt.wantErr {
				t.Errorf("SynchPluginStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShowPluginStorePlugin(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		storeType    string
		pluginName   string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Central store valid plugin",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store", pluginName: "example-plugin"},
			wantErr: false,
		},
		{
			name:    "Central store invalid plugin",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "central-store", pluginName: "invalid-plugin"},
			wantErr: true,
		},
		{
			name:    "Default store valid plugin",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "default", pluginName: "example-plugin"},
			wantErr: false,
		},
		{
			name:    "Default store invalid plugin",
			args:    args{captenConfig: config.CaptenConfig{}, storeType: "default", pluginName: "invalid-plugin"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowPluginStorePlugin(tt.args.captenConfig, tt.args.storeType, tt.args.pluginName); (err != nil) != tt.wantErr {
				t.Errorf("ShowPluginStorePlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShowPluginStorePlugin(tt.args.captenConfig, tt.args.storeType, tt.args.pluginName); (err != nil) != tt.wantErr {
				t.Errorf("ShowPluginStorePlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
