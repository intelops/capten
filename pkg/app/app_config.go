package app

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	folderPrmission os.FileMode = 0755
	filePrmission   os.FileMode = 0644
)

func GetClusterGlobalValues(valuesFilePath string) (map[string]interface{}, error) {
	var values map[string]interface{}
	data, err := os.ReadFile(valuesFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal values file, %s", valuesFilePath)
	}
	return values, nil
}

func GetAppGroups(appGroupsFilePath string) ([]string, error) {
	var values types.AppGroupList
	data, err := os.ReadFile(appGroupsFilePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read file, %s", appGroupsFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal file, %s", appGroupsFilePath)
	}
	return values.Groups, err
}

func GetApps(appListFilePath string) ([]string, error) {
	var values types.AppList
	data, err := os.ReadFile(appListFilePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read values file, %s", appListFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal values file, %s", appListFilePath)
	}
	return values.Apps, err
}

func GetAppConfig(appConfigFilePath string) (types.AppConfig, error) {
	var values types.AppConfig
	data, err := os.ReadFile(appConfigFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read values file, %s", appConfigFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal values file, %s", appConfigFilePath)
	}
	return values, err
}

func WriteAppConfig(captenConfig config.CaptenConfig, appConfig types.AppConfig) error {
	err := os.MkdirAll(captenConfig.PrepareDirPath(captenConfig.AppsTempDirPath), folderPrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to create directory %s", captenConfig.AppsTempDirPath)
	}

	data, err := yaml.Marshal(&appConfig)
	if err != nil {
		return errors.WithMessagef(err, "failed to unmarshal %s app config", appConfig.Name)
	}

	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.AppsTempDirPath, appConfig.Name+".yaml"), data, filePrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to write %s app config to file", appConfig.Name)
	}
	return nil
}

func PrepareGlobalVaules(captenConfig config.CaptenConfig) (map[string]interface{}, error) {
	globalValues, err := GetClusterGlobalValues(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CaptenGlobalValuesFileName))
	if err != nil {
		return nil, err
	}

	err = generateAppGlobalValuesandAppend(globalValues)
	if err != nil {
		return nil, err
	}
	return globalValues, err
}
