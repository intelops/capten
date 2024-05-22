package agent

import (
	"capten/pkg/agent/pb/captenpluginspb"
	"capten/pkg/clog"
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

		if len(projects.Projects) == 0 {
			clog.Logger.Info("No git projects added to cluster")
			return nil
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

		if len(providers.CloudProviders) == 0 {
			clog.Logger.Info("No cloud providers added to cluster")
			return nil
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

		if len(registries.Registries) == 0 {
			clog.Logger.Info("No container registries added to cluster")
			return nil
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

func AddClusterResource(captenConfig config.CaptenConfig, resourceType string, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	switch resourceType {
	case "git-project":
		_, err = client.AddGitProject(context.TODO(), &captenpluginspb.AddGitProjectRequest{
			ProjectUrl:  attributes["git-project-url"],
			Labels:      strings.Split(attributes["labels"], ","),
			AccessToken: attributes["access-token"],
			UserID:      attributes["user-id"],
		})
	case "cloud-provider":
		var cloudAttributes map[string]string
		cloudAttributes, err = prepareCloudAttributes(attributes)
		if err != nil {
			return err
		}

		_, err = client.AddCloudProvider(context.TODO(), &captenpluginspb.AddCloudProviderRequest{
			CloudType:       attributes["cloud-type"],
			Labels:          strings.Split(attributes["labels"], ","),
			CloudAttributes: cloudAttributes,
		})
	case "container-registry":
		_, err = client.AddContainerRegistry(context.TODO(), &captenpluginspb.AddContainerRegistryRequest{
			RegistryUrl:  attributes["registry-url"],
			Labels:       strings.Split(attributes["labels"], ","),
			RegistryType: attributes["registry-type"],
			RegistryAttributes: map[string]string{
				"registry-username": attributes["username"],
				"registry-password": attributes["password"],
			},
		})
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}

	if err != nil {
		return err
	}

	clog.Logger.Infof("%s resource added to cluster", resourceType)
	return nil
}

func UpdateClusterResource(captenConfig config.CaptenConfig, resourceType, id string, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	switch resourceType {
	case "git-project":
		_, err = client.UpdateGitProject(context.TODO(), &captenpluginspb.UpdateGitProjectRequest{
			Id:          id,
			ProjectUrl:  attributes["git-project-url"],
			Labels:      strings.Split(attributes["labels"], ","),
			AccessToken: attributes["access-token"],
			UserID:      attributes["user-id"],
		})
	case "cloud-provider":
		var cloudAttributes map[string]string
		cloudAttributes, err = prepareCloudAttributes(attributes)
		if err != nil {
			return err
		}

		_, err = client.UpdateCloudProvider(context.TODO(), &captenpluginspb.UpdateCloudProviderRequest{
			Id:              id,
			CloudType:       attributes["cloud-type"],
			Labels:          strings.Split(attributes["labels"], ","),
			CloudAttributes: cloudAttributes,
		})
	case "container-registry":
		_, err = client.UpdateContainerRegistry(context.TODO(), &captenpluginspb.UpdateContainerRegistryRequest{
			Id:           id,
			RegistryUrl:  attributes["registry-url"],
			Labels:       strings.Split(attributes["labels"], ","),
			RegistryType: attributes["registry-type"],
			RegistryAttributes: map[string]string{
				"registry-username": attributes["username"],
				"registry-password": attributes["password"],
			},
		})
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}

	if err != nil {
		return err
	}

	clog.Logger.Infof("%s resource updated in cluster", resourceType)
	return nil
}

func prepareCloudAttributes(attributes map[string]string) (map[string]string, error) {
	cloudAttributes := map[string]string{}
	switch attributes["cloud-type"] {
	case "azure":
		cloudAttributes["clientId"] = attributes["client-id"]
		cloudAttributes["clientSecret"] = attributes["client-secret"]
	case "aws":
		cloudAttributes["access-key"] = attributes["accessKey"]
		cloudAttributes["secret-key"] = attributes["secretKey"]
	default:
		return nil, fmt.Errorf("invalid cloud type: %s", attributes["cloud-type"])
	}
	return cloudAttributes, nil
}

func RemoveClusterResource(captenConfig config.CaptenConfig, resourceType, id string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	switch resourceType {
	case "git-project":
		_, err = client.DeleteGitProject(context.TODO(), &captenpluginspb.DeleteGitProjectRequest{
			Id: id,
		})
	case "cloud-provider":
		_, err = client.DeleteCloudProvider(context.TODO(), &captenpluginspb.DeleteCloudProviderRequest{
			Id: id,
		})
	case "container-registry":
		_, err = client.DeleteContainerRegistry(context.TODO(), &captenpluginspb.DeleteContainerRegistryRequest{
			Id: id,
		})
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType)
	}

	if err != nil {
		return err
	}

	clog.Logger.Infof("%s resource removed from cluster", resourceType)
	return nil
}
