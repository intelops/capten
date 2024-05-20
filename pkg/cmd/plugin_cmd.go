package cmd

import (
	"capten/pkg/clog"
	"fmt"

	"github.com/spf13/cobra"
)

func readAndValidDeployPluginFlags(cmd *cobra.Command) (resourceType string, err error) {
	resourceType, _ = cmd.Flags().GetString("type")
	if len(resourceType) == 0 {
		return "", fmt.Errorf("specify the resource type in the command line")
	}
	return
}

var pluginDeployCreateSubCmd = &cobra.Command{
	Use:   "deploy",
	Short: "plugin deploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin Deployed, %s", resourceType)
	},
}

var pluginUnDeployCreateSubCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "plugin undeploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin Undeployed, %s", resourceType)
	},
}

var pluginListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "plugin list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin Listed, %s", resourceType)
	},
}

var pluginShowSubCmd = &cobra.Command{
	Use:   "show",
	Short: "plugin show",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin showed, %s", resourceType)
	},
}

var pluginConfigSubCmd = &cobra.Command{
	Use:   "config",
	Short: "plugin config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin configured, %s", resourceType)
	},
}
