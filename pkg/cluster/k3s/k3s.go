package k3s

import (
	"os"
	"text/template"

	"capten/pkg/config"
	"capten/pkg/terraform"

	"github.com/pkg/errors"
)

func Create(captenConfig config.CaptenConfig) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
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

func Destroy(captenConfig config.CaptenConfig) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
	if err != nil {
		return err
	}

	tf, err := terraform.New(captenConfig, clusterInfo)
	if err != nil {
		return errors.WithMessage(err, "failed to initialise the terraform")
	}
	return tf.Destroy()
}
