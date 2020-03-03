# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

- Interoperability Layer - [OpenHIM](http://openhim.org/)
- FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once installed, we can start the minikube service be executing the below command:

```sh
minikube start --cpus 4 --memory 8192
```

This also updates the VM settings to make use of 4 CPU's and 8GB of RAM, instead of the default 2 CPU's and 4GB of RAM

## AWS Cloud

Some prerequisites are required before we can continue to deploy our Kubernetes infrastructure to an AWS cluster.

* You have created all the various users and permissions as required.
* You have given the users the relevant access to to the AWS services
* You have generated an access token for your AWS user

### Create the cluster

Before we can deploy our Kubernetes infrastructure we need to make sure we have created a cluster for us to deploy to. Execute the below command to create the cluster within AWS EKS

```sh
eksctl create cluster -f cluster.yml
```

### Configure cluster users

Once the cluster has been created successfully, we also need to give access to the various users accessing the cluster. Update the `cluster-auth.yml` file with the users that need access and replace the `data.mapRoles.rolearn` with the arn of the role created to manage this cluster. Execute the below command to find the ARN of the role linked to the cluster:

```sh
kubectl describe configmap -n kube-system aws-auth
```

Once the `cluster.auth.yml` file has been updated, execute the below command to apply the users to cluster.

```sh
kubectl apply -f cluster-auth.yml
```

### Getting Started

Before we proceed with creating our `Core Package` services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes`) we can proceed to create our core instant ohie deployment with the following command:

```bash
./main/k8s-mini.sh up
```

This bash script will enable `ingress` on minikube then proceed to apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR). This script will also implement the HOST mapping which is needed to access the OpenHIM Core and Console locally (on linux).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

#### View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all
```

To view all the ingress resources: (Services that are exposed publicly)

```sh
kubectl get ingress
```

To tear down this deployment use the opposing command:

```bash
./main/k8s-mini.sh down
```

To completely remove all project components use the following option:

```bash
./main/k8s-mini.sh destroy
```

The OpenHIM console will be accessible on <http://openhim-console.instant/> and core will be accessible on:

- API: <http://openhim-core.api.instant/>
- HTTPS routing: <https://openhim-core.ssl.instant/>
- HTTP routing: <http://openhim-core.non-ssl.instant/>

### Initial OpenHIM Config

We have included a useful set of scripts to initialise the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. From the `kubernetes` directory, use the following command to implement:

```bash
./importer/k8s.sh up
```

> These scripts can be duplicated and modified to implement custom imports

To clean up the remaining job and pods from a successful setup run the following:

```bash
./importer/k8s.sh clean
```

Once the config is done HAPI FHIR will be accessible on: <https://openhim-core.ssl.instant/hapi-fhir-jpaserver/fhir/>

You may test that the OpenHIM is routing requests to HAPI FHIR by running:

```bash
./test.sh
```

### Development mode for exposed services

To run in development mode, where the OpenHIM mongo database, HAPI fhir server and the MySQL database can be accessed directly through their urls, run the following command

```bash
./dev/k8s.dev.sh
```

## AWS CLI

Useful links:

- [EKS](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-console.html)

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

### Create Cluster

```sh
eksctl create cluster --name OHIE --version 1.14 --region eu-west-2 --nodegroup-name ohie-workers --node-type t3.medium --nodes 3 --nodes-min 1 --nodes-max 8 --ssh-access --ssh-public-key ~/.ssh/<rsa>.pub --managed
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

### Switch cluster context

```sh
kubectl config get-contexts

kubectl config use context <context-name>
```

### Kill cluster

```sh
eksctl delete cluster --name OHIE
```

### Deploy Scripts via Kubernetes

Run the scripts as normal after checking your cluster context is correct.

```sh
kubectl config get-contexts

./k8s-aws.sh up
```
