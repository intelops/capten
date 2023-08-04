package k8s

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func MakeNamespacePrivilege(kubeconfigPath string, ns string) error {
	clientSet, err := GetK8SClient(kubeconfigPath)
	if err != nil {
		return err
	}

	nsObj, err := clientSet.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if err != nil {
		return err
	}

	nsObj.Labels["pod-security.kubernetes.io/enforce"] = "privileged"
	_, err = clientSet.CoreV1().Namespaces().Update(context.Background(), nsObj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetK8SClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, errors.WithMessage(err, "error while building kubeconfig")
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.WithMessage(err, "error while getting k8s config")
	}
	return clientSet, nil
}

func CreateNamespaceIfNotExists(kubeconfigPath, namespace string) error {
	clientSet, err := GetK8SClient(kubeconfigPath)
	if err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			newNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			_, err := clientSet.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create namespace %s: %v", namespace, err)
			}
			return nil
		}
		return fmt.Errorf("failed to get namespace %s: %v", namespace, err)
	}
	return nil
}
