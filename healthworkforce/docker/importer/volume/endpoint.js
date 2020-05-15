'use strict'

const fs = require('fs')
const http = require('http')
const path = require('path')

const MEDIATOR_HOSTNAME = process.env.MEDIATOR_HOST_NAME || 'mcsdMediator'
const MEDIATOR_API_PORT = process.env.MEDIATOR_API_PORT || 3003


const jsonData = JSON.parse(
  fs.readFileSync(path.resolve(__dirname, 'endpoint.json'))
)

const data = JSON.stringify(jsonData)

const options = {
  protocol: 'http:',
  host: MEDIATOR_HOSTNAME,
  port: MEDIATOR_API_PORT,
  path: '/endpoints',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': data.length,
  }
}

const req = http.request(options, res => {
  if (res.statusCode === 400) {
    let data = ''
    res.on('data', chunk => {
      data += chunk.toString()
    })

    res.on('end', () => {
      if (data) {
        data = JSON.parse(data)
        if (
          data.error && data.error.match(/duplicate key error/).length
          ) {
          console.log('mCSD endpoint already exists')
          return
        }
        throw Error('mCSD endpoint creation failed')
      }
    }
    )

    res.on('error', err => {
      throw Error(err)
    })
    return
  }

  if (res.statusCode != 201) {
    throw Error(`Failed to create mCSD mediator endpoint: ${res.statusCode}`)
  } else {
    console.log('Successfully created mCSD endpoint')
  }
})

req.on('error', error => {
  console.error('Failed to create mCSD mediator endpoint: ', error)
})

req.write(data)
req.end()
