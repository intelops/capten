package types

type AppList struct {
	Apps []string `yaml:"Apps"`
}

type LaunchUIConfig struct {
	RedirectURL string `yaml:"RedirectURL"`
}

type Override struct {
	LaunchUIConfig LaunchUIConfig         `yaml:"LaunchUIConfig"`
	LaunchUIValues map[string]interface{} `yaml:"LaunchUIValues"`
	Values         map[string]interface{} `yaml:"Values"`
}

type AppConfig struct {
	Name                string   `yaml:"Name"`
	ChartName           string   `yaml:"ChartName"`
	RepoName            string   `yaml:"RepoName"`
	RepoURL             string   `yaml:"RepoURL"`
	Namespace           string   `yaml:"Namespace"`
	ReleaseName         string   `yaml:"ReleaseName"`
	Version             string   `yaml:"Version"`
	Override            Override `yaml:"Override"`
	CreateNamespace     bool     `yaml:"CreateNamespace"`
	PrivilegedNamespace bool     `yaml:"PrivilegedNamespace"`
}

type ClusterInfo struct {
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
