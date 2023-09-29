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

func Create(captenConfig config.CaptenConfig) error {
	clog.Logger.Debugf("creating cluster on %s cloud with %s cluster type", captenConfig.CloudService, captenConfig.ClusterType)

	var clusterInfo interface{}
	var err error

	if captenConfig.CloudService == "aws" {
		awsclusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return err
		}
		clusterInfo = awsclusterInfo

	} else {
		azureClusterInfo, err := config.GetClusterInfoAzure(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return err
		}
		clusterInfo = azureClusterInfo

	}

	switch info := clusterInfo.(type) {
	case types.AWSClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)

		err = generateTemplateVarFile(captenConfig, info)
		if err != nil {
			return err
		}

		tf, err := terraform.New(captenConfig, info)
		if err != nil {
			return errors.WithMessage(err, "failed to initialise the terraform")
		}
		return tf.Apply()
	case types.AzureClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)

		err = generateTemplateVarFile(captenConfig, info)
		if err != nil {

			return err
		}
		tf, err := terraform.NewAzure(captenConfig, info)

		if err != nil {

			return errors.WithMessage(err, "failed to initialise the terraform")
		}
		return tf.Apply()
	default:
		return errors.New("unsupported cloud service")
	}

}

func Destroy(captenConfig config.CaptenConfig) error {

	var clusterInfo interface{}
	var err error

	if captenConfig.CloudService == "aws" {
		awsclusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return err
		}
		clusterInfo = awsclusterInfo

	} else {
		azureClusterInfo, err := config.GetClusterInfoAzure(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return err
		}
		clusterInfo = azureClusterInfo

	}

	switch info := clusterInfo.(type) {
	case types.AWSClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)
		err = generateTemplateVarFile(captenConfig, info)
		if err != nil {
			return err
		}

		tf, err := terraform.New(captenConfig, info)
		if err != nil {
			return errors.WithMessage(err, "failed to initialise the terraform")
		}
		return tf.Destroy()
	case types.AzureClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)
		err = generateTemplateVarFile(captenConfig, info)
		if err != nil {

			return err
		}
		tf, err := terraform.NewAzure(captenConfig, info)

		if err != nil {

			return errors.WithMessage(err, "failed to initialise the terraform")
		}
		return tf.Destroy()
	default:
		return errors.New("unsupported cloud service")
	}

}

func generateTemplateVarFile(captenConfig config.CaptenConfig, clusterInfo interface{}) error {
	content, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath, captenConfig.AzureTerraformTemplateFileName))
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