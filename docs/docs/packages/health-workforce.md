---
id: health-workforce
title: Health-Workforce
sidebar_label: Health-Workforce
keywords:
  - Instant OpenHIE
  - Health
  - Workforce
  - Package
description: The health-workforce package of the Instant OpenHIE
---

## Package functionality

This package sets up iHRIS and GOFR applications which are able to be queried for facility and practitioner information. It also sets up a mediator that synchronises (using mCSD) practitioner and facility data with the central FHIR store that is provided by the core package. This allows the user of the HIE to query this data to answer questions such as the following scenario:

* A doctor, Joseph,  at a rural clinic wants to refer a patient, Mousa, to an Oncologist because of a lump that they suspect may be cancerous. They are able to look up a list of specialists that offer that service in their EMR system. They choose a particular specialist and find out the facilities in which they work. A referral can now be produced by the EMR for the patient and they can be sent to that facility for their visit. (PractitionerRole -> Practitioner -> Facility)


## Deployment strategy

TODO - this package is still under development
