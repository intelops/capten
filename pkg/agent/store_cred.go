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
	qtSecretNameVar         = "qtSecretName"
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
				err = putCredentialInVault(vaultClient, config, credential, genericCredentailType)
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
				err = putCredentialInVault(vaultClient, config, credential, genericCredentailType)
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

		secretKeyMapping := map[string][]string{
			"username": {"username"},
			"password": {"password"},
		}

		err = configureSecret(captenConfig, vaultClient, config, secretKeyMapping, nil, serviceCredentailType)
		if err != nil {
			return fmt.Errorf("error while configuring secret: %v", err)
		}
		appGlobalValues[clickhouseSecretNameVar] = config.SecretName
		log.Println("Clickhouse secretName ", appGlobalValues[clickhouseSecretNameVar])

	case "qt-password":

		dbkey := map[string]string{
			"username": config.UserName,
		}
		err := generateAndStoreDBPassword(vaultClient, config, "password", dbkey)
		if err != nil {
			return fmt.Errorf("error while getting and storing password: %v", err)
		}

		secretKeyMapping := map[string][]string{
			"username": {"username"},
			"password": {"password"},
		}

		err = configureSecret(captenConfig, vaultClient, config, secretKeyMapping, nil, serviceCredentailType)
		if err != nil {
			return fmt.Errorf("error while configuring secret: %v", err)
		}
		appGlobalValues[qtSecretNameVar] = config.SecretName
		log.Println("Secret Name", appGlobalValues[qtSecretNameVar])

	case "temporal-password":
		temporaldbuserkey := map[string]string{
			"username": config.UserName,
		}
		err := generateAndStoreDBPassword(vaultClient, config, "password", temporaldbuserkey)
		if err != nil {
			return fmt.Errorf("error while getting and storing password: %v", err)
		}
		postgresconfig := types.CredentialAppConfig{
			Name:       "postgres-cred",
			SecretName: "postgres-admin-secret",
			Namespaces: []string{"observability", "platform", "capten", "quality-trace"}, //	"platform", "capten", "quality-trace"

			CredentialEntity:     "postgres",
			CredentialIdentifier: "postgres-admin",
			CredentialType:       "postgres-password",
			UserName:             "postgres",
		}

		posgresdbuserkey := map[string]string{
			"username": postgresconfig.UserName,
		}
		err = generateAndStoreDBPassword(vaultClient, postgresconfig, "password", posgresdbuserkey)
		if err != nil {
			return fmt.Errorf("error while generating and storing secret in vault: %v", err)

		}

		postgressecretPath := fmt.Sprintf("%s/%s/%s", serviceCredentailType, postgresconfig.CredentialEntity, postgresconfig.CredentialIdentifier)

		tempsecretPath := fmt.Sprintf("%s/%s/%s", serviceCredentailType, config.CredentialEntity, config.CredentialIdentifier)

		secretKey := map[string][]string{
			postgressecretPath: {"password"},
			tempsecretPath:     {"password"},
		}

		secretPropertiesMapping := map[string][]string{
			"password": {"admin-password", "password"},
		}

		err = configureSecret(captenConfig, vaultClient, postgresconfig, secretKey, secretPropertiesMapping, serviceCredentailType)
		if err != nil {
			return err
		}
		appGlobalValues[postgresSecretNameVar] = postgresconfig.SecretName

	default:
		return fmt.Errorf("unknown credential type: %s", config.CredentialType)
	}

	return nil
}

func putCredentialInVault(vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, credential map[string]string, credentialType string) error {
	_, err := vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
		CredentialType: credentialType,
		CredEntityName: config.CredentialEntity,
		CredIdentifier: config.CredentialIdentifier,
		Credential:     credential,
	})
	return err
}

func configureCosignKeysSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {

	secretKeyMapping := map[string][]string{
		"cosign.key": {"cosign.key"},
		"cosign.pub": {"cosign.pub"},
	}

	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping, nil, genericCredentailType)
}

func configureNatsSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig) error {

	secretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, config.CredentialEntity, config.CredentialIdentifier)
	secretKeyMapping := map[string][]string{
		secretPath: {"token"},
	}

	return configureSecret(captenConfig, vaultClient, config, secretKeyMapping, nil, genericCredentailType)
}

func generatePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	source := random.NewSource(time.Now().UnixNano())
	rng := random.New(source)

	password := make([]byte, 11)
	for i := range password {
		password[i] = charset[rng.Intn(len(charset))]
	}
	return string(password)
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

			err := putCredentialInVault(vaultClient, config, credential, serviceCredentailType)
			if err != nil {
				return fmt.Errorf("error storing credentials: %v", err)
			}

		} else {
			return err
		}
	}

	return nil
}

func configureSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient, config types.CredentialAppConfig, secretPathMapping map[string][]string, propertyMapping map[string][]string, credentialType string) error {
	secretPath := fmt.Sprintf("%s/%s/%s", credentialType, config.CredentialEntity, config.CredentialIdentifier)

	for _, namespace := range config.Namespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, namespace, nil)
		if err != nil {
			return err
		}

		if propertyMapping == nil {
			propertyMapping = make(map[string][]string)
			for k, v := range secretPathMapping {
				propertyMapping[k] = v
			}
		}

		secretPathData := make([]*vaultcredpb.SecretPathRef, 0)

		propertyIndex := make(map[string]int)
		if config.CredentialType == "postgres-password" {
			for secretpath, properties := range secretPathMapping {
				for _, property := range properties {
					secretKey := property

					if props, exists := propertyMapping[property]; exists && len(props) > 0 {
						idx := propertyIndex[property]

						if idx < len(props) {
							secretKey = props[idx]
						}

						propertyIndex[property]++
					}

					secretPathData = append(secretPathData, &vaultcredpb.SecretPathRef{
						SecretPath: secretpath,
						SecretKey:  secretKey,
						Property:   property,
					})
				}
			}
		} else {
			for _, properties := range secretPathMapping {
				for _, property := range properties {
					secretKey := property

					if props, exists := propertyMapping[property]; exists && len(props) > 0 {
						idx := propertyIndex[property]

						if idx < len(props) {
							secretKey = props[idx]
						}

						propertyIndex[property]++
					}

					secretPathData = append(secretPathData, &vaultcredpb.SecretPathRef{
						SecretPath: secretPath,
						SecretKey:  secretKey,
						Property:   property,
					})
				}
			}
		}

		request := &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName:     config.SecretName,
			Namespace:      namespace,
			SecretPathData: secretPathData,
			DomainName:     "capten.svc.cluster.local:8200",
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), request)
		if err != nil {
			return fmt.Errorf("failed to configure vault secret: %v", err)
		}
	}
	return nil
}
