package main

import (
	"net/http"
)

// Setup makes the disclaimer and installs dotfolder and clones/updates the repo
func Setup(w http.ResponseWriter, r *http.Request) {
	// makeDisclaimer()
	// http.ServeFile(w, r, "/templates/index.html")
	setup()
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
	SomeStuff()
	// setup()
}
