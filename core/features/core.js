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
      Authorization: `Custom ${CUSTOM_TOKEN_ID}`
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
        Authorization: `Custom ${CUSTOM_TOKEN_ID}`
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

Given('an authorised client, Alice, exists in the OpenHIM', function () {
  return 'pending'
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
