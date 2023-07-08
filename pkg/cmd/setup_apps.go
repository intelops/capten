package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/app"
	"capten/pkg/cert"
	"capten/pkg/config"
	"capten/pkg/k8s"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "sets up apps cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}

		if err := cert.GenerateCerts(captenConfig); err != nil {
			logrus.Errorf("failed to generate certificate, %v", err)
			return
		}
		logrus.Info("Generated Certificates")

		if err := k8s.CreateOrUpdateCertSecrets(captenConfig); err != nil {
			logrus.Errorf("failed to create secret for certs, %v", err)
			return
		}
		logrus.Info("Configured Certificates on Capten Cluster")

		err = app.DeployApps(captenConfig)
		if err != nil {
			logrus.Errorf("applications deployment failed, %v", err)
			return
		}

		if captenConfig.StoreCredOnAgent {
			err = agent.StoreCredential(captenConfig)
			if err != nil {
				logrus.Errorf("store cluster credentials failed, %v", err)
				return
			}
		}

		//push the app config to cluster
		//prepare agent proto to push app config
		//agent store data on cassandra
		logrus.Info("Default Applications Installed")
	},
}
