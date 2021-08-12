Feature: Analytics Infrastructure Test
    Rule: The infrastructure should be up

        Scenario: Infrastructure up
           Then the JS Report service should be up and running
           And the ES analytics service should be up and running
           And the Kibana shpuld be up and running
