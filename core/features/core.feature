Feature: FHIR server available via OpenHIM

  Rule: Authorised clients should be able to access the FHIR server through the OpenHIM

    Scenario: An authorised client searches for a patient
      Given a patient, Jane Doe, exists in the FHIR server
      And an authorised client, Alice, exists in the OpenHIM
      When Alice searches for a patient
      Then Alice is able to get a result

    Scenario: An un-authorised client searches for a patient
      Given a patient, Jane Doe, exists in the FHIR server
      And an authorised client, Alice, exists in the OpenHIM
      When Malice searches for a patient
      Then Malice is NOT able to get a result
