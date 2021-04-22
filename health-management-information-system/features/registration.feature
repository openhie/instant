Feature: Tracked Entity Creation

  Scenario: A patient is registered
    Given that dhis is set up and the metadata import has been done
    When a patient is created
    Then the patient should exist in DHIS
    And the patient should then be deleted
