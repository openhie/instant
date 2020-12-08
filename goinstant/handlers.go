package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func addHandler(router *mux.Router) {
	router.HandleFunc("/events", sseHandler)
	router.HandleFunc("/setup", Setup)
	router.HandleFunc("/disclaimer", Disclaimer)
	router.HandleFunc("/decline", Decline)
	router.HandleFunc("/debugdocker", DebugDocker)
	router.HandleFunc("/debugkubernetes", DebugKubernetes)
	router.HandleFunc("/composeupcoredod", ComposeUpCoreDOD)
}

// Setup makes the disclaimer and installs dotfolder and clones/updates the repo
func Setup(w http.ResponseWriter, r *http.Request) {
	// makeDisclaimer()
	// http.ServeFile(w, r, "/templates/index.html")
	setup()
}

// Disclaimer runs makeDisclaimer
func Disclaimer(w http.ResponseWriter, r *http.Request) {
	makeDisclaimer()
}

// Decline just displays the decline page
func Decline(w http.ResponseWriter, r *http.Request) {
	// makeDisclaimer()
	// setup()
}

// DebugDocker checks that docker is running
func DebugDocker(w http.ResponseWriter, r *http.Request) {
	debugDocker()
	// setup()
}

// DebugKubernetes checks that kubectl has a config
func DebugKubernetes(w http.ResponseWriter, r *http.Request) {
	debugKubernetes()
	// setup()
}

// ComposeUpCoreDOD checks that kubectl has a config
func ComposeUpCoreDOD(w http.ResponseWriter, r *http.Request) {
	consoleSender(server, "Note: Initial setup takes 1-5 minutes. wait for the DONE message")
	SomeStuff()
	// setup()
}
