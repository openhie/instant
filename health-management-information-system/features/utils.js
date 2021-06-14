'use strict'

const axios = require('axios')

const DHIS2_PROTOCOL = process.env.DHIS2_PROTOCOL || 'http'
const DHIS2_API_HOSTNAME = process.env.DHIS2_API_HOSTNAME || 'localhost'
const DHIS2_API_PASSWORD = process.env.DHIS2_API_PASSWORD || 'district'
const DHIS2_API_PORT = process.env.DHIS2_API_PORT || 8081
const DHIS2_API_USERNAME = process.env.DHIS2_API_USERNAME || 'admin'

const authHeader = new Buffer.from(
  `${DHIS2_API_USERNAME}:${DHIS2_API_PASSWORD}`
).toString('base64')

exports.verifyDhis2Running = async () => {
  const options = {
    url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/dhis-web-commons/security/login.action`,
    method: 'GET'
  }

  try {
    const response = await axios(options)

    if (response && response.status === 200) {
      console.log(`DHIS2 running`)
    } else {
      throw new Error('DHIS2 NOT running!')
    }
  } catch (error) {
    throw new Error(`DHIS2 issues: ${error.message}`)
  }
}

exports.getSystemInfo = async () => {
 const response = await axios({
      url: `${DHIS2_PROTOCOL}://${DHIS2_API_HOSTNAME}:${DHIS2_API_PORT}/api/system/info`,
      method: 'GET',
      headers: {
        Authorization: `Basic ${authHeader}`,
        'Content-Type': 'application/json'
      }
  })

  if (response.status !== 200 || !response.data.databaseInfo || response.data.databaseInfo.url != 'jdbc:postgresql://dhis-postgres/dhis2') {
    throw Error(`Get System Info Failed: ${response.data}`)
  }
}
