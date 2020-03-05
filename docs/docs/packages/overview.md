---
id: overview
title: Overview of Packages
sidebar_label: Overview
keywords:
  - InstantHIE
  - Packages
description: An overview of all existing InstantHIE packages
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

So far InstantHIE consists of two packages:

- Core
- Health Worker Force

## Getting Started

To start all the InstantHIE packages, select a deployment platform below.

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

Then to remove the containers run the command below. However, make sure you have stopped all the containers before trying to delete them. This action will also delete any volumes created by the containers.

```sh
./instant.sh destroy docker
```

</TabItem>
<TabItem value="kubernetes">

Before proceeding, ensure you are in the main `/instant` directory containing `instant.sh` script. Then you can execute the following command:

```sh
./instant.sh up k8s
```

</TabItem>
</Tabs>
