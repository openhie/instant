---
id: overview
title: Overview
sidebar_label: Overview
keywords:
  - Instant OpenHIE
  - Overview
  - Architecture
description: An overview of the Instant OpenHIE architecture
---

Instant OpenHIE provides multiple sets of scripts to configure and setup HIE components for particular OpenHIE use cases and workflows. These scripts and configuration are organised into self-contained packages. Each of these packages may depend on other packages which allows highly complex infrastructure to be setup instantly once a number of packages have been created.

Each of these packages will contain scripts which will setup containerised applications, configure them and ensure necessary data is loaded into them. Docker will be used to containerise each of the necessary applications and to enable them to be easily deployed.
