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
	router.HandleFunc("/onefunc", OneFunc)
	router.HandleFunc("/listdocker", ListDocker)
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
}

// DebugDocker checks that docker is running
func DebugDocker(w http.ResponseWriter, r *http.Request) {
	debugDocker()
}

// DebugKubernetes checks that kubectl has a config
func DebugKubernetes(w http.ResponseWriter, r *http.Request) {
	debugKubernetes()
	// setup()
}

// OneFunc runs, stops, starts, destroys everything
func OneFunc(w http.ResponseWriter, r *http.Request) {
	// composeUpCoreDOD()
	SomeStuff(r)
}

// ListDocker gets list of running containers
func ListDocker(w http.ResponseWriter, r *http.Request) {
	listDocker()
}
