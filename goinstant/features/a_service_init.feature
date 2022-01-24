Feature: initialise service
  Scenario: initialise and run OpenHIM service
    Given the OpenHIM service is not instantiated
    When the OpenHIM service is initialised
    Then the service should be reachable