package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gookit/color"
)

func debugKubernetes() {

	fmt.Println("...checking your cluster setup")

	color.Yellow.Println("Running on", runtime.GOOS)

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("kubectl", "config", "view")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
			color.Red.Println("The kubectl cluster manager is not installed.\nUsing a local cluster requires that Docker for Desktop for Windows or Mac or minikube be installed. Please install one of those.")

		} else {
			color.Green.Println("kubectl cluster manager is installed")
		}

	case "windows":
		cmd := exec.Command("cmd", "/C", "kubectl", "config", "view")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
			color.Red.Println("The kubectl cluster manager is not installed.\nUsing a local cluster requires that Docker for Desktop for Windows or Mac or minikube be installed. Please install one of those.")
		} else {
			color.Green.Println("kubectl cluster manager is installed")
		}
	default:
		fmt.Println("What operating system is this?", runtime.GOOS)
	}

	// switch runtime.GOOS {
	// case "linux", "darwin":
	// 	cmd := exec.Command("kubectl", "cluster-info")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		log.Fatalf("cmd.Run() failed with %s\n", err)
	// 		color.Red.Println("error")

	// 	} else {
	// 		color.Green.Println("not error")
	// 	}

	// case "windows":
	// 	cmd := exec.Command("cmd", "/C", "kubectl", "cluster-info")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// 	if err := cmd.Run(); err != nil {
	// 		fmt.Println("Error: ", err)
	// 		color.Red.Println("error")
	// 	} else {
	// 		color.Green.Println("not error")
	// 	}
	// default:
	// 	fmt.Println("What operating system is this?", runtime.GOOS)
	// }

	home, _ := os.UserHomeDir()
	kubeconfig := path.Join(home, ".kube/config")
	fmt.Println("...using cluster config file at:", kubeconfig)

}

func configServerKubernetes() {

	path, _ := os.Getwd()
	// if err != nil {
	// }

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)

	for _, file := range files {
		fmt.Println(file.Name())
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".go" {
				fmt.Println("found some:", file.Name())
			}
		}
	}

}
