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
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}

		err = validateClusterFlags(captenConfig.CloudService, captenConfig.ClusterType)
		if err != nil {
			logrus.Errorf("cluster config not valid, %v", err)
			return
		}

		err = cluster.Destroy(captenConfig)
		if err != nil {
			logrus.Errorf("failed to destroy cluster, %v", err)
			return
		}
		logrus.Info("Cluster Destroyed")
	},
}
