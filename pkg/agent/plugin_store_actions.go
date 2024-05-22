package agent

import (
	"capten/pkg/agent/pb/pluginstorepb"
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func getStoreTypeEnum(storeType string) (pluginstorepb.StoreType, error) {
	switch storeType {
	case "local":
		return pluginstorepb.StoreType_LOCAL_STORE, nil
	case "central":
		return pluginstorepb.StoreType_CENTRAL_STORE, nil
	case "default":
		return pluginstorepb.StoreType_DEFAULT_STORE, nil
	default:
		return pluginstorepb.StoreType_DEFAULT_STORE, fmt.Errorf("invalid store type: %s", storeType)
	}
}

func ListPluginStoreApps(captenConfig config.CaptenConfig, storeType string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	var storeTypeEnum pluginstorepb.StoreType
	storeTypeEnum, err = getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}

	resp, err := client.GetPlugins(context.TODO(), &pluginstorepb.GetPluginsRequest{
		StoreType: storeTypeEnum,
	})
	if err != nil {
		return err
	}

	if len(resp.Plugins) == 0 {
		clog.Logger.Info("No plugins found on plugin store")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Name", "Version"})
	for _, plugin := range resp.Plugins {
		table.Append([]string{plugin.Category, plugin.PluginName, strings.Join(plugin.Versions, ",")})
	}
	table.Render()
	return nil
}

func ConfigPluginStore(captenConfig config.CaptenConfig, storeType, gitProjectId string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	var storeTypeEnum pluginstorepb.StoreType
	storeTypeEnum, err = getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}

	_, err = client.ConfigurePluginStore(context.TODO(), &pluginstorepb.ConfigurePluginStoreRequest{
		Config: &pluginstorepb.PluginStoreConfig{
			StoreType:    storeTypeEnum,
			GitProjectId: gitProjectId,
		},
	})

	return err
}

func SynchPluginStore(captenConfig config.CaptenConfig, storeType string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	var storeTypeEnum pluginstorepb.StoreType
	storeTypeEnum, err = getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}

	_, err = client.SyncPluginStore(context.TODO(), &pluginstorepb.SyncPluginStoreRequest{
		StoreType: storeTypeEnum,
	})

	return err
}

func ShowPluginStorePlugin(captenConfig config.CaptenConfig, storeType, pluginName string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	var storeTypeEnum pluginstorepb.StoreType
	storeTypeEnum, err = getStoreTypeEnum(storeType)
	if err != nil {
		return err
	}

	resp, err := client.GetPluginData(context.TODO(), &pluginstorepb.GetPluginDataRequest{
		StoreType: storeTypeEnum,
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})
	table.Append([]string{"plugin-name", resp.PluginData.PluginName})
	table.Append([]string{"category", resp.PluginData.Category})
	table.Append([]string{"versions", strings.Join(resp.PluginData.Versions, ",")})
	table.Append([]string{"description", resp.PluginData.Description})
	table.Append([]string{"store-type", storeType})
	table.Render()
	return err
}
