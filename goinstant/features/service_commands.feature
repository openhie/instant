Feature: Test Deploy Commands
  Scenario: Initialise Core Service
    When the command "init core -t=docker" is run
    Then check the CLI output is "init -t docker core"

  Scenario: Up Core Service
    When the command "core up -t=docker" is run
    Then check the CLI output is "up -t docker core"

  Scenario: Down Core Service
    When the command "core -t=docker down" is run
    Then check the CLI output is "down -t docker core"

  Scenario: Destroy Core Service
    When the command "core -t=docker destroy" is run
    Then check the CLI output is "destroy -t docker core"
