'use strict'

const axios = require('axios')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'
const CUSTOM_TOKEN_ID = process.env.CUSTOM_TOKEN_ID || 'test'
const MOCK_SERVER_HOST = process.env.MOCK_SERVER_HOST || 'localhost'
const MOCK_SERVER_PORT = process.env.MOCK_SERVER_PORT || '4000'

const testLocation = {
  resourceType: 'Location',
  id: '2test',
  name: 'GoodHealth Clinic',
  identifier: [
    {
      use: 'temp',
      value: 'testLocation'
    }
  ]
}

const testOrganization = {
  resourceType: 'Organization',
  id: '2test',
  identifier: [
    {
      system: 'http://www.acme.org.au/units',
      value: 'testOrganization'
    }
  ],
  name: 'Clinical Lab'
}

const testPractitioner = {
  resourceType: 'Practitioner',
  id: 'P10004test',
  active: true,
  identifier: [
    {
      system: 'http://www.acme.org.au/units',
      value: 'testPractitioner'
    }
  ],
  name: [
    {
      use: 'official',
      text: 'Bob Murray',
      given: ['Bob'],
      family: 'Murray'
    }
  ]
}

const testPractitionerRole = {
  resourceType: 'PractitionerRole',
  id: 'PR10001Test',
  active: true,
  practitioner: {
    reference: 'Practitioner/P10004test'
  },
  identifier: [
    {
      system: 'http://www.acme.org.au/units',
      value: 'testPractitionerRole'
    }
  ]
}

exports.triggerSync = async () => {
  await axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/mcsd-trigger`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then((res) => {
      console.log(`\n Triggered sync`)
    })
    .catch((err) => {
      console.error(`Trigger Failed. Clean up needed. ${err}`)
      throw err
    })
}

exports.ihrisMockServicePractitioner = async () => {
  const name = 'Practitioner'

  await verifyResourceDoesNotExistInFHIR(
    name,
    testPractitioner.identifier[0].value
  )

  // create new practitioner Dr Bob on ihris mock server
  const response = await createResource(name, testPractitioner)

  if (response.status != 201)
    throw Error(`Failed to create ${name} for testing`)

  console.log(
    '\n Created Practitioner resource (Dr Bob) on mock ihris mock server'
  )
}

exports.ihrisMockServicePractitionerRole = async () => {
  const name = 'PractitionerRole'

  await verifyResourceDoesNotExistInFHIR(
    'PractitionerRole',
    testPractitionerRole.identifier[0].value
  )

  // create new practitionerRole for Dr Bob on ihris mock server
  const response = await createResource(name, testPractitionerRole)

  if (response.status != 201)
    throw Error(`Failed to create ${name} for testing`)

  console.log(
    '\n Created PractitionerRole resource (for Dr Bob) on mock ihris mock server'
  )
}

exports.gofrMockServiceLocation = async () => {
  const name = 'Location'

  await verifyResourceDoesNotExistInFHIR(
    'Location',
    testLocation.identifier[0].value
  )

  // create new location on gofr mock server
  const response = await createResource(name, testLocation)

  if (response.status != 201)
    throw Error(`Failed to create ${name} for testing`)

  console.log(
    '\n Created Location resource (GoodHealth Clinic) on mock gofr mock server'
  )
}

exports.gofrMockServiceOrganization = async () => {
  const name = 'Organization'

  await verifyResourceDoesNotExistInFHIR(
    'Organization',
    testOrganization.identifier[0].value
  )

  // create new organization on gofr mock server
  const response = await createResource(name, testOrganization)

  if (response.status != 201)
    throw Error(`Failed to create ${name} for testing`)

  console.log(
    '\n Created Organization resource (Clinical Lab) on mock gofr mock server...'
  )
}

const getResource = (resource, identifierValue) => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/${resource}?identifier:value=${identifierValue}`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
      'Cache-Control': 'no-cache'
    }
  })
}

