package helm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"capten/pkg/config"
	"capten/pkg/types"
)

type Client struct {
	settings       *cli.EnvSettings
	defaultTimeout time.Duration
	captenConfig   config.CaptenConfig
}

func NewClient(captenConfig config.CaptenConfig) (*Client, error) {
	settings := cli.New()
	settings.KubeConfig = captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	return &Client{
		captenConfig:   captenConfig,
		settings:       settings,
		defaultTimeout: time.Second * 900,
	}, nil
}

func (h *Client) Install(ctx context.Context, appConfig types.AppConfig) error {
	repoEntry := &repo.Entry{
		Name: appConfig.RepoName,
		URL:  appConfig.RepoURL,
	}

	settings := cli.New()
	settings.KubeConfig = h.settings.KubeConfig
	settings.SetNamespace(appConfig.Namespace)
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
	err = actionConfig.Init(settings.RESTClientGetter(), appConfig.Namespace, "", debug)
	if err != nil {
		return errors.Wrap(err, "failed to setup actionConfig for helm")
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = appConfig.Namespace
	client.ReleaseName = appConfig.ReleaseName
	client.Version = appConfig.Version
	client.Timeout = h.defaultTimeout
	client.CreateNamespace = appConfig.CreateNamespace
	client.DryRun = h.captenConfig.AppDeployDryRun
	client.Devel = h.captenConfig.AppDeployDebug

	cp, err := client.ChartPathOptions.LocateChart(appConfig.ChartName, settings)
	if err != nil {
		return errors.Wrap(err, "chart locate error")
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "chart load error")
	}

	if len(appConfig.Override.Values) == 0 {
		_, err = client.Run(chartReq, nil)
		return errors.Wrap(err, "chart run error")
	}

	releaseInfo, err := client.Run(chartReq, appConfig.Override.Values)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	logrus.Debug("release info ", releaseInfo)
	return nil
}

func debug(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(format, v))
}

func getValues(values map[string]interface{}) []string {
	vals := make([]string, 0)
	for key, val := range values {
		vals = append(vals, fmt.Sprintf("%v=%v", key, val))
	}

	return vals
}
