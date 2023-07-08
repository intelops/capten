package agent

import (
	"os"

	"capten/pkg/config"

	"github.com/pkg/errors"
)

const (
	k8sCredEntityName            string = "k8s"
	kubeconfigCredIdentifier     string = "kubeconfig"
	s3BucketCredEntityName       string = "s3bucket"
	terraformStateCredIdentifier string = "terraform-state"

	terraformStateBucketNameKey string = "bucketName"
	terraformStateAwsAccessKey  string = "awsAccessKey"
	terraformStateAwsSecretKey  string = "awsSecretKey"
)

func StoreCredential(captenConfig config.CaptenConfig) error {
	err := StoreKubeConfig(captenConfig)
	if err != nil {
		return err
	}
	return StoreTerraformStateConfig(captenConfig)
}

func StoreKubeConfig(captenConfig config.CaptenConfig) error {
	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return err
	}

	_ = map[string]string{
		kubeconfigCredIdentifier: string(configContent),
	}

	// call agent to store cred
	return nil
}

func StoreTerraformStateConfig(captenConfig config.CaptenConfig) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
	if err != nil {
		return err
	}

	if len(clusterInfo.TerraformBackendConfigs) > 0 {
		return errors.New("Terraform backend configs are missing")
	}

	_ = map[string]string{
		terraformStateBucketNameKey: clusterInfo.TerraformBackendConfigs[0],
		terraformStateAwsAccessKey:  clusterInfo.AwsAccessKey,
		terraformStateAwsSecretKey:  clusterInfo.AwsSecretKey,
	}

	// call agent to store cred
	return nil
}
