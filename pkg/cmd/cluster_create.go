package cmd

import (
	"capten/pkg/clog"
	"capten/pkg/cluster"
	"capten/pkg/config"

	"github.com/spf13/cobra"
)

var clusterCreateSubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster create operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cloudService, clusterType, err := readAndValidClusterFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = config.UpdateClusterValues(&captenConfig, cloudService, clusterType)
		if err != nil {
			clog.Logger.Errorf("failed to update capten config, %v", err)
			return
		}

		err = cluster.Create(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to create cluster, %v", err)
			return
		}
		err = config.UpdateClusterEndpoint(&captenConfig, captenConfig.CaptenClusterHost.LoadBalancerHost)
		if err != nil {
			clog.Logger.Errorf("failed to update LoadBalancer Host, %v", err)
			return
		}
		clog.Logger.Info("Cluster Created")
	},
}
