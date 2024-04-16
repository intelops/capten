---
title: "Configuration"
date: 2023-07-24
weight: 2
draft: false
---

## Setting Up Configurations for cluster creation

You can edit the specifications in [azure_config.yaml](https://github.com/intelops/capten/blob/main/config/azure_config.yaml) or [aws_config.yaml](https://github.com/intelops/capten/blob/main/config/aws_config.yaml) as per your requirements.

## Setting Up Configurations for App Deployment

you can specify the domain name and  also you can configure alerts (alerting in case when node goes down or when pod is in crashloopbackoff state and in many cases) in `teams` or `slack` by editing the specifications in [capten.yaml](https://github.com/intelops/capten/blob/main/config/capten.yaml)

## How to configure alerts in teams?

1.Select the channel where the cluster alerts should be send

2.Click the channel and select the connectors option

   ![connector-image](./teamsconnector.png#gh-light-mode-only)

3.you can see the incoming webhook option after clicking the connector

![configure-teams](./teamsconfigure.png#gh-light-mode-only)

3.Once after configuring the webhook,you can see the url created.

![teams-url](./teamsurl.png#gh-light-mode-only)

Copy the created url and Specify the teams url in the [capten.yaml](https://github.com/intelops/capten/blob/main/config/capten.yaml) .


#### Note
you can configure alerts in either teams or slack. If you are using teams for getting cluster alerts,specify 'teams' in the  Social Integration field in the `capten.yaml` or else specify 'slack' .

