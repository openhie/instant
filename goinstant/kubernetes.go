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

	// TODO: choose kubernetes cluster to use
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	kubeconfig := path.Join(home, ".kube/config")
	fmt.Println("...using cluster config file at:")
	fmt.Println(kubeconfig)

	return nil
}
