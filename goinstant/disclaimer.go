package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
)

// existDisclaimer detects if the accept_disclaimer file exists in the $HOME/.instant folder.
// If it exists, then the user has accepted and it returns 'success', otherwise 'fail'.
func existDisclaimer() string {
	home, _ := os.UserHomeDir()
	// must use filepath.join not path.join for windows compat
	dotfiles := filepath.Join(home, ".instant")
	fileName := filepath.Join(dotfiles, "accept_disclaimer")
	status := "success"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		status = "fail"
	}
	return status
}

func cliDisclaimer() {
	status := existDisclaimer()
	if status == "fail" {
		pkgerPrint("/templates/disclaimer.txt", "yellow")
		prompt := promptui.Select{
			Label: "Do you agree to use this application?",
			Items: []string{"Yes", "No", "Quit"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Printf("You chose %q\n", result)

		switch result {
		case "Yes":
			makeDisclaimer()
			setup()
			selectSetup()
		case "No":
			fmt.Println("Understood. Exiting.")
			os.Exit(1)
		case "Quit":
			os.Exit(1)
		}
	} else {
		selectSetup()
	}
}

func makeDisclaimer() string {
	home, _ := os.UserHomeDir()
	dotfiles := filepath.Join(home, ".instant")
	fileName := filepath.Join(dotfiles, "accept_disclaimer")
	_, err := os.Stat(fileName)
	status := "success"
	if os.IsNotExist(err) {
		err := os.MkdirAll(dotfiles, os.ModePerm)
		if err != nil {
			status = "fail"
			log.Fatal(err)
		}
		file, err := os.Create(fileName)
		if err != nil {
			status = "fail"
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(fileName, currentTime, currentTime)
		if err != nil {
			status = "fail"
			fmt.Println(err)
		}
	}
	return status
}
