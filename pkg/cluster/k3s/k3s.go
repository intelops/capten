package k3s

import (
	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/k8s"

	"os"
	"text/template"

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

	} else if captenConfig.CloudService == "azure" {
		azureClusterInfo, err := config.GetClusterInfoAzure(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
		if err != nil {
			return nil, err
		}
		clusterInfo = azureClusterInfo
	} else {
		return nil, errors.Errorf("Unsupported Cloud Service")
	}

	return clusterInfo, nil
}

func createOrDestroyCluster(captenConfig config.CaptenConfig, action string) error {
	clog.Logger.Debugf("%s cluster on %s cloud with %s cluster type", action, captenConfig.CloudService, captenConfig.ClusterType)

	clusterInfo, err := getClusterInfo(captenConfig)
	if err != nil {
		return err
	}

	switch info := clusterInfo.(type) {
	case types.AWSClusterInfo:
		info.ConfigFolderPath = captenConfig.PrepareDirPath(captenConfig.ConfigDirPath)
		info.TerraformModulesDirPath = captenConfig.PrepareDirPath(captenConfig.TerraformModulesDirPath)
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
		info.TerraformModulesDirPath = captenConfig.PrepareDirPath(captenConfig.TerraformModulesDirPath)
		err = generateTemplateVarFile(captenConfig, info, captenConfig.AzureTerraformTemplateFileName)
		if err != nil {
			return err
		}

		tf, err := terraform.NewAzure(captenConfig, info)
		if err != nil {
			return errors.WithMessage(err, "failed to initialize the terraform")
		}

		if action == "create" {
			err = tf.Apply()
			if err != nil {
				return err
			}
			kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
			clientSet, err := k8s.GetK8SClient(kubeconfigPath)
			if err != nil {
				return err
			}
			err = k8s.UpdateNodeLabels(clientSet, &info)

			if err != nil {
				clog.Logger.Debug("Error while updating node label", err)
				return err
			}
		} else if action == "destroy" {
			return tf.Destroy()
		}
	default:
		return errors.New("unsupported cloud service")
	}

	return nil
}

func Create(captenConfig config.CaptenConfig) error {

	return createOrDestroyCluster(captenConfig, "create")

}

func Destroy(captenConfig config.CaptenConfig) error {
	return createOrDestroyCluster(captenConfig, "destroy")
}

func generateTemplateVarFile(captenConfig config.CaptenConfig, clusterInfo interface{}, templateFileName string) error {
	content, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath, templateFileName))
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
