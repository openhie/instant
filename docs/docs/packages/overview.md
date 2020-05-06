---
id: overview
title: Overview of Packages
sidebar_label: Overview
keywords:
  - Instant OpenHIE
  - Packages
description: An overview of all existing Instant OpenHIE packages
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

So far Instant OpenHIE consists of two packages:

- Core
- Health Worker Force

## Getting Started

To start all the Instant OpenHIE packages, select a deployment platform below.

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

Before proceeding, ensure you are in the main `/instant` directory containing `instant.sh` script. Then you can execute the command below. This commands will output urls from which you can access your Instant OpenHIE instance.

```sh
./instant.sh up k8s
```

To delete all the deployment related pods, run the command below. This command will leave services, and volumes intact.

```sh
./instant.sh down k8s
```

To delete the entire Instant OpenHIE system run the command below.

```sh
./instant.sh destroy k8s
```

</TabItem>
</Tabs>
