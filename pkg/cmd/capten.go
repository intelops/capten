package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	vaultcredclient "github.com/intelops/go-common/vault-cred-client"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"capten/pkg/cert"
	"capten/pkg/cluster"
	"capten/pkg/config"
	"capten/pkg/helm"
	"capten/pkg/k8s"
	"capten/pkg/types"
)

type CLIFormatter struct {
}

func (f *CLIFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor *color.Color
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed, color.Bold)
	default:
		levelColor = color.New()
	}
	message := fmt.Sprintf("[%s] %s\n", levelColor.Sprint(strings.ToUpper(entry.Level.String())), entry.Message)
	return []byte(message), nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "capten",
	Short: "",
	Long:  `command line tool for building cluster`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logrus.SetFormatter(&CLIFormatter{})
	cobra.CheckErr(rootCmd.Execute())
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "for creation of resources or cluster",
	Long:  ``,
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy created cluster",
	Long:  ``,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "sets up cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//h := helm.NewHelm()
	},
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "sets up apps cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}

		if err := cert.GenerateCerts(captenConfig); err != nil {
			logrus.Errorf("failed to generate certificate, %v", err)
			return
		}
		logrus.Info("Generated Certificates")

		if err := k8s.CreateOrUpdateAgnetCertSecret(captenConfig); err != nil {
			logrus.Errorf("failed to patch namespace with privilege, %v", err)
			return
		}
		logrus.Info("Configured Certificates on Capten Cluster")

		hc, err := helm.NewClient(captenConfig)
		if err != nil {
			logrus.Errorf("applications installation failed, %v", err)
			return
		}
		err = hc.PrepareAppValues()
		if err != nil {
			logrus.Errorf("applications installation failed, %v", err)
			return
		}
		hc.Install()

		//push kubeconfig and bucket credential to cluster
		configContent, err := ioutil.ReadFile(captenConfig.ConfigDirPath + "/" + captenConfig.KubeConfigFileName)
		if err != nil {
			logrus.Error("error while reading kube config file", err)
			return
		}
		credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
		if err != nil {
			logrus.Error("error in initializing vault credential client", err)
			return
		}
		ctx := context.Background()
		err = credAdmin.PutGenericCredential(ctx, "k8s", "kubeconfig", vaultcredclient.GerericCredentail{
			Credential: map[string]string{
				"kubeconfig": string(configContent),
			},
		})
		if err != nil {
			logrus.Error("error in adding kubeconfig to vault", err)
			return
		}

		awsConfigByte, err := ioutil.ReadFile(captenConfig.ConfigDirPath + "/aws_config.yaml")
		if err != nil {
			logrus.Error("Error reading aws config YAML file", err)
			return
		}
		var awsconfig types.AWSConfig
		err = yaml.Unmarshal(awsConfigByte, &awsconfig)
		if err != nil {
			logrus.Error("Error unmarshaling AWS YAML config file", err)
		}

		err = credAdmin.PutGenericCredential(ctx, "bucket", "terraform-state", vaultcredclient.GerericCredentail{
			Credential: map[string]string{
				"bucketName": awsconfig.TerraformBackendConfigs[0],
				"awsKey":     awsconfig.AwsAccessKey,
				"awsSecrete": awsconfig.AwsSecretKey,
			},
		})
		if err != nil {
			logrus.Error(err, "error in adding bucket credentials to vault", err)
			return
		}

		//push the app config to cluster
		//prepare agent proto to push app config
		//agent store data on cassandra
		logrus.Info("Default Applications Installed")
	},
}

var clusterDestroySubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster destroy operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		clusterType, cloudType, err := readAndValidClusterFlags(cmd)
		if err != nil {
			logrus.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}
		err = cluster.Destroy(captenConfig, clusterType, cloudType)
		if err != nil {
			logrus.Errorf("failed to destroy cluster, %v", err)
			return
		}
		logrus.Info("Cluster Destroyed")
	},
}

var clusterCreateSubCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster create operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		clusterType, cloudType, err := readAndValidClusterFlags(cmd)
		if err != nil {
			logrus.Error(err)
			return
		}

		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Errorf("failed to read capten config, %v", err)
			return
		}
		err = cluster.Create(captenConfig, clusterType, cloudType)
		if err != nil {
			logrus.Errorf("failed to create cluster, %v", err)
			return
		}
		logrus.Info("Cluster Created")
	},
}

func readAndValidClusterFlags(cmd *cobra.Command) (clusterType string, cloudType string, err error) {
	clusterType, _ = cmd.Flags().GetString("type")
	if len(clusterType) == 0 {
		clusterType = "k3s"
	}
	if clusterType != "k3s" {
		err = fmt.Errorf("cluster type '%s' is not supported, supported types: k3s", clusterType)
		return
	}

	cloudType, _ = cmd.Flags().GetString("cloud")
	if len(cloudType) == 0 {
		cloudType = "aws"
	}
	if cloudType != "aws" {
		err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudType)
		return
	}
	return
}

var showClusterInfoCmd = &cobra.Command{
	Use:   "clusterinfo",
	Short: "show the clinster info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			logrus.Error("failed to read capten config", err)
			return
		}
		fmt.Println("Agetnt hostname:", "agent."+captenConfig.DomainName)
		fmt.Println("Agent IP:", "0.0.0.127")
	},
}

func init() {
	clusterCreateSubCmd.PersistentFlags().String("cloud", "", "cloud service (default: aws)")
	clusterDestroySubCmd.PersistentFlags().String("cloud", "", "cloud service (default: aws)")
	clusterCreateSubCmd.PersistentFlags().String("type", "", "type of cluster (default: k3s)")
	clusterDestroySubCmd.PersistentFlags().String("type", "", "type of cluster (default: k3s)")

	createCmd.AddCommand(clusterCreateSubCmd)
	destroyCmd.AddCommand(clusterDestroySubCmd)
	setupCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(registerAgentCmd)
	rootCmd.AddCommand(showClusterInfoCmd)
}
