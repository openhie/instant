package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/markbates/pkger"
)

func pkgerPrint(text string, scolor string) {

	f, err := pkger.Open(text)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	s := string(b)

	if scolor == "green" {
		color.Green.Println(s)
	}

	if scolor == "yellow" {
		color.Yellow.Println(s)
	}
}

func selectSetup() {

	prompt := promptui.Select{
		Label: "Please choose how you want to run Instant. \nChoose Docker if you're running on your PC. \nIf you want to run Instant on a cluster, then you have should been provided credentials. \nOnly choose Cluster on your PC if you're an expert user.",
		Items: []string{"Use Docker on your PC", "Use Cluster on your PC", "Use Cluster on Remote Server", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Use Docker on your PC":
		selectPackageDocker()

	case "Use Cluster on your PC":
		debugKubernetes()
		configServerKubernetes()
		selectPackageClusterLocal()

	case "Use Cluster on Remote Server":
		fmt.Println("Great, but this feature isn't built yet.")
		debugKubernetes()
		configServerKubernetes()
		selectPackageClusterRemote()

	case "Quit":
		os.Exit(1)
	}

}

func selectDocker() {

	debugDocker()

	prompt := promptui.Select{
		Label: "Setup",
		Items: []string{"Check Docker again", "Clean up Docker", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Check Docker again":
		debugDocker()
		selectPackageDocker()
	case "Clean up Docker":
		fmt.Println("This feature needs more work, sorry.")
	case "Quit":
		os.Exit(1)
	}

}

func selectPackageDocker() {

	prompt := promptui.Select{
		Label: "Great, now choose a package",
		Items: []string{"Core", "Core + Facility", "Core + Facility + Workforce", "Stop Services and Cleanup Docker", "Developer Mode", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Core":
		fmt.Println("...Setting up")
		composeUpCore()
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

	case "Core + Facility":
		fmt.Println("Core + Facility")
	case "Core + Facility + Workforce":
		fmt.Println("Core + Facility + Workforce")
	case "Stop Services and Cleanup Docker":
		composeDownCore()
	case "Developer Mode":
		selectPackageDockerDev()
	case "Quit":
		os.Exit(1)
	}

}

func selectPackageClusterLocal() {

	prompt := promptui.Select{
		Label: "Great, now choose a package",
		Items: []string{"Core", "Core + Facility", "Core + Facility + Workforce", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Core":
		fmt.Println("Core")
	case "Core + Facility":
		fmt.Println("Core + Facility")
	case "Core + Facility + Workforce":
		fmt.Println("Core + Facility + Workforce")
	case "Quit":
		os.Exit(1)
	}

}

func selectPackageClusterRemote() {

	prompt := promptui.Select{
		Label: "Great, now choose a package",
		Items: []string{"Core", "Core + Facility", "Core + Facility + Workforce", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Core":
		fmt.Println("Core")
	case "Core + Facility":
		fmt.Println("Core + Facility")
	case "Core + Facility + Workforce":
		fmt.Println("Core + Facility + Workforce")
	case "Quit":
		os.Exit(1)
	}

}

func selectPackageDockerDev() {

	prompt := promptui.Select{
		Label: "Great, now choose a package",
		Items: []string{"Core (dev.yml)", "Facility (w/o Core)", "Facility + Workforce (w/o Core)", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Core (dev.yml)":
		fmt.Println("Core (dev.yml)")
		// composeUpCore()
	case "Facility (w/o Core)":
		fmt.Println("Facility (w/o Core)")
	case "Facility + Workforce (w/o Core)":
		fmt.Println("Facility + Workforce (w/o Core)")
	case "Quit":
		os.Exit(1)
	}

}

// old start menu system
func mainMenu() {

	setup()

	prompt := promptui.Select{
		Label: "Choose Setup if this is your first time",
		Items: []string{"Setup", "Select Packages", "Start Instant OpenHIE", "Stop Instant OpenHIE", "Debug", "Help", "Quit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You chose %q\n", result)

	stack := "https://raw.github.com/openhie/instant/master/core/docker/docker-compose.yml"

	switch result {
	case "Check Setup Again":
		selectSetup()
	case "Select Packages":
		fmt.Println("in-progress")
	case "Start Instant OpenHIE":

		http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, you hit foo!")
		})

		dir := http.FileServer(pkger.Dir("/templates"))
		// use in goroutine to return control
		go http.ListenAndServe(":27517", dir)
		go openBrowser("http://localhost:27517")

		debugDocker()
		stuff := composeGet(stack)
		composeUp(stuff)

		// color.Green.Println("A browser will open http://localhost:27517")
		// color.Red.Println("Enter 'control c' key combination to stop the utility.")
		// color.Println("Then stop containers by running 'goinstant' again and choosing 'Stop OpenHIE")

	case "Stop Instant OpenHIE":
		stuff := composeGet(stack)
		composeDown(stuff)
	case "Debug":
		debugDocker()
	case "Help":
		help()
	case "Quit":
		os.Exit(1)
	}

}
