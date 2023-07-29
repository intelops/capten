package cmd

import (
	"capten/pkg/clog"
	"capten/pkg/cluster"
	"capten/pkg/config"

	"github.com/spf13/cobra"
)

var clusterDestroySubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster destroy operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = validateClusterFlags(captenConfig.CloudService, captenConfig.ClusterType)
		if err != nil {
			clog.Logger.Errorf("cluster config not valid, %v", err)
			return
		}

		err = cluster.Destroy(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to destroy cluster, %v", err)
			return
		}
		clog.Logger.Info("Cluster Destroyed")
	},
}
