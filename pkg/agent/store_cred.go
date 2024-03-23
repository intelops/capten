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
	"strings"

	"os"

	"capten/pkg/agent/vaultcredpb"
	"capten/pkg/clog"
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

func storeNatsCredentials(captenConfig config.CaptenConfig, appGlobalVaules map[string]interface{}, vaultClient vaultcredpb.VaultCredClient) error {

	_, err := vaultClient.GetCredential(context.Background(), &vaultcredpb.GetCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: natsCredEntity,
		CredIdentifier: natsCredIdentifier,
	})

	if err != nil {
		if strings.Contains(err.Error(), "secret not found") {
			val, err := randomTokenGeneration()
			if err != nil {
				return fmt.Errorf("Nats Token generation failed, %v", err)
			}
			credentail := map[string]string{
				tokenAttributeName: val,
			}
			_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
				CredentialType: genericCredentailType,
				CredEntityName: natsCredEntity,
				CredIdentifier: natsCredIdentifier,
				Credential:     credentail,
			})
			if err != nil {
				return fmt.Errorf("store credentails failed, %s", err)

			}
			err = configireNatsSecret(captenConfig, vaultClient)
			if err != nil {
				return err
			}

		} else {

			return fmt.Errorf("Error while getting credential: %s", err)
		}
	} else {

		clog.Logger.Info("Credential already exists in vault")
	}

	appGlobalVaules[natsSecretNameVar] = natsTokenSecretName
	return nil
}

func configireNatsSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {
	natsTokenSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, natsCredEntity, natsCredIdentifier)
	for _, natsTokenNamespace := range natsTokenNamespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, natsTokenNamespace, nil)
		if err != nil {
			return err
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName: natsTokenSecretName,
			Namespace:  natsTokenNamespace,
			SecretPathData: []*vaultcredpb.SecretPathRef{
				&vaultcredpb.SecretPathRef{SecretPath: natsTokenSecretPath, SecretKey: tokenAttributeName},
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
	_, err := vaultClient.GetCredential(context.Background(), &vaultcredpb.GetCredentialRequest{
		CredentialType: genericCredentailType,
		CredEntityName: natsCredEntity,
		CredIdentifier: natsCredIdentifier,
	})
	if err != nil {
		if strings.Contains(err.Error(), "secret not found") {
			privateKeyBytes, publicKeyBytes, err := generateCosignKeyPair()
			if err != nil {
				return fmt.Errorf("cosign key generation failed")
			}
			credentail := map[string]string{
				"cosign.key": string(privateKeyBytes),
				"cosign.pub": string(publicKeyBytes),
			}

			_, err = vaultClient.PutCredential(context.Background(), &vaultcredpb.PutCredentialRequest{
				CredentialType: genericCredentailType,
				CredEntityName: cosignEntity,
				CredIdentifier: cosignCredIdentifier,
				Credential:     credentail,
			})
			if err != nil {
				return err
			}

			err = configireCosignKeysSecret(captenConfig, vaultClient)
			if err != nil {
				return err
			}

		} else {
		
			return fmt.Errorf("Error while getting credential: %s", err)
		}
	} else {

		clog.Logger.Info("Credential already exists in vault")
	}

	appGlobalVaules[cosignKeysSecretNameVar] = cosignKeysSecretName
	return nil
}

func configireCosignKeysSecret(captenConfig config.CaptenConfig, vaultClient vaultcredpb.VaultCredClient) error {
	cosignKeysSecretPath := fmt.Sprintf("%s/%s/%s", genericCredentailType, cosignEntity, cosignCredIdentifier)
	for _, cosignKeysNamespace := range cosignKeysNamespaces {
		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err := k8s.CreateNamespaceIfNotExist(kubeconfigPath, cosignKeysNamespace, nil)
		if err != nil {
			return err
		}

		_, err = vaultClient.ConfigureVaultSecret(context.Background(), &vaultcredpb.ConfigureVaultSecretRequest{
			SecretName: cosignKeysSecretName,
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
