# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

- Interoperability Layer - [OpenHIM](http://openhim.org/)
- FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Kubernetes Host Environment Setup

> If you are working with multiple cluster hosting clients, you may need to switch cluster context for your deployment. See the below commands:

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

The OpenHIM console url will be displayed in the terminal output when the script completes. The Url may take a few minutes to become active as the pod may not be fully initialized yet.

This bash script will apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

### View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all -n core-component
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

We have included a useful set of scripts to initialize the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. From the `kubernetes` directory, use the following command to implement:

```bash
./importer/k8s.sh up
```

> These scripts can be duplicated and modified to implement custom imports

To clean up the remaining job and pods from a successful setup, run the following:

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

To run in development mode, where the OpenHIM mongo database, HAPI FHIR server, and the MySQL database can be accessed directly through their urls, run the following command

```bash
./dev/k8s.dev.sh
```

---
