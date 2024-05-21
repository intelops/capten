package agent

import (
	"capten/pkg/agent/pb/captenpluginspb"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func ListClusterResources(captenConfig config.CaptenConfig, resourceType string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	switch resourceType {
	case "git-project":
		projects, err := client.GetGitProjects(context.TODO(), &captenpluginspb.GetGitProjectsRequest{})
		if err != nil {
			return err
		}
		table.SetHeader([]string{"ID", "Project URL", "Labels"})
		for _, project := range projects.Projects {
			table.Append([]string{project.Id, project.ProjectUrl, strings.Join(project.Labels, ",")})
		}
		table.Render()
	case "cloud-provider":
		providers, err := client.GetCloudProviders(context.TODO(), &captenpluginspb.GetCloudProvidersRequest{})
		if err != nil {
			return err
		}
		table.SetHeader([]string{"ID", "Cloud Type", "Labels"})
		for _, provider := range providers.CloudProviders {
			table.Append([]string{provider.Id, provider.CloudType, strings.Join(provider.Labels, ",")})
		}
		table.Render()
	case "container-registry":
		registries, err := client.GetContainerRegistry(context.TODO(), &captenpluginspb.GetContainerRegistryRequest{})
		if err != nil {
			return err
		}
		table.SetHeader([]string{"ID", "Registry Type", "Registry URL", "Labels"})
		for _, registry := range registries.Registries {
			table.Append([]string{registry.Id, registry.RegistryType, registry.RegistryUrl, strings.Join(registry.Labels, ",")})
		}
		table.Render()
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}
	return nil
}
