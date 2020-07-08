Feature: Basic Health Worker Registry Queries

  Scenario: Query for HW records
     Given fhir server returns a capability statement for practitioner
      When practitioner resources exist on the fhir server
      Then get one practitioner resource