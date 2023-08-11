package types

import (
	"capten/pkg/agent/agentpb"

	"gopkg.in/yaml.v2"
)

type AppGroupList struct {
	Groups []string `yaml:"AppGroups"`
}
type AppList struct {
	Apps []string `yaml:"Apps"`
}

type LaunchUIConfig struct {
	RedirectURL string `yaml:"RedirectURL"`
}

type AppConfig struct {
	Name                string                 `yaml:"Name"`
	ChartName           string                 `yaml:"ChartName"`
	RepoName            string                 `yaml:"RepoName"`
	RepoURL             string                 `yaml:"RepoURL"`
	Namespace           string                 `yaml:"Namespace"`
	ReleaseName         string                 `yaml:"ReleaseName"`
	Version             string                 `yaml:"Version"`
	Category            string                 `yaml:"Category"`
	Description         string                 `yaml:"Description"`
	LaunchURL           string                 `yaml:"LaunchURL"`
	LaunchRedirectURL   string                 `yaml:"LaunchRedirectURL"`
	LaunchUIConfig      LaunchUIConfig         `yaml:"LaunchUIConfig"`
	LaunchUIValues      map[string]interface{} `yaml:"LaunchUIValues"`
	OverrideValues      map[string]interface{} `yaml:"OverrideValues"`
	CreateNamespace     bool                   `yaml:"CreateNamespace"`
	PrivilegedNamespace bool                   `yaml:"PrivilegedNamespace"`
}

type ClusterInfo struct {
	ConfigFolderPath        string   `yaml:"ConfigFolderPath"`
	CloudService            string   `yaml:"CloudService"`
	ClusterType             string   `yaml:"ClusterType"`
	AwsAccessKey            string   `yaml:"AwsAccessKey"`
	AwsSecretKey            string   `yaml:"AwsSecretKey"`
	AlbName                 string   `yaml:"AlbName"`
	PrivateSubnet           string   `yaml:"PrivateSubnet"`
	Region                  string   `yaml:"Region"`
	SecurityGroupName       string   `yaml:"SecurityGroupName"`
	VpcCidr                 string   `yaml:"VpcCidr"`
	VpcName                 string   `yaml:"VpcName"`
	InstanceType            string   `yaml:"InstanceType"`
	NodeMonitoringEnabled   string   `yaml:"NodeMonitoringEnabled"`
	MasterCount             string   `yaml:"MasterCount"`
	WorkerCount             string   `yaml:"WorkerCount"`
	TraefikHttpPort         string   `yaml:"TraefikHttpPort"`
	TraefikHttpsPort        string   `yaml:"TraefikHttpsPort"`
	TalosTg                 string   `yaml:"TalosTg"`
	TraefikTg80Name         string   `yaml:"TraefikTg80Name"`
	TraefikTg443Name        string   `yaml:"TraefikTg443Name"`
	TraefikLbName           string   `yaml:"TraefikLbName"`
	TerraformBackendConfigs []string `yaml:"TerraformBackendConfigs"`
}

func (a AppConfig) ToSyncAppData() (agentpb.SyncAppData, error) {
	marshaledOverride, err := yaml.Marshal(a.OverrideValues)
	if err != nil {
		return agentpb.SyncAppData{}, err
	}

	marshaledLaunchUi, err := yaml.Marshal(a.LaunchUIValues)
	if err != nil {
		return agentpb.SyncAppData{}, err
	}

	return agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         a.ReleaseName,
			AppName:             a.Name,
			Version:             a.Version,
			Category:            a.Category,
			Description:         a.Description,
			ChartName:           a.ChartName,
			RepoName:            a.RepoName,
			RepoURL:             a.RepoURL,
			Namespace:           a.Namespace,
			CreateNamespace:     a.CreateNamespace,
			PrivilegedNamespace: a.PrivilegedNamespace,
			Icon:                []byte{},
			LaunchURL:           a.LaunchURL,
			LaunchRedirectURL:   a.LaunchRedirectURL,
		},
		Values: &agentpb.AppValues{
			OverrideValues: marshaledOverride,
			LaunchUIValues: marshaledLaunchUi,
		},
	}, nil
}
