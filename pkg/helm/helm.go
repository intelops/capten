package helm

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"capten/pkg/config"
	"capten/pkg/k8s"
	"capten/pkg/util"
)

type helm struct {
	config         *viper.Viper
	settings       *cli.EnvSettings
	defaultTimeout time.Duration
}

type SecretInfo struct {
	FilePath string `mapstructure:"FilePath"`
	Key      string `mapstructure:"Key"`
}

type ChartInfo struct {
	Name            string                 `mapstructure:"Name"`
	ChartName       string                 `mapstructure:"ChartName"`
	RepoName        string                 `mapstructure:"RepoName"`
	RepoURL         string                 `mapstructure:"RepoURL"`
	Namespace       string                 `mapstructure:"Namespace"`
	ReleaseName     string                 `mapstructure:"ReleaseName"`
	Version         string                 `mapstructure:"Version"`
	Override        map[string]interface{} `mapstructure:"Override"`
	CreateNamespace bool                   `mapstructure:"CreateNamespace"`
	MakeNsPrivilege bool                   `mapstructure:"CreateNamespace"`
	SecretInfos     []SecretInfo           `mapstructure:"SecretInfos"`
}

func NewHelm(configPath, kubeConfigPath string) (*helm, error) {
	cfg, err := config.GetClusterConfig(configPath)
	if err != nil {
		return nil, err
	}

	settings := cli.New()
	settings.KubeConfig = kubeConfigPath

	return &helm{
		config:         cfg,
		settings:       settings,
		defaultTimeout: time.Second * 900,
	}, nil
}

func (h *helm) Install() {
	chartInfos := make([]ChartInfo, 0)
	if err := h.config.UnmarshalKey("PostInstall", &chartInfos); err != nil {
		log.Println("failed to get postInstall app info from config", err)
		return
	}

	globalValues := h.config.GetStringMap("GlobalValues")
	for _, chartInfo := range chartInfos {
		if err := populateSecretValues(&chartInfo); err != nil {
			logrus.Error("failed to populate secret values", err)
			return
		}

		chartInfo.Override = util.MergeMap(util.ProcessMap(globalValues), util.ProcessMap(chartInfo.Override))
		logrus.Debugf("chart Overrides are %+v", chartInfo.Override)
		if err := h.Run(context.Background(), chartInfo); err != nil {
			logrus.Error("install failed", err)
			return
		}

		if chartInfo.MakeNsPrivilege {
			err := k8s.MakeNamespacePrivilege(h.settings.KubeConfig, chartInfo.Namespace)
			if err != nil {
				logrus.Error("failed to patch namespace with privilege", err)
				return
			}
		}
	}
}

func (h *helm) Run(ctx context.Context, chartInfo ChartInfo) error {
	repoEntry := &repo.Entry{
		Name: chartInfo.RepoName,
		URL:  chartInfo.RepoURL,
	}

	settings := cli.New()
	settings.KubeConfig = h.settings.KubeConfig
	settings.SetNamespace(chartInfo.Namespace)
	r, err := repo.NewChartRepository(repoEntry, getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to create new repo")
	}

	r.CachePath = settings.RepositoryCache
	fmt.Println("cache path", settings.RepositoryCache)
	fmt.Println("repo config", settings.RepositoryConfig)
	_, err = r.DownloadIndexFile()
	if err != nil {
		return errors.Wrap(err, "unable to download the index file")
	}

	var repoFile repo.File
	repoFile.Update(repoEntry)
	err = repoFile.WriteFile(settings.RepositoryConfig, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write the helm-chart path")
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), chartInfo.Namespace, "", debug)
	if err != nil {
		return errors.Wrap(err, "failed to setup actionConfig for helm")
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = chartInfo.Namespace
	client.ReleaseName = chartInfo.ReleaseName
	client.Version = chartInfo.Version
	client.Timeout = h.defaultTimeout
	client.CreateNamespace = chartInfo.CreateNamespace
	cp, err := client.ChartPathOptions.LocateChart(chartInfo.ChartName, settings)
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "chart load error")
	}

	if len(chartInfo.Override) == 0 {
		_, err = client.Run(chartReq, nil)
		return errors.Wrap(err, "chart run error")
	}

	valOptions := &values.Options{
		Values: getValues(chartInfo.Override),
	}

	vals, err := valOptions.MergeValues(getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to merge the helm values")
	}

	_, err = client.Run(chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	//log.Println("release info", releaseInfo)
	return nil
}

func debug(format string, v ...interface{}) {
	//format = fmt.Sprintf("[debug] %v\n", format)
	log.Output(2, fmt.Sprintf(format, v))
}

func getValues(values map[string]interface{}) []string {
	vals := make([]string, 0)
	for key, val := range values {
		vals = append(vals, fmt.Sprintf("%v=%v", key, val))
	}

	return vals
}

// populateSecretValues reads the data from secretInfo converts file content to base64 encoded overrides for helm
func populateSecretValues(info *ChartInfo) error {
	if info.Override == nil {
		info.Override = make(map[string]interface{})
	}

	for _, secretInfo := range info.SecretInfos {
		content, err := os.ReadFile(secretInfo.FilePath)
		if err != nil {
			return errors.Wrapf(err, "unable to read secret file %v", secretInfo.FilePath)
		}

		info.Override[secretInfo.Key] = base64.StdEncoding.EncodeToString(content)
	}

	return nil
}
