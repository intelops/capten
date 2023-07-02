package config

import (
	"io/ioutil"
	"log"

	"capten/pkg/util"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	defaultCaptenConfigPath = "config.yaml"
	configEnvKey            = "CAPTEN_CONFIG"
)

type CaptenConfig struct {
	DomainName           string `envconfig:"DOMAIN_NAME" default:"dev.intelops.app"`
	CaptenNamespace      string `envconfig:"CAPTEN_NAMESPACE" default:"capten"`
	AgentCertSecretName  string `envconfig:"AGENT_CERT_SECRET_NAME" default:"capten-agent-cert"`
	ConfigPath           string `envconfig:"CONFIG_PATH" default:"./config/"`
	CertPath             string `envconfig:"CERT_PATH" default:"./cert/"`
	KubeConfigPath       string `envconfig:"KUBE_CONFIG_PATH" default:"./config/kubeconfig"`
	CaptenValuesFilePath string `envconfig:"CAPTEN_VALUES_FILE_PATH" default:"./config/capten.yaml"`
	// LegacyAppConfigFilePath    string   `envconfig:"CAPTEN_VALUES_FILE_PATH" default:"./config/legacy_app_config.yaml"`
	AgentCertFileName          string   `envconfig:"AGENT_CERT_FILE_NAME" default:"agent.crt"`
	AgentKeyFileName           string   `envconfig:"AGENT_KEY_FILE_NAME" default:"agent.key"`
	ClientCertFileName         string   `envconfig:"CLIENT_CERT_FILE_NAME" default:"client.crt"`
	ClientKeyFileName          string   `envconfig:"CLIENT_KEY_FILE_NAME" default:"client.key"`
	CAFileName                 string   `envconfig:"CA_FILE_NAME" default:"ca.crt"`
	ClientCertExportFileName   string   `envconfig:"CLIENT_CERT_EXPORT_FILE_NAME" default:"capten-client-auth-certs.zip"`
	OrgName                    string   `envconfig:"ORG_NAME" default:"Intelops"`
	RootCACommonName           string   `envconfig:"ROOT_CA_CN" default:"Capten Root CA"`
	IntermediateCACommonName   string   `envconfig:"INTERMEDIATE_CA_CN" default:"Capten Cluster CA"`
	AgentCertCommonName        string   `envconfig:"AGENT_CERT_CN" default:"Capten Agent"`
	AgentDNSNamePrefixes       []string `envconfig:"AGENT_DNS_NAME_PREFIX" default:"*,agent"`
	CaptenClientCertCommonName string   `envconfig:"CAPTEN_CLIENT_CA_CN" default:"Capten Client"`
	AppsFilePath               string   `envconfig:"APPS_FILE_PATH" default:"./config/apps.yaml"`
	AppValuesDir               string   `envconfig:"APP_VALUES_DIR" default:"./config/values/"`
	AgentDNSNames              []string
}

type CaptenClusterValues struct {
	DomainName          string `yaml:"domain"`
	CaptenNamespace     string `yaml:"captenNamespace"`
	AgentCertSecretName string `yaml:"agentCertSecretName"`
}

// GetClusterConfig config for cluster creation
func GetClusterConfig(configPath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(configPath)
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

func GetConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(util.GetEnv(configEnvKey, defaultCaptenConfigPath))
	if err := config.ReadInConfig(); err != nil {
		log.Println("viper cli config read error", err)
		return nil, err
	}

	return config, nil
}

func GetCaptenConfig() (CaptenConfig, error) {
	cfg := CaptenConfig{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	values, err := GetCaptenClusterValues(cfg.CaptenValuesFilePath)
	if err != nil {
		return cfg, err
	}

	if len(values.DomainName) != 0 {
		cfg.DomainName = values.DomainName
	}

	cfg.AgentDNSNames = []string{}
	for _, prefixName := range cfg.AgentDNSNamePrefixes {
		cfg.AgentDNSNames = append(cfg.AgentDNSNames, prefixName+"."+cfg.DomainName)
	}
	return cfg, err
}

func GetCaptenClusterValues(valuesFilePath string) (CaptenClusterValues, error) {
	var values CaptenClusterValues
	data, err := ioutil.ReadFile(valuesFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal values file, %s", valuesFilePath)
	}
	return values, nil
}
