package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
)

func makeSetup() string {

	home, _ := os.UserHomeDir()
	dotfiles := filepath.Join(home, ".instant")

	status := "success"

	// check for dotfolder, create if it doesn't exist
	if _, err := os.Stat(dotfiles); os.IsNotExist(err) {
		os.Mkdir(dotfiles, 0700)
	}

	// check repo and clone or pull
	repo := filepath.Join(dotfiles, "instant")
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		_, err := git.PlainClone(repo, false, &git.CloneOptions{
			URL:      "https://github.com/openhie/instant",
			Progress: os.Stdout,
		})
		if err != nil {
			status = "fail"
		}
	} else {
		const (
			repoURL = "https://github.com/openhie/instant.git"
		)

		dir, _ := ioutil.TempDir("", "temp_dir")

		options := &git.CloneOptions{
			URL: repoURL,
		}

		_, err := git.PlainClone(dir, false, options)
		if err != nil {
			status = "fail"
		}

	}
	return status
}

func setup() {

	color.Green.Println("getting ready")

	home, _ := os.UserHomeDir()
	dotfiles := filepath.Join(home, ".instant")
	fmt.Println("...checking for config folder at:", dotfiles)

	// check for dotfolder, create if it doesn't exist
	if _, err := os.Stat(dotfiles); os.IsNotExist(err) {
		fmt.Println("config folder does not exist, creating it")
		os.Mkdir(dotfiles, 0700)
		color.Green.Println("created", dotfiles)
	} else {
		color.Green.Println("config folder exists")
	}

	// check repo and clone or pull
	fmt.Println("...cloning or pulling latest code from repo")

	repo := filepath.Join(dotfiles, "instant")
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		fmt.Println("...repo folder does not exist, cloning it")
		_, err := git.PlainClone(repo, false, &git.CloneOptions{
			URL:      "https://github.com/openhie/instant",
			Progress: os.Stdout,
		})
		if err != nil {
		}
		color.Green.Println("successfully cloned", dotfiles)
	} else {
		fmt.Println("...repo folder exists, pulling changes")
		const (
			repoURL = "https://github.com/openhie/instant.git"
		)

		dir, _ := ioutil.TempDir("", "temp_dir")

		options := &git.CloneOptions{
			URL: repoURL,
		}

		_, err := git.PlainClone(dir, false, options)
		if err != nil {
		}

	}
	color.Green.Println("git repo is ready")
	color.Green.Println("ready!")

}
