package main

import "net/http"

// Index makes the disclaimer and installs dotfolder and clones/updates the repo
func Index(w http.ResponseWriter, r *http.Request) {
	makeDisclaimer()
	setup()
}

// Decline just displays the decline page
func Decline(w http.ResponseWriter, r *http.Request) {
	makeDisclaimer()
	setup()
}
