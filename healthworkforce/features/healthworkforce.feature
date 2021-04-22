Feature: mCSD support

  Rule: mCSD directories are sync'd to the FHIR server

    Scenario: New changes are sync'd to the FHIR server
      Given a new practitioner, Dr Bob, has been added in iHRIS
      And a new practitionerRole has been added in iHRIS
      And a new location, GoodHealth Clinic, has been added in GOFR
      And a new organization, Clinical Lab, has been added in GOFR
      When the sync is triggered
      Then the new practitionerRole can be found in the FHIR server
      And the new practitioner, Dr Bob, can be found in the FHIR server
      And the new location, GoodHealth Clinic, can be found in the FHIR server
      And the new organization, Clinical Lab, can be found in the FHIR server
