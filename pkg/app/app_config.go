package app

import (
	"capten/pkg/config"
	"html/template"

	//"capten/pkg/k8s"
	"capten/pkg/types"

	"log"
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
		return values, errors.WithMessagef(err, "failed to read cluster values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal cluster values file, %s", valuesFilePath)
	}
	return values, nil
}

func GetApps(appListFilePath string) ([]string, error) {
	var values types.AppList
	data, err := os.ReadFile(appListFilePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read app group file, %s", appListFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal app group file, %s", appListFilePath)
	}
	return values.Apps, err
}

func GetAppConfig(appConfigFilePath string, globalValues map[string]interface{}) (types.AppConfig, error) {
	var values types.AppConfig
	data, err := os.ReadFile(appConfigFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read app config file, %s", appConfigFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal app config file, %s", appConfigFilePath)
	}
	return values, err
}

func GetAppValuesTemplate(captenConfig config.CaptenConfig, appName string) []byte {
	appValuesTemplateFilePath := captenConfig.PrepareFilePath(captenConfig.AppsValuesDirPath, appName+"_template.yaml")
	data, err := os.ReadFile(appValuesTemplateFilePath)
	if err != nil {
		return nil
	}
	return data
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

	cacertpath := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName)
	log.Println("cacertPath", cacertpath)
	readfile, err := os.ReadFile(cacertpath)
	if err != nil {
		log.Println("error while reading file", readfile)
	}
	//
	// cacert,err:=k8s.GetCACert(captenConfig)
	// if err != nil {
	// 	log.Println("error while getting cacer",err)
	// 	return nil, err
	// }
	globalValues["identityTrustAnchorsPEM"] = string(readfile)

	err = generateAppGlobalValuesandAppend(globalValues)
	if err != nil {
		return nil, err
	}
	log.Println("Global values", globalValues)
	return globalValues, err
}

type TemplateData struct {
	FilePath string
}

func LinkerdCreation(captenConfig config.CaptenConfig, appName string) {
	// Set the file path variable
	filePath := captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName)
	log.Println("cacertPath", filePath)
	// filePath := "/path/to/your/file.txt"
	readfile, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("error while reading file", readfile)
	}
	// Create template data with the file path
	data := TemplateData{
		FilePath: string(readfile),
	}
	templatefilepath := captenConfig.PrepareFilePath(captenConfig.AppsConfigDirPath, appName+".yaml")
	// Read the content of the template file
	templateContent, err := os.ReadFile(templatefilepath)
	if err != nil {
		panic(err)
	}

	// Parse the template
	tmpl, err := template.New("yaml").Parse(string(templateContent))
	if err != nil {
		panic(err)
	}

	// Execute the template and print the result to stdout
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
