package k8s

import (
	"capten/pkg/clog"
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels = map[string]string{
	"pod-security.kubernetes.io/enforce": "privileged",
}

func CreateNamespaceIfNotExist(kubeconfigPath, namespaceName string) error {
	clientset, err := GetK8SClient(kubeconfigPath)
	if err != nil {
		return err
	}

	exists, err := namespaceExists(clientset, namespaceName)
	if err != nil {
		return err
	}

	if !exists {
		err := createNamespace(clientset, namespaceName, labels)
		if err != nil {
			return err
		}
		clog.Logger.Debugf("Namespace %s created with label %s\n", namespaceName, labels)

	}
	return nil
}

func CreateorUpdateNamespaceWithLabel(kubeconfigPath, namespaceName string) error {
	clientset, err := GetK8SClient(kubeconfigPath)
	if err != nil {
		return err
	}

	exists, err := namespaceExists(clientset, namespaceName)
	if err != nil {
		return err
	}

	if !exists {
		err := createNamespace(clientset, namespaceName, labels)
		if err != nil {
			return err
		}
		clog.Logger.Debugf("Namespace %s created with label %s\n", namespaceName, labels)

	} else {
		err := updateNamespaceLabel(clientset, namespaceName, labels)
		if err != nil {
			return err
		}
		clog.Logger.Debugf("Namespace %s updated with label %s\n", namespaceName, labels)

	}
	return nil
}

func namespaceExists(clientset *kubernetes.Clientset, name string) (bool, error) {
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err == nil {
		return true, nil
	}
	if errors.IsNotFound(err) {
		return false, nil
	}
	return false, err

}
func createNamespace(clientset *kubernetes.Clientset, name string, labels map[string]string) error {
	newNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
	return err
}

func updateNamespaceLabel(clientset *kubernetes.Clientset, name string, labels map[string]string) error {
	namespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for key, value := range labels {
		namespace.Labels[key] = value
	}

	_, err = clientset.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
	return err
}
