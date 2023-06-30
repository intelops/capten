package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"capten/pkg/api"
	"capten/pkg/cert"
	"capten/pkg/cluster"
	"capten/pkg/config"
	"capten/pkg/helm"
	"capten/pkg/k8s"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "for creation of resources or cluster",
	Long:  ``,
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy created cluster",
	Long:  ``,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "sets up cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//h := helm.NewHelm()
	},
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "sets up apps cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Error("failed to read capten config", err)
			return
		}

		if err := cert.GenerateCerts(captenConfig); err != nil {
			logrus.Errorf("failed to generate certificate. Error - %v", err)
			return
		}
		logrus.Info("Generated Certificates")

		if err := k8s.CreateOrUpdateAgnetCertSecret(captenConfig); err != nil {
			logrus.Error("failed to patch namespace with privilege", err)
			return
		}
		logrus.Info("Configured Certificates on Capten Cluster")

		helmObj, err := helm.NewHelm(captenConfig)
		if err != nil {
			logrus.Error("applications installation failed", err)
			return
		}
		helmObj.Install()
		logrus.Info("Default Applications Installed")
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall intelop's admin apps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var clusterDestroySubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster create/destroy operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		//clusterType, _ := cmd.Flags().GetString("type")
		workingDir, _ := cmd.Flags().GetString("work-dir")
		cluster.Destroy(configPath, workingDir)
	},
}

var clusterCreateSubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster create/destroy operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		clusterType, _ := cmd.Flags().GetString("type")
		workingDir, _ := cmd.Flags().GetString("work-dir")
		cluster.Create(configPath, clusterType, workingDir)
	},
}

var registerAgentCmd = &cobra.Command{
	Use:   "register",
	Short: "registers the endpoint and certs of agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		agentHost, _ := cmd.Flags().GetString("host")
		apps, _ := cmd.Flags().GetBool("apps")
		customerId, _ := cmd.Flags().GetString("customerId")
		if apps {

		}

		api.RegisterAgentInfo(customerId, agentHost)
	},
}

func init() {
	clusterCreateSubCmd.PersistentFlags().String("config", "", "config path")
	clusterCreateSubCmd.PersistentFlags().String("type", "", "type of cluster")
	clusterCreateSubCmd.PersistentFlags().String("work-dir", "", "terraform work directory path")
	_ = clusterCreateSubCmd.MarkPersistentFlagRequired("config")

	clusterDestroySubCmd.PersistentFlags().String("work-dir", "", "terraform work directory path")
	_ = clusterDestroySubCmd.MarkPersistentFlagRequired("config")

	registerAgentCmd.PersistentFlags().String("host", "", "endpoint of agent that needs to be registered")
	registerAgentCmd.PersistentFlags().Bool("apps", true, "endpoint of agent that needs to be registered")
	registerAgentCmd.PersistentFlags().String("customerId", "", "customerId to be registered for")
	//_ = registerAgentCmd.MarkPersistentFlagRequired("host")
	_ = registerAgentCmd.MarkPersistentFlagRequired("customerId")

	createCmd.AddCommand(clusterCreateSubCmd)
	destroyCmd.AddCommand(clusterDestroySubCmd)
	setupCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(registerAgentCmd)
}
