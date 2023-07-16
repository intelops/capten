package k8s

import (
	"capten/pkg/config"
	"context"
	"io/ioutil"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func CreateOrUpdateCertSecrets(captenConfig config.CaptenConfig) error {
	kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	clientSet, err := getK8SClient(kubeconfigPath)
	if err != nil {
		return err
	}
	err = createOrUpdateAgentCertSecret(captenConfig, clientSet)
	if err != nil {
		return err
	}
	return createOrUpdateAgentCACert(captenConfig, clientSet)
}

func createOrUpdateAgentCertSecret(captenConfig config.CaptenConfig, k8sClient *kubernetes.Clientset) error {
	certData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentCertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client cert")
	}
	keyData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentKeyFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client key")
	}
	caCertChainData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading ca cert chain")
	}

	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      captenConfig.AgentCertSecretName,
			Namespace: captenConfig.CaptenNamespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       certData,
			corev1.TLSPrivateKeyKey: keyData,
			"ca.crt":                caCertChainData,
		},
		Type: corev1.SecretTypeTLS,
	}
	return createOrUpdateSecret(k8sClient, secret)
}

func createOrUpdateAgentCACert(captenConfig config.CaptenConfig, k8sClient *kubernetes.Clientset) error {
	caCertChainData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading ca cert chain")
	}

	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      captenConfig.AgentCACertSecretName,
			Namespace: captenConfig.CaptenNamespace,
		},
		Data: map[string][]byte{
			"ca.crt": caCertChainData,
		},
		Type: corev1.SecretTypeOpaque,
	}
	return createOrUpdateSecret(k8sClient, secret)
}

func createOrUpdateSecret(k8sClient *kubernetes.Clientset, secret *corev1.Secret) error {
	_, err := k8sClient.CoreV1().Secrets(secret.Namespace).Get(context.Background(), secret.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		_, err := k8sClient.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating secret")
		}
		return nil
	}

	_, err = k8sClient.CoreV1().Secrets(secret.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error in updating secret")
	}
	return nil
}
