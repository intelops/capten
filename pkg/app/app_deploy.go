package app

import (
	"bytes"
	"capten/pkg/config"
	"capten/pkg/helm"
	"capten/pkg/k8s"
	"capten/pkg/types"
	"context"
	"html/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func DeployApps(captenConfig config.CaptenConfig) error {
	globalValues, err := GetClusterGlobalValues(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CaptenGlobalValuesFileName))
	if err != nil {
		return err
	}

	coreAppGroupAppConfigs, err := prepareAppGroupConfigs(captenConfig, globalValues, captenConfig.CoreAppGroupsFileName)
	if err != nil {
		return err
	}

	defaultAppGroupAppConfigs, err := prepareAppGroupConfigs(captenConfig, globalValues, captenConfig.DefaultAppGroupsFileName)
	if err != nil {
		return err
	}

	hc, err := helm.NewClient(captenConfig)
	if err != nil {
		return err
	}

	status := installAppGroup(captenConfig, hc, coreAppGroupAppConfigs)
	if !status {
		return errors.New("core applications deployment failed")
	}

	status = installAppGroup(captenConfig, hc, defaultAppGroupAppConfigs)
	if !status {
		return errors.New("default applications deployment failed")
	}
	return nil
}

func installAppGroup(captenConfig config.CaptenConfig, hc *helm.Client, appConfigs []types.AppConfig) bool {
	successStatus := true
	for _, appConfig := range appConfigs {
		logrus.Infof("[app: %s] installing", appConfig.Name)
		if err := hc.Install(context.Background(), appConfig); err != nil {
			logrus.Errorf("%s installation failed, %v", appConfig.Name, err)
			successStatus = false
			continue
		}
		logrus.Infof("[app: %s] installed", appConfig.Name)

		if err := WriteAppConfig(captenConfig, appConfig); err != nil {
			logrus.Errorf("failed to write %s config, %v", appConfig.Name, err)
			successStatus = false
			continue
		}
		if appConfig.PrivilegedNamespace {
			err := k8s.MakeNamespacePrivilege(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName),
				appConfig.Namespace)
			if err != nil {
				logrus.Error("failed to patch namespace with privilege", err)
				successStatus = false
				continue
			}
		}
	}
	return successStatus
}

func prepareAppGroupConfigs(captenConfig config.CaptenConfig, globalValues map[string]interface{},
	appGroupNameFile string) (appConfigs []types.AppConfig, err error) {
	var apps []string
	apps, err = GetApps(captenConfig.PrepareFilePath(captenConfig.AppsDirPath, appGroupNameFile))
	if err != nil {
		return
	}

	appConfigs = []types.AppConfig{}
	for _, appName := range apps {
		var appConfig types.AppConfig
		appConfig, err = GetAppConfig(captenConfig.PrepareFilePath(captenConfig.AppsConfigDirPath, appName+".yaml"))
		if err != nil {
			err = errors.WithMessagef(err, "failed load %s config", appName)
			return
		}
		appConfig.Override.Values, err = replaceTemplateValues(appConfig.Override.Values, globalValues)
		if err != nil {
			err = errors.WithMessagef(err, "failed transform %s values", appName)
			return
		}
		appConfigs = append(appConfigs, appConfig)
		logrus.Debug(appName, " : ", appConfig)
	}
	return
}

func replaceTemplateValues(templateData map[string]interface{},
	values map[string]interface{}) (transformedData map[string]interface{}, err error) {
	yamlData, err := yaml.Marshal(templateData)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = map[string]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &transformedData)
	if err != nil {
		return
	}
	return
}
