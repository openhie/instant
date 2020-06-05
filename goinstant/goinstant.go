package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
)

func main() {

	version := "1.0.0-alpha"

	// defaults are not used for package or state
	headlessPtr := flag.Bool("headless", false, "headless mode, no prompts or ui. this is for automated testing.")
	packagePtr := flag.String("package", "", "[headless mode] package(s): core, core+hwf, core+facility, all")
	statePtr := flag.String("state", "", "[headless mode] up or down")

	flag.Parse()
	// packageflag := isFlagPassed(*packagePtr)
	// stateflag := isFlagPassed(*statePtr)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	})
	api.HandleFunc("/existdisclaimer", func(w http.ResponseWriter, r *http.Request) {
		x := existDisclaimer()
		fmt.Fprintln(w, x)
	})

	api.HandleFunc("/makedisclaimer", func(w http.ResponseWriter, r *http.Request) {
		x := makeDisclaimer()
		fmt.Fprintln(w, x)
	})

	api.HandleFunc("/makesetup", func(w http.ResponseWriter, r *http.Request) {
		x := makeSetup()
		fmt.Fprintln(w, x)
	})
	// api.HandleFunc("/debugdocker", func(w http.ResponseWriter, r *http.Request) {
	// 	z := debugDocker()
	// 	fmt.Fprintln(w, z)
	// })
	// api.HandleFunc("/composeup", func(w http.ResponseWriter, r *http.Request) {
	// 	z := composeUp()
	// 	fmt.Fprintln(w, z)
	// })
	// api.HandleFunc("/composedown", func(w http.ResponseWriter, r *http.Request) {
	// 	z := composeDown()
	// 	fmt.Fprintln(w, z)
	// })

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(pkger.Dir("/templates")))
	// Serve index page on all unhandled routes
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/index.html")
	})

	go http.ListenAndServe(":27517", router)

	pkgerPrint("/templates/banner.txt", "green")
	color.Green.Println("Version:", version)
	color.Green.Println("Site: http://localhost:27517")
	color.Green.Println("API: http://localhost:27517/api")
	color.Green.Println("The app can be run in headless mode. Run with -help to see options.\n")

	if *headlessPtr == true {

		color.Green.Println("Starting in headless mode...\n")
		fmt.Println("Options provided:")
		fmt.Println("package:", *packagePtr)
		fmt.Println("state:", *statePtr)
		if *statePtr == "" {
			color.Red.Println("state is empty but required. type goinstant -h for help.")
		}
		if *packagePtr == "" {
			color.Red.Println("package flag is empty but required. type goinstant -h for help.")
		}
	} else {
		color.Green.Println("Welcome to Instant. The tool can be run from the web interface or the prompt below.\n")
		color.Red.Println("Remember to clean up after your work or the app will continue to run in the background and have an adverse impact on performance.")

		c := existDisclaimer()
		fmt.Println(c)
		switch c {
		case "fail":
			go openBrowser("http://localhost:27517/disclaimer.html")
			cliDisclaimer()
			// TODO: disclaimer accept button hits: makeDisclaimer(), makeSetup(),
			//
			// then redirects to index.html
		case "success":
			makeDisclaimer()
			go openBrowser("http://localhost:27517")

		}
		// if c == "fail" {
		// 	go openBrowser("http://localhost:27517/disclaimer.html")
		// 	cliDisclaimer()
		// 	// accept hits makeDisclaimer() then redirects to index.html
		// } else if c == "success" {
		// 	makeDisclaimer()
		// 	go openBrowser("http://localhost:27517")
		// }
	}

	// TODO:
	// on the initial startup, the existdisclaimer runs
	// hit the endpoint for disclaimer which calls the function
	// func checks if the disclaimer exists, then send to index.html,
	// if not then to disclaimer page

	// selectSetup()
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
