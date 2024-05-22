package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/clog"
	"capten/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

func readPluginNameFlags(cmd *cobra.Command) (pluginName string, err error) {
	pluginName, _ = cmd.Flags().GetString("plugin-name")
	if len(pluginName) == 0 {
		return "", fmt.Errorf("specify the plugin name in the command line")
	}
	return
}

func readAndValidDeployPluginBaseFlags(cmd *cobra.Command) (storeType, pluginName string, err error) {
	storeType, _ = cmd.Flags().GetString("store-type")
	if len(storeType) == 0 {
		return "", "", fmt.Errorf("specify the store type in the command line")
	}

	if storeType != "local" && storeType != "central" && storeType != "default" {
		return "", "", fmt.Errorf("invalid store type: %s for list plugin store", storeType)
	}

	pluginName, _ = cmd.Flags().GetString("plugin-name")
	if len(pluginName) == 0 {
		return "", "", fmt.Errorf("specify the plugin name in the command line")
	}
	return
}

func readAndValidDeployPluginFlags(cmd *cobra.Command) (storeType, pluginName, version string, err error) {
	storeType, pluginName, err = readAndValidDeployPluginBaseFlags(cmd)
	if err != nil {
		return "", "", "", err
	}

	version, _ = cmd.Flags().GetString("version")
	if len(version) == 0 {
		return "", "", "", fmt.Errorf("specify the plugin version in the command line")
	}
	return
}

var pluginDeploySubCmd = &cobra.Command{
	Use:   "deploy",
	Short: "plugin deploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		storeType, pluginName, version, err := readAndValidDeployPluginFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.DeployPlugin(captenConfig, storeType, pluginName, version)
		if err != nil {
			clog.Logger.Errorf("failed to trigger deploy plugin, %v", err)
			return
		}

		clog.Logger.Infof("Plugin '%s' deploy triggerred", pluginName)
	},
}

var pluginUnDeploySubCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "plugin undeploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		storeType, pluginName, err := readAndValidDeployPluginBaseFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.UnDeployPlugin(captenConfig, storeType, pluginName)
		if err != nil {
			clog.Logger.Errorf("failed to trigger undeploy plugin, %v", err)
			return
		}

		clog.Logger.Infof("Plugin '%s' un-deploy triggerred", pluginName)
	},
}

var pluginListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "plugin list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}
		err = agent.ListClusterPlugins(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to list cluster plugins, %v", err)
			return
		}
	},
}

var pluginShowSubCmd = &cobra.Command{
	Use:   "show",
	Short: "plugin show",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName, err := readPluginNameFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.ShowClusterPluginData(config.CaptenConfig{}, pluginName)
		if err != nil {
			clog.Logger.Errorf("failed to show cluster plugin data, %v", err)
			return
		}
	},
}

var pluginConfigSubCmd = &cobra.Command{
	Use:   "config",
	Short: "plugin config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName, err := readPluginNameFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.ConfigureClusterPlugin(config.CaptenConfig{}, pluginName)
		if err != nil {
			clog.Logger.Errorf("failed to show cluster plugin data, %v", err)
			return
		}
	},
}
