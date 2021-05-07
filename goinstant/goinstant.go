package main

import (
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/r3labs/sse"
)

var server *sse.Server

func main() {

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
	color.Cyan("Version: 1.02b")
	// color.Green("Site: http://localhost:27517")
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	// mainMenu()
	if len(os.Args) > 1 {
		CLI()
	} else {
		selectSetup()
	}
	// c := existDisclaimer()
	// switch c {
	// case "fail":
	// 	const url = "http://localhost:27517/disclaimer.html"
	// 	browser.OpenURL(url)
	// 	http.ListenAndServe(":27517", router)

	// disable cli
	// mainMenu()

	// case "success":
	// 	const url = "http://localhost:27517/index.html"
	// 	browser.OpenURL(url)
	// 	http.ListenAndServe(":27517", router)

	// disable cli
	// mainMenu()

	// }
}
