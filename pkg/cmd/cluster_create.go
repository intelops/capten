package cmd

import (
	"capten/pkg/clog"
	"capten/pkg/cluster"
	"capten/pkg/config"
	"github.com/olekukonko/tablewriter"
	"os"
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
		clog.Logger.Info("Cluster Created")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetHeader([]string{"Hostname", "Type", "Value"})

		clog.Logger.Println("Before starting the app deployment, please ensure the following record is updated in DNS:")
		table.Append([]string{" *."+captenConfig.DomainName, "CNAME", captenConfig.LoadBalancerHost})
		table.Render()
      
	},
}
