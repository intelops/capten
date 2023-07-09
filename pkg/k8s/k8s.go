package k8s

import (
	"context"
	//	"path/filepath"

	//	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/util/homedir"
)

func MakeNamespacePrivilege(kubeconfigPath string, ns string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	nsObj, err := clientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		return err
	}

	nsObj.Labels["pod-security.kubernetes.io/enforce"] = "privileged"
	nsObj, err = clientSet.CoreV1().Namespaces().Update(context.Background(), nsObj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
