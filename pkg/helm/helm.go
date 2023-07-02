package helm

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"capten/pkg/config"
	"capten/pkg/k8s"
	"capten/pkg/types"
	"capten/pkg/util"
)

type helm struct {
	config         *viper.Viper
	settings       *cli.EnvSettings
	defaultTimeout time.Duration
}

func NewHelm(captenConfig config.CaptenConfig) (*helm, error) {
	// cfg, err := config.GetClusterConfig(captenConfig.LegacyAppConfigFilePath)
	// if err != nil {
	// 	return nil, err
	// }

	settings := cli.New()
	settings.KubeConfig = captenConfig.KubeConfigPath

	return &helm{
		// config:         cfg,
		settings:       settings,
		defaultTimeout: time.Second * 900,
	}, nil
}

func (h *helm) Install(cnf config.CaptenConfig) error {
	appListFile, err := os.Open(cnf.AppsFilePath)
	if err != nil {
		return errors.Wrap(err, "helm install failed")
	}
	appsList, err := getAppsList(appListFile, "apps")
	if err != nil {
		return errors.Wrap(err, "helm install failed")
	}

	chartInfos := readAppValuesDir(cnf.AppValuesDir, appsList)
	// chartInfos := make([]types.ChartInfo, 0)
	// if err := h.config.UnmarshalKey("PostInstall", &chartInfos); err != nil {
	// 	log.Println("failed to get postInstall app info from config", err)
	// 	return err
	// }

	globalStringValues := h.config.GetStringMap("GlobalValues.StringValues")
	globalValues := h.config.GetStringMap("GlobalValues.Values")
	for _, chartInfo := range chartInfos {
		if err := populateSecretValues(&chartInfo); err != nil {
			logrus.Error("failed to populate secret values", err)
			continue
		}

		chartInfo.Override.StringValues = util.MergeMap(util.ProcessMap(globalStringValues), util.ProcessMap(chartInfo.Override.StringValues))
		chartInfo.Override.Values = util.MergeMap(util.ProcessMap(globalValues), util.ProcessMap(chartInfo.Override.Values))
		logrus.Debugf("chart Overrides are %v", chartInfo.Override)
		if err := h.Run(context.Background(), chartInfo); err != nil {
			logrus.Error("install failed", err)
			continue
		}

		if chartInfo.MakeNsPrivilege {
			err := k8s.MakeNamespacePrivilege(h.settings.KubeConfig, chartInfo.Namespace)
			if err != nil {
				logrus.Error("failed to patch namespace with privilege", err)
				continue
			}
		}
	}
	return nil

}

func (h *helm) Run(ctx context.Context, chartInfo types.ChartInfo) error {
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
	//client.DryRun = true
	cp, err := client.ChartPathOptions.LocateChart(chartInfo.ChartName, settings)
	if err != nil {
		return errors.Wrap(err, "locate chart error")
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "chart load error")
	}

	if len(chartInfo.Override.Values) == 0 && len(chartInfo.Override.StringValues) == 0 {
		_, err = client.Run(chartReq, nil)
		return errors.Wrap(err, "chart run error")
	}

	valOptions := &values.Options{
		StringValues: getValues(chartInfo.Override.StringValues),
		Values:       getValues(chartInfo.Override.Values),
	}

	vals, err := valOptions.MergeValues(getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to merge the helm values")
	}

	releaseInfo, err := client.Run(chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	logrus.Debug("release info ", releaseInfo)
	return nil
}

func debug(format string, args ...interface{}) {
	//format = fmt.Sprintf("[debug] %v\n", format)
	log.Output(2, fmt.Sprintf(format, args...))
}

func getValues(values map[string]interface{}) []string {
	vals := make([]string, 0)
	for key, val := range values {
		vals = append(vals, fmt.Sprintf("%v=%v", key, val))
	}

	return vals
}

// populateSecretValues reads the data from secretInfo converts file content to base64 encoded overrides for helm
func populateSecretValues(info *types.ChartInfo) error {
	if info.Override.StringValues == nil {
		info.Override.StringValues = make(map[string]interface{})
	}

	for _, secretInfo := range info.SecretInfos {
		content, err := os.ReadFile(secretInfo.FilePath)
		if err != nil {
			return errors.Wrapf(err, "unable to read secret file %v", secretInfo.FilePath)
		}

		info.Override.StringValues[secretInfo.Key] = base64.StdEncoding.EncodeToString(content)
	}

	return nil
}

// basePath: dir path with all the app values
func readAppValuesDir(basePath string, appsList []string) (ret []types.ChartInfo) {

	for _, app := range appsList {
		filePath := path.Join(basePath, app+".yaml")
		yamlFile, err := os.Open(filePath)
		if err != nil {
			logrus.Errorf("err: %v, yamlfile missing: %s", err, filePath)
			continue
		}
		var charInfo types.ChartInfo
		if err := yaml.NewDecoder(yamlFile).Decode(&charInfo); err != nil {
			logrus.Errorf("err: %v, decoding file: %s", err, filePath)
			continue
		}
		ret = append(ret, charInfo)
	}

	return
}

// useful for reading list like yaml files
func getAppsList(r io.Reader, key string) (ret []string, err error) {

	var appMapping map[string][]string
	err = yaml.NewDecoder(r).Decode(&appMapping)
	if err != nil {
		logrus.Errorf("decode yaml, err: %v", err)
		return nil, err
	}
	ret = append(ret, appMapping[key]...)

	return
}
