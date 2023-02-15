package cluster

import (
	"log"

	"capten/pkg/cluster/k3s"
	"capten/pkg/config"
)

func Create(configPath, clusterType, workingDir string) {
	cfg, err := config.GetClusterConfig(configPath)
	if err != nil {
		log.Println("failed to read config", err)
		return
	}

	if err := k3s.Create(cfg, workingDir); err != nil {
		log.Println("failed to create cluster", err)
	}
}

func Destroy(workingDir string) {
	if err := k3s.Destroy(workingDir); err != nil {
		log.Println("failed to destroy cluster", err)
	}
}
