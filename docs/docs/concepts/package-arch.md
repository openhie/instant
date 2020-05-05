---
id: package-arch
title: Packages
sidebar_label: Packages
keywords:
  - InstantHIE
  - Packages
  - Architecture
description: A description of the instant OpenHIE package architecture
---

import useBaseUrl from '@docusaurus/useBaseUrl';

The fundamental concept of Instant OpenHIE is that it can be extended to support additional use cases and workflows. This is achieved through packages. A [core package](packages/core.md) has been produced which other packages will all derive from. A package will either extend directly from the core package or from another existing package.

Each package will contain the following sorts of technical artefacts:

- Docker compose scripts for setting up the applications required for this package’s use cases and workflows
- Kubernetes scripts for setting up the applications required for this package’s use cases and workflows
- Configuration scripts to setup required configuration metadata
- Extensions to the test harness to test the added use cases with test data

The below diagram shows how packages will extend off each other to add use cases of increasing complexity.

<div class="text--center">
  <img alt="Package architecture" src={useBaseUrl('img/package-arch.png')} />
</div>

A number of essential packages are bundled with Instant OpenHIE, however, it is designed to be extended for implementation specific need. See [creating packages](../how-to/creating-packages) for more information on how to create your own packages to extend instant OpenHIE.
