package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func debugDocker() {

	fmt.Printf("...checking your Docker setup")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Can't get current working directory... this is not a great error.")
		// panic(err)
	} else {
		fmt.Println(cwd)
	}

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Println("Unable to get Docker context. Please ensure that Docker is downloaded and running")
		panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		str1 := "bytes memory is allocated\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		fmt.Println(result)
		fmt.Println("Docker setup looks good")
	}

}

func getPackagePaths(inputArr []string, flags []string) (packagePaths []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				packagePath := strings.Replace(i, flag, "", 1)
				packagePath = strings.Trim(packagePath, "\"")
				packagePaths = append(packagePaths, packagePath)
			}
		}
	}
	return
}

func getEnvironmentVariables(inputArr []string, flags []string) (environmentVariables []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				environmentVariables = append(environmentVariables, strings.SplitN(i, "=", 2)...)
			}
		}
	}
	return
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func extractCommands(startupCommands []string) (environmentVariables []string, deployCommand string, otherFlags []string, deployEnvironment string, packages []string, customPackagePaths []string) {

	for _, option := range startupCommands {
		switch {
		case sliceContains([]string{"init", "up", "down", "destroy"}, option):
			deployCommand = option
		case sliceContains([]string{"docker", "kubernetes", "k8s"}, option):
			deployEnvironment = option
		case strings.HasPrefix(option, "-c=") || strings.HasPrefix(option, "--custom-package="):
			customPackagePaths = append(customPackagePaths, option)
		case strings.HasPrefix(option, "-e=") || strings.HasPrefix(option, "--env-file="):
			environmentVariables = append(environmentVariables, option)
		case strings.HasPrefix(option, "-") || strings.HasPrefix(option, "--"):
			otherFlags = append(otherFlags, option)
		default:
			packages = append(packages, option)
		}
	}

	if len(customPackagePaths) > 0 {
		customPackagePaths = getPackagePaths(customPackagePaths, []string{"-c=", "--custom-package="})
	}

	if len(environmentVariables) > 0 {
		environmentVariables = getEnvironmentVariables(environmentVariables, []string{"-e=", "--env-file="})
	}
	return
}

func RunDirectDockerCommand(startupCommands []string) {
	fmt.Println("Note: Initial setup takes 1-5 minutes.\nWait for the DONE message.\n--------------------------")

	environmentVariables, deployCommand, otherFlags, deployEnvironment, packages, customPackagePaths := extractCommands(startupCommands)

	fmt.Println("Environment:", deployEnvironment)
	fmt.Println("Action:", deployCommand)
	fmt.Println("Package IDs:", packages)
	fmt.Println("Custom package paths:", customPackagePaths)
	fmt.Println("Environment Variables:", environmentVariables)
	fmt.Println("Other Flags:", otherFlags)

	home, _ := os.UserHomeDir()

	if deployCommand == "init" {
		fmt.Println("\n\nDelete a pre-existing instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		RunDockerCommand(commandSlice...)
	}

	fmt.Println("Creating fresh instant container with volumes...")
	commandSlice := []string{
		"create",
		"--rm",
		"--mount=type=volume,src=instant,dst=/instant",
		"--name", "instant-openhie",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", home + "/.kube/config:/root/.kube/config:ro",
		"-v", home + "/.minikube:/home/$USER/.minikube:ro",
		"--network", "host",
	}
	commandSlice = append(commandSlice, environmentVariables...)
	commandSlice = append(commandSlice, []string{"openhie/instant:latest", deployCommand}...)
	commandSlice = append(commandSlice, otherFlags...)
	commandSlice = append(commandSlice, []string{"-t", deployEnvironment}...)
	commandSlice = append(commandSlice, packages...)
	RunDockerCommand(commandSlice...)

	fmt.Println("Adding 3rd party packages to instant volume:")

	for _, c := range customPackagePaths {
		commandSlice = []string{"cp", c, "instant-openhie:instant/"}
		RunDockerCommand(commandSlice...)
	}

	fmt.Println("\nRun Instant OpenHIE Installer Container")
	commandSlice = []string{"start", "-a", "instant-openhie"}
	RunDockerCommand(commandSlice...)

	if deployCommand == "destroy" {
		fmt.Println("Delete instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		RunDockerCommand(commandSlice...)
	}
}

func RunDockerCommand(commandSlice ...string) {
	cmd := exec.Command("docker", commandSlice...)
	cmdReader, err := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd.", stderr.String(), err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd.", stderr.String(), err)
		return
	}
}
