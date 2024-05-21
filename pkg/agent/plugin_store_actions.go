package agent

import (
	"capten/pkg/agent/pb/pluginstorepb"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func ListPluginStoreApps(captenConfig config.CaptenConfig, storeType string) error {
	client, err := GetPluginStoreClient(captenConfig)
	if err != nil {
		return err
	}

	var storeTypeEnum pluginstorepb.StoreType
	switch storeType {
	case "local":
		storeTypeEnum = pluginstorepb.StoreType_LOCAL_STORE
	case "central":
		storeTypeEnum = pluginstorepb.StoreType_CENTRAL_STORE
	case "default":
		storeTypeEnum = pluginstorepb.StoreType_DEFAULT_STORE

	default:
		return fmt.Errorf("invalid store type: %s", storeType)
	}

	resp, err := client.GetPlugins(context.TODO(), &pluginstorepb.GetPluginsRequest{
		StoreType: storeTypeEnum,
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Name", "Version"})
	for _, plugin := range resp.Plugins {
		table.Append([]string{plugin.Category, plugin.PluginName, strings.Join(plugin.Versions, ",")})
	}
	table.Render()
	return nil
}
