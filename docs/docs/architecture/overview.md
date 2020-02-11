---
id: overview
title: Overview
sidebar_label: Overview
keywords:
  - InstantHIE
  - Overview
  - Architecture
description: An overview of the InstantHIE architecture
---

Instant OpenHIE provides multiple sets of scripts to configure and setup HIE components for particular use cases and workflows. Each of these are organised into self contained packages and each of these packages may depend on other packages. This allows highly complex infrastructure to be setup in time once more packages are created.

Each of these packages will contain scripts which will setup containerised applications, configure them and ensure necessary data is loaded into them. Docker will be used to containerise each of the necessary application and to enable them to be easily deployed.
