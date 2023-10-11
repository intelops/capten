package main

import (
	"capten/pkg/agent"
	"capten/pkg/app"
	"capten/pkg/config"

	"github.com/sirupsen/logrus"
)

func main() {
	captenConfig, err := config.GetCaptenConfig()
	if err != nil {
		logrus.Errorf("failed to read capten config, %v", err)
		return
	}
	globalValues, err := app.PrepareGlobalVaules(captenConfig)
	if err != nil {
		logrus.Errorf("applications values preparation failed, %v", err)
		return
	}
	err = agent.StoreCredentials(captenConfig, globalValues)
	if err != nil {
		logrus.Errorf("failed to store, %v", err)
		return
	}
	
	logrus.Info("credential are stored")
}
