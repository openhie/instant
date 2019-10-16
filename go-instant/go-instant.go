package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
)

// func debug does checks on simple but potentially troublesome issues
func debug() {

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		color.Red.Println("Can't get current working directory... this is not a great error.")
		panic(err)
	} else {
		color.Green.Println("Running go-instant [v-alpha] from:", cwd)
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
	// List running containers.
	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}

func getstarted() {
	// Create some places to put downloaded compose files. This command won't wipe out anything.
	if _, err := os.Stat("sandbox/core/"); os.IsNotExist(err) {
		os.MkdirAll("sandbox/core/", 0700)
	}

	// This one will actually remove stuff.
	err := os.Remove("sandbox/core/docker-compose.yml")
	if err != nil {
		fmt.Println(err)
		// return
	}

	// Grab the sandbox core compose file.
	resp, err := grab.Get(
		"sandbox/core/", "https://github.com/openhie/instant/blob/strawperson/sandbox/core/docker-compose.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Download saved to", resp.Filename)
	color.Green.Println("You're ready to go! \nOpen a Windows cmd, PowerShell, or shell in Linux/Mac.")
	color.Green.Println("...then navigate to the sandbox/core directory and run 'docker-compose up'")

}

func help() {
	fmt.Printf("This app is used to troubleshoot and help with Instant OpenHIE. \n")
	fmt.Printf("You don't need this app for Instant OpenHIE if you're comfortable with the command line. \n")
	fmt.Printf("The way this app works is that it creates a folder wherever the app is installed. \n")
	fmt.Printf("Then it downloads docker-compose files from 'github.com/openhie/instant'.\n")
}

func main() {

	prompt := promptui.Select{
		Label: "If you're just starting, choose Get Started, otherwise choose as you wish...",
		Items: []string{"Get Started", "Debug", "Help", "OpenHIE Core", "Kitchen Sink"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Get Started" {
		fmt.Printf("You chose %q\n", result)
		getstarted()
	}

	if result == "Debug" {
		fmt.Printf("You chose %q\n", result)
		debug()
	}

	if result == "Help" {
		fmt.Printf("You chose %q\n", result)
		help()
	}

	if result == "OpenHIE Core" {
		fmt.Printf("You chose %q\n", result)
		debug()
		getstarted()
	}

	if result == "Kitchen Sink" {
		fmt.Printf("You chose %q\n", result)
		debug()
		getstarted()
	}

}
