package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/r3labs/sse"
)

var server *sse.Server

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

	server = sse.New()
	server.AutoReplay = true
	server.CreateStream("messages")
	router.HandleFunc("/events", sseHandler)

	router.HandleFunc("/index", Index)
	router.HandleFunc("/decline", Decline)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(pkger.Dir("/templates")))
	// Serve index page on all unhandled routes
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "/index.html")
	// })

	// go stanleySender(server)

	go http.ListenAndServe(":27517", router)

	pkgerPrint("/templates/banner.txt", "green")
	color.Green.Println("Version:", version)
	color.Green.Println("Site: http://localhost:27517")
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
		switch c {
		case "fail":
			go openBrowser("http://localhost:27517/disclaimer.html")
			cliDisclaimer()
			// TODO: disclaimer accept button hits: makeDisclaimer(), makeSetup(),
			//
			// then redirects to index.html
		case "success":
			makeDisclaimer()
			go openBrowser("http://localhost:27517/index.html")
			// locks into prompts
			setup()
			selectSetup()

		}
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

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("new client", r.Header.Get("X-Forwarded-For"))
	server.HTTPHandler(w, r)
}

func consoleSender(server *sse.Server, text string) {

	server.Publish("messages", &sse.Event{
		Data: []byte(text),
	})

}
