package agent

import (
	"capten/pkg/agent/pb/captenpluginspb"
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func ConfigureClusterPlugin(captenConfig config.CaptenConfig, pluginName, action string,
	actionAttributes map[string]string) error {

	switch pluginName {
	case "crossplane":
		return configureCrossplanePlugin(captenConfig, action, actionAttributes)
	case "tekton":
		return configureTektonPlugin(captenConfig, action)
	case "proact":
		return fmt.Errorf("configure actions for plugin is not implemented yet")
	default:
		return fmt.Errorf("no configure actions for plugin supported")
	}
}

func configureCrossplanePlugin(captenConfig config.CaptenConfig, action string,
	actionAttributes map[string]string) error {
	switch action {
	case "list-actions":
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Action", "Attributes"})
		table.SetRowLine(true)

		table.Append([]string{"show-crossplane-project", ""})
		table.Append([]string{"synch-crossplane-project", ""})
		table.Append([]string{"create-crossplane-provider", "cloud-type, cloud-provider-id"})
		table.Append([]string{"update-crossplane-provider", "crossplane-provider-id, cloud-type, cloud-provider-id"})
		table.Append([]string{"delete-crossplane-provider", "crossplane-provider-id"})
		table.Append([]string{"list-crossplane-providers", ""})
		table.Append([]string{"list-managed-clusters", ""})
		table.Append([]string{"download-kubeconfig", "managed-cluster-id"})
		table.SetAutoMergeCells(true)
		table.Render()
	case "show-crossplane-project":
		return showCrossplaneProject(captenConfig)
	case "synch-crossplane-project":
		return synchCrossplaneProject(captenConfig)
	case "create-crossplane-provider":
		return createCrossplaneProvider(captenConfig, actionAttributes)
	case "update-crossplane-provider":
		return updateCrossplaneProvider(captenConfig, actionAttributes)
	case "delete-crossplane-provider":
		return deleteCrossplaneProvider(captenConfig, actionAttributes)
	case "list-crossplane-providers":
		return listCrossplaneProviders(captenConfig)
	case "list-managed-clusters":
		return listManagedClusters(captenConfig)
	case "download-kubeconfig":
		return downloadKubeconfig(captenConfig, actionAttributes)
	default:
		return fmt.Errorf("action is not supported for plugin")
	}
	return nil
}

func configureTektonPlugin(captenConfig config.CaptenConfig, action string) error {
	switch action {
	case "list-actions":
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Action", "Attributes"})
		table.SetRowLine(true)

		table.Append([]string{"show-tekton-project", ""})
		table.Append([]string{"synch-tekton-project", ""})
		table.SetAutoMergeCells(true)
		table.Render()
	case "show-tekton-project":
		return showTektonProject(captenConfig)
	case "synch-tekton-project":
		return synchTektonProject(captenConfig)
	default:
		return fmt.Errorf("action is not supported for plugin")
	}
	return nil
}

func showCrossplaneProject(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetCrossplaneProject(context.TODO(), &captenpluginspb.GetCrossplaneProjectsRequest{})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Attribute", "Value"})
	table.Append([]string{"git-project-url", resp.Project.GitProjectUrl})
	table.Append([]string{"status", resp.Project.Status})
	table.Render()
	return nil
}

func synchCrossplaneProject(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.RegisterCrossplaneProject(context.TODO(), &captenpluginspb.RegisterCrossplaneProjectRequest{})
	if err != nil {
		return err
	}
	clog.Logger.Info("crossplane project synched")
	return nil
}

func createCrossplaneProvider(captenConfig config.CaptenConfig, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.AddCrossplanProvider(context.TODO(), &captenpluginspb.AddCrossplanProviderRequest{
		CloudType:       attributes["cloud-type"],
		CloudProviderId: attributes["cloud-provider-id"],
	})
	if err != nil {
		return err
	}
	clog.Logger.Info("crossplane provider created")
	return nil
}

func updateCrossplaneProvider(captenConfig config.CaptenConfig, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.UpdateCrossplanProvider(context.TODO(), &captenpluginspb.UpdateCrossplanProviderRequest{
		Id:              attributes["crossplane-provider-id"],
		CloudType:       attributes["cloud-type"],
		CloudProviderId: attributes["cloud-provider-id"],
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteCrossplaneProvider(captenConfig config.CaptenConfig, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.DeleteCrossplanProvider(context.TODO(), &captenpluginspb.DeleteCrossplanProviderRequest{
		Id: attributes["crossplane-provider-id"],
	})
	if err != nil {
		return err
	}
	return nil
}

func listCrossplaneProviders(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetCrossplanProviders(context.TODO(), &captenpluginspb.GetCrossplanProvidersRequest{})
	if err != nil {
		return err
	}

	if len(resp.Providers) == 0 {
		clog.Logger.Info("No crossplane providers added to cluster")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Cloud Type", "Cloud Provider ID", "Status"})
	for _, provider := range resp.Providers {
		table.Append([]string{provider.Id, provider.CloudType, provider.CloudProviderId, provider.Status})
	}
	table.Render()
	return nil
}

func listManagedClusters(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetManagedClusters(context.TODO(), &captenpluginspb.GetManagedClustersRequest{})
	if err != nil {
		return err
	}

	if len(resp.Clusters) == 0 {
		clog.Logger.Info("No managed clusters added to cluster")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Cluster Name", "Cluster Endpoint", "Cluster Deploy Status"})
	for _, cluster := range resp.Clusters {
		table.Append([]string{cluster.Id, cluster.ClusterName, cluster.ClusterEndpoint, cluster.ClusterDeployStatus})
	}
	table.Render()
	return nil
}

func downloadKubeconfig(captenConfig config.CaptenConfig, attributes map[string]string) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetManagedClusterKubeconfig(context.TODO(), &captenpluginspb.GetManagedClusterKubeconfigRequest{
		Id: attributes["managed-cluster-id"],
	})
	if err != nil {
		return err
	}

	fileName := "kubeconfig-" + attributes["managed-cluster-id"] + ".yaml"
	err = os.WriteFile(fileName, []byte(resp.Kubeconfig), 0644)
	if err != nil {
		return err
	}

	clog.Logger.Infof("kubeconfig downloaded to ./%s", fileName)
	return nil
}

func showTektonProject(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetTektonProject(context.TODO(), &captenpluginspb.GetTektonProjectRequest{})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Attribute", "Value"})
	table.Append([]string{"git-project-url", resp.Project.GitProjectUrl})
	table.Append([]string{"status", resp.Project.Status})
	table.Render()
	return nil
}

func synchTektonProject(captenConfig config.CaptenConfig) error {
	client, err := GetCaptenPluginClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.RegisterTektonProject(context.TODO(), &captenpluginspb.RegisterTektonProjectRequest{})
	if err != nil {
		return err
	}
	clog.Logger.Info("tekton project synched")
	return nil
}
