package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/openhie/instant/goinstant/pkg"
)

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
		pkg.Debug()
		stuff := pkg.ComposeGet(stack)
		pkg.ComposeUp(stuff)

		color.Green.Println("A browser will open http://localhost:27517")
		color.Red.Println("Enter 'control c' key combination to stop the utility.")
		color.Println("Then stop containers by running 'goinstant' again and choosing 'Stop OpenHIE")
		pkg.OpenBrowser("http://localhost:27517")

		http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, you hit foo!")
		})

		box := packr.New("someBoxName", "./templates")
		http.Handle("/", http.FileServer(box))

		// this will stay open and block opening a new browser
		http.ListenAndServe(":27517", nil)

	case "Stop Instant OpenHIE":
		stuff := pkg.ComposeGet(stack)
		pkg.ComposeDown(stuff)
	case "Debug":
		pkg.Debug()
	case "Help":
		pkg.Help()
	case "Quit":
		os.Exit(1)
	}

}
