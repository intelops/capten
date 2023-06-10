package types

type SecretInfo struct {
	FilePath string `mapstructure:"filePath"`
	Key      string `mapstructure:"key"`
}

type Override struct {
	StringValues map[string]interface{} `mapstructure:"stringValues"`
	Values       map[string]interface{} `mapstructure:"values"`
}

type ChartInfo struct {
	Name            string       `mapstructure:"name"`
	ChartName       string       `mapstructure:"chartName"`
	RepoName        string       `mapstructure:"repoName"`
	RepoURL         string       `mapstructure:"repoURL"`
	Namespace       string       `mapstructure:"namespace"`
	ReleaseName     string       `mapstructure:"releaseName"`
	Version         string       `mapstructure:"version"`
	Override        Override     `mapstructure:"override"`
	CreateNamespace bool         `mapstructure:"createNamespace"`
	MakeNsPrivilege bool         `mapstructure:"createNamespace"`
	SecretInfos     []SecretInfo `mapstructure:"secretInfos"`
}
