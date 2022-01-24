Feature: bring service down
  Scenario: bring OpenHIM service down
    Given the OpenHIM service is initialised and running
    When the OpenHIM service is brought down
    Then the service should not be reachable