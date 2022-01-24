Feature: destroy service
  Scenario: destroy OpenHIM service                   
    Given the OpenHIM service is initialised and running
    When the OpenHIM service is destroyed
    Then the service should not be reachable