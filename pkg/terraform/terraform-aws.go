package terraform

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/types"
)

type terraform struct {
	azureconfig  types.AzureClusterInfo
	config       types.AWSClusterInfo
	exec         *tfexec.Terraform
	captenConfig config.CaptenConfig
}

func NewAws(captenConfig config.CaptenConfig, config types.AWSClusterInfo) (*terraform, error) {
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.3.7")),
		InstallDir: "./",
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		clog.Logger.Infof("execPath: %s", execPath)

		return nil, errors.WithMessage(err, "error installing Terraform")
	}

	workDir := captenConfig.PrepareDirPath(captenConfig.TerraformModulesDirPath + captenConfig.CloudService + "/" + captenConfig.ClusterType)
	clog.Logger.Debugf("terraform workingDir: %s, execPath: %s", workDir, execPath)
	tf, err := tfexec.NewTerraform(workDir, execPath)
	if err != nil {
		return nil, errors.WithMessage(err, "error running NewTerraform")
	}

	tf.SetLogger(clog.Logger)
	//set the output files, defaulted to terminal
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	return &terraform{config: config, exec: tf, captenConfig: captenConfig}, nil
}

func (t *terraform) initAws() error {

	backendConfigOptionsStr := []string{
		"region=" + t.config.Region,
		"access_key=" + t.config.AwsAccessKey,
		"secret_key=" + t.config.AwsSecretKey,
	}
	backendConfigOptionsStr = append(backendConfigOptionsStr, t.config.TerraformBackendConfigs...)
	initOptions := make([]tfexec.InitOption, 0)
	for _, backendConfigOption := range backendConfigOptionsStr {
		initOptions = append(initOptions, tfexec.BackendConfig(backendConfigOption))
	}

	initOptions = append(initOptions, tfexec.Upgrade(t.captenConfig.TerraformInitUpgrade))
	initOptions = append(initOptions, tfexec.Reconfigure(t.captenConfig.TerraformInitReconfigure))
	err := t.exec.Init(context.Background(), initOptions...)
	if err != nil {
		return errors.WithMessage(err, "terraform init failed")
	}
	return nil
}

func (t *terraform) initCommon() error {
	if t.captenConfig.CloudService == "azure" {
		err := t.initAzure()
		if err != nil {
			return err
		}
	} else {
		err := t.initAws()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *terraform) Apply() error {
	if err := t.initCommon(); err != nil {
		return err
	}

	_, err := t.exec.Show(context.Background())
	if err != nil {
		return errors.WithMessage(err, "error running show")
	}

	varFile := fmt.Sprintf("%s%s%s", t.captenConfig.CurrentDirPath, t.captenConfig.TerraformTemplateDirPath, t.captenConfig.TerraformVarFileName)
	_, err = t.exec.Plan(context.Background(), tfexec.VarFile(varFile))
	if err != nil {
		return errors.WithMessage(err, "error running plan")
	}

	if err := t.exec.Apply(context.Background(), tfexec.VarFile(varFile)); err != nil {
		return errors.WithMessage(err, "error running apply")
	}
	return nil
}

func (t *terraform) Destroy() error {
	if err := t.initCommon(); err != nil {
		return err
	}

	_, err := t.exec.Show(context.Background())
	if err != nil {
		return errors.WithMessage(err, "error running show")
	}

	varFile := fmt.Sprintf("%s%s%s", t.captenConfig.CurrentDirPath, t.captenConfig.TerraformTemplateDirPath, t.captenConfig.TerraformVarFileName)
	return t.exec.Destroy(context.Background(), tfexec.VarFile(varFile))
}

