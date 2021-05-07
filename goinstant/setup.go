package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func setup() {

	consoleSender(server, "getting ready")

	home, _ := os.UserHomeDir()
	dotfiles := filepath.Join(home, ".instant")
	consoleSender(server, "...checking for config folder")

	// check for dotfolder, create if it doesn't exist
	if _, err := os.Stat(dotfiles); os.IsNotExist(err) {
		consoleSender(server, "config folder does not exist, creating it")
		os.Mkdir(dotfiles, 0700)
		consoleSender(server, "created")
	} else {
		consoleSender(server, "config folder exists")
	}

	// check repo and clone or pull
	consoleSender(server, "...cloning or pulling latest code from repo")

	repo := filepath.Join(dotfiles, "instant")
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		consoleSender(server, "...repo folder does not exist, cloning it")
		_, err := git.PlainClone(repo, false, &git.CloneOptions{
			URL:      "https://github.com/openhie/instant",
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Println("error")
		}
		consoleSender(server, "successfully cloned")
	} else {
		consoleSender(server, "...repo folder exists, pulling changes")
		const (
			repoURL = "https://github.com/openhie/instant.git"
		)

		dir, _ := ioutil.TempDir("", "temp_dir")

		options := &git.CloneOptions{
			URL: repoURL,
		}

		_, err := git.PlainClone(dir, false, options)
		if err != nil {
			fmt.Println("error")
		}
	}
	consoleSender(server, "git repo is ready")
}
