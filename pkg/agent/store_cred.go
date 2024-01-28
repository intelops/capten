package agent

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"capten/pkg/agent/agentpb"
	"capten/pkg/config"
	"capten/pkg/k8s"

	"github.com/pkg/errors"
	"github.com/secure-systems-lab/go-securesystemslib/encrypted"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
)

var (
	tokenAttributeName   string = "token"
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

	natsTokenSecretName     = "nats-token"
	cosignKeysSecretName    = "cosign-keys"
	natsSecretNameVar       = "natsTokenSecretName"
	cosignKeysSecretNameVar = "cosignKeysSecretName"

	natsTokenNamespaces  []string = []string{"observability"}
	cosignKeysNamespaces []string = []string{"kyverno", "tekton-pipelines", "tek"}
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

	err = storeCosignKeys(captenConfig, appGlobalVaules, agentClient)
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
	val, err := randomTokenGeneration()
	if err != nil {
		return fmt.Errorf("Nats Token generation failed, %v", err)
	}
	credentail := map[string]string{
		tokenAttributeName: val,
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

	err = configireNatsSecret(captenConfig, agentClient)
	if err != nil {
		return err
	}
	appGlobalVaules[natsSecretNameVar] = natsTokenSecretName
	return nil
}

func configireNatsSecret(captenConfig config.CaptenConfig, agentClient agentpb.AgentClient) error {
	natsTokenSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, natsCredEntity, natsCredIdentifier)
	for _, natsTokenNamespace := range natsTokenNamespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, natsTokenNamespace)
		if err != nil {
			return err
		}

		resp, err := agentClient.ConfigureVaultSecret(context.Background(), &agentpb.ConfigureVaultSecretRequest{
			SecretName: natsTokenSecretName,
			Namespace:  natsTokenNamespace,
			SecretPathData: []*agentpb.SecretPathRef{
				&agentpb.SecretPathRef{SecretPath: natsTokenSecretPath, SecretKey: tokenAttributeName},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to configure nats secret, %v", err)
		}
		if resp.Status != agentpb.StatusCode_OK {
			return fmt.Errorf("failed to configure nats secret, %s", resp.StatusMessage)

		}
	}
	return nil
}

func storeCosignKeys(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, agentClient agentpb.AgentClient) error {
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

	err = configireCosignKeysSecret(captenConfig, agentClient)
	if err != nil {
		return err
	}
	appGlobalVaules[cosignKeysSecretNameVar] = cosignKeysSecretName
	return nil
}

func configireCosignKeysSecret(captenConfig config.CaptenConfig, agentClient agentpb.AgentClient) error {
	cosignKeysSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, cosignEntity, cosignCredIdentifier)
	for _, cosignKeysNamespace := range cosignKeysNamespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, cosignKeysNamespace)
		if err != nil {
			return err
		}

		resp, err := agentClient.ConfigureVaultSecret(context.Background(), &agentpb.ConfigureVaultSecretRequest{
			SecretName: cosignKeysSecretName,
			Namespace:  cosignKeysNamespace,
			SecretPathData: []*agentpb.SecretPathRef{
				&agentpb.SecretPathRef{SecretPath: cosignKeysSecretPath, SecretKey: "cosign.key"},
				&agentpb.SecretPathRef{SecretPath: cosignKeysSecretPath, SecretKey: "cosign.pub"},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to configure cosign keys secret, %v", err)
		}
		if resp.Status != agentpb.StatusCode_OK {
			return fmt.Errorf("failed to configure cosign keys secret, %s", resp.StatusMessage)

		}
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

func randomTokenGeneration() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.WithMessage(err, "error while generating random key")
	}
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)[:32]
	return randomString, nil
}
