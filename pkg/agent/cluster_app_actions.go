package agent

import (
	"capten/pkg/agent/pb/agentpb"
	"capten/pkg/clog"
	"capten/pkg/config"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func ListClusterApplications(captenConfig config.CaptenConfig) error {
	client, err := GetAgentClient(captenConfig)
	if err != nil {
		return err
	}

	resp, err := client.GetDefaultAppsStatus(context.TODO(), &agentpb.GetDefaultAppsStatusRequest{})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Name", "Version", "Status"})
	for _, clusterApp := range resp.DefaultAppsStatus {
		table.Append([]string{clusterApp.Category, clusterApp.AppName, clusterApp.Version, clusterApp.InstallStatus})
	}
	table.Render()
	return nil
}

func DeployApps(captenConfig config.CaptenConfig) error {
	client, err := GetAgentClient(captenConfig)
	if err != nil {
		return err
	}

	_, err = client.DeployDefaultApps(context.TODO(), &agentpb.DeployDefaultAppsRequest{})
	if err != nil {
		return err
	}
	return nil
}

func WaitAndTrackDefaultAppsDeploymentStatus(captenConfig config.CaptenConfig) {
	timeout := time.After(1 * time.Hour)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			clog.Logger.Errorf("Default Apps deployment timed out")
			return
		case <-ticker.C:
			completed, status, err := GetDefaultAppsDeploymentStatus(captenConfig)
			if err != nil {
				clog.Logger.Errorf("failed to get default apps deployment status, %v", err)
			} else {
				clog.Logger.Infof("%s", status)
				if completed {
					return
				}
			}
		}
	}
}

func GetDefaultAppsDeploymentStatus(captenConfig config.CaptenConfig) (bool, string, error) {
	var completed bool
	var defaultAppsDeploymentStatus string
	client, err := GetAgentClient(captenConfig)
	if err != nil {
		return completed, defaultAppsDeploymentStatus, err
	}

	resp, err := client.GetDefaultAppsStatus(context.TODO(), &agentpb.GetDefaultAppsStatusRequest{})
	if err != nil {
		return completed, defaultAppsDeploymentStatus, err
	}

	deployedApps := []string{}
	failedApps := []string{}
	ongoingApps := []string{}
	for _, appStatus := range resp.DefaultAppsStatus {
		if appStatus.InstallStatus == "Installed" {
			deployedApps = append(deployedApps, appStatus.AppName)
		} else if appStatus.InstallStatus == "Installation Failed" {
			failedApps = append(failedApps, appStatus.AppName)
		} else {
			ongoingApps = append(ongoingApps, appStatus.AppName)
		}
	}

	switch resp.DeploymentStatus {
	case agentpb.DeploymentStatus_ONGOING:
		if len(failedApps) > 0 {
			defaultAppsDeploymentStatus = fmt.Sprintf("Deploying applications, %d/%d deployed, %d failed", len(deployedApps), len(resp.DefaultAppsStatus), len(failedApps))
		} else {
			defaultAppsDeploymentStatus = fmt.Sprintf("Deploying applications, %d/%d deployed", len(deployedApps), len(resp.DefaultAppsStatus))
		}
	case agentpb.DeploymentStatus_SUCCESS, agentpb.DeploymentStatus_FAILED:
		if len(failedApps) > 0 {
			defaultAppsDeploymentStatus = fmt.Sprintf("Deployed applications, %d/%d deployed, %d failed", len(deployedApps), len(resp.DefaultAppsStatus), len(failedApps))
		} else {
			defaultAppsDeploymentStatus = fmt.Sprintf("Deployed applications, %d/%d deployed", len(deployedApps), len(resp.DefaultAppsStatus))
		}
		if len(deployedApps) > 0 {
			defaultAppsDeploymentStatus = defaultAppsDeploymentStatus + fmt.Sprintf("\nDeployed Apps: %v", deployedApps)
		}
		if len(failedApps) > 0 {
			defaultAppsDeploymentStatus = defaultAppsDeploymentStatus + fmt.Sprintf("\nFailed Apps: %v", failedApps)
		}
		if len(ongoingApps) > 0 {
			defaultAppsDeploymentStatus = defaultAppsDeploymentStatus + fmt.Sprintf("\nOngoing Apps: %v", ongoingApps)
		}
		completed = true
	default:
		defaultAppsDeploymentStatus = fmt.Sprintf("Deploying applications, %d/%d deployed", len(deployedApps), len(resp.DefaultAppsStatus))
	}
	return completed, defaultAppsDeploymentStatus, nil
}
