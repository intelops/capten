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
		clog.Logger.Info("Cluster Created")
		cfg, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to get capten config, %v", err)
			return
		}
		cfg.AgentDNSNamePrefixes = []string{"gitbridge", "containerbridge", "loki", "grafana", "prometheus", "signoz", "otelcollector", "tracetest"}
		for _, prefixName := range cfg.AgentDNSNamePrefixes {
			cfg.AgentDNSNames = append(cfg.AgentDNSNames, prefixName+"."+cfg.DomainName)
		}
		agenthostname := cfg.AgentHostName + cfg.DomainName

		clog.Logger.Println("Before starting the app deployment, please ensure the following records are updated in DNS:")
		for _, prefixName := range cfg.AgentDNSNames {
			clog.Logger.Printf("Hostname: %s ", prefixName)
			clog.Logger.Println("Type: CNAME")
			clog.Logger.Printf("Value: %s \n", cfg.LoadBalancerHost)
		}
		clog.Logger.Printf("Hostname: %s", agenthostname)
		clog.Logger.Println("Type: CNAME")
		clog.Logger.Printf("Value:%s \n", cfg.LoadBalancerHost)

	},
}
