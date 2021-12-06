Feature: Analytics Infrastructure Test
  Rule: The infrastructure should be up

    Scenario: Infrastructure up
      Then the ES analytics service should be up and running
      And the Kibana shpuld be up and running
