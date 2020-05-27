'use strict'

const fs = require('fs')
const http = require('http')

const config = require('./config').getConfig()

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

const server = http.createServer((req, res) => {
  res.writeHead(200, {'Content-Type': 'application/json'})

  if (req.method === 'GET' && req.url === '/gofr-location-mock/_history' ) {
    res.write(JSON.stringify(locationResource))
  } else if (req.method === 'GET' && req.url === '/ihris-practitioner-mock/_history') {
    res.write(JSON.stringify(practitionerResource))
  } else if (req.method === 'GET' && req.url === '/gofr-organization-mock/_history') {
    res.write(JSON.stringify(organizationResource))
  } else if (req.method === 'GET' && req.url === '/ihris-practitionerRole-mock/_history') {
    res.write(JSON.stringify(practitionerRoleResource))
  } else if (req.method === 'POST' && req.url === '/fhir-mock') {
    res.write(JSON.stringify(fhirTransactionBundleResultResource))
  } else {
    res.writeHead(404, {'Content-Type': 'application/json'})
    res.write(JSON.stringify({error: 'Not Found'}))
  }
  res.end()
})

server.listen(config.port)
console.log('Server started on port: ', config.port)
