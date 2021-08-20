'use strict'

const axios = require('axios')

const { AfterAll, Given, Then, When } = require('@cucumber/cucumber')
const { expect } = require('chai')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'
const OPENHIM_MEDIATOR_API_PORT =
  process.env.OPENHIM_MEDIATOR_API_PORT || '8080'
const CUSTOM_TOKEN_ID = process.env.CUSTOM_TOKEN_ID || 'test'
const BASIC_AUTH_HEADER =
  process.env.BASIC_AUTH_HEADER || 'Basic cm9vdEBvcGVuaGltLm9yZzppbnN0YW50MTAx'

// Save test Patient resource ID for post test cleanup
let hapiFhirPatientID

Given('a patient, Jane Doe, exists in the FHIR server', async function () {
  const checkPatientExistsOptions = {
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/Patient?identifier:value=test`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
      'Cache-Control': 'no-cache'
    }
  }

  const checkPatientExistsResponse = await axios(checkPatientExistsOptions)

  if (checkPatientExistsResponse.data.total === 0) {
    console.log(
      `\nPatient record for Jane Doe does not exist. Creating Patient`
    )
    const options = {
      url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/Patient`,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
        'Cache-Control': 'no-cache'
      },
      data: {
        resourceType: 'Patient',
        name: {
          use: 'temp',
          family: 'Doe',
          given: ['Jane']
        },
        identifier: {
          use: 'temp',
          value: 'test'
        },
        gender: 'male'
      }
    }

    const createPatientResponse = await axios(options)

    expect(createPatientResponse.status).to.eql(201)

    hapiFhirPatientID = createPatientResponse.data.id
  } else if (checkPatientExistsResponse.data.total === 1) {
    console.log(`\nPatient record for Jane Doe already exists`)

    expect(
      checkPatientExistsResponse.data.entry[0].resource.name[0].given
    ).to.eql(['Jane'])
    expect(
      checkPatientExistsResponse.data.entry[0].resource.name[0].family
    ).to.eql('Doe')

    hapiFhirPatientID = checkPatientExistsResponse.data.entry[0].resource.id
  } else {
    // Previous test data should have been cleaned out
    throw new Error(
      `Multiple Patient records for Jane Doe exist: ${checkPatientExistsResponse.data.total}`
    )
  }
})

Given('an authorised client, Alice, exists in the OpenHIM', async function () {
  process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0
  const checkClientExistsOptions = {
    url: `https://${OPENHIM_API_HOSTNAME}:${OPENHIM_MEDIATOR_API_PORT}/clients`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: BASIC_AUTH_HEADER
    }
  }

  const checkClientExistsResponse = await axios(checkClientExistsOptions)

  let createClient = true

  for (let client of checkClientExistsResponse.data) {
    if (client.clientID === 'test-harness-client') {
      expect(client.name).to.eql('Alice')
      createClient = false
      break
    }
  }

  if (createClient) {
    console.log(`\nThe test Harness Client does not exist. Creating Client`)
    const options = {
      url: `https://${OPENHIM_API_HOSTNAME}:${OPENHIM_MEDIATOR_API_PORT}/clients`,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: BASIC_AUTH_HEADER
      },
      data: {
        roles: ['instant'],
        clientID: 'test-harness-client',
        name: 'Alice',
        customTokenID: 'test-harness-token'
      }
    }

    const response = await axios(options)
    expect(response.status).to.eql(201)
  } else {
    console.log(`\nThe Test Harness Client (Alice) already exists`)
  }
})

When('Alice searches for a patient', async function () {
  const checkPatientExistsOptions = {
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/Patient?identifier:value=test`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom test-harness-token`,
      'Cache-Control': 'no-cache'
    }
  }
  const checkPatientExistsResponse = await axios(checkPatientExistsOptions)
  expect(checkPatientExistsResponse.status).to.eql(200)
  this.setTo(checkPatientExistsResponse.data.entry[0])
})

Then('Alice is able to get a result', function () {
  expect(this.searchResults.resource.resourceType).equal('Patient')
  expect(this.searchResults.resource.identifier[0].value).equal('test')
  expect(this.searchResults.resource.name[0].given[0]).equal('Jane')
  expect(this.searchResults.resource.name[0].family).equal('Doe')
})

When('Malice searches for a patient', async function () {
  const checkPatientExistsOptions = {
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/Patient?identifier:value=test`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom invalid-token`,
      'Cache-Control': 'no-cache'
    },
    validateStatus: (status) => {
      return status >= 400
    }
  }
  const checkPatientExistsResponse = await axios(checkPatientExistsOptions)
  expect(checkPatientExistsResponse.status).to.eql(401)
  this.setTo(checkPatientExistsResponse.data)
})

Then('Malice is NOT able to get a result', function () {
  expect(this.searchResults).to.eql('')
})

AfterAll(async function () {
  if (hapiFhirPatientID) {
    console.log(`\nDeleting FHIR test Patient record`)
    const deletePatientOptions = {
      url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/Patient/${hapiFhirPatientID}`,
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
        'Cache-Control': 'no-cache'
      }
    }

    const deletePatientResponse = await axios(deletePatientOptions)
    expect(deletePatientResponse.status).to.eql(200)
  }

  process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0
  const checkClientExistsOptions = {
    url: `https://${OPENHIM_API_HOSTNAME}:${OPENHIM_MEDIATOR_API_PORT}/clients`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: BASIC_AUTH_HEADER
    }
  }

  const checkClientExistsResponse = await axios(checkClientExistsOptions)

  let clientObjectId
  // Previous test data should have been cleaned out
  for (let client of checkClientExistsResponse.data) {
    if (client.clientID === 'test-harness-client') {
      clientObjectId = client._id
      break
    }
  }

  if (clientObjectId) {
    console.log(`\nDeleting OpenHIM test Client record`)
    const deleteClientOptions = {
      url: `https://${OPENHIM_API_HOSTNAME}:${OPENHIM_MEDIATOR_API_PORT}/clients/${clientObjectId}`,
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: BASIC_AUTH_HEADER
      }
    }

    const deleteClientResponse = await axios(deleteClientOptions)
    expect(deleteClientResponse.status).to.eql(200)
  }
})
