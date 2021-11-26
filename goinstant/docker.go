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

func filterAndSplit(items []string, filterTest func(string) bool) (filtered []string, unfiltered []string) {
	for _, s := range items {
		if filterTest(s) {
			filtered = append(filtered, s)
		} else {
			unfiltered = append(unfiltered, s)
		}
	}
	return
}

func getPackagePaths(inputArr []string, flags []string) (packagePaths []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				packagePath := strings.Replace(i, flag, "", 1)
				packagePath = strings.Trim(packagePath, "\"")
				fmt.Print(packagePath + " ")
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
				fmt.Print(i)
				environmentVariables = append(environmentVariables, strings.SplitN(i, "=", 2)...)
			}
		}
	}
	return
}

func extractFlags(customFlags []string) (customPackages []string, environmentVariables []string, otherFlags []string) {
	testCustomPackage := func(s string) bool { return strings.HasPrefix(s, "-c=") || strings.HasPrefix(s, "--custom-package=") }
	customPackages, otherFlags = filterAndSplit(customFlags, testCustomPackage)

	testEnvVars := func(s string) bool { return strings.HasPrefix(s, "-e=") || strings.HasPrefix(s, "--env-file=") }
	environmentVariables, otherFlags = filterAndSplit(otherFlags, testEnvVars)

	if len(customPackages) > 0 {
		fmt.Print("Custom packages requested: ")
		customPackages = getPackagePaths(customPackages, []string{"-c=", "--custom-package="})
	}

	if len(environmentVariables) > 0 {
		fmt.Print("Environment Variables provided: ")
		environmentVariables = getEnvironmentVariables(environmentVariables, []string{"-e=", "--env-file="})
	}
	return
}

func RunDirectDockerCommand(runner string, pk string, action string, customFlags ...string) {
	var customPackages []string
	var otherFlags []string
	var environmentVariables []string

	fmt.Println("Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	fmt.Println("Runner requested: " + runner)
	fmt.Println("Package requested: " + pk)
	fmt.Println("Action requested: " + action)

	if len(customFlags) > 3 {
		customFlags = customFlags[3:]
		customPackages, environmentVariables, otherFlags = extractFlags(customFlags)
	}

	home, _ := os.UserHomeDir()

	if action == "init" {
		fmt.Println("\n\nDelete a pre-existing instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		RunDockerCommand(commandSlice...)
	}

	fmt.Println("Creating fresh instant container with volumes...")
	commandSlice := []string{"create", "--rm",
		"--mount=type=volume,src=instant,dst=/instant",
		"--name", "instant-openhie",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", home + "/.kube/config:/root/.kube/config:ro",
		"-v", home + "/.minikube:/home/$USER/.minikube:ro",
		"--network", "host"}
	commandSlice = append(commandSlice, environmentVariables...)
	commandSlice = append(commandSlice, []string{"openhie/instant:latest", action}...)
	commandSlice = append(commandSlice, otherFlags...)
	commandSlice = append(commandSlice, []string{"-t", runner, pk}...)
	RunDockerCommand(commandSlice...)

	fmt.Println("Adding 3rd party packages to instant volume:")

	for _, c := range customPackages {
		fmt.Print("- " + c)
		commandSlice = []string{"cp", c, "instant-openhie:instant/"}
		RunDockerCommand(commandSlice...)
	}

	fmt.Println("\nRun Instant OpenHIE Installer Container")
	commandSlice = []string{"start", "-a", "instant-openhie"}
	RunDockerCommand(commandSlice...)

	if action == "destroy" {
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
