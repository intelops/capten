package config

import (
	"capten/pkg/types"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type CaptenConfig struct {
	DomainName                 string   `envconfig:"DOMAIN_NAME" default:"dev.intelops.app"`
	CaptenNamespace            string   `envconfig:"CAPTEN_NAMESPACE" default:"capten"`
	AgentCertSecretName        string   `envconfig:"AGENT_CERT_SECRET_NAME" default:"capten-agent-cert"`
	AppsDirPath                string   `envconfig:"APPS_DIR_PATH" default:"/apps/"`
	AppsConfigDirPath          string   `envconfig:"APPS_CONFIG_DIR_PATH" default:"/apps/conf/"`
	AppsTempDirPath            string   `envconfig:"APPS_TEMP_DIR_PATH" default:"/apps/temp/"`
	ConfigDirPath              string   `envconfig:"CONFIG_DIR_PATH" default:"/config/"`
	CertDirPath                string   `envconfig:"CERT_DIR_PATH" default:"/cert/"`
	TerraformModulesDirPath    string   `envconfig:"TERRAFORM_MODULE_DIR_PATH" default:"/terraform_modules/"`
	TerraformTemplateDirPath   string   `envconfig:"TERRAFORM_TEMPLATE_DIR_PATH" default:"/templates/k3s/"`
	AppListFileName            string   `envconfig:"APP_LIST_FILE_NAME" default:"app_list.yaml"`
	CaptenGlobalValuesFileName string   `envconfig:"CAPTEN_VALUES_FILE_PATH" default:"capten.yaml"`
	KubeConfigFileName         string   `envconfig:"KUBE_CONFIG_PATH" default:"kubeconfig"`
	TerraformTemplateFileName  string   `envconfig:"TERRAFORM_TEMPLATE_FILE_NAME" default:"values.tfvars.tmpl"`
	TerraformVarFileName       string   `envconfig:"TERRAFORM_VAR_FILE_NAME" default:"values.tfvars"`
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
	AgentDNSNames              []string
	CurrentDirPath             string
}

type CaptenClusterValues struct {
	DomainName string `yaml:"DomainName"`
}

func GetCaptenConfig() (CaptenConfig, error) {
	cfg := CaptenConfig{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.CurrentDirPath, err = os.Getwd()
	if err != nil {
		return cfg, errors.WithMessage(err, "error getting current directory")
	}
	err = addCurrentDirToPath(cfg.CurrentDirPath)
	if err != nil {
		return cfg, errors.WithMessage(err, "error adding current directory to env")
	}

	values, err := GetCaptenClusterValues(cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenGlobalValuesFileName))
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
	data, err := os.ReadFile(valuesFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal values file, %s", valuesFilePath)
	}
	return values, nil
}

func GetClusterInfo(clusterInfoFilePath string) (types.ClusterInfo, error) {
	var values types.ClusterInfo
	data, err := os.ReadFile(clusterInfoFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read cluster info file, %s", clusterInfoFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal cluster info file, %s", clusterInfoFilePath)
	}
	return values, err
}

func (c CaptenConfig) PrepareFilePath(dir, path string) string {
	return fmt.Sprintf("%s%s%s", c.CurrentDirPath, dir, path)
}

func (c CaptenConfig) PrepareDirPath(dir string) string {
	return fmt.Sprintf("%s%s", c.CurrentDirPath, dir)
}

func addCurrentDirToPath(dir string) error {
	path := os.Getenv("PATH")
	if strings.Contains(path, dir) {
		return nil
	}

	newPath := fmt.Sprintf("%s:%s", dir, path)
	err := os.Setenv("PATH", newPath)
	if err != nil {
		return err
	}
	return nil
}
