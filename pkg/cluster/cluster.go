package cluster

import (
	"capten/pkg/cluster/k3s"
	"capten/pkg/config"
)

func Create(captenConfig config.CaptenConfig) error {
	return k3s.Create(captenConfig)
}

func Destroy(captenConfig config.CaptenConfig) error {
	return k3s.Destroy(captenConfig)
}