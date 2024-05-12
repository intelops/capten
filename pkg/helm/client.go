package helm

import (
	"bytes"
	"context"
	"html/template"

	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/types"
)

const (
	folderPrmission os.FileMode = 0755
	filePrmission   os.FileMode = 0644
)

type Client struct {
	Settings       *cli.EnvSettings
	defaultTimeout time.Duration
	captenConfig   config.CaptenConfig
}

func NewClient(captenConfig config.CaptenConfig) (*Client, error) {
	settings := cli.New()
	settings.KubeConfig = captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)

	err := os.MkdirAll(captenConfig.PrepareDirPath(captenConfig.AppValuesTempDirPath), folderPrmission)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create directory %s", captenConfig.AppsTempDirPath)
	}
	return &Client{
		captenConfig:   captenConfig,
		Settings:       settings,
		defaultTimeout: time.Second * 900,
	}, nil
}

func (h *Client) Install(ctx context.Context, appConfig *types.AppConfig) (alreadyInstalled bool, err error) {
	repoEntry := &repo.Entry{
		Name: appConfig.RepoName,
		URL:  appConfig.RepoURL,
	}

	settings := cli.New()
	settings.KubeConfig = h.Settings.KubeConfig
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
	err = actionConfig.Init(settings.RESTClientGetter(), appConfig.Namespace, "", LogHelmDebug)
	if err != nil {
		err = errors.Wrap(err, "failed to setup actionConfig for helm")
		return
	}

	alreadyInstalled, err = h.IsAppInstalled(actionConfig, appConfig.ReleaseName)

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

func (h *Client) installApp(ctx context.Context, settings *cli.EnvSettings, actionConfig *action.Configuration, appConfig *types.AppConfig) error {
	action.NewList(&action.Configuration{})
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
		return errors.Wrap(err, "failed to locate chart locate")
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "failed load chart")
	}

	if len(appConfig.TemplateValues) == 0 {
		_, err = client.Run(chartReq, nil)
		if err != nil {
			return errors.Wrap(err, "failed chart install run")
		}
		return nil
	}

	appValuesFile, err := h.prepareAppValues(appConfig)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(appValuesFile) }()

	valueOpts := &values.Options{ValueFiles: []string{appValuesFile}}
	vals, err := valueOpts.MergeValues(getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to merge chart values")
	}

	releaseInfo, err := client.Run(chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "failed chart install run with values")
	}

	clog.Logger.Debugf("release info: %v", releaseInfo)
	return nil
}

func (h *Client) upgradeApp(ctx context.Context, settings *cli.EnvSettings, actionConfig *action.Configuration, appConfig *types.AppConfig) error {
	action.NewList(&action.Configuration{})
	client := action.NewUpgrade(actionConfig)
	client.Namespace = appConfig.Namespace
	client.Version = appConfig.Version
	client.Timeout = h.defaultTimeout
	client.DryRun = h.captenConfig.AppDeployDryRun
	client.Devel = h.captenConfig.AppDeployDebug
	client.ResetValues = true

	cp, err := client.ChartPathOptions.LocateChart(appConfig.ChartName, settings)
	if err != nil {
		return errors.Wrap(err, "failed to locate chart locate")
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		return errors.Wrap(err, "failed load chart")
	}

	if len(appConfig.TemplateValues) == 0 {
		_, err = client.Run(appConfig.ReleaseName, chartReq, nil)
		if err != nil {
			return errors.Wrap(err, "failed chart upgrade run")
		}
		return nil
	}

	appValuesFile, err := h.prepareAppValues(appConfig)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(appValuesFile) }()

	valueOpts := &values.Options{ValueFiles: []string{appValuesFile}}
	vals, err := valueOpts.MergeValues(getter.All(settings))
	if err != nil {
		return errors.Wrap(err, "failed to merge chart values")
	}

	releaseInfo, err := client.Run(appConfig.ReleaseName, chartReq, vals)
	if err != nil {
		return errors.Wrap(err, "failed chart upgrade run with values")
	}

	clog.Logger.Debugf("release info: %v", releaseInfo)
	return nil
}

func (h *Client) IsAppInstalled(actionConfig *action.Configuration, releaseName string) (bool, error) {
	releaseClient := action.NewList(actionConfig)
	releases, err := releaseClient.Run()
	if err != nil {
		return false, errors.WithMessage(err, "failed to check helm releases on cluster")
	}

	for _, release := range releases {
		if strings.EqualFold(release.Name, releaseName) {

			if release.Info.Status == "deployed" {
				return true, nil
			} else {
				return false, nil
			}

		}
	}
	return false, nil
}

func LogHelmDebug(format string, v ...interface{}) {
	clog.Logger.Debug(format, v)
}

func (h *Client) prepareAppValues(appConfig *types.AppConfig) (string, error) {
	transformedData, err := executeAppConfigTemplate(appConfig.TemplateValues, appConfig.OverrideValues)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to transform app values")
	}

	tmpValuesPath := h.captenConfig.PrepareFilePath(h.captenConfig.AppValuesTempDirPath, appConfig.Name+"-values.yaml")
	err = os.WriteFile(tmpValuesPath, transformedData, filePrmission)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to write app values to file %s", tmpValuesPath)
	}
	return tmpValuesPath, nil
}

func executeAppConfigTemplate(data []byte, values map[string]interface{}) (transformedData []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}
