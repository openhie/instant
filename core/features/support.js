const { setWorldConstructor } = require('@cucumber/cucumber')

class CoreWorld {
  constructor() {
    this.searchResults = null
  }

  setTo(result) {
    this.searchResults = result
  }
}

setWorldConstructor(CoreWorld)
