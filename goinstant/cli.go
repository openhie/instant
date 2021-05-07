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
	utils		usage: utils ig load <url> <fhirserver>, ig examples <url> <fhirserver> 
	`)

	case "default":
		fmt.Println("default")
	}

}
