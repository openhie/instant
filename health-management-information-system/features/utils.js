'use strict'

const axios = require('axios')
const fs = require('fs')
const path = require('path')
const zlib = require('zlib')

const DHIS2_PROTOCOL = process.env.DHIS2_PROTOCOL || 'http'
const DHIS2_API_HOSTNAME = process.env.DHIS2_API_HOSTNAME || 'localhost'
const DHIS2_API_PASSWORD = process.env.DHIS2_API_PASSWORD || 'district'
const DHIS2_API_PORT = process.env.DHIS2_API_PORT || 8081
const ORG_UNIT = process.env.ORG_UNIT || 'rgNg6HkgJcs'
const DHIS2_API_USERNAME = process.env.DHIS2_API_USERNAME || 'admin'
const DHIS2_METADATA_FILENAME =
  process.env.DHIS2_METADATA_FILENAME || 'metadata.json.gz'

const authHeader = new Buffer.from(
  `${DHIS2_API_USERNAME}:${DHIS2_API_PASSWORD}`
).toString('base64')
const patientId = '12305759751'
const programId = 'eOMz8AkNk2a'
let trackedEntityId

exports.importMetaData = async () => {
  const fileContents = fs.createReadStream(
    path.resolve(__dirname, DHIS2_METADATA_FILENAME)
  )
  const unzip = zlib.createGunzip()

  function streamToString(stream) {
    const chunks = []
    return new Promise((resolve, reject) => {
      stream.on('data', (chunk) => chunks.push(chunk))
      stream.on('error', reject)
      stream.on('end', () => resolve(Buffer.concat(chunks).toString('utf8')))
    })
  }

  const configData = await streamToString(fileContents.pipe(unzip))

  console.log(
    'Importing DHIS2 data... Byte Length: ',
    Buffer.byteLength(configData)
  )

  const options = {
    url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/api/metadata`,
    method: 'post',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': Buffer.byteLength(configData),
      Authorization: `Basic ${authHeader}`
    },
    maxContentLength: 100000000,
    data: configData,
    params: {
      mergeMode: 'REPLACE',
      atomicMode: 'NONE'
    }
  }

  try {
    const response = await axios(options)

    console.log(
      `Successfully Imported DHIS2 Config.\n\nImport summary:${JSON.stringify(
        response.data.stats
      )}`
    )
  } catch (error) {
    throw new Error(`Failed to import DHIS2 config: ${error.message}`)
  }
}

exports.deletePatient = async () => {
 const result = await axios({
      url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/api/trackedEntityInstances/${trackedEntityId}`,
      method: 'DELETE',
      headers: {
      Authorization: `Basic ${authHeader}`
      }
  })

  if (result.status !== 200 && !result.data.response.importCount.deleted) {
    throw Error('Clean up failed - patient')
  }
}

exports.verifyPatientExists = async () => {
  const result = await axios({
    url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/api/trackedEntityInstances`,
    method: 'GET',
    headers: {
      "Content-Type": 'application/json',
      "Cache-Control": 'no-cache',
      Authorization: `Basic ${authHeader}`
    },
    params: {
      ouMode: 'ALL',
      program: programId,
      filter: `VuoMp8yYPYz:EQ:${patientId}`
    }
  })

  if (result.status != 200 || !result.data.trackedEntityInstances.length) {
    throw Error('Patient verification failed')
  }

  trackedEntityId = result.data.trackedEntityInstances[0].trackedEntityInstance
}

exports.createPatient = async () => {
  const data = {
    trackedEntityType: "hc6rGplhGZC",
    orgUnit: ORG_UNIT,
    attributes: [
      {
        attribute: "HyZ808XBvrX",
        value: true
      },
      {
        attribute: "bUXlTnm6OK8",
        value: "Test clinic"
      },
      {
        attribute: 'xIPU67t6kzK',
        value: "Samantha"
      },
      {
        attribute: "VuoMp8yYPYz",
        value: patientId
      }
    ],
    enrollments: [{
      orgUnit: ORG_UNIT,
      program: 'eOMz8AkNk2a',
      enrollmentDate: new Date(),
      incidentDate: new Date()
    }]
  }

  const result = await axios({
    url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/api/trackedEntityInstances`,
    method: 'POST',
    headers: {
      "Content-Type": 'application/json',
      "Cache-Control": 'no-cache',
      Authorization: `Basic ${authHeader}`
    },
    data: JSON.stringify(data)
  })


  if (result.status != 200) {
    throw Error('Creation of patient failed')
  }
}
