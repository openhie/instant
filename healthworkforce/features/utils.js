'use strict'

const axios = require('axios')
const { expect } = require('chai')
const nock = require('nock')

const OPENHIM_PROTOCOL = process.env.OPENHIM_PROTOCOL || 'http'
const OPENHIM_API_HOSTNAME = process.env.OPENHIM_API_HOSTNAME || 'localhost'
const OPENHIM_TRANSACTION_API_PORT =
  process.env.OPENHIM_TRANSACTION_API_PORT || '5001'

const OPENHIM_API_PASSWORD =
  process.env.OPENHIM_API_PASSWORD || 'instant101'
const OPENHIM_API_USERNAME =
  process.env.OPENHIM_API_USERNAME || 'root@openhim.org'

const authHeader = new Buffer.from(
  `${OPENHIM_API_USERNAME}:${OPENHIM_API_PASSWORD}`
).toString('base64')

const practitionerName = "Bob"
const practitionerRoleId = "PR123566"
const locationName = "GoodHealth Clinic"
const organizationName = "Clinical Lab"

exports.triggerSync = async () => {
  await axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/mcsd`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Basic ${authHeader}`
    }
  })
  .then(() => {
    console.log('Data successfully sent to dhis')
  })
  .catch(() => {
    console.log('Error in sending data to dhis')
  })
}

exports.ihrisMockServicePractitioner = async () => {
  await verifyResourceDoesNotExist('Practitioner')

  nock('http://mock-service:4000')
    .get('/ihris-practitioner-mock/_history')
    .reply(
      200,
      {
        "resourceType": "Bundle",
        "type": "transaction",
        "entry": [
          {
            "resource": {
              "resourceType": "Practitioner",
              "meta": {
                "profile": [
                  "http://ihris.org/fhir/StructureDefinition/iHRISPractitioner"
                ]
              },
              "id": "P123456",
              "active": true,
              "name": [
                {
                  "use": "official",
                  "text": "Tekiokio Traifrop",
                  "given": [
                    practitionerName
                  ],
                  "family": "Traifrop"
                }
              ],
              "identifier": [
                {
                  "system": "http://www.acme.org.au/units",
                  "value": "testPractitioner"
                }
              ],
              "gender": "female",
              "birthDate": "1979-01-01"
            },
            "request": {
              "method": "PUT",
              "url": "Practitioner/P123456"
            }
          }
        ]
      })
}

exports.ihrisMockServicePractitionerRole = async () => {
  await verifyResourceDoesNotExist('PractitionerRole')

  nock('http://mock-service:4000')
    .get('/ihris-practitionerRole-mock/_history')
    .reply(
      200,
      {
        "resourceType": "Bundle",
        "type": "transaction",
        "entry": [
          {
            "resource": {
              "resourceType": "PractitionerRole",
              "meta": {
                "profile": [
                  "http://ihris.org/fhir/StructureDefinition/iHRISPractitionerRole"
                ]
              },
              "id": practitionerRoleId,
              "active": true,
              "practitioner": {
                "reference": "Practitioner/P123456"
              },
              "identifier": [
                {
                  "system": "http://www.acme.org.au/units",
                  "value": "testPractitionerRole"
                }
              ],
              "location": [
                {
                  "reference": "Location/TF.S.NYS.11"
                }
              ],
              "period": {
                "start": "1998-01-01"
              }
            },
            "request": {
              "method": "PUT",
              "url": "PractitionerRole/PR123456"
            }
          }
        ]
      })
}

exports.gofrMockServiceLocation = async () => {
  await verifyResourceDoesNotExist('Location')

  nock('http://mock-service:4000')
    .get('/gofr-location-mock/_history')
    .reply(
      200,
      {
        "resourceType": "Bundle",
        "type": "searchset",
        "entry": [
          {
            "fullUrl": "http://hl7.org/fhir/Location/123456",
            "resource": {
              "resourceType": "Location",
              "id": "123456",
              "text": {
                "status": "generated",
                "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">USS Enterprise</div>"
              },
              "identifier": [
                {
                  "system": "http://www.acme.org.au/units",
                  "value": "testLocation"
                }
              ],
              "status": "active",
              "name": locationName,
              "mode": "instance"
            }
          }
        ]
      })
}

exports.gofrMockServiceOrganization = async () => {
  await getResource('Organization')

  nock('http://mock-service:4000')
    .get('/gofr-organization-mock/_history')
    .reply(
      200,
      {
        "resourceType": "Bundle",
        "type": "searchset",
        "entry": [
          {
            "resourceType": "Organization",
            "id": "123456",
            "text": {
              "status": "generated",
              "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">\n      \n      <p>Clinical Laboratory @ Acme Hospital. ph: +1 555 234 1234, email: \n        <a href=\"mailto:contact@labs.acme.org\">contact@labs.acme.org</a>\n      </p>\n    \n    </div>"
            },
            "identifier": [
              {
                "system": "http://www.acme.org.au/units",
                "value": "testOrganization"
              }
            ],
            "name": organizationName
          }
        ]
      })
}

const getResource = resource => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?identifier:value=test${resource}`,
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${authHeader}`
    }
  })
}

const removeResource = resource => {
  return axios({
    url: `${OPENHIM_PROTOCOL}://${OPENHIM_API_HOSTNAME}:${OPENHIM_TRANSACTION_API_PORT}/hapi-fhir-jpaserver/fhir/${resource}?identifier:value=test${resource}`,
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Custom ${authHeader}`
    }
  })
}

const verifyResourceDoesNotExist = async resource => {
  const response = await getResource(resource)

  expect(response.status).to.eql(200)

  if (response.data.total > 0) throw Error(
    `Test aborted! ${resource} resource used in test already exists and will be removed from the FHIR server`
    )
}

exports.verifyPractitionerExistsAndCleanup = async () => {
  const resource = 'Practitioner'
  const response = await getResource(resource)

  expect(response.status).to.eql(200)
  expect(response.body.name[0].given, practitionerName)

  await removeResource(resource)
  .then(() => {console.log(`Test ${resource} successfully removed from FHIR`)})
  .catch(() => {console.error(`Clean up failed, test ${resource} not removed from FHIR`)})
}

exports.verifyPractitionerRoleExistsAndCleanup = async () => {
  const resource = 'PractitionerRole'
  const response = await getResource(resource)

  expect(response.status).to.eql(200)
  expect(response.body.id, practitionerRoleId)

  await removeResource(resource)
    .then(() => {console.log(`Test ${resource} successfully removed from FHIR`)})
    .catch(() => {console.error(`Clean up failed, test ${resource} not removed from FHIR`)})
}

exports.verifyLocationExistsAndCleanup = async () => {
  const resource = 'Location'
  const response = await getResource(resource)

  expect(response.status).to.eql(200)
  expect(response.body.name, locationName)

  await removeResource(resource)
    .then(() => {console.log(`Test ${resource} successfully removed from FHIR`)})
    .catch(() => {console.error(`Clean up failed, test ${resource} not removed from FHIR`)})
}

exports.verifyOrganizationExistsAndCleanup = async () => {
  const resource = 'Organization'
  const response = await getResource(resource)

  expect(response.status).to.eql(200)
  expect(response.body.name, organizationName)

  await removeResource(resource)
  .then(() => {console.log(`Test ${resource} successfully removed from FHIR`)})
  .catch(() => {console.error(`Clean up failed, test ${resource} not removed from FHIR`)})
}
