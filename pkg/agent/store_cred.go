package agent

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"capten/pkg/agent/agentpb"
	"capten/pkg/config"

	"github.com/pkg/errors"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"github.com/theupdateframework/go-tuf/encrypted"
)

const (
	natsCredEntity       string = "nats"
	natsCredIdentifier   string = "auth-token"
	cosignEntity         string = "cosign"
	cosignCredIdentifier string = "signer"

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

	err = storeCosignKeys(agentClient)
	if err != nil {
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

func storeCosignKeys(agentClient agentpb.AgentClient) error {
	privateKeyBytes, publicKeyBytes, err := generateCosignKeyPair()
	if err != nil {
		return fmt.Errorf("Cosign key generation failed")
	}
	credentail := map[string]string{
		"cosign.key": string(privateKeyBytes),
		"cosign.pub": string(publicKeyBytes),
	}

	response, err := agentClient.StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: cosignEntity,
		CredIdentifier: cosignCredIdentifier,
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

func generateCosignKeyPair() ([]byte, []byte, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	keypair := struct {
		private crypto.PrivateKey
		public  crypto.PublicKey
	}{priv, priv.Public()}

	x509Encoded, err := x509.MarshalPKCS8PrivateKey(keypair.private)
	if err != nil {
		return nil, nil, fmt.Errorf("x509 encoding private key: %w", err)
	}

	encBytes, err := encrypted.Encrypt(x509Encoded, []byte{})
	if err != nil {
		return nil, nil, err
	}

	privBytes := pem.EncodeToMemory(&pem.Block{
		Bytes: encBytes,
		Type:  "ENCRYPTED COSIGN PRIVATE KEY",
	})

	pubBytes, err := cryptoutils.MarshalPublicKeyToPEM(keypair.public)
	if err != nil {
		return nil, nil, err
	}

	return privBytes, pubBytes, nil
}
