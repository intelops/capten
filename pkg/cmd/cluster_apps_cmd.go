package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/app"
	"fmt"
	"time"

	"capten/pkg/cert"
	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/k8s"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func readAppsNameFlags(cmd *cobra.Command) (appsName string, err error) {
	appsName, err = cmd.Flags().GetString("app-name")
	if len(appsName) == 0 {
		return "", fmt.Errorf("specify the name of the apps in the command line %v", err)
	}
	return
}

type SetupAppsActionList struct {
	Actions SetupAppsActions `yaml:"actions"`
}

type SetupAppsActions struct {
	CreateNamespaces        map[string]interface{} `yaml:"create-namespaces"`
	InstallCoreAppGroup     map[string]interface{} `yaml:"install-core-app-group"`
	ConfigureAgentCerts     map[string]interface{} `yaml:"configure-agent-certs"`
	ConfigureSecrets        map[string]interface{} `yaml:"configure-secrets"`
	ConfigureCertIssuer     map[string]interface{} `yaml:"configure-cert-issuer"`
	ConfigureCstorPool      map[string]interface{} `yaml:"configure-cstor-pool"`
	InstallDefaultAppGroup  map[string]interface{} `yaml:"install-default-app-group"`
	SynchApps               map[string]interface{} `yaml:"synch-apps"`
	StoreClusterCredentials map[string]interface{} `yaml:"store-cluster-credentials"`
	FetchLoadBalancerHost   map[string]interface{} `yaml:"fetch-loadBalancerHost"`
}

