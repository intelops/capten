package cmd

import (
	"capten/pkg/cluster"
	"capten/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clusterDestroySubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster destroy operations",
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
		err = cluster.Destroy(captenConfig, clusterType, cloudType)
		if err != nil {
			logrus.Errorf("failed to destroy cluster, %v", err)
			return
		}
		logrus.Info("Cluster Destroyed")
	},
}
