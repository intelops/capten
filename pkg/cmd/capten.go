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
	clusterAppsCmd.AddCommand(appsShowSubCmd)

	//cluster resources create options
	resourceCreateSubCmd.PersistentFlags().String("type", "", "type of resource")
	resourceCreateSubCmd.PersistentFlags().String("git-project-url", "", "url of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("access-token", "", "access token of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("user-id", "", "user id of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("labels", "", "labels of resource (e.g. 'crossplane,tekton')")
	resourceCreateSubCmd.PersistentFlags().String("registry-url", "", "registry url of container registry resource")
	resourceCreateSubCmd.PersistentFlags().String("registry-type", "", "registry type of container registry resource")
	resourceCreateSubCmd.PersistentFlags().String("cloud-type", "", "cloud type of cloud provider resource (aws, azure)")
	resourceCreateSubCmd.PersistentFlags().String("registry-username", "", "registry user name of container registry resource")
	resourceCreateSubCmd.PersistentFlags().String("registry-password", "", "registry password of container registry resource")
	resourceCreateSubCmd.PersistentFlags().String("access-key", "", "access key of aws cloud provider resource")
	resourceCreateSubCmd.PersistentFlags().String("secret-key", "", "secret key of aws cloud provider resource")
	resourceCreateSubCmd.PersistentFlags().String("client-id", "", "client id of azure cloud provider resource")
	resourceCreateSubCmd.PersistentFlags().String("client-secret", "", "client secret of azure cloud provider resource")
	clusterResourcesCmd.AddCommand(resourceCreateSubCmd)

	//cluster resources update options
	resourceUpdateSubCmd.PersistentFlags().String("type", "", "type of resource")
	resourceUpdateSubCmd.PersistentFlags().String("id", "", "id of resource")
	resourceUpdateSubCmd.PersistentFlags().String("git-project-url", "", "url of git project resource")
	resourceUpdateSubCmd.PersistentFlags().String("access-token", "", "access token of git project resource")
	resourceUpdateSubCmd.PersistentFlags().String("user-id", "", "user id of git project resource")
	resourceUpdateSubCmd.PersistentFlags().String("labels", "", "labels of resource (e.g. 'crossplane,tekton')")
	resourceUpdateSubCmd.PersistentFlags().String("registry-url", "", "registry url of container registry resource")
	resourceUpdateSubCmd.PersistentFlags().String("registry-type", "", "registry type of container registry resource")
	resourceUpdateSubCmd.PersistentFlags().String("cloud-type", "", "cloud type of cloud provider resource")
	resourceUpdateSubCmd.PersistentFlags().String("registry-username", "", "registry user name of container registry resource")
	resourceUpdateSubCmd.PersistentFlags().String("registry-password", "", "registry password of container registry resource")
	resourceUpdateSubCmd.PersistentFlags().String("access-key", "", "access key of aws cloud provider resource")
	resourceUpdateSubCmd.PersistentFlags().String("secret-key", "", "secret key of aws cloud provider resource")
	resourceUpdateSubCmd.PersistentFlags().String("client-id", "", "client id of azure cloud provider resource")
	resourceUpdateSubCmd.PersistentFlags().String("client-secret", "", "client secret of azure cloud provider resource")
	clusterResourcesCmd.AddCommand(resourceUpdateSubCmd)

	//cluster resources delete options
	resourceDeleteSubCmd.PersistentFlags().String("type", "", "type of resource")
	resourceDeleteSubCmd.PersistentFlags().String("id", "", "id of resource")
	clusterResourcesCmd.AddCommand(resourceDeleteSubCmd)

	//cluster resources list options
	resourceListSubCmd.PersistentFlags().String("type", "", "type of resource")
	clusterResourcesCmd.AddCommand(resourceListSubCmd)

	//plugin deploy options
	pluginDeploySubCmd.PersistentFlags().String("store-type", "", "store type (local, central, default)")
	pluginDeploySubCmd.PersistentFlags().String("plugin-name", "", "name of the plugin")
	pluginDeploySubCmd.PersistentFlags().String("version", "", "version of the plugin")
	pluginCmd.AddCommand(pluginDeploySubCmd)

	//plugin undeploy options
	pluginUnDeploySubCmd.PersistentFlags().String("store-type", "", "store type (local, central, default)")
	pluginUnDeploySubCmd.PersistentFlags().String("plugin-name", "", "name of the plugin")
	pluginCmd.AddCommand(pluginUnDeploySubCmd)

	//plugin list options
	pluginCmd.AddCommand(pluginListSubCmd)

	//plugin show options
	pluginShowSubCmd.PersistentFlags().String("plugin-name", "", "name of the plugin")
	pluginCmd.AddCommand(pluginShowSubCmd)

	//plugin config options
	pluginConfigSubCmd.PersistentFlags().String("plugin-name", "", "name of the plugin")
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
