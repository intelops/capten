package k8s

import (
	"context"

	"fmt"
	"log"
	
	//"capten/pkg/config"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)
const(
	
	labelkey = "pod-security.kubernetes.io/enforce"
	labelValue = "privileged"
)

func CreateorUpdateNamespaceWithLabel(kubeconfigPath,namespaceName string) error{
	// Initialize the Kubernetes client
	//kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	clientset, err := GetK8SClient(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing client: %v", err)
		return err
	}

	// Check if the namespace exists
	exists, err := namespaceExists(clientset, namespaceName)
	
	if err != nil {
		log.Fatalf("Error checking namespace existence: %v", err)
		return err
	}

	if !exists {
		err := createNamespace(clientset, namespaceName, labelkey, labelValue)
		if err != nil {
			log.Fatalf("Error creating namespace: %v", err)
			return err
		}
		fmt.Printf("Namespace %s created with label %s=%s\n", namespaceName, labelkey, labelValue)
	} else {
		err := updateNamespaceLabel(clientset, namespaceName, labelkey, labelValue)
		if err != nil {
			log.Fatalf("Error updating namespace: %v", err)
			return err
		}
		fmt.Printf("Namespace %s updated with label %s=%s\n", namespaceName, labelkey, labelValue)
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
func createNamespace(clientset *kubernetes.Clientset, name, labelkey, labelValue string) error {
	newNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				labelkey: labelValue,
			},

		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
	return err
}


func updateNamespaceLabel(clientset *kubernetes.Clientset, name, labelkey, labelValue string) error {
	namespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	namespace.Labels[labelkey] = labelValue
	_, err = clientset.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
	return err
}
