'use strict'

const { Given, When, Then, setDefaultTimeout } = require('@cucumber/cucumber')

// Set timeout for the steps. The default timeout of 5000 is not enough as the process take a while
setDefaultTimeout(30000)

const {
  gofrMockServiceLocation, gofrMockServiceOrganization,
  ihrisMockServicePractitioner, ihrisMockServicePractitionerRole, triggerSync,
  verifyLocationExistsAndCleanup, verifyOrganizationExistsAndCleanup,
  verifyPractitionerExistsAndCleanup, verifyPractitionerRoleExistsAndCleanup
} = require('./utils')

Given('a new practitioner, Dr Bob, has been added in iHRIS', ihrisMockServicePractitioner)
Given('a new practitionerRole has been added in iHRIS', ihrisMockServicePractitionerRole)
Given('a new location, GoodHealth Clinic, has been added in GOFR', gofrMockServiceLocation)
Given('a new organization, Clinical Lab, has been added in GOFR', gofrMockServiceOrganization)

When('the sync is triggered', triggerSync)

Then('the new practitioner, Dr Bob, can be found in the FHIR server', verifyPractitionerExistsAndCleanup)
Then('the new practitionerRole can be found in the FHIR server', verifyPractitionerRoleExistsAndCleanup)
Then('the new location, GoodHealth Clinic, can be found in the FHIR server', verifyLocationExistsAndCleanup)
Then('the new organization, Clinical Lab, can be found in the FHIR server', verifyOrganizationExistsAndCleanup)
