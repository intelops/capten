package agent

import (
	"capten/pkg/agent/pb/clusterpluginspb"
	"capten/pkg/agent/pb/pluginstorepb"
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func ListClusterPlugins(captenConfig config.CaptenConfig) error {
	client, err := GetClusterPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetClusterPlugins(context.TODO(), &clusterpluginspb.GetClusterPluginsRequest{})
	if err != nil {
		return err
	}

	if len(resp.Plugins) == 0 {
		clog.Logger.Info("No plugins found on cluster")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Name", "Version", "Store Type", "Status"})
	for _, plugin := range resp.Plugins {
		table.Append([]string{plugin.Category, plugin.PluginName, plugin.Version, plugin.StoreType.String(), plugin.InstallStatus})
	}
	table.Render()
	return nil
}

func DeployPlugin(captenConfig config.CaptenConfig, storeType, pluginName, version string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	storeTypeEnum, err := getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}
	_, err = client.DeployPlugin(context.TODO(), &pluginstorepb.DeployPluginRequest{
		StoreType:  storeTypeEnum,
		PluginName: pluginName,
		Version:    version,
	})
	if err != nil {
		return err
	}
	return nil
}

func UnDeployPlugin(captenConfig config.CaptenConfig, storeType, pluginName string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	storeTypeEnum, err := getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}

	_, err = client.UnDeployPlugin(context.TODO(), &pluginstorepb.UnDeployPluginRequest{
		StoreType:  storeTypeEnum,
		PluginName: pluginName,
	})
	if err != nil {
		return err
	}
	return nil
}

func ShowClusterPluginData(captenConfig config.CaptenConfig, pluginName string) error {
	return fmt.Errorf("show plugin data is not implemented yet")
}
