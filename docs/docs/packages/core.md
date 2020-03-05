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

### Local installation

To be able to deploy a Kubernetes clusters, you either need a Cloud Services provider account that provides a Kubernetes service or you could deploy to a locally run Kubernetes instance called `Minikube`

Follow the steps outlined [here](https://kubernetes.io/docs/tasks/tools/install-minikube/) to install minikube on your machine.

Once minikube has been installed, we need to start it so that we can deploy out Kubernetes scripts to it. Execute the below command to start minikube

```sh
minikube start --cpus 4 --memory 8192
```

Once the cluster is created, you will be able to see the Kubernetes cluster you are deploying to be executing the below command

```sh
kubectl config current-context
```

### Cloud Services installation

You could create this Kubernetes cluster on a Cloud Services provider like, AWS, Google Cloud and Azure. This will however incur some costs for making use of these services.

If you do have an account with once of these, you will need to ensure you setup a cluster using the relevant provider SDK.

AWS Example: Once the relevant SDK's have been downloaded for AWS, you can execute the below command to create the `instanthie` cluster. Creating cluster can take some time to complete.

```sh
eksctl create cluster -f cluster.yml
```

Once the cluster is created, you will be able to see the Kubernetes cluster you are deploying to be executing the below command

```sh
kubectl config current-context
```

### Deploying Kubernetes scripts

Before we proceed with creating our Core package components, we need to ensure we are in the correct directory containing our `kubernetes` script.

Once you are in the correct working directory (`core/kubernetes/`) we can proceed to execute our `k8s` script by running the below command which will create all the services and output the results.

```bash
./main/k8s.sh up
```

We can also execute the impoert script which will configure some OpenHIM settings for us to be able to communicate successfully with the FHIR server

```sh
./importer/k8s.sh up
```

Once the Kubernetes scripts have finished executing we are able to connect to our Core package through the OpenHIM.

The link to the OpenHIM will be printed in the output of the `./main/k8s.sh` script

</TabItem>
</Tabs>

## Accessing the services

- **OpenHIM**
  - Console: <http://localhost:9000> - or displayed in the output of the startup script
  - Username: **root@openhim.org**
  - Password: **instant101**
- **HAPI FHIR**
  - This service should not be publicly accessible and only accessed via the Interoperability Layer

## Testing the Core package

As part of the Core package setup we also do some initial importation of config for connecting the services together.

- OpenHIM: Import a public channel configuration that routes requests to the HAPI FHIR services
- HAPI FHIR: _Not config import yet_

For testing this Core package we will be making use of `curl` for sending our request, but any client could be used to achieve the same result.

Execute the below `curl` request to successfully route a request through the OpenHIM to query the HAPI FHIR server.

```bash
curl http://localhost:5001/hapi-fhir-jpaserver/fhir/Patient
```
