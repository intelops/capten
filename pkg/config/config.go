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
	CaptenClusterValues
	AgentHostName              string   `envconfig:"AGENT_HOST_NAME" default:"captenagent"`
	AgentHostPort              string   `envconfig:"AGENT_HOST_PORT" default:":443"`
	AgentSecure                bool     `envconfig:"AGENT_SECURE" default:"true"`
	CaptenNamespace            string   `envconfig:"CAPTEN_NAMESPACE" default:"capten"`
	CertManagerNamespace       string   `envconfig:"CERT_MANAGER_NAMESPACE" default:"cert-manager"`
	ClusterCACertSecretName    string   `envconfig:"INTER_CERT_SECRET_NAME" default:"capten-ca-cert"`
	InterCACertFileName        string   `envconfig:"INTER_CERT_FILE_NAME" default:"inter-ca.crt"`
	InterCAKeyFileName         string   `envconfig:"INTER_CERT_KEY_FILE_NAME" default:"inter-ca.key"`
	AgentCertSecretName        string   `envconfig:"AGENT_CERT_SECRET_NAME" default:"kad-agent-cert"`
	AgentCACertSecretName      string   `envconfig:"AGENT_CA_CERT_SECRET_NAME" default:"kad-agent-ca-cert"`
	AppsDirPath                string   `envconfig:"APPS_DIR_PATH" default:"/apps/"`
	AppsConfigDirPath          string   `envconfig:"APPS_CONFIG_DIR_PATH" default:"/apps/conf/"`
	AppsValuesDirPath          string   `envconfig:"APPS_VALUES_DIR_PATH" default:"/apps/conf/values/"`
	AppsTempDirPath            string   `envconfig:"APPS_TEMP_DIR_PATH" default:"/apps/tmp/"`
	AppValuesTempDirPath       string   `envconfig:"APPS_TEMPVAL_DIR_PATH" default:"/apps/tmp/val/"`
	AppIconsDirPath            string   `envconfig:"APPS_ICON_DIR_PATH" default:"/apps/icons/"`
	ConfigDirPath              string   `envconfig:"CONFIG_DIR_PATH" default:"/config/"`
	CertDirPath                string   `envconfig:"CERT_DIR_PATH" default:"/cert/"`
	TerraformModulesDirPath    string   `envconfig:"TERRAFORM_MODULE_DIR_PATH" default:"/terraform_modules/"`
	TerraformTemplateDirPath   string   `envconfig:"TERRAFORM_TEMPLATE_DIR_PATH" default:"/templates/k3s/"`
	CoreAppGroupsFileName      string   `envconfig:"CORE_APP_GROUPS_FILE_NAME" default:"core_group_apps.yaml"`
	DefaultAppGroupsFileName   string   `envconfig:"DEFAULT_APP_GROUPS_FILE_NAME" default:"default_group_apps.yaml"`
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
	AppDeployDryRun            bool     `envconfig:"APP_DEPLOY_DRYRUN" default:"false"`
	AppDeployDebug             bool     `envconfig:"APP_DEPLOY_DEBUG" default:"false"`
	StoreCredOnAgent           bool     `envconfig:"STORE_CRED_ON_AGENT" default:"true"`
	SkipAppsDeploy             bool     `envconfig:"SKIP_APPS_DEPLOY" default:"false"`
	ForceGenerateCerts         bool     `envconfig:"FORCE_GENERATE_CERTS" default:"false"`
	UpgradeAppIfInstalled      bool     `envconfig:"UPGRADE_APP_IF_INSTALLED" default:"false"`
	TerraformInitReconfigure   bool     `envconfig:"TERRAFORM_INIT_RECONFIGURE" default:"true"`
	TerraformInitUpgrade       bool     `envconfig:"TERRAFORM_INIT_UPGRADE" default:"true"`
	AgentDNSNames              []string
	CurrentDirPath             string
	PoolClusterName            string `envconfig:"POOL_CLUSTER_NAME" default:"cstor-disk-pool"`
	PoolClusterNamespace       string `envconfig:"POOL_CLUSTER_NAMESPACE" default:"openebs-cstor"`
}

type CaptenClusterValues struct {
	DomainName       string `yaml:"DomainName" envconfig:"DOMAIN_NAME" default:"dev.intelops.app"`
	LoadBalancerHost string `yaml:"LoadBalancerHost" envconfig:"CLUSTER_LB_HOST"`
	CloudService     string `yaml:"CloudService" envconfig:"CLOUD_SERVICE"`
	ClusterType      string `yaml:"ClusterType" envconfig:"CLUSTER_TYPE"`
	ClusterCAIssuer  string `yaml:"ClusterCAIssuer" envconfig:"CLUSTER_CA_ISSUER" default:"capten-issuer"`
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
	if len(values.CloudService) != 0 {
		cfg.CloudService = values.CloudService
	}
	if len(values.ClusterType) != 0 {
		cfg.ClusterType = values.ClusterType
	}
	if len(values.LoadBalancerHost) != 0 {
		cfg.LoadBalancerHost = values.LoadBalancerHost
	}

	cfg.AgentDNSNames = []string{}
	for _, prefixName := range cfg.AgentDNSNamePrefixes {
		cfg.AgentDNSNames = append(cfg.AgentDNSNames, prefixName+"."+cfg.DomainName)
	}
	return cfg, err
}

func (c CaptenConfig) GetCaptenAgentEndpoint() string {
	if c.AgentSecure {
		return fmt.Sprintf("%s.%s%s", c.AgentHostName, c.DomainName, c.AgentHostPort)
	}
	return fmt.Sprintf("%s.%s:80", c.AgentHostName, c.DomainName)
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

func UpdateClusterValues(cfg *CaptenConfig, cloudService, clusterType string) error {
	clusterValuesPath := cfg.PrepareFilePath(cfg.ConfigDirPath, cfg.CaptenGlobalValuesFileName)
	clusterValues, err := GetCaptenClusterValues(clusterValuesPath)
	if err != nil {
		return err
	}
	clusterValues.CloudService = cloudService
	clusterValues.ClusterType = clusterType
	clusterValuesData, err := yaml.Marshal(&clusterValues)
	if err != nil {
		return err
	}

	err = os.WriteFile(clusterValuesPath, clusterValuesData, 0644)
	if err != nil {
		return err
	}
	cfg.CloudService = cloudService
	cfg.ClusterType = clusterType
	return nil
}
