'use strict'

const axios = require('axios')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'
const CUSTOM_TOKEN_ID = process.env.CUSTOM_TOKEN_ID || 'test'

const practitionerId = 'P10004'
const practitionerRoleId = 'PR10001'
const locationId = '2'
const organizationId = '3'

exports.triggerSync = async () => {
  await axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/mcsd-trigger`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  })

  // Allow the syncing to finish. The data will only be available in FHIR after a few seconds
  await new Promise ((resolve, _reject) => {
    setTimeout(() => resolve(), 60000)
  })
}

exports.ihrisMockServicePractitioner = async () => {
  // Mock ihris service running in docker container
  await verifyResourceDoesNotExistInFHIR('Practitioner', practitionerId)
}

exports.ihrisMockServicePractitionerRole = async () => {
  // Mock ihris service running in docker container 
  await verifyResourceDoesNotExistInFHIR('PractitionerRole', practitionerRoleId)
}

exports.gofrMockServiceLocation = async () => {
  // Mock gofr service running in docker container
  await verifyResourceDoesNotExistInFHIR('Location', locationId)
}

exports.gofrMockServiceOrganization = async () => {
  // Mock gofr service running in docker container
  await verifyResourceDoesNotExistInFHIR('Organization', organizationId)
}

const getResource = (resource, id) => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?_id=${id}`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${CUSTOM_TOKEN_ID}`,
      'Cache-Control': 'no-cache'
    }
  })
}

const removeResource = (resource, id) => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?_id=${id}`,
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
    `Test aborted! ${resource} resource (id: ${resourceId}) used in test already exists and will be removed from the FHIR server`
    )
}

exports.verifyPractitionerExistsAndCleanup = async () => {
  const resource = 'Practitioner'

  const response = await getResource(resource, practitionerId)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name[0].text === 'Tekiokio Traifrop')
    ) throw Error(`${resource} with id ${practitionerId} does not exist`)

  const result = await removeResource(resource, practitionerId)

  if (
    !JSON.stringify(result.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (id: ${practitionerId}) not removed from FHIR`)
}

exports.verifyPractitionerRoleExistsAndCleanup = async () => {
  const resource = 'PractitionerRole'

  const response = await getResource(resource, practitionerRoleId)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.id === practitionerRoleId)
    ) throw Error(`${resource} with id ${practitionerRoleId} does not exist`)

  const result = await removeResource(resource, practitionerRoleId)

  if (
    !JSON.stringify(result.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (id: ${practitionerRoleId}) not removed from FHIR`)
}

exports.verifyLocationExistsAndCleanup = async () => {
  const resource = 'Location'

  const response = await getResource(resource, locationId)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'USSS Enterprise-D Sickbay')
    ) throw Error(`${resource} with id ${locationId} does not exist`)

  const result = await removeResource(resource, locationId)

  if (
    !JSON.stringify(result.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (id: ${locationId}) not removed from FHIR`)
}

exports.verifyOrganizationExistsAndCleanup = async () => {
  const resource = 'Organization'

  const response = await getResource(resource, organizationId)

  if (
    !(response.status === 200) ||
    !(response.data.total === 1) ||
    !(response.data.entry[0].resource.name === 'Clinical Lab')
    ) throw Error(`${resource} with id ${organizationId} does not exist`)

  const result = await removeResource(resource, organizationId)

  if (
    !JSON.stringify(result.data).match(/Successfully deleted 1 resource/)
    ) throw(`Clean up failed, test ${resource} (id: ${organizationId}) not removed from FHIR`)
}
