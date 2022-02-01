Feature: Test Deploy Options Commands
  Scenario: Initialise Core Service in Dev Mode
    When the command "init core -t=docker --dev" is run
    Then check the CLI output is "init --dev -t docker core"

  Scenario: Initialise Custom Service in Dev Mode
    When the command "init template -t=docker -c=https://github.com/jembi/instant-openhie-template-package.git" is run
    Then check the CLI output is "3rd party packages to instant volume:- https://github.com/jembi/instant-openhie-template-package.git"
    Then check the CLI output is "init -t docker template"
