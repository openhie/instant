package main

import (
	"fmt"
	"os"
)

func CLI() {
	argsWOProg := os.Args[1:]
	// fmt.Println(argsWOProg[0])

	switch argsWOProg[0] {
	case "docker", "k8s", "kubernetes":
		SomeStuffDirect(argsWOProg[0], argsWOProg[1], argsWOProg[2])

	case "help":
		fmt.Println(`
Commands: 
	help 		this menu
	docker		usage: docker <package> <state> e.g. docker core init
	kubernetes	usage: k8s/kubernetes <package> <state>, e.g. k8s core init
	install		usage: install <ig_url> <fhir_server>, e.g. install https://intrahealth.github.io/simple-hiv-ig/ http://hapi.fhir.org/baseR4
	`)

	case "install":
		// loadIGexamples(ig_url, fhir_server)
		loadIGpackage(argsWOProg[1], argsWOProg[2])

	case "default":
		fmt.Println("default")
	}

}
