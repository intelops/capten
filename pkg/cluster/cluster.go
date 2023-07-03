package cluster

import (
	"capten/pkg/cluster/k3s"
	"capten/pkg/config"
)

func Create(captenConfig config.CaptenConfig, clusterType, cloudType string) error {
	return k3s.Create(captenConfig, clusterType, cloudType)
}

func Destroy(captenConfig config.CaptenConfig, clusterType, cloudType string) error {
	return k3s.Destroy(captenConfig, clusterType, cloudType)
}
