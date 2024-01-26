package agent

import (
	"context"
	"fmt"
	"log"

	"os"

	"capten/pkg/agent/agentpb"
	"capten/pkg/config"

	"github.com/pkg/errors"
)

const (
	natsCredEntity     string = "nats"
	natsCredIdentifier string = "auth-token"

	genericCredentailType        string = "generic"
	k8sCredEntityName            string = "k8s"
	captenConfigEntityName       string = "capten-config"
	globalValuesCredIdentifier   string = "global-values"
	kubeconfigCredIdentifier     string = "kubeconfig"
	s3BucketCredEntityName       string = "s3bucket"
	terraformStateCredIdentifier string = "terraform-state"

	terraformStateBucketNameKey string = "bucketName"
	terraformStateAwsAccessKey  string = "awsAccessKey"
	terraformStateAwsSecretKey  string = "awsSecretKey"
	externalsecretns            string = "external-secrets"
	externalsecretname          string = "nats-token"
	natsTokenSecretPath         string = "generic/nats/auth-token"
	vaultTokenSecretName        string = "vaultToken"
	clusterId                   string = "1"
)

func StoreCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}) error {
	agentClient, err := GetAgentClient(captenConfig)
	if err != nil {
		return err
	}

	err = storeKubeConfig(captenConfig, agentClient)
	if err != nil {
		return err
	}

	err = storeClusterGlobalValues(captenConfig, agentClient)
	if err != nil {
		return err
	}

	err = storeNatsCredentials(captenConfig, appGlobalVaules, agentClient)
	if err != nil {
		return err
	}
	err = RetrieveNatsCredential(captenConfig, appGlobalVaules, agentClient)

	if err != nil {
		log.Println("Error while retrieving nats credential", err)
		return err
	}

	return nil
}

func StoreClusterCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}) error {
	agentClient, err := GetAgentClient(captenConfig)
	if err != nil {
		return err
	}
	err = storeTerraformStateConfig(captenConfig, agentClient)
	if err != nil {
		return err
	}

	return nil
}

func storeKubeConfig(captenConfig config.CaptenConfig, agentClient agentpb.AgentClient) error {
	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return err
	}

	credentail := map[string]string{
		kubeconfigCredIdentifier: string(configContent),
	}

	response, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: k8sCredEntityName,
		CredIdentifier: kubeconfigCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
	}

	if response.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("store credentails failed, %s", response.StatusMessage)
	}
	return nil
}

func storeClusterGlobalValues(captenConfig config.CaptenConfig, agentClient agentpb.AgentClient) error {
	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CaptenGlobalValuesFileName))
	if err != nil {
		return err
	}
	hostValues, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CaptenHostValuesFileName))
	if err != nil {
		return err
	}

	credentail := map[string]string{
		globalValuesCredIdentifier: string(configContent) + "\n" + string(hostValues),
	}

	response, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: captenConfigEntityName,
		CredIdentifier: globalValuesCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
	}

	if response.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("store credentails failed, %s", response.StatusMessage)
	}
	return nil
}

func storeTerraformStateConfig(captenConfig config.CaptenConfig, agentClient agentpb.AgentClient) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
	if err != nil {
		return err
	}

	if len(clusterInfo.TerraformBackendConfigs) > 0 {
		return errors.New("Terraform backend configs are missing")
	}

	credentail := map[string]string{
		terraformStateBucketNameKey: clusterInfo.TerraformBackendConfigs[0],
		terraformStateAwsAccessKey:  clusterInfo.AwsAccessKey,
		terraformStateAwsSecretKey:  clusterInfo.AwsSecretKey,
	}

	response, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: s3BucketCredEntityName,
		CredIdentifier: terraformStateCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
	}

	if response.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("store credentails failed, %s", response.StatusMessage)
	}
	return nil
}

func storeNatsCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, agentClient agentpb.AgentClient) error {
	val, ok := appGlobalVaules["NatsToken"]
	if !ok {
		return fmt.Errorf("NatsToken is missing")
	}
	credentail := map[string]string{
		natsCredEntity: val.(string),
	}

	response, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: natsCredEntity,
		CredIdentifier: natsCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
	}

	if response.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("store credentails failed, %s", response.StatusMessage)
	}
	return nil
}

func RetrieveNatsCredential(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, agentClient agentpb.AgentClient) error {
	resp, err := agentClient.ConfigureVaultSecret(context.Background(), &agentpb.ConfigureVaultSecretRequest{
		SecretName: externalsecretname,
		Namespace:  externalsecretns,
		SecretPathData: []*agentpb.SecretPathRef{
			&agentpb.SecretPathRef{SecretPath: natsTokenSecretPath, SecretKey: "nats"},
		},
	})
	if err != nil {
		fmt.Errorf("Unable to configure vault secret, %s", err)
		return err
	}
	if resp.Status != agentpb.StatusCode_OK {
		return fmt.Errorf("retrieve credentails failed, %s", resp.StatusMessage)

	}
	fmt.Println("Retrieved the nats token successfully")

	return nil
}
