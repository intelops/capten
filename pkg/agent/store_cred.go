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
	"capten/pkg/agent/vaultcredpb"
	"capten/pkg/config"
	"capten/pkg/k8s"
	"capten/pkg/types"

	"github.com/pkg/errors"
	"github.com/secure-systems-lab/go-securesystemslib/encrypted"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"gopkg.in/yaml.v2"
)

var (
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

	natsSecretNameVar       = "natsTokenSecretName"
	cosignKeysSecretNameVar = "cosignKeysSecretName"
)

func StoreCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}) error {
	vaultClient, err := GetVaultClient(captenConfig)
	if err != nil {
		return err
	}
	err = storeKubeConfig(captenConfig, vaultClient)
	if err != nil {
		return err
	}

	err = storeClusterGlobalValues(captenConfig, vaultClient)
	if err != nil {
		return err
	}

	err = storeNatsCredentials(captenConfig, appGlobalVaules, vaultClient)
	if err != nil {
		return err
	}

	err = storeCosignKeys(captenConfig, appGlobalVaules, vaultClient)
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

func storeKubeConfig(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {
	configContent, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return err
	}

	credentail := map[string]string{
		kubeconfigCredIdentifier: string(configContent),
	}

	_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: k8sCredEntityName,
		CredIdentifier: kubeconfigCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {

		return fmt.Errorf("store credentails failed, %s", err)

	}

	return nil
}

func storeClusterGlobalValues(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {
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

	_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: captenConfigEntityName,
		CredIdentifier: globalValuesCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return fmt.Errorf("store credentails failed, %s", err)

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

func storeNatsCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, vaultClient vaultcredpb.VaultCredClient) error {
	val, err := randomTokenGeneration()
	if err != nil {
		return fmt.Errorf("nats Token generation failed, %v", err)
	}

	config, err := readCredAppConfig(captenConfig, "nats-cred.yaml")
	if err != nil {
		return fmt.Errorf("error reading credential config YAML file: %v", err)
	}

	credentail := map[string]string{
		config.TokenAttributeName: val,
	}

	_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: config.CredentialEntity,
		CredIdentifier: config.CredentialIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return fmt.Errorf("store credentails failed, %s", err)

	}

	err = configireNatsSecret(captenConfig, vaultClient)
	if err != nil {
		return err
	}
	appGlobalVaules[natsSecretNameVar] = config.SecretName

	return nil
}

func configireNatsSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {

	config, err := readCredAppConfig(captenConfig, "nats-cred.yaml")
	if err != nil {
		return fmt.Errorf("Error reading Credential Configuration YAML file: %v", err)
	}

	natsTokenSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)

	for _, natsTokenNamespace := range config.Namespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, natsTokenNamespace)
		if err != nil {
			return err
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName: config.SecretName,
			Namespace:  natsTokenNamespace,

			SecretPathData: []*vaultcredpb.SecretPathRef{
				&vaultcredpb.SecretPathRef{SecretPath: natsTokenSecretPath, SecretKey: config.TokenAttributeName},
			},
			DomainName: captenConfig.DomainName,
		})
		if err != nil {
			return fmt.Errorf("failed to configure nats secret, %v", err)
		}

	}
	return nil
}

func storeCosignKeys(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, vaultClient vaultcredpb.VaultCredClient) error {
	privateKeyBytes, publicKeyBytes, err := generateCosignKeyPair()
	if err != nil {
		return fmt.Errorf("cosign key generation failed")
	}

	config, err := readCredAppConfig(captenConfig, "cosign-cred.yaml")
	if err != nil {
		return fmt.Errorf("Error reading Credential Config YAML file: %v", err)
	}

	credentail := map[string]string{
		"cosign.key": string(privateKeyBytes),
		"cosign.pub": string(publicKeyBytes),
	}

	_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: config.CredentialEntity,
		CredIdentifier: config.CredentialIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
	}

	err = configireCosignKeysSecret(captenConfig, vaultClient)
	if err != nil {
		return err
	}
	appGlobalVaules[cosignKeysSecretNameVar] = config.SecretName
	return nil
}

func configireCosignKeysSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {

	config, err := readCredAppConfig(captenConfig, "cosign-cred.yaml")
	if err != nil {
		return fmt.Errorf("Error while reading Credential Config YAML file: %v", err)
	}
	cosignKeysSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
	for _, cosignKeysNamespace := range config.Namespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, cosignKeysNamespace)
		if err != nil {
			return err
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName: config.SecretName,
			Namespace:  cosignKeysNamespace,
			SecretPathData: []*vaultcredpb.SecretPathRef{
				&vaultcredpb.SecretPathRef{SecretPath: cosignKeysSecretPath, SecretKey: "cosign.key"},
				&vaultcredpb.SecretPathRef{SecretPath: cosignKeysSecretPath, SecretKey: "cosign.pub"},
			},
			DomainName: captenConfig.DomainName,
		})
		if err != nil {
			return fmt.Errorf("failed to configure cosign keys secret, %v", err)
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

func readCredAppConfig(captenConfig config.CaptenConfig, filename string) (types.CredentialAppConfig, error) {
	var config types.CredentialAppConfig
	dirpath := captenConfig.PrepareDirPath(captenConfig.AppsConfigDirPath + captenConfig.AppsCredentialDirPath)

	filePath := dirpath + filename

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing YAML: %v", err)
	}

	return config, nil
}
