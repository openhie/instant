package main

import (
	"fmt"
	"io/ioutil"
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
		Label: "Please choose how you want to run Instant. \nChoose Docker if you're running on your PC. \nIf you want to run Instant on Kubernetes, then you have should been provided credentials or have Kubernetes running on your PC.",
		Items: []string{"Use Docker on your PC", "Use a Kubernetes Cluster", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Use Docker on your PC":
		debugDocker()
		listDocker()
		selectPackageDocker()

	case "Use a Kubernetes Cluster":
		debugKubernetes()
		// configServerKubernetes()
		selectPackageCluster()

	// case "Use Cluster on Remote Server":
	// 	fmt.Println("Great, but this feature isn't built yet.")
	// 	debugKubernetes()
	// 	configServerKubernetes()
	// 	selectPackageClusterRemote()

	case "Quit":
		os.Exit(1)
	}

}

// func selectDocker() {

// 	prompt := promptui.Select{
// 		Label: "Setup",
// 		Items: []string{"Check Docker again", "Clean up Docker", "Quit"},
// 	}

// 	_, result, err := prompt.Run()

// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return
// 	}

// 	fmt.Printf("You choose %q\n", result)

// 	switch result {
// 	case "Check Docker again":
// 		debugDocker()
// 		selectPackageDocker()
// 	case "Clean up Docker":
// 		fmt.Println("This feature needs more work, sorry.")
// 	case "Quit":
// 		os.Exit(1)
// 	}

// }

func selectPackageDocker() {

	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{"Launch Core (Required, Start Here)", "Launch Facility Registry", "Launch Workforce", "Stop and Cleanup Core", "Stop and Cleanup Facility Registry", "Stop and Cleanup Workforce", "Stop All Services and Cleanup Docker", "Developer Mode", "Quit", "Back"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Launch Core (Required, Start Here)":
		println("...Setting up Core Package")
		SomeStuffDirect("docker", "core", "init")
		SomeStuffDirect("docker", "core", "up")
		println("OpenHIM Console: http://localhost:9000/\nUser: root@openhim.org password: openhim-password")
		println("HAPI FHIR base URL: http://localhost:3447/")
		selectPackageDocker()

		// fmt.Print("Press 'Enter' to continue...")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')

	case "Launch Facility Registry":
		println("...Setting up Facility Registry Package")
		SomeStuffDirect("docker", "facility", "up")
		selectPackageDocker()

	case "Launch Workforce":
		println("...Setting up Workforce Package")
		SomeStuffDirect("docker", "healthworker", "up")
		selectPackageDocker()

	case "Stop and Cleanup Core":
		println("Stopping and Cleaning Up Core...")
		SomeStuffDirect("docker", "core", "destroy")
		selectPackageDocker()

	case "Stop and Cleanup Facility Registry":
		println("Stopping and Cleaning Up Facility Registry...")
		SomeStuffDirect("docker", "facility", "destroy")
		selectPackageDocker()

	case "Stop and Cleanup Workforce":
		println("Stopping and Cleaning Up Workforce...")
		SomeStuffDirect("docker", "healthworker", "destroy")
		selectPackageDocker()

	case "Stop All Services and Cleanup Docker":
		// composeDownCore()
		println("Stopping and Cleaning Up Everything...")
		SomeStuffDirect("docker", "core", "destroy")
		SomeStuffDirect("docker", "facility", "destroy")
		SomeStuffDirect("docker", "healthworker", "destroy")
		selectPackageDocker()

	case "Developer Mode":
		selectPackageDockerDev()
		// selectPackageDocker()

	case "Quit":
		os.Exit(1)

	case "Back":
		selectSetup()
	}

}