var appsInstallSubCmd = &cobra.Command{
	Use:   "install",
	Short: "install capten stack apps on cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		globalValues, err := app.PrepareGlobalVaules(captenConfig)
		if err != nil {
			clog.Logger.Errorf("applications values preparation failed, %v", err)
			return
		}

		actions, err := loadSetupAppsActions(captenConfig)
		if err != nil {
			clog.Logger.Errorf("loading setup apps actions failed, %v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.CreateNamespaces, func() error {
			kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
			return k8s.CreateNamespaceIfNotExists(kubeconfigPath, captenConfig.CaptenNamespace)
		})
		if err != nil {
			clog.Logger.Errorf("capten namespace creation failed, %v", err)
		}

		err = execActionIfEnabled(actions.Actions.InstallCoreAppGroup, func() error {
			return app.DeployApps(captenConfig, globalValues, captenConfig.CoreAppGroupsFileName)
		})
		if err != nil {
			clog.Logger.Errorf(" Error while deploying apps %v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.FetchLoadBalancerHost, func() error {
			lbhostName, err := k8s.FetchClusterLoadBalancerHost(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath,
				captenConfig.KubeConfigFileName), "traefik", captenConfig.LBServiceName)
			if err != nil {
				clog.Logger.Error("failed to get LoadBalancerService ", err)
			}

			err = retry(10, 30*time.Second, func() error {
				if err := config.UpdateLBEndpointFile(&captenConfig, lbhostName, ""); err != nil {
					clog.Logger.Infof("LB is not updated in the capten_lb_endpoint.yaml ")
					return errors.WithMessage(err, "failed to update LB ")
				}
				return nil
			})
			if err != nil {
				return err
			}
			clog.Logger.Info("Fetched cluster agent address")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.ConfigureAgentCerts, func() error {
			if err := cert.PrepareCerts(captenConfig); err != nil {
				return errors.WithMessage(err, "failed to generate certificate")
			}
			if err := k8s.CreateOrUpdateCertSecrets(captenConfig); err != nil {
				return errors.WithMessage(err, "failed to create secret for certs")
			}
			clog.Logger.Info("Configured Certificates for Cluster Agent")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.InstallDefaultAppGroup, func() error {
			if err = k8s.CreateOrUpdateClusterIssuer(captenConfig); err != nil {
				return errors.WithMessage(err, "failed to create cstorPoolCluster")
			}
			clog.Logger.Info("Configured Certificate Issuer on Cluster")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.ConfigureCstorPool, func() error {
			err = k8s.CreateCStorPoolClusterWithRetries(captenConfig)
			if err != nil {
				clog.Logger.Errorf("Failed to configure storage pool, %v", err)
				return err
			}
			clog.Logger.Info("Configured storage pool")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.StoreClusterCredentials, func() error {
			clog.Logger.Info("Storing credentails on cluster")
			err = retry(10, 30*time.Second, func() error {
				err = agent.StoreCredentials(captenConfig, globalValues)
				if err != nil {
					clog.Logger.Infof("Vault is not ready")
					return errors.WithMessage(err, "failed to store credentials")
				}
				if captenConfig.CloudService == "aws" {
					err = agent.StoreClusterCredentials(captenConfig, globalValues)
					if err != nil {
						return errors.WithMessage(err, "failed to store cluster credentials")
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
			clog.Logger.Info("Stored credentails on cluster")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.InstallDefaultAppGroup, func() error {
			return app.DeployApps(captenConfig, globalValues, captenConfig.DefaultAppGroupsFileName)
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.FetchLoadBalancerHost, func() error {

			natslbhostname, err := k8s.FetchClusterLoadBalancerHost(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath,
				captenConfig.KubeConfigFileName), "observability", captenConfig.NatsLBServiceName)

			if err != nil {
				clog.Logger.Error("failed to get NatsLoadBalancerService ", err)
			}

			err = retry(10, 30*time.Second, func() error {
				if err := config.UpdateLBEndpointFile(&captenConfig, "", natslbhostname); err != nil {
					clog.Logger.Infof("LB is not updated in the capten_lb_endpoint.yaml ")
					return errors.WithMessage(err, "failed to update LB ")
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = agent.StoreCredentials(captenConfig, globalValues)
			if err != nil {

				return errors.WithMessage(err, "failed to store lbip credentials")
			}
			clog.Logger.Info("Fetched nats loadbalancer host and stored in vault")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

		err = execActionIfEnabled(actions.Actions.SynchApps, func() error {
			clog.Logger.Info("Synchonizing Applications with Cluster Agent")
			err = retry(12, 30*time.Second, func() error {
				if err := agent.SyncInstalledAppConfigsOnAgent(captenConfig); err != nil {
					clog.Logger.Infof("Capten Agent is not ready")
					return errors.WithMessage(err, "failed to sync installed apps config in cluster")
				}
				return nil
			})
			if err != nil {
				return err
			}
			clog.Logger.Info("Applications Synchonized with Cluster Agent")
			return nil
		})
		if err != nil {
			clog.Logger.Errorf("%v", err)
			return
		}

	},
}

var appsListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "list deployed apps on cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}
		err = agent.ListClusterApplications(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to fetch applications from capten cluster, %v", err)
		}
	},
}

var appsShowSubCmd = &cobra.Command{
	Use:   "show",
	Short: "show deployed apps on cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		appsName, err := readAppsNameFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.ShowClusterAppData(captenConfig, appsName)
		if err != nil {
			clog.Logger.Errorf("failed to fetch application from capten cluster, %v", err)
		}
	},
}

func loadSetupAppsActions(captenConfig config.CaptenConfig) (*SetupAppsActionList, error) {
	actionsFile := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.SetupAppsConfigFile)
	yamlFile, err := os.ReadFile(actionsFile)
	if err != nil {
		return nil, err
	}

	var actions SetupAppsActionList
	err = yaml.Unmarshal(yamlFile, &actions)
	if err != nil {
		return nil, err
	}
	return &actions, err
}

func isEnabled(actionConfig map[string]interface{}) bool {
	enabledVal, ok := actionConfig["enabled"]
	if !ok {
		return false
	}
	enabled, ok := enabledVal.(bool)
	if !ok {
		return false
	}
	return enabled
}

func execActionIfEnabled(actionConfig map[string]interface{}, f func() error) error {
	if isEnabled(actionConfig) {
		return f()
	}
	return nil
}

func retry(retries int, interval time.Duration, f func() error) (err error) {
	for i := 0; i <= retries; i++ {
		if err = f(); err == nil {
			return nil
		}
		time.Sleep(interval)
	}
	return
}
