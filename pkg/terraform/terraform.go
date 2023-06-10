package terraform

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/viper"

	"capten/pkg/cluster/types"
)

type terraform struct {
	config     *viper.Viper
	exec       *tfexec.Terraform
	workingDir string
}

func New(config *viper.Viper, workingDir string) (*terraform, error) {
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.3.7")),
		InstallDir: "./",
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Printf("error installing Terraform: %s", err)
		return nil, err
	}

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Printf("error running NewTerraform: %s", err)
		return nil, err
	}

	//set the output files, defaulted to terminal
	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	return &terraform{config: config, exec: tf, workingDir: workingDir}, nil
}

func (t *terraform) Apply() error {
	backendConfigs := t.config.GetStringSlice(types.TerraformBackendConfigs)
	backendConfigOptionsStr := []string{
		"region=" + t.config.GetString(types.Region),
		"access_key=" + t.config.GetString(types.AwsAccessKey),
		"secret_key=" + t.config.GetString(types.AwsSecretKey),
	}
	backendConfigOptionsStr = append(backendConfigOptionsStr, backendConfigs...)
	fmt.Println(backendConfigOptionsStr)
	initOptions := make([]tfexec.InitOption, 0)
	for _, backendConfigOption := range backendConfigOptionsStr {
		initOptions = append(initOptions, tfexec.BackendConfig(backendConfigOption))
	}

	initOptions = append(initOptions, tfexec.Upgrade(true))
	err := t.exec.Init(context.Background(), initOptions...)
	if err != nil {
		log.Printf("error running Init: %s", err)
		return err
	}

	_, err = t.exec.Show(context.Background())
	if err != nil {
		log.Printf("error running show: %s", err)
		return err
	}

	varFile := fmt.Sprintf("%s/%s", t.workingDir, "values.tfvars")
	_, err = t.exec.Plan(context.Background(), tfexec.VarFile(varFile))
	if err != nil {
		log.Printf("error running plan: %s", err)
		return err
	}

	log.Println("terraform plan is completed")
	if err := t.exec.Apply(context.Background(), tfexec.VarFile(varFile)); err != nil {
		log.Printf("error running apply: %s", err)
		return err
	}

	return nil
}

func (t *terraform) Destroy() error {
	varFile := fmt.Sprintf("%s/%s", t.workingDir, "values.tfvars")
	return t.exec.Destroy(context.Background(), tfexec.VarFile(varFile))
}
