Feature: bring service up
  Scenario: bring OpenHIM service up
    Given the service should not be reachable
    When the OpenHIM service is brought up
    Then the service should be reachable