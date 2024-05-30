package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/clog"
	"capten/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

func readAndValidatePluginStoreTypeFlags(cmd *cobra.Command) (storeType string, err error) {
	storeType, _ = cmd.Flags().GetString("store-type")
	if len(storeType) == 0 {
		return "", fmt.Errorf("specify the store type in the command line")
	}

	if storeType != "local" && storeType != "central" && storeType != "default" {
		return "", fmt.Errorf("invalid store type: %s for list plugin store", storeType)
	}
	return
}

func readAndValidatePluginStoreShowFlags(cmd *cobra.Command) (pluginName, storeType string, err error) {
	storeType, err = readAndValidatePluginStoreTypeFlags(cmd)
	if err != nil {
		return "", "", err
	}

	pluginName, _ = cmd.Flags().GetString("plugin-name")
	if len(pluginName) == 0 {
		return "", "", fmt.Errorf("specify the plugin name in the command line")
	}

	return
}

func readAndValidatePluginStoreConfigFlags(cmd *cobra.Command) (storeType, gitProjectId string, err error) {
	storeType, _ = cmd.Flags().GetString("type")
	if len(storeType) == 0 {
		storeType = "local"
	}

	if storeType != "local" {
		return "", "", fmt.Errorf("invalid store type: %s for config plugin store", storeType)
	}

	gitProjectId, _ = cmd.Flags().GetString("git-projec-id")
	if len(gitProjectId) == 0 {
		return "", "", fmt.Errorf("specify the git project identifier in the command line")
	}

	return
}

var pluginStoreListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "plugin store list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		storeType, err := readAndValidatePluginStoreTypeFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenconfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.ListPluginStoreApps(captenconfig, storeType)
		if err != nil {
			clog.Logger.Errorf("failed to list plugin store apps, %v", err)
			return
		}
	},
}

var pluginStoreShowSubCmd = &cobra.Command{
	Use:   "show",
	Short: "plugin store show",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pluginName, storeType, err := readAndValidatePluginStoreShowFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenconfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.ShowPluginStorePlugin(captenconfig, storeType, pluginName)
		if err != nil {
			clog.Logger.Errorf("failed to config plugin store, %v", err)
			return
		}
	},
}

var pluginStoreSynchSubCmd = &cobra.Command{
	Use:   "synch",
	Short: "plugin store synch",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		storeType, err := readAndValidatePluginStoreTypeFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenconfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.SynchPluginStore(captenconfig, storeType)
		if err != nil {
			clog.Logger.Errorf("failed to config plugin store, %v", err)
			return
		}
		clog.Logger.Infof("Plugin store synched, %s", storeType)
	},
}

var pluginStoreConfigSubCmd = &cobra.Command{
	Use:   "config",
	Short: "plugin store config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitProjectId, storeType, err := readAndValidatePluginStoreConfigFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		captenconfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Error(err)
			return
		}

		err = agent.ConfigPluginStore(captenconfig, storeType, gitProjectId)
		if err != nil {
			clog.Logger.Errorf("failed to config plugin store, %v", err)
			return
		}
		clog.Logger.Infof("Plugin store configured")
	},
}
