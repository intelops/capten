package cmd

import (
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

		if err := k8s.CreateOrUpdateAgnetCertSecret(captenConfig); err != nil {
			logrus.Errorf("failed to patch namespace with privilege, %v", err)
			return
		}
		logrus.Info("Configured Certificates on Capten Cluster")

		err = app.DeployApps(captenConfig)
		if err != nil {
			logrus.Errorf("setup applications failed, %v", err)
			return
		}

		//push kubeconfig and bucket credential to cluster

		//push the app config to cluster
		//prepare agent proto to push app config
		//agent store data on cassandra
		logrus.Info("Default Applications Installed")
	},
}
