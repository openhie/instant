'use strict'

const axios = require('axios')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'
const CUSTOM_TOKEN_ID = process.env.CUSTOM_TOKEN_ID || 'test'
const MOCK_SERVER_PORT = process.env.MOCK_SERVER_PORT || '4000'

const testLocationBundle = {
  resourceType: 'Bundle',
  type: 'searchset',
  entry: [
    {
      resource: {
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
    }
  ]
}

const testOrganizationBundle = {
  resourceType: 'Bundle',
  type: 'searchset',
  entry: [
    {
      resource: {
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
    }
  ]
}

const testPractitionerBundle = {
  resourceType: 'Bundle',
  type: 'transaction',
  entry: [
    {
      resource: {
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
      },
      request: {
        method: 'PUT',
        url: 'Practitioner/P10004test'
      }
    }
  ]
}

const testPractitionerRoleBundle = {
  resourceType: 'Bundle',
  type: 'transaction',
  entry: [
    {
      resource: {
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
        ],
      },
      request: {
        method: 'PUT',
        url: 'PractitionerRole/PR10001test'
      }
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

  console.log('\n Triggered sync')
}

exports.ihrisMockServicePractitioner = async () => {
  const name = 'Practitioner'

  await verifyResourceDoesNotExistInFHIR(name, testPractitionerBundle.entry[0].resource.identifier[0].value)

  // create new practitioner Dr Bob on ihris mock server
  const response = await createResourceBundle(name, testPractitionerBundle)

  if (response.status != 201) throw Error(`Failed to create ${name} for testing`)

  console.log('\n Created Practitioner resource (Dr Bob) on mock ihris mock server')
}

exports.ihrisMockServicePractitionerRole = async () => {
  const name = 'PractitionerRole'

  await verifyResourceDoesNotExistInFHIR('PractitionerRole', testPractitionerRoleBundle.entry[0].resource.identifier[0].value)

  // create new practitionerRole for Dr Bob on ihris mock server
  const response = await createResourceBundle(name, testPractitionerRoleBundle)

  if (response.status != 201) throw Error(`Failed to create ${name} for testing`)

  console.log('\n Created PractitionerRole resource (for Dr Bob) on mock ihris mock server')
}

exports.gofrMockServiceLocation = async () => {
  const name = 'Location'

  await verifyResourceDoesNotExistInFHIR('Location', testLocationBundle.entry[0].resource.identifier[0].value)

  // create new location on gofr mock server
  const response = await createResourceBundle(name, testLocationBundle)

  if (response.status != 201) throw Error(`Failed to create ${name} for testing`)

  console.log('\n Created Location resource (GoodHealth Clinic) on mock gofr mock server')
}

exports.gofrMockServiceOrganization = async () => {
  const name = 'Organization'

  await verifyResourceDoesNotExistInFHIR('Organization', testOrganizationBundle.entry[0].resource.identifier[0].value)

  // create new organization on gofr mock server
  const response = await createResourceBundle(name, testOrganizationBundle)

  if (response.status != 201) throw Error(`Failed to create ${name} for testing`)

  console.log('\n Created Organization resource (Clinical Lab) on mock gofr mock server...')
}

const getResource = (resource, identifierValue) => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?identifier:value=${identifierValue}`,
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
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?identifier:value=${identifierValue}`,
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

  if (response.data.total > 0) throw Error(
    `Test aborted! ${resource} resource (identifier-value: ${resourceId}) used in test already exists and will be removed from the FHIR server`
    )
}

const createResourceBundle = (resourceName, resourceBundle) => {
  return axios({
    url: `http://localhost:${MOCK_SERVER_PORT}/create-resource/${resourceName}`,
    method: 'POST',
    data: resourceBundle
  })
}

const deleteResourceBundle = resourceName => {
  return axios({
    url: `http://localhost:${MOCK_SERVER_PORT}/cleanup/${resourceName}`,
    method: 'DELETE'
  })
}

exports.verifyPractitionerExistsAndCleanup = async () => {
  const resource = 'Practitioner'
  const identifierValue = testPractitionerBundle.entry[0].resource.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name[0].text === 'Bob Murray')
    ) throw Error(`${resource} with identifier value ${identifierValue} does not exist in FHIR`)

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(resource, identifierValue)

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`)

  // Remove the practitioner bundle created on the ihris server
  const resultOfDeleteResourceMockServer = await deleteResourceBundle(resource)

  if (resultOfDeleteResourceMockServer.status != 200) throw Error(`${resource} not removed from the mock-service`)

  console.log(`\n ${resource} resource existence on FHIR verified, and clean up done...`)
}

exports.verifyPractitionerRoleExistsAndCleanup = async () => {
  const resource = 'PractitionerRole'
  const identifierValue = testPractitionerRoleBundle.entry[0].resource.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1)
    ) throw Error(`${resource} with identifier-value ${identifierValue} does not exist in FHIR`)

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(resource, identifierValue)

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`)

  // Remove the practitionerRole bundle created on the ihris server
  const resultOfDeleteResourceMockServer = await deleteResourceBundle(resource)

  if (resultOfDeleteResourceMockServer.status != 200) throw Error(`${resource} not removed from the mock-service`)

  console.log(`\n ${resource} resource existence on FHIR verified, and clean up done...`)
}

exports.verifyLocationExistsAndCleanup = async () => {
  const resource = 'Location'
  const identifierValue = testLocationBundle.entry[0].resource.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'GoodHealth Clinic')
    ) throw Error(`${resource} with identifier-value ${identifierValue} does not exist`)

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(resource, identifierValue)

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (identifier-value: ${identifierValue}) not removed from FHIR`)

  // Remove the location bundle created on the gofr server
  const resultOfDeleteResourceMockServer = await deleteResourceBundle(resource)

  if (resultOfDeleteResourceMockServer.status != 200) throw Error(`${resource} not removed from the mock-service`)

  console.log(`\n ${resource} resource existence on FHIR verified, and clean up done...`)
}

exports.verifyOrganizationExistsAndCleanup = async () => {
  const resource = 'Organization'
  const identifierValue = testOrganizationBundle.entry[0].resource.identifier[0].value

  const response = await getResource(resource, identifierValue)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'Clinical Lab')
    ) throw Error(`${resource} with identifier-value ${identifierValue} does not exist`)

  // Remove resource from FHIR
  const resultOfDeleteResourceFHIR = await removeResource(resource, identifierValue)

  if (
    !JSON.stringify(resultOfDeleteResourceFHIR.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (id: ${identifierValue}) not removed from FHIR`)

  // Remove the organization bundle created on the gofr server
  const resultOfDeleteResourceMockServer = await deleteResourceBundle(resource)

  if (resultOfDeleteResourceMockServer.status != 200) throw Error(`${resource} not removed from the mock-service`)

  console.log(`\n ${resource} resource existence on FHIR verified, and clean up done...`)
}
