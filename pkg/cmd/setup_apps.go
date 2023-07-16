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

		if err := cert.PrepareCerts(captenConfig); err != nil {
			logrus.Errorf("failed to generate certificate, %v", err)
			return
		}
		logrus.Info("Certificates prepared for cluster")

		if err := k8s.CreateOrUpdateCertSecrets(captenConfig); err != nil {
			logrus.Errorf("failed to create secret for certs, %v", err)
			return
		}

		err = k8s.CreateOrUpdateClusterIssuer(captenConfig)
		if err != nil {
			logrus.Errorf("failed to create cluster issuer, %v", err)
			return
		}
		logrus.Info("Configured Certificates on Capten Cluster")

		globalValues, err := app.PrepareGlobalVaules(captenConfig)
		if err != nil {
			logrus.Errorf("applications values preparation failed, %v", err)
			return
    }
    
		if captenConfig.SKipAppsDeploy {
			err = app.DeployApps(captenConfig)
			if err != nil {
				logrus.Errorf("applications deployment failed, %v", err)
				return
			}
		}

		if !captenConfig.SKipAppsDeploy {
			err = app.DeployApps(captenConfig, globalValues)
			if err != nil {
				logrus.Errorf("applications deployment failed, %v", err)
				return
			}
		}

		if captenConfig.StoreCredOnAgent {
			err = agent.StoreCredentials(captenConfig, globalValues)
			if err != nil {
				logrus.Errorf("failed to store cluster credentials, %v", err)
				return
			}
		}

		//push the app config to cluster
		//prepare agent proto to push app config
		//agent store data on cassandra
		logrus.Info("Default Applications Installed")
	},
}
