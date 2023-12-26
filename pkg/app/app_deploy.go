package app

import (
	"bytes"
	"capten/pkg/config"
	"capten/pkg/helm"
	"capten/pkg/k8s"
	"capten/pkg/types"
	"context"
	"html/template"

	"capten/pkg/clog"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func DeployApps(captenConfig config.CaptenConfig, globalValues map[string]interface{}, groupFile string) error {
	appGroupAppConfigs, err := prepareAppGroupConfigs(captenConfig, globalValues, groupFile)
	if err != nil {
		return err
	}

	hc, err := helm.NewClient(captenConfig)
	if err != nil {
		return err
	}

	status := installAppGroup(captenConfig, hc, appGroupAppConfigs)
	if !status {
		return errors.New("applications deployment failed")
	}
	return nil
}

func installAppGroup(captenConfig config.CaptenConfig, hc *helm.Client, appConfigs []types.AppConfig) bool {
	successStatus := true
	for _, appConfig := range appConfigs {
		if appConfig.PrivilegedNamespace {
			err := k8s.CreateorUpdateNamespaceWithLabel(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName),
				appConfig.Namespace)
			if err != nil {
				clog.Logger.Error("failed to patch namespace with privilege", err)
				successStatus = false
				continue
			}
		}
		alreadyInstalled, err := hc.Install(context.Background(), &appConfig)
		if err != nil {
			clog.Logger.Errorf("[app: %s] installation failed, %v", appConfig.Name, err)
			successStatus = false
			continue
		}
		if alreadyInstalled {
			clog.Logger.Infof("[app: %s] already installed", appConfig.Name)
		} else {
			clog.Logger.Infof("[app: %s] installed", appConfig.Name)
		}

		appConfig.TemplateValues = nil
		if err := WriteAppConfig(captenConfig, appConfig); err != nil {
			clog.Logger.Errorf("failed to write %s config, %v", appConfig.Name, err)
			successStatus = false
			continue
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
		appConfig, err = GetAppConfig(captenConfig.PrepareFilePath(captenConfig.AppsConfigDirPath, appName+".yaml"), globalValues)
		if err != nil {
			err = errors.WithMessagef(err, "failed load %s config", appName)
			return
		}

		appConfig.TemplateValues = GetAppValuesTemplate(captenConfig, appName)
		appConfig.OverrideValues, err = replaceOverrideTemplateValues(appConfig.OverrideValues, globalValues)
		if err != nil {
			err = errors.WithMessagef(err, "failed transform '%s' override values", appName)
			return
		}

		appConfig.LaunchURL, err = replaceTemplateStringValues(appConfig.LaunchURL, globalValues)
		if err != nil {
			err = errors.WithMessagef(err, "failed transform '%s' string value", appName)
			return
		}
		appConfigs = append(appConfigs, appConfig)
		clog.Logger.Debug(appName, " : ", appConfig)
	}
	return
}

func replaceOverrideTemplateValues(templateData map[string]interface{},
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

func replaceTemplateStringValues(templateStringData string,
	values map[string]interface{}) (transformedStringData string, err error) {
	tmpl, err := template.New("templateVal").Parse(templateStringData)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}
	return buf.String(), nil
}
