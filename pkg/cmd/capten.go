package cmd

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"capten/pkg/api"
	"capten/pkg/cluster"
	"capten/pkg/helm"
	"capten/pkg/util"
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
		configPath, _ := cmd.Flags().GetString("config")
		kubeCfgPath, _ := cmd.Flags().GetString("kubeconfig")
		if err := util.OsExec("bash", "./generate.sh"); err != nil {
			logrus.Errorf("failed to generate certificate %v", err)
			return
		}

		helmObj, err := helm.NewHelm(configPath, kubeCfgPath)
		if err != nil {
			log.Println("failed to setup apps", err)
			return
		}

		helmObj.Install()
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
		//configPath, _ := cmd.Flags().GetString("config")
		//clusterType, _ := cmd.Flags().GetString("type")
		workingDir, _ := cmd.Flags().GetString("work-dir")
		cluster.Destroy(workingDir)
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
		customerId, _ := cmd.Flags().GetString("customerId")
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

	appsCmd.PersistentFlags().String("config", "", "config path")
	appsCmd.PersistentFlags().String("kubeconfig", "", "kube config path")
	_ = appsCmd.MarkPersistentFlagRequired("config")
	_ = appsCmd.MarkPersistentFlagRequired("kubeconfig")

	registerAgentCmd.PersistentFlags().String("host", "", "endpoint of agent that needs to be registered")
	registerAgentCmd.PersistentFlags().String("customerId", "", "customerId to be registered for")
	_ = registerAgentCmd.MarkPersistentFlagRequired("host")
	_ = registerAgentCmd.MarkPersistentFlagRequired("customerId")

	createCmd.AddCommand(clusterCreateSubCmd)
	destroyCmd.AddCommand(clusterDestroySubCmd)
	setupCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(registerAgentCmd)
}
