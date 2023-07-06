package k3s

import (
	"os"
	"text/template"

	"capten/pkg/config"
	"capten/pkg/terraform"
	"capten/pkg/types"

	"github.com/pkg/errors"
)

func Create(captenConfig config.CaptenConfig, clusterType, cloudType string) error {
	clusterInfo, err := prepareClusterInfo(captenConfig, clusterType, cloudType)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath, captenConfig.TerraformTemplateFileName))
	if err != nil {
		return errors.WithMessage(err, "failed to read template file")
	}

	contentStr := string(content)
	templateObj, err := template.New("terraformTemplate").Parse(contentStr)
	if err != nil {
		return err
	}

	templateFile, err := os.Create(captenConfig.PrepareFilePath(captenConfig.TerraformTemplateDirPath, captenConfig.TerraformVarFileName))
	if err != nil {
		return err
	}

	if err := templateObj.Execute(templateFile, clusterInfo); err != nil {
		return err
	}

	tf, err := terraform.New(captenConfig, clusterInfo)
	if err != nil {
		return errors.WithMessage(err, "failed to initialise the terraform")
	}
	return tf.Apply()
}

func Destroy(captenConfig config.CaptenConfig, clusterType, cloudType string) error {
	clusterInfo, err := prepareClusterInfo(captenConfig, clusterType, cloudType)
	if err != nil {
		return err
	}

	tf, err := terraform.New(captenConfig, clusterInfo)
	if err != nil {
		return errors.WithMessage(err, "failed to initialise the terraform")
	}

	return tf.Destroy()
}

func prepareClusterInfo(captenConfig config.CaptenConfig, clusterType, cloudType string) (clusterInfo types.ClusterInfo, err error) {
	clusterInfo, err = config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, cloudType+"_config.yaml"))
	if err != nil {
		return clusterInfo, err
	}
	clusterInfo.ClusterType = clusterType
	clusterInfo.CloudService = cloudType
	return clusterInfo, nil
}
