package agent

import (
	"capten/pkg/config"
	"testing"
)

func TestListClusterApplications(t *testing.T) {
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
			if err := ListClusterApplications(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("ListClusterApplications() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShowClusterAppData(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
		appName      string
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
			if err := ShowClusterAppData(tt.args.captenConfig, tt.args.appName); (err != nil) != tt.wantErr {
				t.Errorf("ShowClusterAppData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeployDefaultApps(t *testing.T) {
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
			if err := DeployDefaultApps(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("DeployDefaultApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitAndTrackDefaultAppsDeploymentStatus(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WaitAndTrackDefaultAppsDeploymentStatus(tt.args.captenConfig)
		})
	}
}

func TestGetDefaultAppsDeploymentStatus(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetDefaultAppsDeploymentStatus(tt.args.captenConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultAppsDeploymentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultAppsDeploymentStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetDefaultAppsDeploymentStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