func selectPackageCluster() {

	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{"Launch Core (Required, Start Here)", "Launch Facility Registry", "Launch Workforce", "Stop and Cleanup Core", "Stop and Cleanup Facility Registry", "Stop and Cleanup Workforce", "Stop All Services and Cleanup Kubernetes", "Developer Mode", "Quit", "Back"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Launch Core (Required, Start Here)":
		println("...Setting up Core Package")
		SomeStuffDirect("k8s", "core", "init")
		SomeStuffDirect("k8s", "core", "up")
		selectPackageCluster()

	case "Launch Facility Registry":
		println("...Setting up Facility Registry Package")
		SomeStuffDirect("k8s", "facility", "up")
		selectPackageCluster()

	case "Launch Workforce":
		println("...Setting up Workforce Package")
		SomeStuffDirect("k8s", "healthworker", "up")
		selectPackageCluster()

	case "Stop and Cleanup Core":
		println("Stopping and Cleaning Up Core...")
		SomeStuffDirect("k8s", "core", "destroy")
		selectPackageCluster()

	case "Stop and Cleanup Facility Registry":
		println("Stopping and Cleaning Up Facility Registry...")
		SomeStuffDirect("k8s", "facility", "destroy")
		selectPackageCluster()

	case "Stop and Cleanup Workforce":
		println("Stopping and Cleaning Up Workforce...")
		SomeStuffDirect("k8s", "healthworker", "destroy")
		selectPackageCluster()

	case "Stop All Services and Cleanup Kubernetes":
		// composeDownCore()
		println("Stopping and Cleaning Up Everything...")
		SomeStuffDirect("k8s", "core", "destroy")
		SomeStuffDirect("k8s", "facility", "destroy")
		SomeStuffDirect("k8s", "healthworker", "destroy")
		selectPackageCluster()

	case "Developer Mode":
		selectPackageDockerDev()
		// selectPackageCluster()

	case "Quit":
		os.Exit(1)

	case "Back":
		selectSetup()
	}

}

// func selectPackageClusterLocal() {

// 	prompt := promptui.Select{
// 		Label: "Great, now choose a package",
// 		Items: []string{"Core", "Core + Facility", "Core + Facility + Workforce", "Quit"},
// 	}

// 	_, result, err := prompt.Run()

// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return
// 	}

// 	fmt.Printf("You choose %q\n", result)

// 	switch result {
// 	case "Core":
// 		fmt.Println("Core")
// 	case "Core + Facility":
// 		fmt.Println("Core + Facility")
// 	case "Core + Facility + Workforce":
// 		fmt.Println("Core + Facility + Workforce")
// 	case "Quit":
// 		os.Exit(1)
// 	}

// }

// func selectPackageClusterRemote() {

// 	prompt := promptui.Select{
// 		Label: "Great, now choose a package",
// 		Items: []string{"Core", "Core + Facility", "Core + Facility + Workforce", "Quit"},
// 	}

// 	_, result, err := prompt.Run()

// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return
// 	}

// 	fmt.Printf("You choose %q\n", result)

// 	switch result {
// 	case "Core":
// 		fmt.Println("Core")
// 	case "Core + Facility":
// 		fmt.Println("Core + Facility")
// 	case "Core + Facility + Workforce":
// 		fmt.Println("Core + Facility + Workforce")
// 	case "Quit":
// 		os.Exit(1)
// 	}

// }

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
// func mainMenu() {

// 	prompt := promptui.Select{
// 		Label: "Developer Mode",
// 		Items: []string{"Setup", "Select Packages", "Start Instant OpenHIE", "Stop Instant OpenHIE", "Debug", "Help", "Quit"},
// 	}

// 	_, result, err := prompt.Run()
// 	if err != nil {
// 		fmt.Printf("Prompt failed %v\n", err)
// 		return
// 	}
// 	fmt.Printf("You chose %q\n", result)

// 	stack := "https://raw.github.com/openhie/instant/master/core/docker/docker-compose.yml"

// 	switch result {
// 	case "Setup":
// 		setup()
// 		selectSetup()
// 	case "Select Packages":
// 		fmt.Println("in-progress")
// 	case "Start Instant OpenHIE":

// 		// http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
// 		// 	fmt.Fprintln(w, "Hello, you hit foo!")
// 		// })

// 		// dir := http.FileServer(pkger.Dir("/templates"))
// 		// // use in goroutine to return control
// 		// go http.ListenAndServe(":27517", dir)
// 		// go openBrowser("http://localhost:27517")

// 		debugDocker()
// 		stuff := composeGet(stack)
// 		composeUp(stuff)

// 		// color.Green.Println("A browser will open http://localhost:27517")
// 		// color.Red.Println("Enter 'control c' key combination to stop the utility.")
// 		// color.Println("Then stop containers by running 'goinstant' again and choosing 'Stop OpenHIE")

// 	case "Stop Instant OpenHIE":
// 		stuff := composeGet(stack)
// 		composeDown(stuff)
// 	case "Debug":
// 		debugDocker()
// 	case "Help":
// 		help()
// 	case "Quit":
// 		os.Exit(1)
// 	}

// }
