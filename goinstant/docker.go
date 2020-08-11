package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gookit/color"
	"golang.org/x/net/context"
)

func debugDocker() {

	consoleSender(server, "...checking your Docker setup")

	cwd, err := os.Getwd()
	if err != nil {
		consoleSender(server, "Can't get current working directory... this is not a great error.")
		// panic(err)
	} else {
		consoleSender(server, cwd)
	}

	// ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		// panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		consoleSender(server, "Unable to list Docker containers. Please ensure that Docker is downloaded and running")
		// return
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		consoleSender(server, "Unable to get Docker context. Please ensure that Docker is downloaded and running")
		// panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		fmt.Printf("%d bytes memory is allocated.\n", info.MemTotal)
		str1 := " bytes memory is allocated.\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		consoleSender(server, result)
		consoleSender(server, "Docker setup looks good")
	}

	// fmt.Println(reflect.TypeOf(containers).String())
	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}

// ComposeGet gets a docker-compose from github
func composeGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		color.Red.Println("Are you connected to the Internet? Error:", err)
		consoleSender(server, "Are you connected to the Internet? Error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red.Println("Strange error reading the downloaded body. Error:", err)
		consoleSender(server, "Strange error reading the downloaded body.")
	}
	fmt.Println(string(body))
	return (string(body))
}

// SomeStuff tests some stuff - new way to stream real-time, only test core package
func SomeStuff() {
	home, _ := os.UserHomeDir()
	cmd := exec.Command("docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"/.kube/config:/root/.kube/config:ro", "-v", home+"/.minikube:/home/$USER/.minikube:ro", "--mount=type=volume,src=instant,dst=/instant", "--network", "host", "openhie/instant:latest", "init", "-t", "docker", "core")
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

func composeUpCoreDOD() {

	home, _ := os.UserHomeDir()
	color.Yellow.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		// cmd := exec.Command("docker-compose", "-f", composefile, "up", "-d")
		cmd := exec.Command("docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"/.kube/config:/root/.kube/config:ro", "-v", home+"/.minikube:/home/$USER/.minikube:ro", "--mount=type=volume,src=instant,dst=/instant", "--network", "host", "openhie/instant:latest", "init", "-t", "docker")

		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)

		}
		consoleSender(server, outb.String())
		fmt.Println("out:", outb.String(), "err:", errb.String())

	case "windows":
		// cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "up", "-d")
		cmd := exec.Command("cmd", "/C", "docker", "run", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", home+"\\.kube:/root/.kube/config:ro", "--mount=type=volume,src=instant,dst=/instant", "openhie/instant:latest", "init", "-t", "docker")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		consoleSender(server, "What operating system is this?")
	}

}

func composeUpCore() {

	home, _ := os.UserHomeDir()
	composefile := path.Join(home, ".instant/instant/core/docker/docker-compose.yml")
	color.Yellow.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", composefile, "up", "-d")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "up", "-d")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		consoleSender(server, "What operating system is this?")
	}

}

func composeDownCore() {

	home, _ := os.UserHomeDir()
	composefile := path.Join(home, ".instant/instant/core/docker/docker-compose.yml")
	color.Yellow.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", composefile, "down")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "down")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		consoleSender(server, "What operating system is this?")
	}

}

// ComposeUp brings up based on docker-compose
func composeUp(composeFile string) {

	color.Yellow.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", "-", "up", "-d")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "up", "-d")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		consoleSender(server, "What operating system is this?")
	}

}

// ComposeDown stops containers based on docker-compose
func composeDown(composeFile string) {
	color.Yellow.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", "-", "down")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "down")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		consoleSender(server, "What operating system is this?")
	}
}
