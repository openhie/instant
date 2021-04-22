'use strict'

const {Given, When, Then, setDefaultTimeout} = require('@cucumber/cucumber')

const {
  importMetaData,
  createPatient,
  verifyPatientExists,
  deletePatient
} = require('./utils')

setDefaultTimeout(20000)

Given('that dhis is set up and the metadata import has been done', importMetaData)

When('a patient is created', createPatient)

Then('the patient should exist in DHIS', verifyPatientExists)

Then('the patient should then be deleted', deletePatient)
