Feature: DHIS2 Responsive

  Scenario: DHIS2 up and running
    Given that dhis is set up
    Then DHIS2 should respond to an authenticated API request
