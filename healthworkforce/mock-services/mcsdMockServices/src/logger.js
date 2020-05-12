const pino = require('pino')

const config = require('./config')
const configOptions = config.getConfig()

const logger = pino({
  level: configOptions.logLevel,
  prettyPrint: {
    colorize: true,
    translateTime: 'sys:UTC:yyyy-mm-dd"T"HH:MM:ss:l"Z"',
    ignore: 'pid,hostname'
  },
  serializers: {
    err: pino.stdSerializers.err
  },
  enabled: configOptions.enableLogging
})

module.exports = exports = logger
