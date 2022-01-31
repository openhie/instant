Feature: Test Deploy Options Commands
  Scenario: Initialise Core Service in Dev Mode
    When the command "init core -t=docker --dev" is run
    Then check the CLI output is "init --dev -t docker core"

  Scenario: Initialise Custom Service in Dev Mode
    When the command "init core -t=docker --dev" is run
    Then check the CLI output is "init --dev -t docker core"
