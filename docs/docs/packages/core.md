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

The InstantHIE Core Package is the common base of the InstantHIE architecture.

This package consists of two components that support all other packages, these are:

- Interoperability Layer - [OpenHIM](http://openhim.org/)
- FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Package functionality

This package sets up a number of containers that standup the OpenHIM as well as the HAPI FHIR server. It also configures the OpenHIM with a channel so that the HAPI FHIR server may be accessed through the OpenHIM.

To use the HAPI FHIR server from an external point of care application you may access it through the OpenHIM at a URL like the following:

```bash
curl <openhim_core_transaction_api_url>/hapi-fhir-jpaserver/fhir/Patient
```

> The **openhim_core_transaction_api_url** is displayed in the output of the startup script

Accessing the services created by this package:
- **OpenHIM**
  - Console: Displayed in the output of the startup script
  - Username: **root@openhim.org**
  - Password: **instant101**
- **HAPI FHIR**
  - This service should not be publicly accessible and only accessed via the Interoperability Layer

## Deployment strategy

The OpenHIM was already dockerised so we were able to re-use those images for our work in the core package. HAPI FHIR didn't have official dockerfiles available, however, several community contributed optiosn existed. We chose to use what seemed like the most robust option.

We supplied docker compose files for the setup and configuration of these application. We chose to split the docker-compose file into three file:

1. a main docker-compose.yml file that setup the base applications
2. a config docker-compose.config.yml file that when executed configures the OpenHIM with a channel route to HAPI FHIR
3. a dev docker-compose.dev.yml file that exposes all open port to the host for easy debuging, this would be insecure in a production environment

For Kubernetes we created deployment and service resource files for each of the components of each application and hooked these up with a kustomization.yml file for easy deployment. To import configuration into the OpenHIM we use job resources that only execute when the OpenHIM core is up. This is done by using init container to wait for the OpenHIM core port to become available.

For importing config we use a custom image which is just a node.js container that can run node.js scripts that we define. It also has a `wait-on` module installed to allow it to wait on certain prot being available before executing.

## Core Package Dev guide

For testing purposes this package can be run independently. The below are some notes of how to do this. The recommended way to run Instant OpenHIE is described [here](../introduction/getting-started).

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

## Useful compose flags

Some additional flags can be passed to the `docker-compose` command making it a bit easier to work with.

- `-d`: Run the services in a detached mode. This means that when you close or exit your terminal, the services will still be running in the background.
- `-f`: Specify the location of the `docker-compose` file to be executed. Omitting this flag will look for the default `docker-compose.yml` file.
- `--force-recreate`: This will force the container/image to be re-created if a newer version is found. This is useful when a new image has been released but not yet pulled onto the host machine.

```bash
docker-compose up -d --force-recreate
```

## Environment configuration

By running the above command to get started with the Core package we create all the services that need to be defined, but this script might have some limitations depending on the type of environment you want to run the configuration

Additional `docker-compose` files are available for extra environment configuration

- **docker-compose.yml**: Main `docker-compose` script to create the services
- **docker-compose.dev.yml**: Development `docker-compose` script to override some of the default configurations to be used in a development environment (Open service ports for access etc)

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

> If you are working with multiple cluster hosting clients you may need to switch cluster context for your deployment. See the below commands:

```sh
kubectl config get-contexts
kubectl config use-context <context-name>
```

## Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once installed, we can start the minikube service by executing the below command:

```sh
minikube start --cpus 4 --memory 8192
```

This also updates the VM settings to make use of 4 CPU's and 8GB of RAM, instead of the default 2 CPU's and 4GB of RAM

## Getting Started

Useful Links:

- [Kubectl Cheat Sheet](https://www.digitalocean.com/community/cheatsheets/getting-started-with-kubernetes-a-kubectl-cheat-sheet)

Before we proceed with creating our `Core Package` services, we need to ensure we are in the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes`) we can proceed to create our core instant ohie deployment with the following command:

```bash
./main/k8s.sh up
```

The OpenHIM console url will be displayed in the terminal output when the script completes. The Url may take a few minutes to become active as the pod may not be fully initialized yet.

This bash script will apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

## View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all --all-namespaces
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
</TabItem>
</Tabs>
