package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "capten",
	Short: "",
	Long:  `command line tool for building cluster`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

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
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show cluster details",
	Long:  ``,
}

func readAndValidClusterFlags(cmd *cobra.Command) (cloudService string, clusterType string, err error) {
	cloudService, _ = cmd.Flags().GetString("cloud")
	if len(cloudService) == 0 {
		cloudService = "azure"
	} else {
		cloudService = "aws"
	}

	clusterType, _ = cmd.Flags().GetString("type")
	if len(clusterType) == 0 {
		clusterType = "talos"
	}
	err = validateClusterFlags(cloudService, clusterType)
	return
}

func validateClusterFlags(cloudService, clusterType string) (err error) {
	if cloudService == "aws" {
		//err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudService)
		return
	} else if cloudService == "azure" {
		return
	} else {
		err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudService)
		//  return
	}

	if clusterType != "talos" {
		err = fmt.Errorf("cluster type '%s' is not supported, supported types: talos", clusterType)
		return
	}
	return
}

func init() {
	clusterCreateSubCmd.PersistentFlags().String("cloud", "", "cloud service (default: azure)")
	clusterCreateSubCmd.PersistentFlags().String("type", "", "type of cluster (default: talos)")

	createCmd.AddCommand(clusterCreateSubCmd)
	destroyCmd.AddCommand(clusterDestroySubCmd)
	showCmd.AddCommand(showClusterInfoCmd)
	setupCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(showCmd)
}
