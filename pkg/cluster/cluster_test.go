package cluster

import (
	"capten/pkg/config"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{
					// Add valid captenConfig properties here
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid captenConfig",
			args: args{
				captenConfig: config.CaptenConfig{
					// Add invalid captenConfig properties here
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Create(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestDestroy(t *testing.T) {
	type args struct {
		captenConfig config.CaptenConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful destruction of a cluster",
			args: args{
				captenConfig: config.CaptenConfig{},
			},
			wantErr: false,
		},
		{
			name: "Error handling when trying to destroy a non-existing cluster",
			args: args{
				captenConfig: config.CaptenConfig{
					// Add invalid captenConfig properties here
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Destroy(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
