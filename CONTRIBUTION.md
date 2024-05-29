# Contribution Guidelines
Please read this guide if you plan to contribute to the Capten. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with the Capten, please [open an issue on GitHub](https://github.com/intelops/capten/issues/new?labels=kind%2Fbug&template=bug-report.md&title=Bug:) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/intelops/capten/issues/new?labels=kind%2Ffeature&template=feature-request.md&title=Feature%20Request:).


## Developing 

Development can be conducted using  GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode).

For contributing Capten,First you need to understand the folder structure.Kindly refer below to understand he folder structure

This document deals with the detailed description on how to contribute capten

### Directory Structure 


```
capten/
│
├── apps/                   
│   ├── conf/
│   │   ├── credentials/         # Config files for generating and storing app credentials.
│   │   ├── values/              # Helm chart values for applications.
│   ├── icons/                   # Icons for applications.
│   ├── tmp/                     # Temporary folder created during app deployment with app values.
│   └── core_group_apps.yaml     # YAML file for core group application configurations.
│   └── default_group_apps.yaml  # YAML file for default group application configurations.
│
├── cert/                        # Public certificates and assets.
│
├── config/
│   ├── aws_config.yaml/        # Configuration for AWS cluster.
│   ├── azure_config.yaml/      # Configuration for Azure cluster.
│   ├── capten-lb-endpoint.yaml # Capten load balancer endpoint.
│   └── capten.yaml             # Main Capten configuration file.
│   └── setup_apps.yaml         # Setup configuration for applications.
│
├── cmd/                        # Main entry point for command line commands.
│
├── pkg/                        # Package directory for various components and services.
│   ├── agent/                  # Code related to the agent component.
│   │   ├── pb/                 # Protocol buffer files and generated code.
│   └── app/                    # Application-related code.
│   ├── cert/                   # Certificate management code.
│   ├── clog/                   # Custom logging utilities.
│   └── cluster/                # Cluster management code.
│       ├── k3s/                # K3s specific configurations and code.
│   ├── cmd/                    # Command-related code.
│   └── config/                 # Configuration management code.
│   └── helm/                   # Helm chart management.
│   ├── k8s/                    # Kubernetes-specific utilities and configurations.
│   └── terraform/              # Terraform configurations and modules.
│   └── types/                  # Common types used across the project.
│
├── templates/                  # Template files for various configurations.
│   ├── values.aws.tmpl         # Template for AWS-specific Helm chart values.
│   └── values.azure.tmpl/      # Template for Azure-specific Helm chart values.
│   ├── values.tfvars/          # Template for Terraform variable files.
│
├── README.md                   # Project readme file.
└── .gitignore                  # Git ignore file to exclude specified files and directories from version control.

```

## How to Contribute Capten

You can generally contribute capten in 4 ways,that is given below:

1. Managing cluster creation

2. Supporting additional tools to be deployed in control plane cluster

3. Enhancing CLI

4. Additional Plugin Support


## Managing cluster creation

For supporting capten with any additional cluster or to enhance additional support to the existing cluster,you can follow the below steps:

