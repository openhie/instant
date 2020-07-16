'use strict'

const axios = require('axios')

const { Given, When, Then } = require('cucumber')
const { expect } = require('chai')

const OPENHIM_PROTOCOL = process.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_API_PORT = process.OPENHIM_API_PORT || '5001'


Given('a patient, Jane Doe, exists in the FHIR server', async function () {
  const options = {
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_API_PORT}/hapi-fhir-jpaserver/fhir/Patient`,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom test`
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
