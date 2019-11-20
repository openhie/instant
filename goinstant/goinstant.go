package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gobuffalo/packr/v2"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
)

// also see: https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
func debug() {
	cwd, err := os.Getwd()
	if err != nil {
		color.Red.Println("Can't get current working directory... this is not a great error.")
		panic(err)
	} else {
		color.Green.Println("Running goinstant [v-alpha] from:", cwd)
	}

	// ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		color.Red.Println("Unable to list Docker containers. Please ensure that Docker is downloaded and running!")
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		color.Red.Println("Unable to get Docker context. Please ensure that Docker is downloaded and running!")
		panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		fmt.Printf("%d bytes memory is allocated.\n", info.MemTotal)
	}

	// fmt.Println(reflect.TypeOf(containers).String())
	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}

func composeGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		color.Red.Println("Are you connected to the Internet? Error:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red.Println("Strange error reading the downloaded body. Error:", err)
	}
	fmt.Println(string(body))
	return (string(body))
}

func composeUp(composeFile string) {

	fmt.Println("Running on", runtime.GOOS)
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
		fmt.Println("What operating system is this?", runtime.GOOS)
	}

	color.Green.Println("A browser will open http://localhost:27517")
	color.Red.Println("Enter 'Control C' key combination to stop the utility.")
	color.Println("Then stop containers by running 'goinstant' again and choosing 'Stop OpenHIE")
	openBrowser("http://localhost:27517")
}

func composeDown(composeFile string) {
	fmt.Println("Running on", runtime.GOOS)
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
		fmt.Println("What operating system is this?", runtime.GOOS)

	}
}

func help() {
	fmt.Printf("This is a utility to help run Instant OpenHIE on your personal computer.. \n")
	fmt.Printf("You don't need this app for Instant OpenHIE if you're comfortable with the command line. \n")
	fmt.Printf("The utility downloads docker-compose files from 'github.com/openhie/instant'.\n")
}

func main() {

	prompt := promptui.Select{
		Label: "Choose Start Instant OpenHIE if this is your first time",
		Items: []string{"Start Instant OpenHIE", "Stop Instant OpenHIE", "Debug", "Help", "Quit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You chose %q\n", result)

	stack := "https://raw.github.com/openhie/instant/strawperson/core/docker/docker-compose.yml"

	switch result {
	case "Start Instant OpenHIE":
		debug()
		stuff := composeGet(stack)
		composeUp(stuff)

		box := packr.New("someBoxName", "./templates")
		http.Handle("/", http.FileServer(box))
		// this will stay open and block opening a new browser
		http.ListenAndServe(":27517", nil)

	case "Stop Instant OpenHIE":
		stuff := composeGet(stack)
		composeDown(stuff)
	case "Debug":
		debug()
	case "Help":
		help()
	case "Quit":
		os.Exit(1)
	}

}
