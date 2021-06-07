package main

import (
	"embed"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/r3labs/sse"
)

var server *sse.Server

//go:embed banner.txt
var f embed.FS

func main() {

	router := mux.NewRouter()
	server = sse.New()
	server.AutoReplay = true
	server.CreateStream("messages")
	addHandler(router)
	data, _ := f.ReadFile("banner.txt")
	color.Green(string(data))

	color.Cyan("Version: 1.02b")
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	// mainMenu()
	if len(os.Args) > 1 {
		CLI()
	} else {
		selectSetup()
	}
}
