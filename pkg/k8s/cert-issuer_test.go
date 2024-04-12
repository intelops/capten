package k8s

import (
	"capten/pkg/config"
	"testing"
)

func TestCreateOrUpdateClusterIssuer(t *testing.T) {
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
			if err := CreateOrUpdateClusterIssuer(tt.args.captenConfig); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateClusterIssuer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
