package cmd

import (
	"capten/pkg/clog"
	"capten/pkg/config"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showClusterInfoCmd = &cobra.Command{
	Use:   "clusterinfo",
	Short: "show cluster info",
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
