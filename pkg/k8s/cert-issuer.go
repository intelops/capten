package k8s

import (
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"

	"github.com/pkg/errors"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateOrUpdateClusterIssuer(captenConfig config.CaptenConfig) error {
	kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	issuer := &certmanagerv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: captenConfig.ClusterCAIssuer,
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: certmanagerv1.IssuerConfig{
				CA: &certmanagerv1.CAIssuer{
					SecretName: captenConfig.ClusterCACertSecretName,
				},
			},
		},
	}

	serverIssuer, err := cmClient.CertmanagerV1().ClusterIssuers().Get(context.Background(), issuer.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		result, err := cmClient.CertmanagerV1().ClusterIssuers().Create(context.Background(), issuer, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating cert issuer")
		}
		clog.Logger.Debugf("ClusterIssuer %s created successfully", result.Name)
		return nil
	}

	serverIssuer.Spec.IssuerConfig.CA.SecretName = captenConfig.ClusterCACertSecretName
	issuerClient := cmClient.CertmanagerV1().ClusterIssuers()
	result, err := issuerClient.Update(context.TODO(), serverIssuer, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error while updating cluster issuer")
	}
	clog.Logger.Debugf("ClusterIssuer %s updated successfully", result.Name)
	return nil
}
