package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/pkg/errors"
)

func debugKubernetes() error {
	fmt.Println("...checking your cluster setup")

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("kubectl", "config", "view")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "cmd.Run() failed. Ensure Docker for Desktop for Windows or Mac or minikube is installed")
		}
		fmt.Println("\nkubectl cluster manager is installed, you're ready to use kubernetes")
		fmt.Println("\nThe current kubernetes cluster active will be used. Change now if you wish to use another.")

	case "windows":
		cmd := exec.Command("cmd", "/C", "kubectl", "config", "view")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "cmd.Run() failed. Ensure Docker for Desktop for Windows or Mac or minikube is installed")
		}
		fmt.Println("kubectl cluster manager is installed, you're ready to use kubernetes")
		fmt.Println("The current kubernetes cluster active will be used. Change now if you wish to use another.")

	default:
		fmt.Println("What operating system is this?")
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

	// TODO: choose kubernetes cluster to use
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	kubeconfig := path.Join(home, ".kube/config")
	fmt.Println("...using cluster config file at:")
	fmt.Println(kubeconfig)
}

// func configServerKubernetes() {

// 	path, _ := os.Getwd()
// 	// if err != nil {
// 	// }

// 	files, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

	return nil
}
