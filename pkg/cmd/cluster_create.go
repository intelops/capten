package cmd

import (
	"capten/pkg/cluster"
	"capten/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clusterCreateSubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster create operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		clusterType, cloudType, err := readAndValidClusterFlags(cmd)
		if err != nil {
			logrus.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}
		err = cluster.Create(captenConfig, clusterType, cloudType)
		if err != nil {
			logrus.Errorf("failed to create cluster, %v", err)
			return
		}
		logrus.Info("Cluster Created")
	},
}
