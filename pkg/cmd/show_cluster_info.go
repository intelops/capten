package cmd

import (
	"capten/pkg/config"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var showClusterInfoCmd = &cobra.Command{
	Use:   "clusterinfo",
	Short: "show the clinster info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Error("failed to read capten config", err)
			return
		}
		fmt.Println("Capten Agent Hostname :", captenConfig.AgentHostName)
		fmt.Println("Cluster LB Host :", captenConfig.ClusterLBHost)
	},
}
