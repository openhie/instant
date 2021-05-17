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
		Items: []string{"Use Docker on your PC", "Use a Kubernetes Cluster", "Install FHIR package", "Quit"},
		Size:  12,
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

	case "Install FHIR package":
		selectUtil()

	case "Quit":
		os.Exit(0)
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
// 		os.Exit(0)
// 	}

// }

func selectUtil() {
	fmt.Println("Enter URL for the published package")
	// prompt for url
	prompt := promptui.Prompt{
		Label: "URL",
	}

	ig_url, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fhir_server, params := selectFHIR()
	fmt.Println("FHIR Server target:", fhir_server)
	loadIGpackage(ig_url, fhir_server, params)
	selectSetup()
}

func selectPackageDocker() {

	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{"Launch Core (Required, Start Here)", "Launch Facility Registry", "Launch Workforce", "Stop and Cleanup Core", "Stop and Cleanup Facility Registry", "Stop and Cleanup Workforce", "Stop All Services and Cleanup Docker", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Launch Core (Required, Start Here)":
		fmt.Println("...Setting up Core Package")
		SomeStuffDirect("docker", "core", "init")
		SomeStuffDirect("docker", "core", "up")
		fmt.Println("OpenHIM Console: http://localhost:9000/\nUser: root@openhim.org password: openhim-password")
		// now working
		// fmt.Printlnntln("HAPI FHIR base URL: http://localhost:3447/")
		selectPackageDocker()

		// fmt.Print("Press 'Enter' to continue...")
		// bufio.NewReader(os.Stdin).ReadBytes('\n')

	case "Launch Facility Registry":
		fmt.Println("...Setting up Facility Registry Package")
		SomeStuffDirect("docker", "facility", "up")
		selectPackageDocker()

	case "Launch Workforce":
		fmt.Println("...Setting up Workforce Package")
		SomeStuffDirect("docker", "healthworker", "up")
		selectPackageDocker()

	case "Stop and Cleanup Core":
		fmt.Println("Stopping and Cleaning Up Core...")
		SomeStuffDirect("docker", "core", "destroy")
		selectPackageDocker()

	case "Stop and Cleanup Facility Registry":
		fmt.Println("Stopping and Cleaning Up Facility Registry...")
		SomeStuffDirect("docker", "facility", "destroy")
		selectPackageDocker()

	case "Stop and Cleanup Workforce":
		fmt.Println("Stopping and Cleaning Up Workforce...")
		SomeStuffDirect("docker", "healthworker", "destroy")
		selectPackageDocker()

	case "Stop All Services and Cleanup Docker":
		// composeDownCore()
		fmt.Println("Stopping and Cleaning Up Everything...")
		SomeStuffDirect("docker", "core", "destroy")
		SomeStuffDirect("docker", "facility", "destroy")
		SomeStuffDirect("docker", "healthworker", "destroy")
		selectPackageDocker()

	// case "Developer Mode":
	// selectPackageDockerDev()
	// selectPackageDocker()

	case "Quit":
		os.Exit(0)

	case "Back":
		selectSetup()
	}

}

