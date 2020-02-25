'use strict'

const fs = require('fs')
const https = require('https')
const path = require('path')

const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'openhim-core'
const OPENHIM_API_PASSWORD =
  process.env.OPENHIM_API_PASSWORD || 'openhim-password'
const OPENHIM_API_PORT = process.env.OPENHIM_API_PORT || 8080
const OPENHIM_API_USERNAME =
  process.env.OPENHIM_API_USERNAME || 'root@openhim.org'

const authHeader = new Buffer.from(
  `${OPENHIM_API_USERNAME}:${OPENHIM_API_PASSWORD}`
).toString('base64')

const jsonData = JSON.parse(
  fs.readFileSync(path.resolve(__dirname, 'openhim-import.json'))
)

const data = JSON.stringify(jsonData)

const options = {
  protocol: 'https:',
  hostname: OPENHIM_API_HOSTNAME,
  port: OPENHIM_API_PORT,
  path: '/metadata',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': data.length,
    Authorization: `Basic ${authHeader}`
  }
}

const req = https.request(options, res => {
  if (res.statusCode == 401) {
    throw new Error(`Incorrect OpenHIM API credentials`)
  }

  if (res.statusCode != 201) {
    throw new Error(`Failed to import OpenHIM config: ${res.statusCode}`)
  }

  console.log('Successfully Imported OpenHIM Config')
})

req.on('error', error => {
  console.error('Failed to import OpenHIM config: ', error)
})

req.write(data)
req.end()
