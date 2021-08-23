'use strict'

const axios = require('axios')
const { Then } = require('@cucumber/cucumber')

const PIPELINE_PROTOCOL = process.env.PIPELINE_PROTOCOL || 'http'
const FHIR_EXTRACTOR_HOSTNAME = process.env.FHIR_EXTRACTOR_HOSTNAME  || 'fhir-extractor'
const LOGSTASH_HOSTNAME = process.env.LOGSTASH_HOSTNAME  || 'logstash'
const FHIR_EXTRACTOR_API_PORT =
  process.env.FHIR_EXTRACTOR_API_PORT || '3000'
const LOGSTASH_API_PORT =
  process.env.LOGSTASH_API_PORT || '5055'

const sendRequest = (hostname, port) => axios({
  url: `${PIPELINE_PROTOCOL}://${hostname}:${port}`
})

Then('the logstash service should be up and running', async () => {
  console.log('\n\nChecking the status of the Logstash Service\n\n')

  const result = await sendRequest(LOGSTASH_HOSTNAME, LOGSTASH_API_PORT)

  console.log(`The Logstash service is up and responded with status - ${result.status}\n`)
})

/*
 The following test only checks whether the service is up and running.
 A request is send to the service and a positive result is one in which the
 http response status is not 500.
*/
Then('the fhir extractor service should be up and running', async () => {
  console.log('\n\nChecking the status of the Fhir Extractor Service\n\n')

  try {
    const result = await sendRequest(FHIR_EXTRACTOR_HOSTNAME ,FHIR_EXTRACTOR_API_PORT)

    console.log(`The Fhir Extractor service is up and responded with a status - ${result.status}\n`)
  } catch (error) {
    if (error.response.status != 500) {
      console.log(`The Fhir Extractor service is up and responded with a status - ${error.response.status}\n`)
    } else {
      throw error
    }
  }
})
