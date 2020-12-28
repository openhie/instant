package main

import (
	"net/http"

	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/pkg/browser"
	"github.com/r3labs/sse"
)

var server *sse.Server

func main() {

	version := "1.0.0-beta"

	router := mux.NewRouter()
	server = sse.New()
	server.AutoReplay = true
	server.CreateStream("messages")
	addHandler(router)
	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(pkger.Dir("/templates")))
	// Serve index page on all unhandled routes
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "/index.html")
	// })

	pkgerPrint("/templates/banner.txt", "green")
	color.Green.Println("Version:", version)
	color.Green.Println("Site: http://localhost:27517")
	color.Green.Println("Welcome to Instant.\n")
	color.Red.Println("Remember to clean up after your work or the app will continue to run in the background and have an adverse impact on performance.")

	c := existDisclaimer()
	switch c {
	case "fail":
		const url = "http://localhost:27517/disclaimer.html"
		browser.OpenURL(url)
		http.ListenAndServe(":27517", router)
		// disable cli
		// mainMenu()
	case "success":
		const url = "http://localhost:27517/index.html"
		browser.OpenURL(url)
		http.ListenAndServe(":27517", router)
		// disable cli
		// mainMenu()

	}
}
