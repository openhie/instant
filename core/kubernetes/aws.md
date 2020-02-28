# AWSCLI

Useful links:

- [EKS](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-console.html)

## Install AWS Cli

```sh
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

unzip awscliv2.zip

sudo ./aws/install
```

## Configure with your AWS token details

```sh
aws configure
```

## Install EksCtl

```sh
curl --silent --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp

sudo mv /tmp/eksctl /usr/local/bin

eksctl version
```

## Create Cluster

```sh
eksctl create cluster --name OHIE --version 1.14 --region eu-west-2 --nodegroup-name ohie-workers --node-type t3.medium --nodes 3 --nodes-min 1 --nodes-max 8 --ssh-access --ssh-public-key ~/.ssh/<rsa>.pub --managed
```

## Access an existing cluster

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

## Switch cluster context

```sh
kubectl config get-contexts

kubectl config use context <context-name>
```

## Kill cluster

```sh
eksctl delete cluster --name OHIE
```

## Deploy Scripts via Kubernetes

Run the scripts as normal after checking your cluster context is correct.

```sh
kubectl config get-contexts

./k8s-aws.sh up
```
