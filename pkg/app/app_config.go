package app

import (
	"capten/pkg/config"
	"capten/pkg/types"
	"io/ioutil"
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
	data, err := ioutil.ReadFile(valuesFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal values file, %s", valuesFilePath)
	}
	return values, nil
}

func GetAppList(appListFilePath string) ([]string, error) {
	var values types.AppList
	data, err := ioutil.ReadFile(appListFilePath)
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
	data, err := ioutil.ReadFile(appConfigFilePath)
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

	err = ioutil.WriteFile(captenConfig.PrepareFilePath(captenConfig.AppsTempDirPath, appConfig.Name+".yaml"), data, filePrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to write %s app config to file", appConfig.Name)
	}
	return nil
}
