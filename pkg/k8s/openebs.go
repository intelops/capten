package k8s

import (
	"capten/pkg/config"
	"context"
	
	"fmt"
	"log"

	"github.com/pkg/errors"


	//"k8s.io/client-go/kubernetes"
	//"k8s.io/kubernetes/plugin/pkg/auth/authorizer/node"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	clientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	
)


func getOpenEBSClient(captenConfig config.CaptenConfig) (*clientset.Clientset, error) {
	kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, errors.WithMessage(err, "error while building kubeconfig")
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}




func getOpenEBSBlockDevices(openebsClientset *clientset.Clientset, captenConfig config.CaptenConfig ) ([]string, error) {

	bdList, err := openebsClientset.OpenebsV1alpha1().BlockDevices(captenConfig.PoolClusterNamespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var blockDevices []string
	for _, bd := range bdList.Items {
		if bd.Name != "" {

			blockDevices = append(blockDevices, bd.Name)
		}

	}

	return blockDevices, nil
}

func getWorkerNodeForBlockDevice(openebsClientset *clientset.Clientset, blockDeviceName string,captenConfig config.CaptenConfig) (string, error) {
	bd, err := openebsClientset.OpenebsV1alpha1().BlockDevices(captenConfig.PoolClusterNamespace).Get(context.TODO(), blockDeviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	nodeName := bd.Spec.NodeAttributes.NodeName

	return nodeName, nil
}



func CreateCStorPoolClusters(captenConfig config.CaptenConfig) error {
    
	openebsClientset, err := getOpenEBSClient(captenConfig)
	if err != nil {
		log.Println("Error while creating openebs client", openebsClientset)
	}
	
	if err != nil {
		log.Println("Error while creating openebs client", openebsClientset)
	}

	blockdevice, err := getOpenEBSBlockDevices(openebsClientset,captenConfig)
	
    if (err!=nil) {
		return fmt.Errorf("failed to retrieve blockdevices %v",  err)
	}

	var poolSpecs []v1.PoolSpec
	
	for _, bd := range blockdevice {
		bdname, err := getWorkerNodeForBlockDevice(openebsClientset, bd,captenConfig)
		
		if err != nil {
			log.Println("Error while retrieve node for bd", err)
		}

		instancebd := []v1.CStorPoolInstanceBlockDevice{{
			BlockDeviceName: bd,
		},
		}
		poolspec := v1.PoolSpec{
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": bdname,
			},
			DataRaidGroups: []v1.RaidGroup{
				{
					CStorPoolInstanceBlockDevices: instancebd,
				},
			},
			PoolConfig: v1.PoolConfig{
				DataRaidGroupType: "stripe",
			},

		}
		poolSpecs = append(poolSpecs, poolspec)
	
	}
	poolClusterClient := openebsClientset.CstorV1().CStorPoolClusters(captenConfig.PoolClusterNamespace)
	poolCluster := &v1.CStorPoolCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      captenConfig.PoolClusterName,
			Namespace: captenConfig.PoolClusterNamespace,
		},
		Spec: v1.CStorPoolClusterSpec{
			Pools: poolSpecs,
		},
	}
	_, err = poolClusterClient.Create(context.TODO(), poolCluster, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("CStorPoolCluster %s created successfully in namespace %s.\n", poolCluster.Name, poolCluster.Namespace)
	}


	return nil
}
