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
	random "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"capten/pkg/agent/vaultcredpb"
	"capten/pkg/clog"

	//"capten/pkg/clog"
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

	natsTokenNamespaces   []string = []string{"observability"}
	cosignKeysNamespaces  []string = []string{"kyverno", "tekton-pipelines", "tek"}
	postgresSecretNameVar          = "postgresSecretName"
)

func StoreCredentials(captenConfig config.CaptenConfig, appGlobalValues map[string]interface{}) error {
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
	err = StoreCredAppConfig(captenConfig, appGlobalValues, vaultClient)
	if err != nil {
		return err
	}

	return nil
}

func StoreClusterCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}) error {

	vaultClient, err := GetVaultClient(captenConfig)
	if err != nil {
		return err
	}
	err = storeTerraformStateConfig(captenConfig, vaultClient)
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

func storeTerraformStateConfig(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {
	clusterInfo, err := config.GetClusterInfo(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CloudService+"_config.yaml"))
	if err != nil {
		return err
	}

	credentail := map[string]string{
		terraformStateBucketNameKey: clusterInfo.TerraformBackendConfigs[0],
		terraformStateAwsAccessKey:  clusterInfo.AwsAccessKey,
		terraformStateAwsSecretKey:  clusterInfo.AwsSecretKey,
	}

	_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: s3BucketCredEntityName,
		CredIdentifier: terraformStateCredIdentifier,
		Credential:     credentail,
	})
	if err != nil {
		return err
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

	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)
	randomString = strings.ReplaceAll(randomString, "-", "")

	if len(randomString) > 32 {
		randomString = randomString[:32]
	}

	return randomString, nil
}

func StoreCredAppConfig(captenConfig config.CaptenConfig, appGlobalValues map[string]interface{}, vaultClient vaultcredpb.VaultCredClient) error {
	var credConfigs types.CredentialAppConfig
	dirpath := captenConfig.PrepareDirPath(captenConfig.AppsConfigDirPath + captenConfig.AppsCredentialDirPath)

	files, err := os.ReadDir(dirpath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {

		filePath := filepath.Join(dirpath, file.Name())
		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file %s: %v", filePath, err)
		}

		err = yaml.Unmarshal(yamlFile, &credConfigs)
		if err != nil {
			return fmt.Errorf("error parsing YAML file %s: %v", filePath, err)
		}
		err = storeCredentials(captenConfig, appGlobalValues, vaultClient, credConfigs)
		if err != nil {
			return fmt.Errorf("error while storing app credentials  %s: %v", filePath, err)
		}

	}
	return nil
}

func storeCredentials(captenConfig config.CaptenConfig, appGlobalValues map[string]interface{}, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {

	var credential map[string]string
	switch config.CredentialType {
	case "cosign":

		_, err := vaultClient.GetCredential(context.Background(), &vaultcredpb.GetCredentialRequest{
			CredentialType: genericCredentailType,
			CredEntityName: config.CredentialEntity,
			CredIdentifier: config.CredentialIdentifier,
		})
		if err != nil {
			if strings.Contains(err.Error(), "secret not found") {
				privateKeyBytes, publicKeyBytes, err := generateCosignKeyPair()
				if err != nil {
					return fmt.Errorf("cosign key generation failed")
				}
				credential := map[string]string{
					"cosign.key": string(privateKeyBytes),
					"cosign.pub": string(publicKeyBytes),
				}
				err = putCredentialInVault(vaultClient, config, credential)
				if err != nil {
					return fmt.Errorf("error storing credentials: %v", err)
				}

				err = configureCosignKeysSecret(captenConfig, vaultClient, config)
				if err != nil {
					return err
				}

			} else {

				return fmt.Errorf("Error while getting credential: %s", err)
			}
		} else {

			clog.Logger.Debug("Credential already exists in vault")
		}

		appGlobalValues[cosignKeysSecretNameVar] = config.SecretName

	case "randomkey":

		_, err := vaultClient.GetCredential(context.Background(), &vaultcredpb.GetCredentialRequest{
			CredentialType: genericCredentailType,
			CredEntityName: config.CredentialEntity,
			CredIdentifier: config.CredentialIdentifier,
		})

		if err != nil {
			if strings.Contains(err.Error(), "secret not found") {
				val, err := randomTokenGeneration()
				if err != nil {
					return fmt.Errorf("Nats Token generation failed, %v", err)
				}
				credential = map[string]string{
					config.TokenAttributeName: val,
				}
				err = putCredentialInVault(vaultClient, config, credential)

				if err != nil {
					return fmt.Errorf("store credentails failed, %s", err)

				}
				err = configureNatsSecret(captenConfig, vaultClient, config)
				if err != nil {
					return err
				}

			} else {

				return fmt.Errorf("Error while getting credential: %s", err)
			}
		} else {

			clog.Logger.Debug("Credential already exists in vault")
		}

		err = configureNatsSecret(captenConfig, vaultClient, config)
		if err != nil {
			return fmt.Errorf("error while configuring cosign key: %v", err)
		}
		appGlobalValues[natsSecretNameVar] = config.SecretName

	case "postgres-password", "temporal-password":
		err := generateAndStorePassword(vaultClient, config)
		if err != nil {
			return fmt.Errorf("error while getting and storing password: %v", err)
		}

	// err = configureSecret(captenConfig, vaultClient, config, secretKeyMapping)
	// if err != nil {
	// 	return fmt.Errorf("error while configuring secret: %v", err)
	// }
	//appGlobalValues[postgresSecretNameVar] = config.SecretName

	default:

		return fmt.Errorf("unknown credential type: %s", config.CredentialType)
	}

	return nil
}

func putCredentialInVault(vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, credential map[string]string) error {

	_, err := vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: config.CredentialEntity,
		CredIdentifier: config.CredentialIdentifier,
		Credential:     credential,
	})
	return err
}

func configureSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, secretKeyMapping map[string]string) error {
	secretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
	for _, namespace := range config.Namespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, namespace, nil)
		if err != nil {
			return err
		}

		secretPathData := make([]*vaultcredpb.SecretPathRef, 0, len(secretKeyMapping))
		for _, vaultKey := range secretKeyMapping {
			secretPathData = append(secretPathData, &vaultcredpb.SecretPathRef{SecretPath: secretPath, SecretKey: vaultKey})
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName:     config.SecretName,
			Namespace:      namespace,
			SecretPathData: secretPathData,
			DomainName:     captenConfig.DomainName,
		})
		if err != nil {
			return fmt.Errorf("failed to configure secret in vault, %v", err)
		}
	}
	return nil
}

func configureCosignKeysSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {
	secretKeyMapping := map[string]string{
		"cosign.key": "cosign.key",
		"cosign.pub": "cosign.pub",
	}

	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping)
}

func configureNatsSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {
	secretKeyMapping := map[string]string{
		config.TokenAttributeName: config.TokenAttributeName,
	}

	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping)
}

func generatePassword() string {
	// Define the characters to choose from for the password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a new source for random number generation
	source := random.NewSource(time.Now().UnixNano())
	rng := random.New(source)

	password := make([]byte, 11)
	for i := range password {
		password[i] = charset[rng.Intn(len(charset))]
	}
	return string(password)
}

func generateAndStorePassword(vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {
	val := generatePassword()

	credential := map[string]string{
		"password": val,
	}
	err := putCredentialInVault(vaultClient, config, credential)
	if err != nil {
		return fmt.Errorf("error storing credentials: %v", err)
	}

	return nil
}
