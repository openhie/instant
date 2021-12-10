package main

import (
	"embed"
	"os"

	"github.com/fatih/color"
)

//go:embed banner.txt
var f embed.FS

type customOption struct {
	startupAction              string
	startupPackages            []string
	envVarFileLocation         string
	envVars                    []string
	customPackageFileLocations []string
	onlyFlag                   bool
	instantVersion             string
}

var customOptions = customOption{
	startupAction:      "init",
	envVarFileLocation: "",
	onlyFlag:           false,
	instantVersion:     "latest",
}

func main() {
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
