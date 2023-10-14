package k8s

import (
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"
	"strings"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	clientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	"github.com/pkg/errors"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func getOpenEBSBlockDevices(openebsClientset *clientset.Clientset, captenConfig config.CaptenConfig) ([]map[string]string, error) {
	bdList, err := openebsClientset.OpenebsV1alpha1().BlockDevices(captenConfig.PoolClusterNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var blockDevicesMappings []map[string]string
	for _, bd := range bdList.Items {
		if strings.EqualFold(bd.Spec.Details.DeviceType, "disk") {
			nodename := bd.Spec.NodeAttributes.NodeName
			blockDeviceMapping := map[string]string{
				"blockDevice": bd.Name,
				"nodeName":    nodename,
			}
			blockDevicesMappings = append(blockDevicesMappings, blockDeviceMapping)
			clog.Logger.Debugf("Added Block device %s", bd.Name)
		} else {
			clog.Logger.Infof("Skipped Block device %s", bd.Name)
		}
	}
	return blockDevicesMappings, nil
}

func CreateCStorPoolClusters(captenConfig config.CaptenConfig) error {
	openebsClientset, err := getOpenEBSClient(captenConfig)
	if err != nil {
		return errors.WithMessage(err, "error while creating openebsClientset")
	}

	bdNodeMap, err := getOpenEBSBlockDevices(openebsClientset, captenConfig)
	if err != nil {
		return errors.WithMessage(err, "failed to retrieve blockdevices")
	}

	var poolSpecs []v1.PoolSpec
	for _, bd := range bdNodeMap {
		blockDevice := bd["blockDevice"]
		nodeName := bd["nodeName"]
		if err != nil {
			return errors.WithMessage(err, "failed to retrieve node")

		}
		instancebd := []v1.CStorPoolInstanceBlockDevice{{
			BlockDeviceName: blockDevice,
		},
		}
		poolspec := v1.PoolSpec{
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": nodeName,
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
	poolCluster, err := poolClusterClient.Get(context.Background(), captenConfig.PoolClusterName, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
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
			return err
		}
	}

	poolCluster.Spec = v1.CStorPoolClusterSpec{
		Pools: poolSpecs,
	}
	_, err = poolClusterClient.Update(context.TODO(), poolCluster, metav1.UpdateOptions{})
	return err
}
