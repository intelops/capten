package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/clog"
	"capten/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
)

func readAndValidPluginStoreFlags(cmd *cobra.Command) (resourceType string, err error) {
	resourceType, _ = cmd.Flags().GetString("type")
	if len(resourceType) == 0 {
		return "", fmt.Errorf("specify the resource type in the command line")
	}
	return
}

var pluginStoreListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "plugin store list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		storeType, err := readAndValidPluginStoreFlags(cmd)
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
		resourceType, err := readAndValidPluginStoreFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin store showed, %s", resourceType)
	},
}

var pluginStoreSynchSubCmd = &cobra.Command{
	Use:   "synch",
	Short: "plugin store synch",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, err := readAndValidPluginStoreFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		clog.Logger.Infof("Plugin store Synched, %s", resourceType)
	},
}
