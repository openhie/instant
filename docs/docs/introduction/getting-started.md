---
id: getting-started
title: Getting started
sidebar_label: Getting started
keywords:
  - InstantHIE
  - Running Instant OpenHIE
description: How to run Instant OpenHIE
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Instant OpenHIE is designed to be easily downloaded and run by anyone on any platform. To get started you will need to first install docker and docker-compose which Instant OpenHIE uses to run all of the necessary services and applications. The following links guide you on how to do this for your platform:

- [Docker engine](https://docs.docker.com/install/)
- [Docker compose](https://docs.docker.com/compose/install/)

## Running using the interactive app (recommended)

Once these are installed you can download the latest Instant OpenHIE executable from github. Look for the latest release and download the executable that matches your platform under the assets section of the release:

- [Instant OpenHIE](https://github.com/openhie/instant/releases)

Once the download is complete, move it somewhere memorable. From there you may execute it and you will be presented with a command-line interface. Use the arrow keys to navigate it and choose 'Start Instant OpenHIE'. This will initiate a download that may take some time. Once completed it will then startup the default packages which each include a number of services/application. Have a look at the existing packages on the sidebar to see what packages are available. A web page will be opened with links to the started services/application as well as links to their documentation.

Congratulations, you have successfully started up Instant OpenHIE. Follow the instruction on the web page to stop the services when you are done.

## Running using the bash scripts directly (useful for automation)

To start all the default the InstantHIE packages, select a deployment platform below. Have a look at the existing packages on the sidebar to see what packages are available.

<Tabs
  defaultValue="dockerCompose"
  values={[
    { label: 'Docker Compose', value: 'dockerCompose' },
    { label: 'Kubernetes', value: 'kubernetes' }
  ]}
>
<TabItem value="dockerCompose">

Before proceeding, ensure you are in the main `/instant` directory containing `instant.sh` script. To start up the system, execute the following command:

```sh
./instant.sh up docker
```

To stop the running containers execute the following:

```sh
./instant.sh down docker
```

Then to remove the containers run the command below. However, make sure you have **stopped all the containers before trying to delete** them. This action will also delete any volumes created by the containers.

```sh
./instant.sh destroy docker
```

</TabItem>
<TabItem value="kubernetes">

Before proceeding, ensure you are in the main `/instant` directory containing `instant.sh` script. Then you can execute the command below. This commands will output urls from which you can access your instantHIE instance.

```sh
./instant.sh up k8s
```

To delete all the deployment related pods, run the command below. This command will leave services, and volumes intact.

```sh
./instant.sh down k8s
```

To delete the entire instantHIE system run the command below.

```sh
./instant.sh destroy k8s
```

</TabItem>
</Tabs>
