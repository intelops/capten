package agent

import (
	"capten/pkg/agent/pb/vaultcredpb"
	"capten/pkg/clog"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	random "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	serviceCredentailType        string = "service-cred"
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

	postgresSecretNameVar   = "postgresSecretName"
	clickhouseSecretNameVar = "clickkhouseSecretName"
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

	log.Println("Cluster Global Values", appGlobalValues)

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
	secretKey := make(map[string]string)

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
				return fmt.Errorf("error while getting credential: %s", err)
			}
		} else {
			err = configureCosignKeysSecret(captenConfig, vaultClient, config)
			if err != nil {
				return err
			}
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
					return fmt.Errorf("nats Token generation failed, %v", err)
				}
				credential = map[string]string{
					"token": val,
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
				return fmt.Errorf("error while getting credential: %s", err)
			}
		} else {
			clog.Logger.Debug("Credential already exists in vault")
			err = configureNatsSecret(captenConfig, vaultClient, config)
			if err != nil {
				return fmt.Errorf("error while configuring cosign key: %v", err)
			}
		}
		appGlobalValues[natsSecretNameVar] = config.SecretName

	case "clickhouse-password":

		dbkey := map[string]string{
			"username": config.UserName,
		}
		err := generateAndStoreDBPassword(vaultClient, config, "password", dbkey)
		if err != nil {
			return fmt.Errorf("error while getting and storing password: %v", err)
		}

		secretKeyMapping := map[string]string{
			"username": "username",
			"password": "password",
		}

		err = configureDBSecret(captenConfig, vaultClient, config, secretKeyMapping)
		if err != nil {
			return fmt.Errorf("error while configuring secret: %v", err)
		}
		appGlobalValues[clickhouseSecretNameVar] = config.SecretName

		log.Println("Secret Name", appGlobalValues[clickhouseSecretNameVar])

	case "temporal-password":
		temporaldbuserkey := map[string]string{
			"username": config.UserName,
		}
		err := generateAndStoreDBPassword(vaultClient, config, "password", temporaldbuserkey)
		if err != nil {
			return fmt.Errorf("error while getting and storing password: %v", err)
		}
		postgresconfig := types.CredentialAppConfig{
			Name:                 "postgres-cred",
			SecretName:           "postgres-admin-secret",
			Namespaces:           []string{"observability", "platform", "capten", "quality-trace"},
			CredentialEntity:     "postgres",
			CredentialIdentifier: "postgres-admin",
			CredentialType:       "postgres-password",
			UserName:             "postgres",
		}

		posgresdbuserkey := map[string]string{
			"username": postgresconfig.UserName,
		}
		err = generateAndStoreDBPassword(vaultClient, postgresconfig, "admin-password", posgresdbuserkey)
		if err != nil {
			log.Println("Errror while gen admin password", err)
		}

		postgressecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, postgresconfig.CredentialEntity, postgresconfig.CredentialIdentifier)

		tempsecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
		secretKey = map[string]string{
			postgressecretPath: "admin-password",
			tempsecretPath:     "password",
		}

		log.Println("Secret Key in postgrespassword", secretKey)

		err = configureSecret(captenConfig, vaultClient, postgresconfig, secretKey)
		if err != nil {
			log.Printf("failed while configuring, %v", err)
			return err
		}
		appGlobalValues[postgresSecretNameVar] = postgresconfig.SecretName

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
		log.Println("Cred type", config.CredentialType)
		if config.CredentialType == "postgres-password" {

			for secretPath, secretkey := range secretKeyMapping {
				secretPathData = append(secretPathData, &vaultcredpb.SecretPathRef{SecretPath: secretPath, SecretKey: secretkey})
			}
			log.Println("Secret Path Data", secretPathData)
		} else {
			for _, vaultKey := range secretKeyMapping {
				secretPathData = append(secretPathData, &vaultcredpb.SecretPathRef{SecretPath: secretPath, SecretKey: vaultKey})
			}
			log.Println("Secret Path Data", secretPathData)
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName:     config.SecretName,
			Namespace:      namespace,
			SecretPathData: secretPathData,

			DomainName: "capten.svc.cluster.local:8200",
		})

		if err != nil {
			return fmt.Errorf("failed to configure secret in vault, %v", err)
		}
	}
	return nil
}

func configureCosignKeysSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {
	// secretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
	secretKeyMapping := map[string]string{
		"cosign.key": "cosign.key",
		"cosign.pub": "cosign.pub",
	}
	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping)
}

func configureNatsSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {

	secretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
	secretKeyMapping := map[string]string{
		secretPath: "token",
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

func configureDBSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, secretKeyMapping map[string]string) error {

	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping)
}

func generateAndStoreDBPassword(vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, passwordKey string, credential map[string]string) error {
	_, err := vaultClient.GetCredential(context.Background(), &vaultcredpb.GetCredentialRequest{
		CredentialType: serviceCredentailType,
		CredEntityName: config.CredentialEntity,
		CredIdentifier: config.CredentialIdentifier,
	})

	if err != nil {
		if strings.Contains(err.Error(), "secret not found") {
			val := generatePassword()

			credential[passwordKey] = val

			log.Println("Credential", credential)
			err := putCredentialInVault(vaultClient, config, credential)
			if err != nil {
				return fmt.Errorf("error storing credentials: %v", err)
			}

		} else {
			log.Printf("Error while getting credential: %s", err)
			return err
		}
	}

	return nil
}
