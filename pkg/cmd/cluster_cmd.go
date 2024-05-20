package cmd

import (
	"capten/pkg/clog"
	"capten/pkg/cluster"
	"capten/pkg/config"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var clusterCreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "cluster create operation",
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

		clog.Logger.Info("Cluster Created")
	},
}

var clusterDestroySubCmd = &cobra.Command{
	Use:   "destroy",
	Short: "cluster destroy operation",
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

var showClusterInfoSubCmd = &cobra.Command{
	Use:   "info",
	Short: "cluster show info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Error("failed to read capten config", err)
			return
		}
		fmt.Println(color.New(color.FgGreen).Sprint("Cluster LB Host:"), captenConfig.LoadBalancerHost)
		fmt.Println(color.New(color.FgGreen).Sprint("Capten Agent Hostname:"), captenConfig.AgentHostName)
	},
}
