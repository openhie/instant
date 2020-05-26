'use strict'

exports.parseStringToBoolean = (value, defaultValue) => {
  if (!value) return defaultValue

  switch (value.toLowerCase()) {
    case 'true':
      return true
    case 'false':
      return false
    default:
      return defaultValue
  }
}
