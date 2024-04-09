package terraform

import (
	"context"

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

func NewAzure(captenConfig config.CaptenConfig, config types.AzureClusterInfo) (*terraform, error) {
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
	return &terraform{azureconfig: config, exec: tf, captenConfig: captenConfig}, nil
}

func (t *terraform) initAzure() error {
	backendConfigOptionsStr := []string{}
	//	backendConfigOptionsStr = append(backendConfigOptionsStr)

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
