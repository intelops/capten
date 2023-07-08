package agent

import (
	"context"
	"os"

	vaultcredclient "github.com/intelops/go-common/vault-cred-client"

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

func StoreKubeConfig(captenConfig config.CaptenConfig) error {
	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return err
	}

	credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
	if err != nil {
		return err
	}

	credential := map[string]string{
		kubeconfigCredIdentifier: string(configContent),
	}

	return credAdmin.PutCredential(context.Background(), k8sCredEntityName, kubeconfigCredIdentifier, credential)
}

func StoreTerraformStateConfig(captenConfig config.CaptenConfig, cloudType string) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, cloudType+"_config.yaml"))
	if err != nil {
		return err
	}

	credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
	if err != nil {
		return err
	}

	if len(clusterInfo.TerraformBackendConfigs) > 0 {
		return errors.New("Terraform backend configs are missing")
	}

	credential := map[string]string{
		terraformStateBucketNameKey: clusterInfo.TerraformBackendConfigs[0],
		terraformStateAwsAccessKey:  clusterInfo.AwsAccessKey,
		terraformStateAwsSecretKey:  clusterInfo.AwsSecretKey,
	}

	return credAdmin.PutCredential(context.Background(), s3BucketCredEntityName, terraformStateCredIdentifier, credential)
}
