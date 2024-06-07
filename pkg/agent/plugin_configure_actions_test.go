package agent

import (
	"capten/pkg/config"

	"fmt"
	"os"
	"testing"
)

func TestConfigureClusterPlugin(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	tests := []struct {
		pluginName       string
		action           string
		actionAttributes map[string]string
		expectedError    error
	}{
		{"crossplane", "list-actions", map[string]string{}, nil},
		{"tekton", "some-action", map[string]string{}, nil},
		{"proact", "any-action", map[string]string{}, fmt.Errorf("configure actions for plugin is not implemented yet")},
		{"unknown", "any-action", map[string]string{}, fmt.Errorf("no configure actions for plugin supported")},
	}

	for _, test := range tests {
		err := ConfigureClusterPlugin(captenConfig, test.pluginName, test.action, test.actionAttributes)
		if err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("For plugin %v and action %v: Expected error: %v, got: %v", test.pluginName, test.action, test.expectedError, err)
		}
		if err == nil && test.expectedError != nil {
			t.Errorf("For plugin %v and action %v: Expected error: %v, got: nil", test.pluginName, test.action, test.expectedError)
		}
	}
}

func TestConfigureCrossplanePlugin(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	tests := []struct {
		action           string
		actionAttributes map[string]string
		expectedError    error
	}{
		{"list-actions", map[string]string{}, nil},
		{"unknown-action", map[string]string{}, fmt.Errorf("action is not supported for plugin")},
	}

	for _, test := range tests {
		err := ConfigureCrossplanePlugin(captenConfig, test.action, test.actionAttributes)
		if err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("For action %v: Expected error: %v, got: %v", test.action, test.expectedError, err)
		}
		if err == nil && test.expectedError != nil {
			t.Errorf("For action %v: Expected error: %v, got: nil", test.action, test.expectedError)
		}
	}
}

func TestConfigureTektonPlugin(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	tests := []struct {
		action        string
		expectedError error
	}{
		{"some-action", nil},
		{"another-action", nil},
	}

	for _, test := range tests {
		err := configureTektonPlugin(captenConfig, test.action)
		if err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("For action %v: Expected error: %v, got: %v", test.action, test.expectedError, err)
		}
		if err == nil && test.expectedError != nil {
			t.Errorf("For action %v: Expected error: %v, got: nil", test.action, test.expectedError)
		}
	}
}

func TestShowCrossplaneProject_Success(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	err := showCrossplaneProject(captenConfig)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShowCrossplaneProject_Error(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	err := showCrossplaneProject(captenConfig)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestSynchCrossplaneProject_Success(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	err := synchCrossplaneProject(captenConfig)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSynchCrossplaneProject_Error(t *testing.T) {
	captenConfig := config.CaptenConfig{}

	err := synchCrossplaneProject(captenConfig)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
