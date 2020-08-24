package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func addHandler(router *mux.Router) {
	// swagger:operation GET /events admin repoList
	// ---
	// summary: Stream of console events and stdout
	// description: Implements unidirectional streaming of console and stdout events for debugging.
	router.HandleFunc("/events", sseHandler)
	// swagger:operation GET /setup admin Setup
	// ---
	// summary: Clones or updates local Instant repository.
	// description: Clones or updates Instant GitHub repo in `$HOME/.instant` Does not require git to be installed.
	router.HandleFunc("/setup", Setup)
	// swagger:operation GET /disclaimer admin Disclaimer
	// ---
	// summary: Adds a disclaimer to the local repository.
	// description: Adds accept_disclaimer file in `$HOME/.instant/` Avoids visiting the disclaimer page.
	router.HandleFunc("/disclaimer", Disclaimer)
	// decline has not rest-api endpoint
	router.HandleFunc("/decline", Decline)
	// swagger:operation GET /debugdocker docker debugDocker
	// ---
	// summary: Confirms Docker is running.
	// description: Runs equivalent of a `docker ps` command to ensure Docker is running and installed locally.
	router.HandleFunc("/debugdocker", DebugDocker)
	// swagger:operation GET /debugkubernetes kubernetes DebugKubernetes
	// ---
	// summary: Confirms Kubernetes is running.
	// description: Runs equivalent of a `kubectl config view` command to ensure kubectl is installed and configured.
	router.HandleFunc("/debugkubernetes", DebugKubernetes)
	// swagger:operation GET /composeupcoredod docker ComposeUpCoreDOD
	// ---
	// summary: Launches docker DooD container and runs core package.
	// description: Runs `docker run /var/run/docker.sock:/var/run/docker.sock openhie/instant:latest init -t docker core`
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
