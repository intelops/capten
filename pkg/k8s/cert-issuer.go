package k8s

import (
	"capten/pkg/config"
	"context"
	"encoding/json"
	"log"

	//	"path/filepath"

	cmacme "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pkg/errors"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
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
	// kubeconfig := flag.String("kubeconfig", "", "Path to the kubeconfig file")
	// flag.Parse()

	// // Build the config from the kubeconfig file
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Create a Kubernetes client
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Define the ClusterIssuer
	issuer := &certmanagerv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "capten-issuer",
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: certmanagerv1.IssuerConfig{
				ACME: &cmacme.ACMEIssuer{
					PrivateKey: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "capten-ca-cert",
						},
						Key: "privatekey",
					},
				},
			},
		},
		// Spec: certmanagerv1.IssuerSpec{
		// 	SecretName: "capten-ca-cert",
		// },
	}
	issuerJSON, err := json.Marshal(issuer)

	if err != nil {
		return errors.WithMessage(err, "error while marshaling ClusterIssuer to JSON")
	}
	log.Println("Issuer Json is", string(issuerJSON))
	// Create the ClusterIssuer

	_, err = clientset.RESTClient().
		Post().
		AbsPath("/apis/certmanager.k8s.io/v1/clusterissuers").
		Body(issuerJSON).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Do(context.TODO()).
		Get()
	if err != nil {
		return errors.WithMessage(err, "error while creating ClusterIssuer")
	}

	// _, err = clientset.RESTClient().
	// 	Post().
	// 	AbsPath("/apis/certmanager.k8s.io/v1/clusterissuers").
	// 	Body(issuer).
	// 	Do(context.TODO()).
	// 	Get()
	// if err != nil {
	// 	return errors.WithMessage(err, "error while creating clusterIssuer")
	// }

	log.Println("ClusterIssuer created successfully.")
	return nil
}
