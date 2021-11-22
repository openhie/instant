package main

import (
	"fmt"
	"os"
)

func CLI() {
	argsWOProg := os.Args[1:]

	switch argsWOProg[0] {
	case "docker", "k8s", "kubernetes":
		if len(argsWOProg) < 3 {
			panic("Incorrect arguments list passed to CLI. Requires at least 3 arguments when in non-interactive mode.")
		}

		RunDirectDockerCommand(argsWOProg[0], argsWOProg[1], argsWOProg[2], argsWOProg...)

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
		switch argsWOProg[3] {
		case "none", "None":
			params.TypeAuth = "None"
			loadIGpackage(argsWOProg[1], argsWOProg[2], params)
		case "basic", "Basic":
			params.TypeAuth = "Basic"
			params.BasicUser = argsWOProg[4]
			params.BasicPass = argsWOProg[5]
			loadIGpackage(argsWOProg[1], argsWOProg[2], params)
		case "token", "Token":
			params.TypeAuth = "Token"
			params.Token = argsWOProg[4]
			loadIGpackage(argsWOProg[1], argsWOProg[2], params)
		case "custom", "Custom":
			params.TypeAuth = "Custom"
			params.Token = argsWOProg[4]
			loadIGpackage(argsWOProg[1], argsWOProg[2], params)
		}
	case "default":
		fmt.Println("The command is not recognized.")
	}

}
