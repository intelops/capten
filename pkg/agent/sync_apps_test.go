package agent

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"reflect"
	"testing"
)

func TestSyncInstalledAppConfigsOnAgent(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncInstalledAppConfigsOnAgent(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("SyncInstalledAppConfigsOnAgent() error = %v, wantErr %v", err, tt.wantErr)
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
		// TODO: Add test cases.
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
