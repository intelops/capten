package config

import (
	"log"

	"capten/pkg/util"

	"github.com/spf13/viper"
)

const (
	defaultCaptenConfigPath = "config.yaml"
	configEnvKey            = "CAPTEN_CONFIG"
)

// GetClusterConfig config for cluster creation
func GetClusterConfig(configPath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(configPath)
	if err := config.ReadInConfig(); err != nil {
		log.Println("viper cluster config read error", err)
		return nil, err
	}

	return config, nil
}

func GetConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(util.GetEnv(configEnvKey, defaultCaptenConfigPath))
	if err := config.ReadInConfig(); err != nil {
		log.Println("viper cli config read error", err)
		return nil, err
	}

	return config, nil
}
