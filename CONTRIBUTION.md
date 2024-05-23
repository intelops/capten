# Contribution Guidelines
Please read this guide if you plan to contribute to the Capten. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with the Capten, please [open an issue on GitHub](https://github.com/intelops/capten/issues/new?labels=kind%2Fbug&template=bug-report.md&title=Bug:) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/intelops/capten/issues/new?labels=kind%2Ffeature&template=feature-request.md&title=Feature%20Request:).


## Developing 

Development can be conducted using  GoLang compatible IDE/editor (e.g., Jetbrains GoLand, VSCode).

There are 3 places where you develop new things on Capten: on the CLI and  on the kad and on the Documentation website.

### Folder Structure 

## Directory Structure

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

## How to Contribute 

Written in Golang, the CLI code is stored in the folder `./pkg/cmd`. You can add any additional CLI options here .
For eg if you wish to cluster creation for any cloud,you can also modify the terraform related changes in `./pkg/terraform` 

And also you 

To test your modification,you can just build the CLI artifact with the below command

```sh
make build.release-linux
cd capten
```
Then with the build binary,you can test your changes

For bringing up any additional apps or tools,you can work on `./apps`.
For eg,create a yaml file for the tool,with the below specifications


```sh
Name: "name of the application"
ChartName: "The path or name of the Helm chart within the repository"
Category: "Helps to group and identify the application type (e.g., Security)"
RepoName: "The name of the repository where the Helm chart is located"
RepoURL: "The URL of the Helm repository from which the chart will be fetched"
Namespace: The Kubernetes namespace in which the Helm release will be deployed. 
ReleaseName: "The name given to this Helm release, uniquely identifying the release within the namespace"
Version: " Specifies the version of the Helm chart to use"
CreateNamespace: A boolean value that determines if the namespace should be created if it doesn't exist

```

If any values needs to be overrided,you can create a sample `_template.yaml` in `./apps/conf/values`.And pass the override values in this yaml file.

Then add the application name in the `./apps/conf/core_group_apps.yaml` or `./apps/conf/default_group_apps.yaml`


## General Instructions for contributing Code
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

