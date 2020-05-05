---
id: setup-kubernetes-in-cloud
title: How to setup managed Kubernetes in the cloud
sidebar_label: Setup managed Kubernetes
keywords:
  - InstantHIE
  - How to
  - Kubernetes
  - AWS
  - Google cloud
  - Azure
  - Digital ocean
description: How to setup managed Kubernetes in the cloud
---

Below are some useful link on getting started with managed Kubernetes in a number of different cloud service providers along with how to get setup with their own CLI applciation that help you setup the managed Kubernetes infrastructure.

## Amazon Web Services

Useful links:

- [AWS Cli Setup guide](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- [Amazon Elastic Kubernetes Service(EKS) Quickstart](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-console.html)

### Install AWS Cli

```sh
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

unzip awscliv2.zip

sudo ./aws/install
```

## Google Cloud

Useful Links:

- [Google Cloud Cli Setup Guide](https://cloud.google.com/sdk/docs#deb)
- [Google Cloud Cluster Quickstart](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl)

## Azure

Useful Links:

- [Azure Cli Setup Guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-apt?view=azure-cli-latest)
- [Quickstart](https://docs.microsoft.com/en-us/cli/azure/aks?view=azure-cli-latest)

### Install Azure Cli

```sh
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

az login
```

## Digital Ocean

Useful Links:

- [Digital Ocean Cli Setup Guide](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/)

### Install Digital Ocean Ctl

```sh
curl -OL https://github.com/digitalocean/doctl/releases/download/v<version>/doctl-<version>-linux-amd64.tar.gz

sudo mv ./doctl /usr/local/bin
```

> [See here for latest version](https://github.com/digitalocean/doctl/releases)