package main

import (
	"embed"
	"log"
	"os"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v2"
)

//go:embed banner.txt
var f embed.FS

//go:embed config.yml
var yamlConfig []byte
var cfg Config

type Package struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type Config struct {
	Image              string    `yaml:"image"`
	DefaultEnvironment string    `yaml:"defaultEnvironment"`
	Packages           []Package `yaml:"packages"`
	DisableKubernetes  bool      `yaml:"disableKubernetes"`
}

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

func stopContainer() {
	commandSlice := []string{"stop", "instant-openhie"}
	suppressErrors := []string{"Error response from daemon: No such container: instant-openhie"}
	_, err := runCommand("docker", suppressErrors, commandSlice...)
	if err != nil {
		gracefulPanic(err, "")
	}
}

//Gracefully shut down the instant container and then kill the go cli with the panic error or message passed.
func gracefulPanic(err error, message string) {
	stopContainer()
	if message != "" {
		panic(message)
	}
	panic(err)
}

func loadConfig() {
	err := yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadConfig()

	data, err := f.ReadFile("banner.txt")
	if err != nil {
		log.Println(err)
	}
	color.Green(string(data))

	color.Cyan("Version: 1.02b")
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	// mainMenu()
	if len(os.Args) > 1 {
		err = CLI()
		if err != nil {
			gracefulPanic(err, "")
		}
	} else {
		err = selectSetup()
		if err != nil {
			gracefulPanic(err, "")
		}
	}
}
