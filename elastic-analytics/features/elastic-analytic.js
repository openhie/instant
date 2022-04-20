'use strict'

const axios = require('axios')
const { Then } = require('@cucumber/cucumber')

const ANALYTICS_PROTOCOL = process.env.ANALYTICS_PROTOCOL || 'http'
const ES_ANALYTICS_HOSTNAME = process.env.ES_ANALYTICS_HOSTNAME || 'es-analytics'
const KIBANA_HOSTNAME = process.env.KIBANA_HOSTNAME || 'kibana'
const KIBANA_API_PORT =
  process.env.KIBANA_API_PORT || '5601'
const ES_ANALYTICS_API_PORT =
  process.env.ES_API_PORT || '9200'

const sendRequest = (port, hostname) => axios({
  url: `${ANALYTICS_PROTOCOL}://${hostname}:${port}`
})

/*
 The following test only checks whether the service is up and running.
 A request is send to the service and a positive result is one in which the
 http response status is not 500.
*/
Then('the ES analytics service should be up and running', async () => {
  console.log('\n\nChecking the status of the ES Analytics Service\n\n')

  try {
    const result = await sendRequest(ES_ANALYTICS_API_PORT, ES_ANALYTICS_HOSTNAME)

    console.log(`The ES Analytics service is up and responded with status ${result.status}\n`)
  } catch (error) {
    if (error.response.status != 500) {
      console.log(`The ES Analytics service is up and responded with a status - ${error.response.status}\n`)
    } else {
      throw error
    }
  }
})

Then('the Kibana should be up and running', async () => {
  console.log('\n\nChecking the status of the Kibana Service\n\n')

  const result = await sendRequest(KIBANA_API_PORT, KIBANA_HOSTNAME)

  console.log(`The Kibana service is up and responded with status - ${result.status}\n`)
})
