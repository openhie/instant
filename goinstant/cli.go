package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func CLI() {
	startupCommands := os.Args[1:]

	switch startupCommands[0] {
	case "docker":
		if len(startupCommands) < 3 {
			panic("Incorrect arguments list passed to CLI. Requires at least 3 arguments when in non-interactive mode.")
		}

		RunDirectDockerCommand(startupCommands)
	case "k8s", "kubernetes":
		color.Red("\nKubernetes not supported for now :(")
	case "help":
		fmt.Println(`
Commands: 
	help 		this menu
	docker		manage package in docker
				usage: docker <package> <state>

				docker core init
				docker core up
				docker core destroy

				note: only one package can be instantiated at a time using the CLI

	kubernetes	manage package in kubernetes, can also use k8s
				usage: k8s/kubernetes <package> <state>

				k8s core init
				kubernetes core up
				kubernetes core destroy

	install		install fhir npm package on fhir server
				usage: install <ig_url> <fhir_server> <authtype> <user/token> <pass>

				install https://intrahealth.github.io/simple-hiv-ig/ http://hapi.fhir.org/baseR4 none
				install <ig_url> <fhir_server> basic smith stuff
				install <ig_url> <fhir_server> token "123"
				install <ig_url> <fhir_server> custom test
	`)

	case "install":
		params := &Params{}
		switch startupCommands[3] {
		case "none", "None":
			params.TypeAuth = "None"
			loadIGpackage(startupCommands[1], startupCommands[2], params)
		case "basic", "Basic":
			params.TypeAuth = "Basic"
			params.BasicUser = startupCommands[4]
			params.BasicPass = startupCommands[5]
			loadIGpackage(startupCommands[1], startupCommands[2], params)
		case "token", "Token":
			params.TypeAuth = "Token"
			params.Token = startupCommands[4]
			loadIGpackage(startupCommands[1], startupCommands[2], params)
		case "custom", "Custom":
			params.TypeAuth = "Custom"
			params.Token = startupCommands[4]
			loadIGpackage(startupCommands[1], startupCommands[2], params)
		}
	default:
		fmt.Println("The deploy command is not recognized: ", startupCommands)
	}

}
