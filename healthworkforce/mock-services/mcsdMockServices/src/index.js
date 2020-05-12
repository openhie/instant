'use strict'

const express = require('express')
const fs = require('fs')

const config = require('./config').getConfig()
const logger = require('./logger')

const practitionerResource = JSON.parse(
  fs.readFileSync('./resources/bundle-Practitioner.json')
)
const practitionerRoleResource = JSON.parse(
  fs.readFileSync('./resources/bundle-PractitionerRole.json')
)
const locationResource = JSON.parse(
  fs.readFileSync('./resources/bundle-Location.json')
)
const organizationResource = JSON.parse(
  fs.readFileSync('./resources/bundle-Organization.json')
)
const fhirTransactionBundleResultResource = JSON.parse(
  fs.readFileSync('./resources/bundle-transaction-response.json')
)

const app = express()

app.use((req, res, next) => {
  logger.info(`Request url": ${req.url}`)
  next()
})

app.get('/gofr-location-mock/_history', function(req, res) {
  logger.info(`Fetching GOFR Location data with query paramter`)
  res.json(locationResource)
})

app.get('/gofr-organization-mock/_history', function(req, res) {
  logger.info(`Fetching GOFR Organization data with query paramter`)
  res.json(organizationResource)
})

app.get('/ihris-practitioner-mock/_history', function(req, res) {
  logger.info(`Fetching iHRIS practitioner data with query paramter`)
  res.json(practitionerResource)
})

app.get('/ihris-practitionerRole-mock/_history', function(req, res) {
  logger.info(`Fetching iHRIS practitionerRole data with query paramter`)
  res.json(practitionerRoleResource)
})

app.post('/fhir-mock', function(req, res) {
  logger.info(`Processing FHIR transaction bundle...`)
  res.json(fhirTransactionBundleResultResource)
})

app.listen(config.port, () => {
  logger.info(`Server listening on port ${config.port}...`)
})
