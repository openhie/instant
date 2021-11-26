package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
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

	// ctx := context.Background()
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

// TODO: change printf to consoleSender
// listDocker may be used in future
func listDocker() {

	fmt.Println("Listing containers...")

	// ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println("Unable to list Docker containers. Please ensure that Docker is downloaded and running")
		// return
	}

	if len(containers) == 0 {
		fmt.Println("No containers are running.")
	} else {
		for _, container := range containers {
			items := fmt.Sprintf("ContainerID: %s Status: %s Image: %s Names: %s", container.ID[:10], container.State, container.Image, container.Names)
			fmt.Println(items)
		}
		fmt.Println("\nContainers are already running.\nCleanup running containers in the Docker dashboard before continuing.")
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

func extractFlags(customFlags []string) (customPackages []string, otherFlags []string) {
	testCustomPackage := func(s string) bool { return strings.HasPrefix(s, "-c=") || strings.HasPrefix(s, "--custom-package=") }
	customPackages, otherFlags = filterAndSplit(customFlags, testCustomPackage)

	if len(customPackages) > 0 {
		fmt.Print("Custom packages requested: ")
		customPackages = getPackagePaths(customPackages, []string{"-c=", "--custom-package="})
	}
	return
}

func RunDirectDockerCommand(runner string, pk string, action string, customFlags ...string) {
	var customPackages []string
	var otherFlags []string

	fmt.Println("Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	fmt.Println("Runner requested: " + runner)
	fmt.Println("Package requested: " + pk)
	fmt.Println("Action requested: " + action)

	if len(customFlags) > 3 {
		customFlags = customFlags[3:]
		customPackages, otherFlags = extractFlags(customFlags)
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
		"--network", "host",
		"openhie/instant:latest",
		action, "-t", runner, pk}
	commandSlice = append(commandSlice, otherFlags...)
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
	//fmt.Println(cmd.String()) //TODO: remove, use for debugging to see docker command
	// create a pipe for the output of the script
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
		/*TODO: The error handling here is terrible, we just get the exit code logged to the output,
		we need more than that to know if it's a valid error. Example, require docker login gives 125 exit code*/
		return
	}
}

// DockerCommand HTTP API-only
func RunHTTPDockerCommand(r *http.Request) {
	fmt.Printf("Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	runner := r.URL.Query().Get("runner")
	pk := r.URL.Query().Get("package")
	state := r.URL.Query().Get("state")
	home, _ := os.UserHomeDir()

	cmd := exec.Command("docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"/.kube/config:/root/.kube/config:ro", "-v", home+"/.minikube:/home/$USER/.minikube:ro", "--mount=type=volume,src=instant,dst=/instant", "--network", "host", "openhie/instant:latest", state, "-t", runner, pk)
	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
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

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}

}

// func composeUpCoreDOD() {

// 	home, _ := os.UserHomeDir()
// 	color.Yellow.Println("Running on", runtime.GOOS)
// 	switch runtime.GOOS {
// 	case "linux", "darwin":
// 		// cmd := exec.Command("docker-compose", "-f", composefile, "up", "-d")
// 		cmd := exec.Command("docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"/.kube/config:/root/.kube/config:ro", "-v", home+"/.minikube:/home/$USER/.minikube:ro", "--mount=type=volume,src=instant,dst=/instant", "--network", "host", "openhie/instant:latest", "init", "-t", "docker")

// 		var outb, errb bytes.Buffer
// 		cmd.Stdout = &outb
// 		cmd.Stderr = &errb
// 		// cmd.Stdout = os.Stdout
// 		// cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			log.Fatalf("cmd.Run() failed with %s\n", err)

// 		}
// 		consoleSender(server, outb.String())
// 		fmt.Println("out:", outb.String(), "err:", errb.String())

// 	case "windows":
// 		// cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "up", "-d")
// 		cmd := exec.Command("cmd", "/C", "docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"\\.kube:/root/.kube/config:ro", "--mount=type=volume,src=instant,dst=/instant", "openhie/instant:latest", "init", "-t", "docker")
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	default:
// 		consoleSender(server, "What operating system is this?")
// 	}

// }
