package k8s

import (
	"context"

	"capten/pkg/config"
	"fmt"

	"strings"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"

	clientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	labelKey         = "topology.cstor.openebs.io/nodeName"
	labelValuePrefix = "ip-192"
	namespace        = "openebs-cstor"
)

// CStorPoolSpec represents the specification for a CStorPool.
type CStorPoolSpec struct {
	NodeSelector map[string]string
	BlockDevices []string
	PoolConfig   PoolConfig
}

// PoolConfig represents the configuration for the CStorPool.
type PoolConfig struct {
	DataRaidGroupType string
}

// getClient creates a new Kubernetes clientset based on the kubeconfig file provided.
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

func getOpenEBSBlockDevices(openebsClientset *clientset.Clientset, node string) ([]string, error) {

	bdList, err := openebsClientset.OpenebsV1alpha1().BlockDevices(namespace).List(context.TODO(), metav1.ListOptions{})

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

func isBlockDeviceClaimed(clientset *clientset.Clientset, blockDevice string) (bool, error) {
	bd, err := clientset.OpenebsV1alpha1().BlockDevices(namespace).Get(context.TODO(), blockDevice, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	res := bd.Status.ClaimState

	if res == "Claimed" {
		return true, nil
	} else {
		return false, nil
	}

}

var sequenceNumber = 1

func generateSequentialName() string {
	name := fmt.Sprintf("cstor-disk-pool-%d", sequenceNumber)
	sequenceNumber++
	return name
}

func createCStorPoolCluster(clientset *clientset.Clientset, poolClusterSpec CStorPoolSpec, blockDevice string) error {
	isClaimed, err := isBlockDeviceClaimed(clientset, blockDevice)
	if err != nil {
		return err
	}
	if isClaimed {
		fmt.Printf("Skipping claimed block device %s.\n", blockDevice)
		return nil
	} else {

		poolClusterClient := clientset.CstorV1().CStorPoolClusters(namespace)

		raid := v1.RaidGroup{
			CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{
				{BlockDeviceName: blockDevice},
			},
		}
		poolCluster := &v1.CStorPoolCluster{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cstor.openebs.io/v1",
				Kind:       "CStorPoolCluster",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      generateSequentialName(),
				Namespace: namespace,
			},
			Spec: v1.CStorPoolClusterSpec{
				Pools: []v1.PoolSpec{
					{
						NodeSelector:   poolClusterSpec.NodeSelector,
						DataRaidGroups: []v1.RaidGroup{raid},
						PoolConfig: v1.PoolConfig{
							DataRaidGroupType: poolClusterSpec.PoolConfig.DataRaidGroupType,
						},
					},
				},
			},
		}

		_, err := poolClusterClient.Create(context.TODO(), poolCluster, metav1.CreateOptions{})
		if err == nil {
			fmt.Printf("CStorPoolCluster %s created successfully in namespace %s.\n", poolCluster.Name, poolCluster.Namespace)
		} else {
			if strings.Contains(err.Error(), "admission-webhook") && strings.Contains(err.Error(), "blockdevice") && strings.Contains(err.Error(), "doesn't belongs to node") {
				fmt.Printf("Skipping block device %s that doesn't belong to the specified node.\n", blockDevice)
			} else {
				return err
			}
		}
	}
	return nil
}

func createCStorPoolClustersForAllBlockDevices(clientset *clientset.Clientset, poolClusterSpecs []CStorPoolSpec) error {
	for _, spec := range poolClusterSpecs {
		for _, blockDevice := range spec.BlockDevices {
			if err := createCStorPoolCluster(clientset, spec, blockDevice); err != nil {
				if !strings.Contains(err.Error(), "admission-webhook") || !strings.Contains(err.Error(), "blockdevice") || !strings.Contains(err.Error(), "doesn't belongs to node") {
					return err
				}
			}
		}
	}

	return nil
}

// createCStorPoolClusterSpecs generates CStorPoolClusterSpecs based on the block devices available on each node.
func createCStorPoolClusterSpecs(blockDevicesByNode map[string][]string) []CStorPoolSpec {
	var poolClusterSpecs []CStorPoolSpec

	for node, blockDevices := range blockDevicesByNode {

		poolClusterSpecs = append(poolClusterSpecs, CStorPoolSpec{
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": node,
			},
			BlockDevices: blockDevices,
			PoolConfig: PoolConfig{
				DataRaidGroupType: "stripe",
			},
		})
	}

	return poolClusterSpecs
}

func OpenEBS(captenConfig config.CaptenConfig) error {
	kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
	clientset, err := GetK8SClient(kubeconfigPath)
	if err != nil {

		fmt.Println("Error creating Kubernetes client:", err)
		return err

	}

	openebsClientset, err := getOpenEBSClient(captenConfig)
	if err != nil {

		fmt.Println("Error creating OpenEBS client:", err)
		return err

	}

	workerNodes, err := getWorkerNodes(clientset, labelKey, labelValuePrefix)
	if err != nil {
		fmt.Println("Error fetching worker nodes:", err)
		return err

	}
	blockDevicesByNode := make(map[string][]string)

	for _, node := range workerNodes.Items {
		nodeName := node.GetName()

		blockDevices, err := getOpenEBSBlockDevices(openebsClientset, nodeName)

		if err != nil {
			fmt.Printf("Error fetching block devices for node %s: %v\n", nodeName, err)
			continue
		}

		blockDevicesByNode[nodeName] = blockDevices

	}

	// Generate CStorPoolClusterSpecs based on the block devices available on each node.
	poolClusterSpecs := createCStorPoolClusterSpecs(blockDevicesByNode)

	err = createCStorPoolClustersForAllBlockDevices(openebsClientset, poolClusterSpecs)
	if err != nil {
		fmt.Println("Error creating CStorPoolClusters:", err)
		return err

	}
	return nil
}

func getWorkerNodes(clientset *kubernetes.Clientset, labelKey, labelValuePrefix string) (*corev1.NodeList, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var filteredNodes []corev1.Node
	for _, node := range nodes.Items {
		for key, value := range node.GetLabels() {
			if key == labelKey && strings.HasPrefix(value, labelValuePrefix) {
				filteredNodes = append(filteredNodes, node)
				break
			}
		}
	}

	filteredNodeList := &corev1.NodeList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeList",
			APIVersion: "v1",
		},
		Items: filteredNodes,
	}

	return filteredNodeList, nil
}
