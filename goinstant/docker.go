package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func debugDocker() {

	consoleSender(server, "...checking your Docker setup")
	println("...checking your Docker setup\n")

	cwd, err := os.Getwd()
	if err != nil {
		consoleSender(server, "Can't get current working directory... this is not a great error.")
		println("Can't get current working directory... this is not a great error.")
		// panic(err)
	} else {
		consoleSender(server, cwd)
		println(cwd)
	}

	// ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		consoleSender(server, "Unable to get Docker context. Please ensure that Docker is downloaded and running")
		println("Unable to get Docker context. Please ensure that Docker is downloaded and running")
		panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		str1 := "bytes memory is allocated\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		consoleSender(server, result)
		consoleSender(server, "Docker setup looks good\n")
		println(result)
		println("Docker setup looks good\n")
	}

}

// TODO: change printf to consoleSender
// listDocker may be used in future
func listDocker() {

	println("Listing containers...\n")
	consoleSender(server, "Listing running containers...")

	// ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		consoleSender(server, "Unable to list Docker containers. Please ensure that Docker is downloaded and running")
		println("Unable to list Docker containers. Please ensure that Docker is downloaded and running")
		// return
	}

	if len(containers) == 0 {
		println("No containers are running.\n")
	} else {
		for _, container := range containers {
			items := fmt.Sprintf("ContainerID: %s Status: %s Image: %s Names: %s", container.ID[:10], container.State, container.Image, container.Names)
			consoleSender(server, items)
			println(items)
		}
		println("\nContainers are already running.\nCleanup running containers in the Docker dashboard before continuing.")
	}

}

func SomeStuffDirect(runner string, pk string, state string) {
	consoleSender(server, "Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	println("Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	// runner := runner
	// pk := pk
	// state := state
	consoleSender(server, "Runner requested: "+runner)
	consoleSender(server, "Package requested: "+pk)
	consoleSender(server, "State requested: "+state)
	println("Runner requested: " + runner)
	println("Package requested: " + pk)
	println("State requested: " + state)

	home, _ := os.UserHomeDir()

	// args := []string{runner, "ever", "you", "like"}
	// cmd := exec.Command(app, args...)
	// consoleSender(server, args[0])

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
			println("\t > %s\n", scanner.Text())
			consoleSender(server, scanner.Text())
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

// SomeStuff HTTP API-only
func SomeStuff(r *http.Request) {
	consoleSender(server, "Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	runner := r.URL.Query().Get("runner")
	pk := r.URL.Query().Get("package")
	state := r.URL.Query().Get("state")
	consoleSender(server, "Runner requested: "+runner)
	consoleSender(server, "Package requested: "+pk)
	consoleSender(server, "State requested: "+state)
	home, _ := os.UserHomeDir()

	// args := []string{runner, "ever", "you", "like"}
	// cmd := exec.Command(app, args...)
	// consoleSender(server, args[0])

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
			println("\t > %s\n", scanner.Text())
			consoleSender(server, scanner.Text())
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
