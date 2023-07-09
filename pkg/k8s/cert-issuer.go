package k8s

import (
	"capten/pkg/config"
	"context"
	"encoding/json"
	"log"

	cmacme "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pkg/errors"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateOrUpdateClusterSecret(captenConfig config.CaptenConfig) error {

	config, err := clientcmd.BuildConfigFromFlags("", captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.WithMessage(err, "error while getting k8s config")
	}

	// Define the ClusterIssuer
	issuer := &certmanagerv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "capten-issuer",
		},
		Spec: certmanagerv1.IssuerSpec{

			IssuerConfig: certmanagerv1.IssuerConfig{
				ACME: &cmacme.ACMEIssuer{
					Server:         "https://acme-v02.api.letsencrypt.org/directory",
					PreferredChain: "",
					PrivateKey: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "capten-ca-cert",
						},
						Key: "privatekey",
					},
				},
			},
		},
	}
	issuerJSON, err := json.Marshal(issuer)

	if err != nil {
		return errors.WithMessage(err, "error while marshaling ClusterIssuer to JSON")
	}

	// Create the ClusterIssuer

	res := clientset.RESTClient().
		Post().
		AbsPath("/apis/certmanager.k8s.io/v1/clusterissuers").
		Body(issuerJSON).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Do(context.TODO())
	log.Println("Res is", res)
	if err != nil {
		return errors.WithMessage(err, "error while creating ClusterIssuer")
	}

	log.Println("ClusterIssuer created successfully.")
	return nil
}
