package agent

import (
	"context"
	"os"

	vaultcredclient "github.com/intelops/go-common/vault-cred-client"

	"capten/pkg/config"
)

const (
	k8sCredEntityName            string = "k8s"
	kubeconfigCredIdentifier     string = "kubeconfig"
	bucketCredEntityName         string = "bucket"
	terraformStateCredIdentifier string = "terraform-state"
)

func PushKubeConfigToVault(captenConfig config.CaptenConfig) error {

	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return err
	}

	credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = credAdmin.PutGenericCredential(ctx, k8sCredEntityName, kubeconfigCredIdentifier, vaultcredclient.GerericCredentail{
		Credential: map[string]string{
			"kubeconfig": string(configContent),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func PushBucketConfigToVault(captenConfig config.CaptenConfig, cloudType string) error {

	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, cloudType+"_config.yaml"))
	if err != nil {
		return err
	}

	credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = credAdmin.PutGenericCredential(ctx, bucketCredEntityName, terraformStateCredIdentifier, vaultcredclient.GerericCredentail{
		Credential: map[string]string{
			"bucketName": clusterInfo.TerraformBackendConfigs[0],
			"awsKey":     clusterInfo.AwsAccessKey,
			"awsSecrete": clusterInfo.AwsSecretKey,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
