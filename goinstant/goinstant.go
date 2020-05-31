package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/markbates/pkger"
)

func main() {

	// defaults are not used for package or state
	headlessPtr := flag.Bool("headless", false, "headless mode, no prompts or ui. this is for automated testing.")
	packagePtr := flag.String("package", "", "package(s) to install: core, core+hwf, core+facility, all")
	statePtr := flag.String("state", "", "up or down")

	flag.Parse()
	// packageflag := isFlagPassed(*packagePtr)
	// stateflag := isFlagPassed(*statePtr)

	fmt.Println("headless:", *headlessPtr)
	fmt.Println("package:", *packagePtr)
	fmt.Println("state:", *statePtr)

	switch {
	case *headlessPtr == true:
		color.Green.Println("headless flag is true")
		if *statePtr == "" {
			color.Red.Println("state is empty but required. type goinstant -h for help.")
		}
		if *packagePtr == "" {
			color.Red.Println("package flag is empty but required. type goinstant -h for help.")
		}
	case *headlessPtr == false:
		pkgerPrint("/banner.txt", "green")
		checkDisclaimer()
	}

}

// // https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go
// func isFlagPassed(name string) bool {
// 	found := false
// 	flag.Visit(func(f *flag.Flag) {
// 		if f.Name == name {
// 			found = true
// 		}
// 	})
// 	return found
// }

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

func checkDisclaimer() {
	home, _ := os.UserHomeDir()
	// must use filepath.join not path.join for windows compat
	dotfiles := filepath.Join(home, ".instant")
	fileName := filepath.Join(dotfiles, "accept_disclaimer")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		pkgerPrint("/disclaimer.txt", "yellow")
		prompt := promptui.Select{
			Label: "Do you agree to use this application?",
			Items: []string{"Yes", "No", "Quit"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Printf("You chose %q\n", result)

		switch result {
		case "Yes":
			makeDisclaimer()
			setup()
			selectSetup()
		case "No":
			fmt.Println("Understood. Exiting.")
			os.Exit(1)
		case "Quit":
			os.Exit(1)
		}
	} else {
		selectSetup()
	}
}

func makeDisclaimer() {
	home, _ := os.UserHomeDir()
	// must use filepath.join not path.join for windows compat
	dotfiles := filepath.Join(home, ".instant")
	fileName := filepath.Join(dotfiles, "accept_disclaimer")
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dotfiles, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(fileName, currentTime, currentTime)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func setup() {

	color.Green.Println("getting ready")

	home, _ := os.UserHomeDir()
	dotfiles := filepath.Join(home, ".instant")
	fmt.Println("...checking for config folder at:", dotfiles)

	// check for dotfolder, create if it doesn't exist
	if _, err := os.Stat(dotfiles); os.IsNotExist(err) {
		fmt.Println("config folder does not exist, creating it")
		os.Mkdir(dotfiles, 0700)
		color.Green.Println("created", dotfiles)
	} else {
		color.Green.Println("config folder exists")
	}

	// check repo and clone or pull
	fmt.Println("...cloning or pulling latest code from repo")

	repo := filepath.Join(dotfiles, "instant")
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		fmt.Println("...repo folder does not exist, cloning it")
		_, err := git.PlainClone(repo, false, &git.CloneOptions{
			URL:      "https://github.com/openhie/instant",
			Progress: os.Stdout,
		})
		if err != nil {
		}
		color.Green.Println("successfully cloned", dotfiles)
	} else {
		fmt.Println("...repo folder exists, pulling changes")
		const (
			repoURL = "https://github.com/openhie/instant.git"
		)

		dir, _ := ioutil.TempDir("", "temp_dir")

		options := &git.CloneOptions{
			URL: repoURL,
		}

		_, err := git.PlainClone(dir, false, options)
		if err != nil {
		}

	}

	color.Green.Println("git repo is ready")

	color.Green.Println("ready!")

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
