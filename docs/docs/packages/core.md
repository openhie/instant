---
id: core
title: Core Package
sidebar_label: Core
keywords:
  - InstantHIE
  - Core
  - Package
description: The core package of the InstantHIE
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two components:

- Interoperability Layer - [OpenHIM](http://openhim.org/)
- FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Getting Started

Select a deployment platform below to follow the getting started steps in setting up this package.

<Tabs
  defaultValue="dockerCompose"
  values={[
    { label: 'Docker Compose', value: 'dockerCompose' },
    { label: 'Kubernetes', value: 'kubernetes' }
  ]
}>
<TabItem value="dockerCompose">

Before we proceed with creating our Core package components, we need to ensure we are in the correct directory containing our `docker-compose` script.

Once you are in the correct working directory (`core/docker/`) we can proceed to execute our `docker-compose` script by running the below command which will create all the services and print out their logs in the terminal.

```bash
docker-compose up
```

### Useful compose flags

Some additional flags can be passed to the `docker-compose` command making it a bit easier to work with.

- `-d`: Run the services in a detached mode. This means that when you close or exit your terminal, the services will still be running in the background.
- `-f`: Specify the location of the `docker-compose` file to be executed. Omitting this flag will look for the default `docker-compose.yml` file.
- `--force-recreate`: This will force the container/image to be re-created if a newer version is found. This is useful when a new image has been released but not yet pulled onto the host machine.

```bash
docker-compose up -d --force-recreate
```

### Environment configuration

By running the above command to get started with the Core package we create all the services that need to be defined, but this script might have some limitations depending on the type of environment you want to run the configuration

Additional `docker-compose` files are available for extra environment configuration

- **docker-compose.yml**: Main `docker-compose` script to create the services
- **docker-compose.dev.yml**: Development `docker-compose` script to override some of the default configuration to be used in a development environment (Open service ports for access etc)

The below command specifies the three `docker-compose` files that need to be executed for the development configuration

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml -f docker-compose.config.yml up -d
```

The below command specifies the two `docker-compose` files that need to be executed for a production-like configuration

```bash
docker-compose -f docker-compose.yml -f docker-compose.config.yml up -d
```

</TabItem>
<TabItem value="kubernetes">

> If you are working with multiple cluster hosting clients you may need to switch cluster contexts for you deployments. See the below commands:

```sh
kubectl config get-contexts
kubectl config use-context <context-name>
```

### Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once installed, we can start the minikube service be executing the below command:

```sh
minikube start --cpus 4 --memory 8192
```

This also updates the VM settings to make use of 4 CPU's and 8GB of RAM, instead of the default 2 CPU's and 4GB of RAM

### Amazon Web Services

Useful links:

- [AWS Cli Setup guide](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- [Amazon Elastic Kubernetes Service(EKS) Quickstart](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-console.html)

Some prerequisites are required before we can continue to deploy our Kubernetes infrastructure to an AWS cluster.

- You have created all the various users and permissions as required.
- You have given the users the relevant access to the AWS services
- You have generated an access token for your AWS user
- You have installed all the relevant CLI tools

#### Install AWS Cli

```sh
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

unzip awscliv2.zip

sudo ./aws/install
```

#### Configure with your AWS token details

```sh
aws configure
```

#### Install EksCtl

```sh
curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp

sudo mv /tmp/eksctl /usr/local/bin

eksctl version
```

#### Install Kubectl

```sh
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt` /bin/linux/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

kubectl version --client
```

#### Create the cluster

Before we can deploy our Kubernetes infrastructure we need to make sure we have created a cluster for us to deploy to. Execute the below command to create the cluster within AWS EKS

```sh
eksctl create cluster -f cluster.yml
```

#### Configure cluster users

Once the cluster has been created successfully, we also need to give access to the various users accessing the cluster. Update the `cluster-auth.yml` file with the users that need access and replace the `data.mapRoles.rolearn` with the arn of the role created to manage this cluster. Execute the below command to find the ARN of the role linked to the cluster:

```sh
kubectl describe configmap -n kube-system aws-auth
```

Once the `cluster.auth.yml` file has been updated, execute the below command to give the users access to the cluster.

```sh
kubectl replace -f cluster-auth.yml
```

#### Access an existing cluster

1. See the available clusters

   ```sh
   eksctl get clusters
   ```

1. Create config file locally to reference existing cluster

   ```sh
   eksctl utils write-kubeconfig --cluster <cluster-name>
   ```

1. Check current cluster context

   ```sh
   kubectl config get-contexts
   ```

#### Kill cluster

```sh
eksctl delete cluster -f cluster.yml
```

### Google Cloud

Useful Links:

- [Google Cloud Cli Setup Guide](https://cloud.google.com/sdk/docs#deb)
- [Google Cloud Cluster Quickstart](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl)

### Azure

Useful Links:

- [Azure Cli Setup Guide](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-apt?view=azure-cli-latest)
- [Quickstart](https://docs.microsoft.com/en-us/cli/azure/aks?view=azure-cli-latest)

#### Install Azure Cli

```sh
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

az login
```

### Digital Ocean

Useful Links:

- [Digital Ocean Cli Setup Guide](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/)

#### Install Digital Ocean Ctl

```sh
curl -OL https://github.com/digitalocean/doctl/releases/download/v<version>/doctl-<version>-linux-amd64.tar.gz

sudo mv ./doctl /usr/local/bin
```

> [See here for latest version](https://github.com/digitalocean/doctl/releases)

---

## Getting Started

Useful Links:

- [Kubectl Cheat Sheet](https://www.digitalocean.com/community/cheatsheets/getting-started-with-kubernetes-a-kubectl-cheat-sheet)

Before we proceed with creating our `Core Package` services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes`) we can proceed to create our core instant ohie deployment with the following command:

```bash
./main/k8s.sh up
```

The OpenHIM console url will be displayed in the terminal output when the script completes. The Url may take a few minutes to become active as the pod may not be fully initialised yet.

This bash script will apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

### View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all
```

To tear down this deployment use the opposing command:

```bash
./main/k8s.sh down
```

To completely remove all project components use the following option:

```bash
./main/k8s.sh destroy
```

---

## Initial OpenHIM Config

We have included a useful set of scripts to initialise the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. From the `kubernetes` directory, use the following command to implement:

```bash
./importer/k8s.sh up
```

> These scripts can be duplicated and modified to implement custom imports

To clean up the remaining job and pods from a successful setup run the following:

```bash
./importer/k8s.sh clean
```

Once the config is done HAPI FHIR will be accessible on: <https://OPENHIM_CORE_HOSTNAME/hapi-fhir-jpaserver/fhir/>

You may test that the OpenHIM is routing requests to HAPI FHIR by running:

```bash
./test.sh <OPENHIM_CORE_HOSTNAME>
```

---

## Development mode for exposed services

To run in development mode, where the OpenHIM mongo database, HAPI fhir server and the MySQL database can be accessed directly through their urls, run the following command

```bash
./dev/k8s.dev.sh
```

---
</TabItem>
</Tabs>

## Accessing the services

- **OpenHIM**
  - Console: Displayed in the output of the startup script
  - Username: **root@openhim.org**
  - Password: **instant101**
- **HAPI FHIR**
  - This service should not be publicly accessible and only accessed via the Interoperability Layer

## Testing the Core package

As part of the Core package setup we also do some initial imports of config for connecting the services together.

- OpenHIM: Import a public channel configuration that routes requests to the HAPI FHIR services
- HAPI FHIR: _Not config import yet_

For testing this Core package we will be making use of `curl` for sending our request, but any client could be used to achieve the same result.

Execute the below `curl` request to successfully route a request through the OpenHIM to query the HAPI FHIR server.

```bash
curl <openhim_core_transaction_api_url>/hapi-fhir-jpaserver/fhir/Patient
```

> The **openhim_core_transaction_api_url** is displayed in the output f the startup script
