package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
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
		consoleSender(server, "Docker setup looks good")
	}

	// fmt.Println(reflect.TypeOf(containers).String())
	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}

// ComposeGet gets a docker-cmpose from github
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
