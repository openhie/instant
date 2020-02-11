---
id: package-arch
title: Packages
sidebar_label: Packages
keywords:
  - InstantHIE
  - Overview
  - Architecture
description: A description of the instant OpenHIE package architecture
---

The fundamental concept of Instant OpenHIE is that it can be extended to support additional use cases and workflows. This will be achieved through packages. A **core package** will be produced in this first phase which other packages will all derive from. A package will either extend directly from the core package or from another package.

Each package will contain the following sorts of technical artefacts:
* Docker compose scripts for setting up the applications required for this package’s use cases and workflows
* Kubernetes scripts for setting up the applications required for this package’s use cases and workflows
* Configuration scripts to setup required configuration metadata
* Extensions to the test harness to test the added use cases with test data

The below diagram shows how packages will extend off each other to add use cases of increasing complexity.

![](package-arch.png)
