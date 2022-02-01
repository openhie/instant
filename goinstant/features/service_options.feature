Feature: Test Deploy Options Commands
  Scenario: Initialise Core Service in Dev Mode
    When the command "init core -t=docker --dev" is run
    Then check the CLI output is "init --dev -t docker core"

  Scenario: Down only Core Service
    When the command "down core -t=docker --only" is run
    Then check the CLI output is "down --only -t docker core"

  Scenario: Initialise Template Custom Service
    When the command "init template -t=docker -c=https://github.com/jembi/instant-openhie-template-package.git" is run
    Then check that the CLI added custom packages
      | directory | location |
      | instant-openhie-template-package | https://github.com/jembi/instant-openhie-template-package.git |
    Then check the CLI output is "init -t docker template"

  Scenario: Initialise Multiple Services
    When the command "init covid19immunization client covid19surveillance -t=docker --custom-package=https://github.com/jembi/covid19-immunization-tracking-package.git -c=https://github.com/jembi/who-covid19-surveillance-package.git" is run
    Then check that the CLI added custom packages
      | directory | location |
      | covid19-immunization-tracking-package | https://github.com/jembi/covid19-immunization-tracking-package.git |
      | who-covid19-surveillance-package | https://github.com/jembi/who-covid19-surveillance-package.git |
    Then check the CLI output is "init -t docker covid19immunization client covid19surveillance"

  Scenario: Initialise Core Service with .env file
    When the command "init core -t=docker --env-file=features/.env.test" is run
    Then check the CLI output is "--env-file features/.env.test"

  Scenario: Initialise Core Service with Environment Variable
    When the command "init core -t=docker -e=TEST=ME" is run
    Then check the CLI output is "-e TEST=ME"
