![CAPTEN](.readme_assets/captenlogo.png)

The open-source platform for creating the cluster,deploying the application and destroying the cluster.
[![Docker Image CI](https://github.com/intelops/capten/actions/workflows/cli_release.yaml/badge.svg)](https://github.com/intelops/capten/actions/workflows/cli_release.yaml)
[![CodeQL](https://github.com/intelops/capten/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/intelops/capten/actions/workflows/github-code-scanning/codeql)


<hr>

## CAPTEN

A tool to create cluster, set up application in the cluster and destroy cluster.

## How Capten works

With Capten,cluster can be created of type talos in the cloud aws .

After cluster is being created,capten supports certain applications to be deployed.So far capten supports some of the applications given below: 

* [Openebs-cstor](https://openebs.io/docs/user-guides/cstor)-Configure cStor storage and use cStor Volumes for running stateful workloads.
* [Cert-Manager](https://github.com/cert-manager/cert-manager#cert-manager)- adds certificates and certificate issuers as resource types in Kubernetes clusters, and simplifies the process of obtaining, renewing and using those certificates.
* [traefik](https://traefik.io/traefik/)- a modern HTTP reverse proxy and load balancer that makes deploying microservices easy.
* pre-install
* [prometheus](https://prometheus.io/)-collects metrics from configured targets at given intervals, evaluates rule expressions, displays the results, and can trigger alerts when specified conditions are observed.
* [vault](https://www.vaultproject.io/)- Manage secrets and protect sensitive data
* [vault-cred](https://github.com/intelops/vault-cred)- Automate Vault unsealing,continuous monitoring of ConfigMap to create vault policy and vault role.Stores service based credential,certificate and any generic credential.
* [external-secrets](https://github.com/external-secrets/external-secrets)-a K8s operator that integrates external secret management systems like AWS Secrets Manager, HashiCorp Vault and many more.The operator reads information from external APIs and automatically injects the values into a Kubernetes Secret.
* [k8ssandra-operator](https://docs.k8ssandra.io/components/k8ssandra-operator/)-Kubernetes-based distribution of Apache Cassandra that includes several tools and components that automate and simplify configuring, managing, and operating a Cassandra cluster.
* [loki](https://grafana.com/oss/loki/)-log aggregation system designed to store and query logs from all your applications and infrastructure.
* [Kyverno](https://kyverno.io/)- a policy engine designed for Kubernetes
* k8ssandra-cluster
* monitoring
* [kubviz-client and kubviz-agent](https://github.com/intelops/kubviz)-Visualize Kubernetes & DevSecOps Workflows. Tracks changes/events real-time across your entire K8s clusters, git repos, container registries, etc. , analyzing their effects and providing you with the context you need to troubleshoot efficiently.
* [signoz](https://signoz.io/)- open-source observability tool that helps you monitor your applications and troubleshoot problems.
* [temporal](https://temporal.io/)- open source programming model that can simplify code, make applications more reliable
* [kad](https://github.com/kube-tarian/kad)-Extensible open-source framework that Integrates & Scales your DevSecOps and MLOps stacks as you need
* policy-reporter
* [kubescape](https://www.armosec.io/kubescape/)- an open-source Kubernetes security platform. It includes risk analysis, security compliance, and misconfiguration scanning.  
* [falco](https://falco.org/)-cloud native runtime security tool for Linux operating systems. It is designed to detect and alert on abnormal behavior and potential security threats in real-time.
* [tracetest](https://tracetest.io/)-a trace-based testing tool for building integration and end-to-end tests in minutes using your OpenTelemetry traces
* [velero](https://velero.io/)- an open source tool to safely backup and restore, perform disaster recovery, and migrate Kubernetes cluster resources and persistent volumes.

Then capten also suuports cluster to be destroyed.

## How to install and run capten

#### Prerequisites
* go binary 


1.Clone the repo
```bash
git clone git@github.com:intelops/capten.git
```
or 
```bash
git clone https://github.com/intelops/capten.git
```
2.Give the below command to create capten zip folder

```bash
make build.release
```
3.Unzip the capten folder
```bash
unzip capten.zip
```
4.Navigate to the unzipped capten folder
```bash
cd capten
```
Then below commands can be used for creating the cluster,setting up the application in the cluster and also for deleting the cluster.
#### Create Cluster
```
./capten create cluster --cloud=aws --type=talos
```

#### Destroy Cluster
```
./capten destroy cluster
```

#### Setup Apps
```
./capten setup apps
```
#### Show Cluster Info
```
./capten show cluster info
```