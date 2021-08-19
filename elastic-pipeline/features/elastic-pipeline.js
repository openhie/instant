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

  try {
    await sendRequest(LOGSTASH_HOSTNAME, LOGSTASH_API_PORT)

    console.log('The Logstash service is up\n')
  } catch (error) {
    if (error.response.status != 500) {
      console.log('The Logstash service is up\n') 
    } else {
      throw error
    }
  }
})

Then('the fhir extractor service should be up and running', async () => {
  console.log('\n\nChecking the status of the Fhir Extractor Service\n\n')

  try {
    await sendRequest(FHIR_EXTRACTOR_HOSTNAME ,FHIR_EXTRACTOR_API_PORT)

    console.log('The Fhir Extractor service is up\n')
  } catch (error) {
    if (error.response.status != 500) {
      console.log('The Fhir Extractor service is up\n') 
    } else {
      throw error
    }
  }
})
