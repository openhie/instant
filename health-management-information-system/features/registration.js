'use strict'

const {Given, When, Then, setDefaultTimeout} = require('@cucumber/cucumber')

const {
  verifyDhis2Running,
  getSystemInfo,
} = require('./utils')

setDefaultTimeout(10000)

Given('that dhis is set up', verifyDhis2Running)

Then('DHIS2 should respond to an authenticated API request', getSystemInfo)
