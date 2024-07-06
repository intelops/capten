package agent

import (
	"capten/pkg/config"
	"log"

	//"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureClusterPlugin(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}
	captenConfig := config.CaptenConfig{
		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsagent.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	tests := []struct {
		name             string
		pluginName       string
		action           string
		actionAttributes map[string]string
		expectedError    bool
	}{
		{
			name:             "Valid tekton Plugin",
			pluginName:       "tekton",
			action:           "show-tekton-project",
			actionAttributes: map[string]string{},
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConfigureClusterPlugin(captenConfig, tt.pluginName, tt.action, tt.actionAttributes)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
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
	}

	for _, test := range tests {
		err := configureCrossplanePlugin(captenConfig, test.action, test.actionAttributes)
		if err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("For action %v: Expected error: %v, got: %v", test.action, test.expectedError, err)
		}

	}
}

func TestConfigureTektonPlugin(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}
	captenConfig := config.CaptenConfig{

		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsagent.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	tests := []struct {
		action        string
		expectedError error
	}{
		{"show-tekton-project", nil},
	}

	for _, test := range tests {
		err := configureTektonPlugin(captenConfig, test.action)

		if err == nil && test.expectedError != nil {
			t.Errorf("For action %v: Expected error: %v, got: nil", test.action, test.expectedError)
		}
	}
}

func TestShowCrossplaneProject_Success(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	captenConfig := config.CaptenConfig{
		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsagent.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	err = showCrossplaneProject(captenConfig)

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShowCrossplaneProject_Error(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	captenConfig := config.CaptenConfig{
		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsagent.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	err = showCrossplaneProject(captenConfig)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestSynchCrossplaneProject_Success(t *testing.T) {
	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	captenConfig := config.CaptenConfig{
		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsagent.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	err = synchCrossplaneProject(captenConfig)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSynchCrossplaneProject_Error(t *testing.T) {

	currentdir, err := os.Getwd()
	if err != nil {
		log.Println("Error while getting cuerent dir", err)
	}
	presentdir, err := getRelativePathUpTo(currentdir)

	if err != nil {
		log.Println("Error while getting working dir", err)
	}

	captenConfig := config.CaptenConfig{
		CertDirPath:        "/" + presentdir + "/cert/",
		ConfigDirPath:      "/" + presentdir + "/config/",
		AgentHostName:      "captenagent",
		KubeConfigFileName: "kubeconfig",
		ClientKeyFileName:  "client.key",
		ClientCertFileName: "client.crt",
		CAFileName:         "ca.crt",
		CaptenClusterValues: config.CaptenClusterValues{
			DomainName: "awsdemo.optimizor.app",
		},
		CaptenClusterHost: config.CaptenClusterHost{
			LoadBalancerHost: "a084c23852d0b428e98f363457fc8f8b-5ee99283c8b044fa.elb.us-west-2.amazonaws.com",
		},
	}

	err = synchCrossplaneProject(captenConfig)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
