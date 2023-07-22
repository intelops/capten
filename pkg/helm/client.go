package helm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
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

func (h *Client) Install(ctx context.Context, appConfig types.AppConfig) (alreadyInstalled bool, err error) {
	repoEntry := &repo.Entry{
		Name: appConfig.RepoName,
		URL:  appConfig.RepoURL,
	}

	settings := cli.New()
	settings.KubeConfig = h.settings.KubeConfig
	settings.SetNamespace(appConfig.Namespace)
	r, err := repo.NewChartRepository(repoEntry, getter.All(settings))
	if err != nil {
		err = errors.Wrap(err, "failed to create new repo")
		return
	}

	r.CachePath = settings.RepositoryCache
	_, err = r.DownloadIndexFile()
	if err != nil {
		err = errors.Wrap(err, "unable to download the index file")
		return
	}

	var repoFile repo.File
	repoFile.Update(repoEntry)
	err = repoFile.WriteFile(settings.RepositoryConfig, 0644)
	if err != nil {
		err = errors.Wrap(err, "failed to write the helm-chart path")
		return
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), appConfig.Namespace, "", logHelmDebug)
	if err != nil {
		err = errors.Wrap(err, "failed to setup actionConfig for helm")
		return
	}

	alreadyInstalled, err = h.isAppInstalled(actionConfig, appConfig.ReleaseName)
	if err != nil {
		return
	}

	if !alreadyInstalled {
		err = h.installApp(ctx, settings, actionConfig, appConfig)
		return
	}

	if h.captenConfig.UpgradeAppIfInstalled {
		err = h.upgradeApp(ctx, settings, actionConfig, appConfig)
		return
	}
	return
}

func (h *Client) installApp(ctx context.Context, settings *cli.EnvSettings, actionConfig *action.Configuration, appConfig types.AppConfig) error {
	action.NewList(&action.Configuration{})
	client := action.NewInstall(actionConfig)
	client.Namespace = appConfig.Namespace
	client.ReleaseName = appConfig.ReleaseName
	client.Version = appConfig.Version
	client.Timeout = h.defaultTimeout
	client.CreateNamespace = appConfig.CreateNamespace
	client.DryRun = h.captenConfig.AppDeployDryRun
	client.Devel = h.captenConfig.AppDeployDebug
	client.ClientOnly = true

	cp, err := client.ChartPathOptions.LocateChart(appConfig.ChartName, settings)
	if err != nil {
		return errors.Wrap(err, "chart locate error")
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "chart load error")
	}

	if len(appConfig.OverrideValues) == 0 {
		_, err = client.Run(chartReq, nil)
		return errors.Wrap(err, "chart run error")
	}

	stringValues := getValues(appConfig.OverrideValues)
	valueOpts := &values.Options{StringValues: stringValues}
	vals, err := valueOpts.MergeValues(nil)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}
	appConfig.OverrideValues = vals

	releaseInfo, err := client.Run(chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	logrus.Debug("release info ", releaseInfo)
	return nil
}

func (h *Client) upgradeApp(ctx context.Context, settings *cli.EnvSettings, actionConfig *action.Configuration, appConfig types.AppConfig) error {
	action.NewList(&action.Configuration{})
	client := action.NewUpgrade(actionConfig)
	client.Namespace = appConfig.Namespace
	client.Version = appConfig.Version
	client.Timeout = h.defaultTimeout
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

	if len(appConfig.OverrideValues) == 0 {
		_, err = client.Run(appConfig.ReleaseName, chartReq, nil)
		return errors.Wrap(err, "chart run error")
	}

	stringValues := getValues(appConfig.OverrideValues)
	valueOpts := &values.Options{StringValues: stringValues}
	vals, err := valueOpts.MergeValues(nil)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}
	appConfig.OverrideValues = vals

	releaseInfo, err := client.Run(appConfig.ReleaseName, chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "chart run error")
	}

	logrus.Debug("release info ", releaseInfo)
	return nil
}

func (h *Client) isAppInstalled(actionConfig *action.Configuration, releaseName string) (bool, error) {
	releaseClient := action.NewList(actionConfig)
	releases, err := releaseClient.Run()
	if err != nil {
		return false, errors.WithMessage(err, "failed to check helm releases on cluster")
	}

	for _, release := range releases {
		if strings.EqualFold(release.Name, releaseName) {
			return true, nil
		}
	}
	return false, nil
}

func logHelmDebug(format string, v ...interface{}) {
	logrus.Debug(format, v)
}

func getValues(values map[string]interface{}) []string {
	vals := make([]string, 0)
	for key, val := range values {
		vals = append(vals, fmt.Sprintf("%v=%v", key, val))
	}

	return vals
}
