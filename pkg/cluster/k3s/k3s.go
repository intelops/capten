package k3s

import (
	"os"
	"text/template"

	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/terraform"
	"capten/pkg/types"

	"github.com/pkg/errors"
)


func getClusterInfo(captenConfig config.CaptenConfig) (interface{}, error) {
	var clusterInfo interface{}


	if captenConfig.CloudService == "aws" {
		awsclusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return nil, err
		}
		clusterInfo = awsclusterInfo

	} else {
		azureClusterInfo, err := config.GetClusterInfoAzure(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return nil, err
		}
		clusterInfo = azureClusterInfo
	}

	return clusterInfo, nil
}

func createOrUpdateCluster(captenConfig config.CaptenConfig, action string) error {
	clog.Logger.Debugf("%s cluster on %s cloud with %s cluster type", action, captenConfig.CloudService, captenConfig.ClusterType)

	clusterInfo, err := getClusterInfo(captenConfig)
	if err != nil {
		return err
	}

	switch info := clusterInfo.(type) {
	case types.AWSClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)
		err = generateTemplateVarFile(captenConfig, info, captenConfig.AWSTerraformTemplateFileName)
		if err != nil {
			return err
		}

		tf, err := terraform.NewAws(captenConfig, info)
		if err != nil {
			return errors.WithMessage(err, "failed to initialize the terraform")
		}

		if action == "create" {
			return tf.Apply()
		} else if action == "destroy" {
			return tf.Destroy()
		}
	case types.AzureClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)
		err = generateTemplateVarFile(captenConfig, info, captenConfig.AzureTerraformTemplateFileName)
		if err != nil {
			return err
		}

		tf, err := terraform.NewAzure(captenConfig, info)
		if err != nil {
			return errors.WithMessage(err, "failed to initialize the terraform")
		}

		if action == "create" {
			return tf.Apply()
		} else if action == "destroy" {
			return tf.Destroy()
		}
	default:
		return errors.New("unsupported cloud service")
	}

	return nil
}

func Create(captenConfig config.CaptenConfig) error {
	return createOrUpdateCluster(captenConfig, "create")
}

func Destroy(captenConfig config.CaptenConfig) error {
	return createOrUpdateCluster(captenConfig, "destroy")
}









func generateTemplateVarFile(captenConfig config.CaptenConfig, clusterInfo interface{},templateFileName string) error {
	content, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath,templateFileName ))
	if err != nil {
		return errors.WithMessage(err, "failed to read template file")
	}

	contentStr := string(content)
	templateObj, err := template.New("terraformTemplate").Parse(contentStr)

	if err != nil {
		clog.Logger.Error("Error while creating templateObj", err)
		return err
	}

	templateFile, err := os.Create(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath, captenConfig.TerraformVarFileName))

	if err != nil {
		clog.Logger.Error("Error while creating templateFile", err)
		return err
	}

	if err := templateObj.Execute(templateFile, clusterInfo); err != nil {
		clog.Logger.Error("Error while executing templateObj", err)
		return err
	}
	return nil
}