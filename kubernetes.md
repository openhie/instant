# Kubernetes Host Environment Setup

> If you are working with multiple cluster hosting clients, you may need to switch cluster context for your deployment. See the below commands:

```sh
kubectl config get-contexts
kubectl config use-context <context-name>
```

## Minikube (local Kubernetes)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once installed, we can start the minikube service be executing the below command:

```sh
minikube start --cpus 4 --memory 8192 --disk-size 8g
```

This also updates the VM settings to make use of 4 CPU's and 8GB of RAM, instead of the default 2 CPU's and 4GB of RAM. Also the disk space is reduced from the default 20GB to 8GB to save space.

## Amazon Web Services

Useful links:

- [AWS Cli Setup guide](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- [Amazon Elastic Kubernetes Service(EKS) Quickstart](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-console.html)

Some prerequisites are required before we can continue to deploy our Kubernetes infrastructure to an AWS cluster.

- You have created all the various users and permissions as required.
- You have given the users the relevant access to the AWS services
- You have generated an access token for your AWS user
- You have installed all the relevant CLI tools

### Install AWS Cli

```sh
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

unzip awscliv2.zip

sudo ./aws/install
```

### Configure with your AWS token details

```sh
aws configure
```

### Install EksCtl

```sh
curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp

sudo mv /tmp/eksctl /usr/local/bin

eksctl version
```

### Install Kubectl

```sh
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt` /bin/linux/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

kubectl version --client
```

### Create the cluster

Before we can deploy our Kubernetes infrastructure we need to make sure we have created a cluster for us to deploy to. Execute the below command to create the cluster within AWS EKS

```sh
yarn eks:cluster:create
```

Once this completes you may execute a Kubernetes deployment, kubectl should already be configured to point to the newly created cluster. E.g.

```sh
yarn docker:instant init -t k8s
```

### Access an existing cluster

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

### Kill cluster

```sh
yarn eks:cluster:destroy
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

---
