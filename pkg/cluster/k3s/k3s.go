package k3s

import (
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"capten/pkg/cluster/types"
	"capten/pkg/terraform"

	"github.com/spf13/viper"
)

type ClusterInfo struct {
	ClusterType           string
	AwsAccessKey          string
	AwsSecretKey          string
	AlbName               string
	PrivateSubnet         string
	Region                string
	SecurityGroupName     string
	VpcCidr               string
	VpcName               string
	InstanceType          string
	NodeMonitoringEnabled string
	MasterCount           string
	WorkerCount           string
	TraefikHttpPort       string
	TraefikHttpsPort      string
	TalosTg               string
	TraefikTg80Name       string
	TraefikTg443Name      string
	TraefikLbName         string
}

func Create(config *viper.Viper, workDir string) error {
	clusterInfo := ClusterInfo{
		ClusterType:           "k3s",
		AwsAccessKey:          config.GetString(types.AwsAccessKey),
		AwsSecretKey:          config.GetString(types.AwsSecretKey),
		AlbName:               config.GetString(types.AlbName),
		PrivateSubnet:         config.GetString(types.PrivateSubnet),
		Region:                config.GetString(types.Region),
		SecurityGroupName:     config.GetString(types.SecurityGroupName),
		VpcCidr:               config.GetString(types.VpcCidr),
		VpcName:               config.GetString(types.VpcName),
		InstanceType:          config.GetString(types.InstanceType),
		NodeMonitoringEnabled: config.GetString(types.NodeMonitoringEnabled),
		MasterCount:           config.GetString(types.MasterCount),
		WorkerCount:           config.GetString(types.WorkerCount),
		TraefikHttpPort:       config.GetString(types.TraefikHttpPort),
		TraefikHttpsPort:      config.GetString(types.TraefikHttpsPort),
		TalosTg:               config.GetString(types.TalosTg),
		TraefikTg80Name:       config.GetString(types.TraefikTg80Name),
		TraefikTg443Name:      config.GetString(types.TraefikTg443Name),
		TraefikLbName:         config.GetString(types.TraefikLbName),
	}

	content, err := os.ReadFile("./templates/k3s/values.tfvars.tmpl")
	if err != nil {
		log.Printf("failed to read %s template file\n", clusterInfo.ClusterType)
		return err
	}

	contentStr := string(content)
	templateObj, err := template.New(clusterInfo.ClusterType).Parse(contentStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	templateFile, err := os.Create(path.Join(workDir, "values.tfvars"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := templateObj.Execute(templateFile, clusterInfo); err != nil {
		fmt.Println(err)
		return err
	}

	tf, err := terraform.New(config, workDir)
	if err != nil {
		log.Println("failed to initialise the terraform", err)
		return err
	}

	return tf.Apply()
}

func Destroy(config *viper.Viper, workDir string) error {
	tf, err := terraform.New(config, workDir)
	if err != nil {
		log.Println("failed to initialise the terraform", err)
		return err
	}

	return tf.Destroy()
}
