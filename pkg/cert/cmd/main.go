package main

import (
	"capten/pkg/cert"
	"capten/pkg/config"

	"github.com/sirupsen/logrus"
)

func main() {
	captenConfig, err := config.GetCaptenConfig()
	if err != nil {
		logrus.Error("failed to read capten config", err)
		return
	}
	if err := cert.GenerateCerts(captenConfig); err != nil {
		logrus.Errorf("failed to generate certificate. Error - %v", err)
		return
	}
	logrus.Info("Generated Certificates")
}
