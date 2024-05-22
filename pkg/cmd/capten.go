package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "capten",
	Short: "",
	Long:  `command line tool for building cluster`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster operations",
	Long:  ``,
}

var clusterShowCmd = &cobra.Command{
	Use:   "show",
	Short: "cluster show details",
	Long:  ``,
}

var clusterAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "cluster apps operations",
	Long:  ``,
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "plugin operations",
	Long:  ``,
}

var pluginStoreCmd = &cobra.Command{
	Use:   "store",
	Short: "plugin store operations",
	Long:  ``,
}

var clusterResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "cluster resources operations",
	Long:  ``,
}

func readAndValidClusterFlags(cmd *cobra.Command) (cloudService string, clusterType string, err error) {
	cloudService, err = cmd.Flags().GetString("cloud")
	if len(cloudService) == 0 {
		return "", "", fmt.Errorf("specify the cloud service either azure or aws in the command line %v", err)
	}
	clusterType, _ = cmd.Flags().GetString("type")
	if len(clusterType) == 0 {
		clusterType = "talos"
	}
	err = validateClusterFlags(cloudService, clusterType)
	return
}

func validateClusterFlags(cloudService, clusterType string) (err error) {

	if cloudService != "aws" && cloudService != "azure" {
		err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudService)
		return
	}

	if clusterType != "talos" {
		err = fmt.Errorf("cluster type '%s' is not supported, supported types: talos", clusterType)
		return
	}
	return
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(pluginCmd)

	//cluster optons
	clusterCmd.AddCommand(clusterShowCmd)
	clusterCmd.AddCommand(clusterAppsCmd)
	clusterCmd.AddCommand(clusterResourcesCmd)

	//cluster create options
	clusterCreateSubCmd.PersistentFlags().String("cloud", "", "cloud service (default: azure)")
	clusterCreateSubCmd.PersistentFlags().String("type", "", "type of cluster (default: talos)")
	clusterCmd.AddCommand(clusterCreateSubCmd)

	//cluster destroy options
	clusterCmd.AddCommand(clusterDestroySubCmd)

	//cluster show options
	clusterShowCmd.AddCommand(showClusterInfoSubCmd)

	//cluster apps options
	clusterAppsCmd.AddCommand(appsInstallSubCmd)
	clusterAppsCmd.AddCommand(appsListSubCmd)

	//cluster resources create options
	resourceCreateSubCmd.PersistentFlags().String("type", "", "type of resource")
	clusterResourcesCmd.AddCommand(resourceCreateSubCmd)

	//cluster resources delete options
	resourceDeleteSubCmd.PersistentFlags().String("type", "", "type of resource")
	clusterResourcesCmd.AddCommand(resourceDeleteSubCmd)

	//cluster resources list options
	resourceListSubCmd.PersistentFlags().String("type", "", "type of resource")
	clusterResourcesCmd.AddCommand(resourceListSubCmd)

	//cluster resources show options
	resourceShowSubCmd.PersistentFlags().String("type", "", "type of resource")
	clusterResourcesCmd.AddCommand(resourceShowSubCmd)

	//plugin deploy options
	pluginDeployCreateSubCmd.PersistentFlags().String("type", "", "type of resource")
	pluginCmd.AddCommand(pluginDeployCreateSubCmd)

	//plugin undeploy options
	pluginUnDeployCreateSubCmd.PersistentFlags().String("type", "", "type of resource")
	pluginCmd.AddCommand(pluginUnDeployCreateSubCmd)

	//plugin list options
	pluginListSubCmd.PersistentFlags().String("type", "", "type of resource")
	pluginCmd.AddCommand(pluginListSubCmd)

	//plugin show options
	pluginShowSubCmd.PersistentFlags().String("type", "", "type of resource")
	pluginCmd.AddCommand(pluginShowSubCmd)

	//plugin config options
	pluginConfigSubCmd.PersistentFlags().String("type", "", "type of resource")
	pluginCmd.AddCommand(pluginConfigSubCmd)

	//plugin store options
	pluginCmd.AddCommand(pluginStoreCmd)

	//plugin store list options
	pluginStoreListSubCmd.PersistentFlags().String("store-type", "", "store type (local, central, default)")
	pluginStoreCmd.AddCommand(pluginStoreListSubCmd)

	//plugin store show options
	pluginStoreShowSubCmd.PersistentFlags().String("store-type", "", "store type (local, central, default)")
	pluginStoreShowSubCmd.PersistentFlags().String("plugin-name", "", "name of the plugin")
	pluginStoreCmd.AddCommand(pluginStoreShowSubCmd)

	//plugin store synch options
	pluginStoreSynchSubCmd.PersistentFlags().String("store-type", "", "store type (local, central, default)")
	pluginStoreCmd.AddCommand(pluginStoreSynchSubCmd)

	//plugin store config options
	pluginStoreConfigSubCmd.PersistentFlags().String("store-type", "", "store type (local)")
	pluginStoreConfigSubCmd.PersistentFlags().String("git-project-id", "", "git project identifier")
	pluginStoreCmd.AddCommand(pluginStoreConfigSubCmd)

}