func selectPackageCluster() {

	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{"Launch Core (Required, Start Here)", "Launch Facility Registry", "Launch Workforce", "Stop and Cleanup Core", "Stop and Cleanup Facility Registry", "Stop and Cleanup Workforce", "Stop All Services and Cleanup Kubernetes", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	switch result {
	case "Launch Core (Required, Start Here)":
		fmt.Println("...Setting up Core Package")
		SomeStuffDirect("k8s", "core", "init")
		SomeStuffDirect("k8s", "core", "up")
		selectPackageCluster()

	case "Launch Facility Registry":
		fmt.Println("...Setting up Facility Registry Package")
		SomeStuffDirect("k8s", "facility", "up")
		selectPackageCluster()

	case "Launch Workforce":
		fmt.Println("...Setting up Workforce Package")
		SomeStuffDirect("k8s", "healthworker", "up")
		selectPackageCluster()

	case "Stop and Cleanup Core":
		fmt.Println("Stopping and Cleaning Up Core...")
		SomeStuffDirect("k8s", "core", "destroy")
		selectPackageCluster()

	case "Stop and Cleanup Facility Registry":
		fmt.Println("Stopping and Cleaning Up Facility Registry...")
		SomeStuffDirect("k8s", "facility", "destroy")
		selectPackageCluster()

	case "Stop and Cleanup Workforce":
		fmt.Println("Stopping and Cleaning Up Workforce...")
		SomeStuffDirect("k8s", "healthworker", "destroy")
		selectPackageCluster()

	case "Stop All Services and Cleanup Kubernetes":
		// composeDownCore()
		fmt.Println("Stopping and Cleaning Up Everything...")
		SomeStuffDirect("k8s", "core", "destroy")
		SomeStuffDirect("k8s", "facility", "destroy")
		SomeStuffDirect("k8s", "healthworker", "destroy")
		selectPackageCluster()

	// case "Developer Mode":
	// 	selectPackageDockerDev()
	// 	// selectPackageCluster()

	case "Quit":
		os.Exit(0)

	case "Back":
		selectSetup()
	}

}

func selectFHIR() (result_url string, params *Params) {

	prompt := promptui.Select{
		Label: "Select or enter URL for a FHIR Server",
		Items: []string{"Docker Default", "Kubernetes Default", "Use Public HAPI Server", "Enter a Server URL", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	fmt.Printf("You choose %q\n", result)

	switch result {

	case "Docker Default":
		result_url := "http://localhost:8080/fhir"
		params := &Params{}
		params.TypeAuth = "Custom"
		params.Token = "test"
		return result_url, params

	case "Kubernetes Default":
		result_url := "http://localhost:8080/fhir"
		params := &Params{}
		params.TypeAuth = "Custom"
		params.Token = "test"
		return result_url, params

	case "Use Public HAPI Server":
		result_url := "http://hapi.fhir.org/baseR4"
		params := &Params{}
		params.TypeAuth = "None"
		return result_url, params

	case "Enter a Server URL":
		prompt := promptui.Prompt{
			Label: "URL",
		}
		result_url, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		// TODO: validate URL
		// params.TypeAuth =
		params := selectParams()
		return result_url, params

	case "Quit":
		os.Exit(0)
		params := &Params{}
		return "", params

	case "Back":
		selectUtil()
		params := &Params{}
		return "", params

	}
	return result_url, params

}

type Params struct {
	// none, token, basic, custom
	TypeAuth  string
	Token     string
	BasicUser string
	BasicPass string
}

func selectParams() *Params {

	a := &Params{}

	prompt := promptui.Select{
		Label: "Choose authentication type",
		Items: []string{"None", "Basic", "Token", "Custom", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	fmt.Printf("You choose %q\n", result)
	switch result {

	case "None":
		a.TypeAuth = "None"
		return a

	case "Basic":
		a.TypeAuth = "Basic"

		// basic user
		prompt_basic_user := promptui.Prompt{
			Label: "Basic User",
		}
		result_basic_user, err := prompt_basic_user.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		a.BasicUser = result_basic_user

		// basic pass
		prompt_basic_pass := promptui.Prompt{
			Label: "Basic Password",
		}
		result_basic_pass, err := prompt_basic_pass.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		a.BasicPass = result_basic_pass

		return a

	case "Token":
		a.TypeAuth = "Token"

		// bearer token
		prompt_token := promptui.Prompt{
			Label: "Bearer Token",
		}
		result_token, err := prompt_token.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		a.Token = result_token
		return a

	case "Custom":
		a.TypeAuth = "Custom"

		// custom token
		prompt_ctoken := promptui.Prompt{
			Label: "Custom Token",
		}
		result_ctoken, err := prompt_ctoken.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		a.Token = result_ctoken
		return a

	case "Quit":
		os.Exit(0)
		return a

	case "Back":
		selectUtil()
		return a
	}
	return a

}
