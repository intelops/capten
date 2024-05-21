package agent

import (
	"capten/pkg/agent/pb/clusterpluginspb"
	"capten/pkg/config"
	"context"
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Name", "Version", "Store Type", "Status"})
	for _, plugin := range resp.Plugins {
		table.Append([]string{plugin.Category, plugin.PluginName, plugin.Version, plugin.StoreType.String(), plugin.InstallStatus})
	}
	table.Render()
	return nil
}