const removeResource = (resource, identifierValue) => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/fhir/${resource}?identifier:value=${identifierValue}`,
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
      'Cache-Control': 'no-cache'
    }
  })
}

const verifyResourceDoesNotExistInFHIR = async (resource, resourceId) => {
  const response = await getResource(resource, resourceId)

  if (response.data.total > 0)
    throw Error(
      `Test aborted! ${resource} resource (identifier-value: ${resourceId}) used in test already exists and will be removed from the FHIR server`
    )
}

const createResource = (resourceName, resource) => {
  return axios({
    url: `http://${MOCK_SERVER_HOST}:${MOCK_SERVER_PORT}/create-resource/${resourceName}`,
    method: 'POST',
    data: resource
  })
}

const deleteResource = (resourceName) => {
  return axios({
    url: `http://${MOCK_SERVER_HOST}:${MOCK_SERVER_PORT}/cleanup/${resourceName}`,
    method: 'DELETE'
  })
}

exports.verifyPractitionerExistsAndCleanup = async () => {
  const resource = 'Practitioner'
  const identifierValue = testPractitioner.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name[0].text === 'Bob Murray')
  )
    throw Error(
      `${resource} with identifier value ${identifierValue} does not exist in FHIR`
    )

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(
    resource,
    identifierValue
  )

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(
      /Successfully deleted 1 resource/
    )
  )
    throw `Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`

  // Remove the practitioner created on the ihris server
  const resultOfDeleteResourceMockServer = await deleteResource(resource)

  if (resultOfDeleteResourceMockServer.status != 200)
    throw Error(`${resource} not removed from the mock-service`)

  console.log(
    `\n ${resource} resource existence on FHIR verified, and clean up done...`
  )
}

exports.verifyPractitionerRoleExistsAndCleanup = async () => {
  const resource = 'PractitionerRole'
  const identifierValue = testPractitionerRole.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (!(response.status === 200) || !(response.data.total === 1))
    throw Error(
      `${resource} with identifier-value ${identifierValue} exists ${response.data.total} times in FHIR`
    )

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(
    resource,
    identifierValue
  )

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(
      /Successfully deleted 1 resource/
    )
  )
    throw `Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`

  // Remove the practitionerRole created on the ihris server
  const resultOfDeleteResourceMockServer = await deleteResource(resource)

  if (resultOfDeleteResourceMockServer.status != 200)
    throw Error(`${resource} not removed from the mock-service`)

  console.log(
    `\n ${resource} resource existence on FHIR verified, and clean up done...`
  )
}

exports.verifyLocationExistsAndCleanup = async () => {
  const resource = 'Location'
  const identifierValue = testLocation.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'GoodHealth Clinic')
  )
    throw Error(
      `${resource} with identifier-value ${identifierValue} does not exist`
    )

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(
    resource,
    identifierValue
  )

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(
      /Successfully deleted 1 resource/
    )
  )
    throw `Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`

  // Remove the location created on the gofr server
  const resultOfDeleteResourceMockServer = await deleteResource(resource)

  if (resultOfDeleteResourceMockServer.status != 200)
    throw Error(`${resource} not removed from the mock-service`)

  console.log(
    `\n ${resource} resource existence on FHIR verified, and clean up done...`
  )
}

exports.verifyOrganizationExistsAndCleanup = async () => {
  const resource = 'Organization'
  const identifierValue = testOrganization.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'Clinical Lab')
  )
    throw Error(
      `${resource} with identifier-value ${identifierValue} does not exist`
    )

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(
    resource,
    identifierValue
  )

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(
      /Successfully deleted 1 resource/
    )
  )
    throw `Clean up failed, test ${resource} (id: ${identifierValue}) not removed from FHIR`

  // Remove the organization created on the gofr server
  const resultOfDeleteResourceMockServer = await deleteResource(resource)

  if (resultOfDeleteResourceMockServer.status != 200)
    throw Error(`${resource} not removed from the mock-service`)

  console.log(
    `\n ${resource} resource existence on FHIR verified, and clean up done...`
  )
}
