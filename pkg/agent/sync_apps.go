package agent

import (
	"bytes"
	"capten/pkg/agent/agentpb"
	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/types"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func SyncInstalledAppConfigsOnAgent(captenConfig config.CaptenConfig) error {
	client, err := GetAgentClient(captenConfig)
	if err != nil {
		return err
	}

	appConfigs, err := readInstalledAppConfigs(captenConfig)
	if err != nil {
		return err
	}

	for _, appConfig := range appConfigs {
		syncAppData, err := appConfig.ToSyncAppData()
		if err != nil {
			clog.Logger.Errorf("failed to parse '%s' app config, %v", appConfig.ReleaseName, err)
			continue
		}

		if len(syncAppData.Config.Icon) != 0 {
			iconBytes, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.AppIconsDirPath, string(syncAppData.Config.Icon)))
			if err != nil {
				clog.Logger.Errorf("failed loading icon for app '%s', %v", appConfig.ReleaseName, err)
			}
			syncAppData.Config.Icon = iconBytes
			clog.Logger.Debugf("'%s' app icon added", appConfig.ReleaseName)
		}

		syncAppData.Config.InstallStatus = "Installed"
		res, err := client.SyncApp(context.TODO(), &agentpb.SyncAppRequest{Data: &syncAppData})
		if err != nil {
			clog.Logger.Errorf("failed to synch '%s' app config to synch with agent, %v", appConfig.ReleaseName, err)
			continue
		}

		if res != nil && res.Status != agentpb.StatusCode_OK {
			clog.Logger.Errorf("failed to synch '%s' app config to synch with agent, %v", appConfig.ReleaseName, res.GetStatusMessage())
			continue
		}
		clog.Logger.Debugf("'%s' app synchronized with agent", appConfig.ReleaseName)
	}
	return nil
}

func readInstalledAppConfigs(config config.CaptenConfig) (ret []types.AppConfig, err error) {
	configDir := config.PrepareDirPath(config.AppsTempDirPath)
	err = filepath.Walk(configDir, func(appConfigFilePath string, info os.FileInfo, ferr error) error {
		if ferr != nil {
			return errors.Wrapf(ferr, "in file %s", appConfigFilePath)
		}

		if info.IsDir() || (filepath.Dir(appConfigFilePath) != strings.TrimRight(configDir, "/")) ||
			!strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}

		data, err := os.ReadFile(appConfigFilePath)
		if err != nil {
			return errors.Wrapf(err, "in file: %s", appConfigFilePath)
		}

		var appConfig types.AppConfig
		if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(&appConfig); err != nil {
			return errors.Wrapf(err, "in file %s", appConfigFilePath)
		}

		ret = append(ret, appConfig)
		return nil
	})

	return
}
