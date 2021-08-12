'use strict'

const axios = require('axios')

const { Given, Then } = require('@cucumber/cucumber')
const { expect } = require('chai')

const PROTOCOL = process.env.PROTOCOL || 'http'
const HOSTNAME = process.env.HOSTNAME || 'localhost'
const KIBANA_API_PORT =
  process.env.KIBANA_API_PORT || '5601'
const ES_ANALYTICS_API_PORT =
  process.env.ES_API_PORT || '9201'
const JS_REPORT_API_PORT =
  process.env.JS_REPORT_API_PORT || '5488'

const sendRequest = port => axios({
    url: `${PROTOCOL}://${HOSTNAME}:${port}`
})

Then('the JS Report service should be up and running', async () => {
    console.log('\n\nChecking the status of the JS Report service\n\n')

    try {
        await sendRequest(JS_REPORT_API_PORT)

        console.log('JS Report service is up\n')
    } catch (error) {
        if (error.response.status != 500) {
            console.log('JS Report service is up\n') 
        } else {
            throw error
        }
    }
    
})

Then('the ES analytics service should be up and running', async () => {
    console.log('\n\nChecking the status of the ES Analytics Service\n\n')

    try {
        const result = await sendRequest(ES_ANALYTICS_API_PORT)

        console.log('The ES Analytics service is up\n')
    } catch (error) {
        if (error.response.status != 500) {
            console.log('The ES Analytics service is up\n') 
        } else {
            throw error
        }
    }
    
})

Then('the Kibana shpuld be up and running', async () => {
    console.log('\n\nChecking the status of the Kibana Service\n\n')

    try {
        const result = await sendRequest(KIBANA_API_PORT)

        console.log('The Kibana service is up\n')
    } catch (error) {
        if (error.response.status != 500) {
            console.log('The Kibana service is up\n') 
        } else {
            throw error
        }
    }
})
