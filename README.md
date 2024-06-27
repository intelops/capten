![CAPTEN](.readme_assets/captenlogo.png)

The open-source DevSecOps platform for manging cloud infrastructure and cloud-native applications.
[![Docker Image CI](https://github.com/intelops/capten/actions/workflows/cli_release.yaml/badge.svg)](https://github.com/intelops/capten/actions/workflows/cli_release.yaml)
[![CodeQL](https://github.com/intelops/capten/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/intelops/capten/actions/workflows/github-code-scanning/codeql)

<hr>

## CAPTEN By INTELOPS

Capten streamlines the management of Kubernetes clusters, making it an ideal solution for teams or individuals who require automated cluster provisioning, application deployment, and lifecycle management in their development and testing workflows.

The all-in-one DevSecOps platform facilitates close collaboration to build and manage cloud-native ecosystems for application and infrastructure modernization, automation, and security.

## How to Install Capten Controlplane Cluster

Capten controlplane cluster creation supported with Capten CLI, Capten CLI distribution available for Linux, Winodws and MacOS.
Capten controlplane cluster creation supported on public cloud providers like AWS and Azure.

#### Prerequisites

- AWS or Azure clound provider account

- Azure CLI (Needed in case of using Azure cloud for cluster setup)

- Docker (Needed in case of using Capten CLI distribution on MacOS)

- kubectl tool to access Capten controlplane cluster

#### Setting up the cluster

1. Download and Extract the latest Capten package from Capten github repoistory [release page](https://github.com/intelops/capten/releases).

```bash
wget https://github.com/intelops/capten/releases/download/<latest-release>/capten_linux.zip
unzip capten_linux.zip && cd capten
```

2. Preparted the cluster installation parameters

Update cluster installation parameters:
For AWS cluster, update cluster installation parameters in the `aws_config.yaml` in `config` folder.

| Parameter               | Description                                                                        |
| ----------------------- | ---------------------------------------------------------------------------------- |
| AwsAccessKey            | Access key for AWS authentication                                                  |
| AwsSecretKey            | Secret key for AWS authentication                                                  |
| AlbName                 | Name of the Application Load Balancer (ALB)                                        |
| PrivateSubnet           | CIDR block for the private subnet(s)                                               |
| Region                  | AWS region where the resources will be deployed                                    |
| SecurityGroupName       | Name of the security group that controls inbound and outbound traffic              |
| VpcCidr                 | CIDR block for the Virtual Private Cloud (VPC)                                     |
| VpcName                 | Name of the Virtual Private Cloud (VPC)                                            |
| InstanceType            | Type of EC2 instance                                                               |
| NodeMonitoringEnabled   | Flag indicating whether node monitoring is enabled or not (true/false)             |
| MasterCount             | Number of master nodes                                                             |
| WorkerCount             | Number of worker nodes                                                             |
| TraefikHttpPort         | Port number for HTTP traffic handled by Traefik load balancer                      |
| TraefikHttpsPort        | Port number for HTTPS traffic handled by Traefik load balancer                     |
| TalosTg                 | Name of the target group for Talos instances                                       |
| TraefikTg80Name         | Name of the target group for port 80 traffic handled by Traefik                    |
| TraefikTg443Name        | Name of the target group for port 443 traffic handled by Traefik                   |
| TraefikLbName           | Name of the Elastic Load Balancer (ELB) used by Traefik                            |
| TerraformBackendConfigs | Configuration settings for Terraform backend (bucket name and DynamoDB table name) |

**Note:**
For a terraform backend,create the bucket and dynamoDB table in aws console.

For Azure cluster, update cluster installation parameters in the `azure_config.yaml` in `config` folder.

| Parameter            | Description                                                       |
| -------------------- | ----------------------------------------------------------------- |
| Region               | The Azure region where resources will be deployed                 |
| MasterCount          | Number of master nodes                                            |
| WorkerCount          | Number of worker nodes                                            |
| NICs                 | Network Interface Controllers (NICs) for master nodes             |
| WorkerNics           | Network Interface Controllers (NICs) for worker nodes             |
| InstanceType         | Type of virtual machine instance used for nodes                   |
| PublicIpName         | Names of public IP addresses assigned to the nodes                |
| TraefikHttpPort      | Port number for HTTP traffic handled by Traefik load balancer     |
| TraefikHttpsPort     | Port number for HTTPS traffic handled by Traefik load balancer    |
| Talosrgname          | Resource group name for the Talos deployment                      |
| Storagergname        | Resource group name for storage-related resources                 |
| Storage_account_name | Name of the storage account used for storing images               |
| Talos_imagecont_name | Name of the container within the storage account for Talos images |
| Talos_cluster_name   | Name of the Talos cluster                                         |
| Nats_client_port     | Port number for NATS client communication                         |

3. Prepare cluster application deployment parameters

Update cluster application deployment parameters in the `capten.yaml` in `config` folder.

| Parameter         | Description                                                                      |
| ----------------- | -------------------------------------------------------------------------------- |
| DomainName        | Name of the domain needed for exposing the application                           |
| ClusterCAIssuer   | The issuer of the Cluster Certificate Authority (CA) for cluster security        |
| SocialIntegration | The social platform like teams or slack integrated for alerting purpose          |
| SlackURL          | Slack channel url (needs to be provided if slack is used for social integration) |
| SlackChannel      | Name of the slack channel                                                        |
| TeamsURL          | Teams channel url (needs to be provided if teams is used for social integration) |

4. Create cluster

For creating the cluster, execute below command

```bash
./capten create cluster --cloud=<cloudtype> --type=talos
```

Note: Cloud type supported are 'aws' and 'azure'

- Cluster Creation through Docker Container:

For creating the cluster through docker container (needed in case of using Capten CLI distribution on Windows or MacOS ), run the below command

```bash
docker run -v /path/to/aws_config.yaml:/app/config/awsorazure_config.yaml -it ghcr.io/intelops/capten:<latest-image-tag>  create cluster --cloud=aws --type=talos
```

Post cluster creation, `kubeconfig` will be generated to `./config/kubeconfig`.
Access cluster using generated kubeconfig with kubectl

```bash
export KUBECONFIG=/home/capten/config/kubeconfig
kubectl get nodes
```

#### Setting up the cluster applications

For deploying the cluster applications, execute below command

```bash
./capten setup apps
```

Capten CLI will deploy Capten application suite and Capten Agent on Controlplane cluster.
Post application deployment, mTLS certificates are generated to access Capten Agent. mTLS certificates `capten-client-auth-certs.zip` generated in `cert` folder.

Deployed applications can be listed with helm tool

```bash
helm list -A
```

#### Show the cluster info

```bash
./capten show cluster info
```


#### Destroying the cluster

Cluster destruction command initiates the process of removing all components associated with the cluster, such as virtual machines, instances, nodes, networking configurations, and any other resources provisioned for the cluster. It effectively undoes the setup and configuration of the cluster, deallocating resources and ensuring they are no longer in use. This command can be used when the cluster is no longer needed or to clean up resources in cloud computing, distributed systems, or container orchestration environments.

```bash
./capten destroy cluster
```


# Capten Plugin SDK

## Overview

Capten Plugin SDK to develop and deploy applications on Capten cluster

## Features

- Onboard Plugin applications in to Capten Plugin Store
- Deploy Plugin applications on Capten cluster
- Deploy Plugin applications on Business cluster
- Capten Cluster Storage service access for Plugin application
- Capten Cluster Vault service access for Plugin application
- Capten Cluster SDK API service access for Plugin application
- Enabling Single SignOn for Plugin application UI

## Capten Plugin Store

Capten SDK provides a simple way to onboard the Capten Plugin Applications (Plugins or Plugin Apps) into the Capten Plugin Store.

### Central Plugin Store

Intelops manages the central plugin store Git repository (https://github.com/intelops/capten-plugin) to onboard and manage Capten Plugin applications. This is an open-source repository for the Capten community to onboard and manage Capten Plugin applications.

This Central Plugin Store is integrated with Capten Stack by default, and Plugins available in the Central Plugin Store can be deployed on the Capten cluster from the Capten UI.

### Local Plugin Store

Capten SDK supports managing a local plugin store, i.e., in your Git repository. A Git repository needs to be integrated into the Capten from Capten UI, which can then be used to onboard and manage Plugin applications.



## Depoly Plugin Application in Central Plugin Store

Deploy Plugin application on Capten cluster 


Capten SDK creates resources for plugin applications for plugin-configured capabilities before deploying plugin applications to the Capten cluster.

For deploying the application like crossplane,argocd and tekton,use the below commands:

1. Sync the applications in the application store with the below command

```bash
 ./capten plugin store synch --store-type central
```

2. For viewing the applications in the central store ,use the below command.

```bash
 ./capten plugin store list --store-type central
```
This command will list the applications in the central store and its respective version


3. For deploying the applications in the central store,use the below sample commands.Use the right version and the application name 

* For deploying tekton,use the below command

```bash
  ./capten plugin deploy --plugin-name tekton --store-type central --version v0.1.9
```
* For deploying crossplane,use the below command

```bash
 ./capten plugin deploy --plugin-name crossplane --store-type central --version v1.0.3
```

* For deploying argocd, use the below command
```bash
 ./capten plugin deploy --plugin-name argo-cd --store-type central --version v1.0.2
```



### Capten Pluguin Resources/Environment

Capten SDK creates resources for plugin application for the configured capabilities before deploying plugin application to Capten cluster

List of supported capabilities:

```
- Capten SDK
- Vault Store
- Postgres Store
```

#### Capten SDK

- This capability provides MTLS certificate for server and client authentications.
- Plugin application uses client certificate to communicate with Capten agent.
- Plugin application can server certificate to enable mtls for communication.

#### Vault Store

- This capability provides access to key secrets.

#### Postgres Store

- This capability provides postgres DB setup required for the plugin application.




# Capten Crossplane Plugin

## Onboard cluster resources:

### Git Project:

### Add Git Project via Capten CLI

For adding the git project,use the below command,

```bash
./capten cluster resources create --resource-type "git-project" --access-token "<gitbub-pat-token>" --user-id "<github-user-id>" --labels="crossplane" --git-project-url="<repo-url>"
```
**Note:**:

For creating Business cluster,add git project with the label "crossplane"

Provide the empty git repository to the repo-url and also provide the necessary github-pat token and userid 


### List Git Project via Capten CLI

```bash
./capten cluster resources list --resource-type="git-project"
```
the sample output of above command is given below

```bash
+--------------------------------------+-----------------------------------------+------------+
|                  ID                  |               PROJECT URL               |   LABELS   |
+--------------------------------------+-----------------------------------------+------------+
| 952001a4-7d6c-4cf7-adf3-83543ddea2bb | https://gitproject/git/sample.git       | crossplane |
+--------------------------------------+-----------------------------------------+------------+
```
### Update Git Project via Capten CLI

Below command is used for updating the git project
```bash
./capten cluster resources create --resource-type "git-project" --access-token "<gitbub-pat-token>" --user-id "<github-user-id>" --labels="crossplane" --git-project-url="<repo-url>" --id "<ID>"
```
**Note:**:

The ID can be retrieved while listing the git project.

For example,
```bash
./capten cluster resources create --resource-type "git-project" --access-token "msdbb" --user-id "samp" --labels="crossplane" --git-project-url="https://gitproject/git/sample.git" --id 952001a4-7d6c-4cf7-adf3-83543ddea2bb
```

### Delete Git Project via Capten 

```bash
/capten cluster resources delete --id "<ID>" --resource-type git-project
```

Below is the sample command for deleting the git project,

```bash
/capten cluster resources delete --id 952001a4-7d6c-4cf7-adf3-83543ddea2bb --resource-type git-project
```
**Note:**:

The id is retrieved by listing the git project


=======

### Cloud Provider:

### Add Cloud Provider via Capten CLI

```bash

./capten cluster resources create --resource-type="cloud-provider" --cloud-type="aws" --labels="crossplane" --access-key="accesskey" --secret-key="secretkey"
```
**Note:** The label _crossplane_ is used by the crossplane plugin to reference both the repository and provider.


### List Cloud Provider via Capten CLI


```bash

./capten cluster resources list  --resource-type="cloud-provider"
```
Below is the sample output

```bash

+--------------------------------------+------------+------------+
|                  ID                  | CLOUD TYPE |   LABELS   |
+--------------------------------------+------------+------------+
| 5a68240f-3f87-4c41-a589-5f32c1040f1e | aws        | crossplane |
+--------------------------------------+------------+------------+
```


### Update Cloud Provider via Capten CLI

```bash

./capten1 cluster resources update --resource-type="cloud-provider" --cloud-type="aws" --labels="tekton" --access-key="accesskey" --secret-key="secretkey" --id "<ID>"
```

For example,sample command will be,

```bash

./capten1 cluster resources update --resource-type="cloud-provider" --cloud-type="aws" --labels="tekton" --access-key="accesskey" --secret-key="secretkey" --id "5a68240f-3f87-4c41-a589-5f32c1040f1e"
```


### Delete Cloud Provider via Capten CLI

Below command is used for deleting the cloud provider

```bash
./capten cluster resources delete --resource-type cloud-provider --id "<ID>"
```


## Create Crossplane provider:

1. Add the git project by referring the above command

2. Sync the crossplane project with the beloe command

```bash
 ./capten plugin config --action synch-crossplane-project --plugin-name crossplane
```

3. Then list the cloud provider with the below command

```bash
./capten cluster resources list  --resource-type="cloud-provider"
```

The sample output is shown below:

```bash
+--------------------------------------+------------+------------+
|                  ID                  | CLOUD TYPE |   LABELS   |
+--------------------------------------+------------+------------+
| 5a68240f-3f87-4c41-a589-5f32c1040f1e | aws        | crossplane |
+--------------------------------------+------------+------------+

```
Then use the above lister id  for creating the crossplane provider,
```bash

./capten plugin config --action create-crossplane-provider --cloud-type aws --cloud-provider-id 5a68240f-3f87-4c41-a589-5f32c1040f1e
```

4. Then again sync the crossplane plugin

```bash
./capten plugin config --action synch-crossplane-project --plugin-name crossplane
```

## Create Business cluster

1. After the sync is successful, the crossplane objects and its argocd applications are added to the empty repository under the infra directory.
2. Go to infra/clusters/cluster-configs/cluster-claim.yaml
3. Uncomment the cluster-claim.yaml file (or add any required changes)
4. Go to argocd UI page and sync all crossplane related applications
5. After the clusterclaim is created, the business cluster creation will get triggered.

![BusinessClusterCreation](.readme_assets/bc-create.png)

Make sure to sync all crossplane related apps_

## Delete Business cluster

1. To delete the business cluster, remove all applications from the business cluster.
2. Go to infra/clusters/cluster-configs and remove cluster-claim.yaml
3. Now prune sync the cluster-config-app application (watching the cluster-claim.yaml).
4. This will trigger the business cluster deletion

## Delete Crossplane provider

1. List the crossplane provider with the below command

```bash
./capten plugin config --action list-crossplane-providers --plugin-name crossplane
```

+--------------------------------------+------------+--------------------------------------+--------+
|                  ID                  | CLOUD TYPE |          CLOUD PROVIDER ID           | STATUS |
+--------------------------------------+------------+--------------------------------------+--------+
| ee162729-6d34-4481-8a85-e3e92694511b | aws        | e2f4886e-c201-45c7-8813-36449af67e09 | Ready  |
+--------------------------------------+------------+--------------------------------------+--------+


2.. Use the below command for deleting the crossplane provider (use the  id for deleting the crossplane provider)

```bash

./capten plugin config --action delete-crossplane-provider --crossplane-provider-id ee162729-6d34-4481-8a85-e3e92694511b --plugin-name crossplane
```

3. Then synch the project with the below command

```bash
./capten plugin config --action synch-crossplane-project --plugin-name crossplane
```

With the above 3 commands ,crossplane provider will be deleted


# Tekton Flow For Capten

## Tekton CI/CD Pipeline

### Add Git Project via Capten CLI

Onboard the git project by providing repo url,access token and label for the customer repo (label is tekton) and for the tekton ci/cd repo (label is IntelopsCi) using the below steps:

1. First onboard the empty git project  

  ```bash
  ./capten cluster resources create --resource-type="git-project" --access-token "<accesskey>" --user-id "<user-id>" --labels "tekton" --git-project-url "<repo-url>"
  ```


2. Then onboard the application repo-url ,access-token and userId with the label IntelopsCi

```bash

./capten cluster resources create --resource-type="git-project" --access-token="access-token" --user-id="git-user-id" --labels="IntelopsCi" --git-project-url="https://github.com/sample/qt-test-application.git"
```


### Onboarding container registry in to capten


1. Add Container registry

Add the registry url,username,access token and label to which the built image needs to be pushed (label is "tekton")

```bash

./capten cluster resources create --resource-type="container-registry" --registry-url="ghcr.io" --registry-type="GitHub Registry" --registry-username="registry-user-name" --registry-password="registry-password"
```

List Container Registry

```bash
./capten cluster resources list  --resource-type="container-registry"
```

```bash
+--------------------------------------+-----------------+--------------+--------+
|                  ID                  |  REGISTRY TYPE  | REGISTRY URL | LABELS |
+--------------------------------------+-----------------+--------------+--------+
| 6d7ec739-168d-472e-84d0-66971fb29bb5 | GitHub Registry | ghcr.io      |        |
+--------------------------------------+-----------------+--------------+--------+
```


Update Container Registry

```bash

./capten cluster resources create --resource-type="container-registry" --registry-url="ghcr.io" --registry-type="GitHub Registry" --registry-username="registry-user-name" --registry-password="registry-password" --id "<id>"
```
 The id can be retrieved by listing the container-registry


Delete Container Registry

```bash
./capten cluster resources delete --id "<id>" --resource-type container-registry
```


# Configuring Tekton

## Configuring Capten Tekton Plugin


Use the below commad for synchronizing the tekton to the customer empty repository.This will configure the tekton and  neccessary folders will be created in the customer's repo.

```bash
 ./capten plugin config --action synch-tekton-project --plugin-name tekton
```

The status can be viewed by using below command

```bash
./capten plugin config --action show-tekton-project --plugin-name tekton
```

The sample output can be shown below

```bash
+-----------------+-----------------------------------------+
|    ATTRIBUTE    |                  VALUE                  |
+-----------------+-----------------------------------------+
| git-project-url | https://github.com/sample/infra.git |
| status          | configured                              |
+-----------------+-----------------------------------------+

```

# Pre-requisite For Tekton CI/CD Pipeline Creation

- Use the already created **tekton-pipelines** namespace for the creation of pipeline.

- * Create a secretstore from the yaml given below.Replace the server with the url which can be obtained from the **kubectl** command given below.

  ```bash
  kubectl get ingress -n capten
  ```

       apiVersion: external-secrets.io/v1beta1
       kind: SecretStore
       metadata:
         name: vault-root-store
       spec:
         provider:
           vault:
             server: <"replace with the ingress host obtained from above command">
             path: "secret"
             version: "v2"
             auth:
               tokenSecretRef:
                 name: "tekton-vault-token"
                 key: "token"
                 

Here, the **tekton-vault-token** is the secret created in tekton namespace to access the vault.Copy-paste the **tekton-vault-token** secret in the required namespace where the tekton pipeline will be present and then create the above secretstore.

- Git secret

  Go to _onboarding-->git_ under the respective git project the path of the vault where the credentials of git stored can be viewed.copy the path and add it to the path in the external secret yaml as given below

## Note

  Annotate the external-secret to specify the domains for which Tekton can use the credentials.

  A credential annotation key must begin with tekton.dev/git- or tekton.dev/docker- and its value is the URL of the host for which Tekton will be using that credential.
  eg-tekton.dev/git-0: https://gitlab.com , tekton.dev/git-0: https://github.com , tekton.dev/docker-0: https://gcr.io

          apiVersion: external-secrets.io/v1beta1
          kind: ExternalSecret
          metadata:
            annotations:
              tekton.dev/git-0: "https://github.com"
            name: gitcred-external
            namespace: tekton-pipelines
          spec:
            refreshInterval: "10s"
            secretStoreRef:
              name: vault-root-store
              kind: SecretStore
            target:
              name: gitcred-capten-pipeline
            data:
            - secretKey: password
              remoteRef:
                key: <vault path cpoied from ui>
                property: accessToken
            - secretKey: username
              remoteRef:
                key: <vault path copied from ui>
                property: userID

- Container registry secret

  Go to _onboarding-->container registry_ under the respective container registry, where the path of the vault where the credentials of container registry stored can be viewed.copy the path and add it to the path in the external secret yaml as given below

         apiVersion: external-secrets.io/v1beta1
         kind: ExternalSecret
         metadata:
           annotations:
             tekton.dev/git-0: "https://github.com"
           name: docker-external
           namespace: tekton-pipelines
         spec:
           refreshInterval: "10s"
           secretStoreRef:
             name: vault-root-store
             kind: SecretStore
           target:
             name: docker-credentials-capten-pipeline
           data:
           - secretKey: config.json
             remoteRef:
               key: <vault path copied from ui>
               property: config.json

- Cosign docker login secret

  Go to _onboarding-->conatainer registry_ under the respective container registry where the path of the vault in which the credentials of container registry stored can be viewed.copy the path and add it to the path in the external secret yaml as given below

      apiVersion: external-secrets.io/v1beta1
      kind: ExternalSecret
      metadata:
        annotations:
          tekton.dev/git-0: "https://github.com"
        name: cosign-docker-external
        namespace: tekton-pipelines
      spec:
        refreshInterval: "10s"
        secretStoreRef:
          name: vault-root-store
          kind: SecretStore
        target:
          name: cosign-docker-secret-capten-pipeline
        data:
        - secretKey: password
          remoteRef:
            key: <vault path copied from ui>
            property: password
        - secretKey: registry
          remoteRef:
            key: <vault path copied from ui>
            property: registry
        - secretKey: username
          remoteRef:
            key: <vault path copied from ui>
            property: username

- Argocd secret
  Use the below secret yaml and replace the password with the encoded argocd password which can be obtained by using the **kubectl** command and the server url is obtained from the capten ui under _capten-->platform-engineering_ .Copy the repo url from the argocd setup ,encoded it and add it to the server url.Username is admin ,add the encoded username to the yaml given below

```bash
 kubectl get secrets argocd-initial-admin-secret -n argo-cd
```

      apiVersion: v1
      data:
        PASSWORD: <replace with encoded argocd secret>
        SERVER_URL: <repo url copied from ui>
        USERNAME: <encoded username>
      kind: Secret
      metadata:
        name: argocd-capten-pipeline
        namespace: tekton-pipelines
      type: Opaque

- cosign-keys

  Now the cosign keys secret is automatically created in tekton-pipelines namespace.

- Extra-config secret

  Go to _onboarding-->git_ under the respective git project where the path of the vault in which the credentials of git stored can be viewed.copy the path and add it to the path in the external secret yaml as given below

      apiVersion: external-secrets.io/v1beta1
      kind: ExternalSecret
      metadata:
        annotations:
          tekton.dev/git-0: "https://github.com"
        name: extraconfig-external
        namespace: tekton-pipelines
      spec:
        refreshInterval: "10s"
        secretStoreRef:
          name: vault-root-store
          kind: SecretStore
        target:
          name: extraconfig-capten-pipeline
        data:
        - secretKey: GIT_TOKEN
          remoteRef:
            key: <vault path copied from ui>
            property: accessToken
        - secretKey: GIT_USER_NAME
          remoteRef:
            key: <vault path copied from ui>
            property: userID

# Prepare Pipeline Resources For The Tekton Pipeline

Now commit the required pipeline,rbac,triggers and ingress in the customer repo under the directory *cicd-->tekton-pipelines-->templates*.Now go to the **argocd ui** and sync the tekton-pipelines application manually. 
once done the argocd will update this changes to the cluster and the pipeline,triggers,rbac and ingress will be created in the controlplane cluster

![PipelineResource](.readme_assets/infra.png)

# QT in tekton

  1. create a new folder **qt_test** in the root directory
   
  2. place the test.yaml file which contains the testcases to test the deployed application inside the folder
   
  3. Note that the folder name should be **qt_test** and file name should be **test.yaml**.The sample testcase is given below,

     
         type: Test
         spec:
           id: zSn7HGzRR
           name: user service
           trigger:
             type: http
             httpRequest:
               method: GET
               url: http://user.svc.local/api/UserService
           specs:
             - selector: span[tracetest.span.type="http" name="GET api/UserService"]
               assertions:
                 - attr:http.request.method = "GET"
                 - attr:http.route = "api/UserService"
                 - attr:name = "GET api/UserService"
                 - attr:tracetest.span.name = "GET api/UserService"
                 - attr:tracetest.span.type = "http"
                 - attr:url.path = "/api/UserService"
                 - attr:url.scheme = "http"
                 - attr:user_agent.original = "Go-http-client/1.1"
             - selector: span[tracetest.span.type="general" name="Tracetest trigger"]
               assertions:
                 - attr:tracetest.span.type = "general"
                 - attr:tracetest.span.name = "Tracetest trigger"

  4. Also ensure whether the  applicatin which is deployed in the business cluster is exposed
   
  5. Then add the url of the application in the test.yaml file
   
  6. After the execution of the pipeline we can check the qt task success and failure in the tekton-dashboard
   
  7. when the tekton pipeline fails due to quality-trace test case, it may be due to which assertions in the test yaml getting failed or configuration error between the test app and quality-trace. We can check the quality-trace server pod logs and quality-trace otel collector logs for errors and trace ids.If the test fails due to configuration errors, then in the quality-trace server log, it can be viewed as "Trace not found" error


* If all the assertions pass,below logs is shown

![AssertionPass](.readme_assets/assertionspass.png)

* If the assertion fails,below logs is shown

![AssertionFail](.readme_assets/assertionfail.png)

# Triggering Tekton Pipeline

  1. Now add the **webhook url** to the tekton ci/cd repo and select the **event** on which the tekton pipeline needs to be executed upon trigger.
   
  ![WebhookImage](.readme_assets/webhook-img.png)

  2. If needed one can protect their branch using the branch protection rule which will be present under *settings-->branches-->add rule*

  3. In the add rule select the Require status checks to pass before merging and Require branches to be up to date before merging
     
  4. Then in the search box that appers under Require branches to be up to date before merging ,search for tekton-pipelines and add it,now whenever a pull_request is raised the check will run and once the checks is success the **merge** option will be visible
 
  ![TektonStatus](.readme_assets/tekton-status.png)

  5. once all the setup is done and now when a pull_request event is triggered (when a pull_request is raised), the tekton pipeline will get executed and the image gets built and pushed to the container registry ,the built image is then signed using cosign and finally once the application is deployed in the bussiness cluster the qt task in the pipeline will run the testcase to test the application and the success/failure task will get executed depending upon the result of pipeline.similarly the pipeline can trigger for event such as push,tag and release

  6. Also the success and failure status will be notified back to the git repo in the case of pull_request event.Note for this the branch protection rule needs to be added
     
  7. The tekton related pipelines and tasks can be viewed in the tekton-dashboard by clicking on the details option where the check is running or by clicking on the tekton icon present under Capten-->Controlplane Cluster in the ui

  ![DeployQT](.readme_assets/deployqt.png)

  ![TektonQTLogs](.readme_assets/tekton-qt-logs.png)

# What is Proact and how to use it?

Proact is CLI/CI Tool for Automating Vulnerability Management for Enhancing Software Supply Chain Security Measures.

* To create a new schedule click on Schedule Scan repo.


```bash
curl -X 'POST' \
  'http://proact-scheduler.awsagents.optimizor.app/api/v1/schedule/' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "schedule_name": "<schedule-name>",
  "container_registry_id": "ghcr.io",
  "container_registry_url": "ghcr.io",
  "cron_schedule": "* * * * *",
  "scan_configs": [
    {
      "docker_image_name": "<docker-image-name>",
      "pyroscope_enabled": "true",
      "pyroscope_url": "http://pyroscope.<DomainName>",
      "pyroscope_app_name": "user",
      "rebuild_image": "true",
      "docker_file_folder_path": "<dockerfile-repo-url>"
    }
  ]
}'

```

**Note**

Use the domain name that is configured in **capten.yaml** in the pyroscope_url.

```bash
 "pyroscope_url": "http://pyroscope.<Domain-Name>",
```


* To List all schedules

```bash
curl -X 'GET' \
  'http://proact-scheduler.<Domain-Name>/api/v1/schedule' \
  -H 'accept: application/json'
```
* To pause schedule:

```bash
curl -X 'PUT' \
  'https://proact-scheduler.<DomainName>/api/v1/schedule/{scheduleId}/pause' \
  -H 'accept: application/json'
```

* To resume schedule:

```bash
curl -X 'PUT' \
  'https://proact-scheduler.<DomainName>/api/v1/schedule/{scheduleId}/resume' \
  -H 'accept: application/json'
```

* To delete schedule:

```bash
curl -X 'DELETE' \
  'https://proact-scheduler.<Domain-name>/api/v1/schedule/{scheduleId}' \
  -H 'accept: application/json'
```


## CAPTEN UI

Follow the documentation [here](https://github.com/intelops/capten/blob/main/capten-ui.md) to access the capten UI 