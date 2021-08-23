Feature: Pipeline Infrastructure Test
  Rule: The infrastructure should be up

    Scenario: Infrastructure up
      Then the fhir extractor service should be up and running
      And the logstash service should be up and running
