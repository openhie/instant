'use strict'

const axios = require('axios')

const { Given, When, Then } = require('cucumber')
const { expect } = require('chai')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'
const OPENHIM_MEDIATOR_API_PORT =
  process.env.OPENHIM_MEDIATOR_API_PORT || '8080'
const CUSTOM_TOKEN_ID = process.env.CUSTOM_TOKEN_ID || 'test'

Given('a patient, Jane Doe, exists in the FHIR server', async function () {
  const checkPatientExistsOptions = {
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/Patient?identifier:value=test`,
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
      `Patient record for Jane Doe does not exist. Creating Patient...`
    )
    const options = {
      url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/Patient`,
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

    const response = await axios(options)
    expect(response.status).to.eql(201)
  } else if (checkPatientExistsResponse.data.total === 1) {
    console.log(
      `Patient record for Jane Doe already exists...`
    )
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
      Authorization: `Basic cm9vdEBvcGVuaGltLm9yZzppbnN0YW50MTAx`
    }
  }

  const checkClientExistsResponse = await axios(checkClientExistsOptions)

  let createClient = true
  // Previous test data should have been cleaned out
  for (let client of checkClientExistsResponse.data) {
    if (client.clientID === 'test-harness-client') {
      createClient = false
      break
    }
  }

  if (createClient) {
    console.log(`The test Harness Client does not exist. Creating Client...`)
    const options = {
      url: `https://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/clients`,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Basic cm9vdEBvcGVuaGltLm9yZzppbnN0YW50MTAx`
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
    console.log(`The Test Harness Client (Alice) already exists...`)
  }
})

When('Alice searches for a patient', function () {
  return 'pending'
})

Then('Alice is able to get a result', function () {
  return 'pending'
})

When('Malice searches for a patient', function () {
  return 'pending'
})

Then('Malice is NOT able to get a result', function () {
  return 'pending'
})