1. Add the cloud configurations yaml file in `./config` directory.you can refer [aws_config.yaml](https://github.com/intelops/capten/blob/main/config/aws_config.yaml) 

2. Modify the below function for supporting the specified cloud service.You can check this function [here](https://github.com/intelops/capten/blob/main/pkg/cmd/capten.go) 

```bash
func validateClusterFlags(cloudService, clusterType string) (err error) {

	if cloudService != "aws" && cloudService != "azure" {
		err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudService)
		return
	}

	if clusterType != "talos" {
		err = fmt.Errorf("cluster type '%s' is not supported, supported types: talos", clusterType)
		return
	}
	return
}
```
3. Add any configurations if needed in ./config/config.go

4. Also for supporting additional cluster,you can add structure in ./pkg/types/types.go 

For example,awsclusterInfo struct is given below

```bash
type AWSClusterInfo struct {
	ConfigFolderPath        string   `yaml:"ConfigFolderPath"`
	TerraformModulesDirPath string   `yaml:"TerraformModulesDirPath"`
	CloudService            string   `yaml:"CloudService"`
	ClusterType             string   `yaml:"ClusterType"`
	AwsAccessKey            string   `yaml:"AwsAccessKey"`
	AwsSecretKey            string   `yaml:"AwsSecretKey"`
	AlbName                 string   `yaml:"AlbName"`
	PrivateSubnet           string   `yaml:"PrivateSubnet"`
	Region                  string   `yaml:"Region"`
	SecurityGroupName       string   `yaml:"SecurityGroupName"`
	VpcCidr                 string   `yaml:"VpcCidr"`
	VpcName                 string   `yaml:"VpcName"`
	InstanceType            string   `yaml:"InstanceType"`
	NodeMonitoringEnabled   string   `yaml:"NodeMonitoringEnabled"`
	MasterCount             string   `yaml:"MasterCount"`
	WorkerCount             string   `yaml:"WorkerCount"`
	TraefikHttpPort         string   `yaml:"TraefikHttpPort"`
	TraefikHttpsPort        string   `yaml:"TraefikHttpsPort"`
	TalosTg                 string   `yaml:"TalosTg"`
	TraefikTg80Name         string   `yaml:"TraefikTg80Name"`
	TraefikTg443Name        string   `yaml:"TraefikTg443Name"`
	TraefikLbName           string   `yaml:"TraefikLbName"`
	TerraformBackendConfigs []string `yaml:"TerraformBackendConfigs"`
	Nats_client_port        string   `yaml:"Nats_client_port"`
	Nats_tg_4222_name       string   `yaml:"Nats_tg_4222_name"`
}

```
5. For supporting additional cloud,kindly go through and understand the code in `./pkg/cluster`. 

Refer the below code snippets and modify the code accordingly for supporting additional cloud

(https://github.com/intelops/capten/blob/main/pkg/cluster/cluster.go)
(https://github.com/intelops/capten/blob/main/pkg/cluster/k3s/k3s.go)


6. Add the template file in required format in `./templates/k3s` directory.You can refer sample template file in [./templates/k3s/values.aws.tmpl](https://github.com/intelops/capten/blob/main/templates/k3s/values.aws.tmpl)

7. Contribute your terraform code  for supporting additional cloud in [controlplane-dataplane repo](https://github.com/kube-tarian/controlplane-dataplane). 

8. Terraform is packaged in the Capten Artifact.So for triggering the terraform  create a separate file in `./pkg/terraform` directory and add the code logic for supporting additional cloud by using terrform  go package `github.com/hashicorp/terraform-exec/tfexec`.

For reference,you can understand the code in [./pkg/terraform/terraform-aws.go](https://github.com/intelops/capten/blob/main/pkg/terraform/terraform-aws.go) 


### How to test the changes:

To test your modification,you can just build the CLI artifact with the below command

First Navigate to the Capten directory and run the below command

```sh

make build.release-linux
cd capten
```
Then with the build binary,you can test your changes

For example for creating cluster with the provided cloud-type,you refer the below command

```sh
./capten cluster create --cloud=<coud-type> --type=talos
```

## Supporting additional tools to be deployed in control plane cluster

For bringing up any additional apps or tools in control plane cluster,you can work on `./apps` directory.

1. create a yaml file `./apps/conf` for the additional tool,with the below specifications


```sh

Name: "name of the application"
ChartName: "The path or name of the Helm chart within the repository"
Category: "Helps to group and identify the application type (e.g., Security)"
RepoName: "The name of the repository where the Helm chart is located"
RepoURL: "The URL of the Helm repository from which the chart will be fetched"
Namespace: The Kubernetes namespace in which the Helm release will be deployed. 
ReleaseName: "The name given to this Helm release, uniquely identifying the release within the namespace"
Version: " Specifies the version of the Helm chart to use"
CreateNamespace: "A boolean value that determines if the namespace should be created if it doesn't exist"
OverrideValues: (Optional)"The OverrideValues section is used to specify configuration values that can be dynamically replaced at runtime or during the deployment of an application"

```

Refer the [sample application](https://github.com/intelops/capten/blob/main/apps/conf/falco.yaml)

2. For Passing the values in the application ,create a `_template.yaml` file in `./apps/conf/values`

You can refer [here](https://github.com/intelops/capten/blob/main/apps/conf/values/falco_template.yaml)


3. Then add the application name in the  `./apps/conf/default_group_apps.yaml`.Capten Supports two varieties of cluster Type 
  
     1. Cloud Managed-If the user has predefined inbuilt cluster like eks or aks
	 2. Talos- Capten supports cluster creation of cluster type talos

  So based on the cluster type,add the applications in `./apps/conf/default_group_apps.yaml`

4. If any app credentials needs to be stored in vault or any external secret that needs to be created with the credentials in the vault or any override values such as secret-name that needs to be passed dynamically (similar to [kubviz-client](https://github.com/intelops/capten/blob/main/apps/conf/kubviz-client.yaml]) )to the application, you can refer the code in `./pkg/agent/store_cred.go` .

For storing any application based credentials,create a yaml file in `./apps/conf/credentials`

```sh
name: This field specifies a human-readable name for the credential configuration.

secretName: This field indicates the name of the Kubernetes secret where the credentials will be stored.

namespaces: This field lists the Kubernetes namespaces where the secret will be available. 

credentialEntity: This field denotes the entity that the credential is associated with. 

credentialIdentifier: This field specifies the identifier used to access the credential within the secret. 

credentialType: This field describes the type of credential being used.

```
For sample reference,[ClickHere](https://github.com/intelops/capten/blob/main/apps/conf/credentials/nats-cred.yaml)

5. For contributing any changes related to app deployment,Go through and understand the code in `./pkg/app` and `./pkg/helm`

### How to test

1. Build the binary with your code changes,with the below command

```sh
make build
```

2. Run the below command,your application will be deployed

```sh
./capten cluster apps install
```

### CLI enhancemenet

For enhancing CLI,you can refer and understand the code in `./pkg/cmd` 

Adding CLI for cluster onboarding git project is shown in the below code snippet.

```sh
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster operations",
	Long:  ``,
}
```

```sh
var clusterResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "cluster resources operations",
	Long:  ``,
}
```

```sh

func init() {
		rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterResourcesCmd)
	resourceCreateSubCmd.PersistentFlags().String("resource-type", "", "type of resource")
	resourceCreateSubCmd.PersistentFlags().String("git-project-url", "", "url of git project resource")

	resourceCreateSubCmd.PersistentFlags().String("git-project-url", "", "url of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("access-token", "", "access token of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("user-id", "", "user id of git project resource")
	resourceCreateSubCmd.PersistentFlags().String("labels", "", "labels of resource (e.g. 'crossplane,tekton')")


}
```
Additional CLI commands can be added in [./pkg/cmd/capten.go](https://github.com/intelops/capten/blob/main/pkg/cmd/capten.go)


```sh
var resourceCreateSubCmd = &cobra.Command{
	Use:   "create",
	Short: "cluster resource create",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resourceType, attributes, err := readAndValidCreateResourceFlags(cmd)
		if err != nil {
			clog.Logger.Error(err)
			return
		}
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		err = agent.AddClusterResource(captenConfig, resourceType, attributes)
		if err != nil {
			clog.Logger.Errorf("failed to create cluster resource, %v", err)
		}
	},
}
```
The CLI commands are executed in a separate file,you can refer in [./pkg/cmd/cluster_resource_cmd.go](https://github.com/intelops/capten/blob/main/pkg/cmd/cluster_resource_cmd.go) to understand the code.


In the above mentioned ways,you can contribute for Capten CLI enhancement.

### How to test

1. Build the binary with your code changes,with the below command

```sh
make build
```

2. Run your CLI command. In the above example,running below CLI command will onboard the git project successfully

```sh
./capten1 cluster resources create --resource-type="git-project" --access-token="abcd" --user-id="xxx" --labels="tekton" --git-project-url="sample-repo-url"

```

## Additional Plugin Support

For supporting additional plugin,you can contribute in this repo [ClickHere](https://github.com/intelops/capten-plugins)

1. Add Plugin Application name

Add the plugin application name in `plugin-store/plugin-list.yaml` file

```
plugins:
  - argo-cd
  - crossplane
  - tekton
```


2. Add Plugin Application Configuration

Create a folder with plugin name `plugin-store/<plugin-name>`, and add plugin metadata files

- Add plugin application configuration in `plugin-store/<plugin-name>/plugin-config.yaml` file
- Add plugin application Icon file in `plugin-store/<plugin-name>/icon.svg`

| Attribute   | Description                           |
| ----------- | ------------------------------------- |
| pluginName  | Plugin application name               |
| description | Plugin application description        |
| category    | Plugin application category           |
| icon        | Plugin application icon               |
| versions    | Plugin application supported versions |


```
pluginName: "argo-cd"
description: "GitOps continuous delivery tool for Kubernetes"
category: "CI/CD"
icon: "icon.svg"
versions:
  - "v1.0.2"
  - "v1.0.5"
```

you can refer [here](https://github.com/intelops/capten-plugins/blob/main/plugin-store/argo-cd/plugin.yaml)

3. Add Plugin Application Version Deployment Configuration

For each supported version, create version folder `plugin-store/<plugin-name>/<version>` and add plugin version deployment metadata files

- add plugin application version deployment configuration in `plugin-store/<plugin-name>/<version>/plugin-config.yaml` file
- add plugin application version values file in `plugin-store/<plugin-name>/<version>/values.yaml` file

\*\* plugin application version deployment configuration attributes \*\*
| Attribute | Description |
| ------------------- | ---------------------------------------- |
| chartName | Plugin application chart name |
| chartRepo | Plugin application chart repo |
| version | Plugin application version |
| defaultNamespace | Plugin application default namespace |
| privilegedNamespace | Plugin application privileged namespace |
| valuesFile | Plugin application values file |
| apiEndpoint | Plugin application API endpoint |
| uiEndpoint | Plugin application ui endpoint |
| capabilities | Plugin application required capabilities |



Table below shows an example of plugin application version deployment configuration

| Capabilities                | Description                                                      |
| --------------------------- | ---------------------------------------------------------------- |
| deploy-controlplane-cluster | Capability to deploy plugin application on control plane cluster |
| deploy-bussiness-cluster    | Capability to deploy plugin application on business cluster      |
| capten-sdk                  | Capability to access Capten cluster SDK                          |
| ui-sso-oauth                | Capability to enable Single SignOn for plugin application UI     |
| postgress-store             | Capability to access Capten cluster storage service              |
| vault-store                 | Capability to access Capten cluster vault service                |

```
deployment:
  controlplaneCluster:
    chartName: "argo-cd"
    chartRepo: "https://kube-tarian.github.io/helmrepo-supporting-tools"
    version: "v1.0.2"
    defaultNamespace: "argo-cd"
    privilegedNamespace: false
    valuesFile: "values.yaml"
apiEndpoint: https://argo.{{.DomainName}}
uiEndpoint: https://argo.{{.DomainName}}
capabilities:
  - deploy-controlplane-cluster
  - capten-sdk
  - ui-sso-oauth
```

You can have sample reference of argocd plugin [here](https://github.com/intelops/capten-plugins/tree/main/plugin-store/argo-cd)

## General Instructions 
This project is written in Golang 

To contribute code.
1. Ensure you are running golang version 1.21 or greater for go module support
2. Set the following environment variables:
    ```
    GO111MODULE=on
    GOFLAGS=-mod=vendor
    ```
3. Fork the project.
4. Clone the project: `git clone https://github.com/[YOUR_USERNAME]/capten && cd capten`
5. kindly refer capten.md file to know the structure of the project.
6. Commit changes *([Please refer the commit message conventions](https://www.conventionalcommits.org/en/v1.0.0/))*
7. Push commits.
8. Open pull request.
