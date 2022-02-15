package main

import (
	"embed"
	"log"
	"os"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v3"
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
	Image                        string    `yaml:"image"`
	DefaultTargetLauncher        string    `yaml:"defaultTargetLauncher"`
	Packages                     []Package `yaml:"packages"`
	DisableKubernetes            bool      `yaml:"disableKubernetes"`
	DisableIG                    bool      `yaml:"disableIG"`
	DisableCustomTargetSelection bool      `yaml:"disableCustomTargetSelection"`
}

type customOption struct {
	startupAction              string
	startupPackages            []string
	envVarFileLocation         string
	envVars                    []string
	customPackageFileLocations []string
	onlyFlag                   bool
	instantVersion             string
	targetLauncher             string
	devMode                    bool
}

var customOptions = customOption{
	startupAction:      "init",
	envVarFileLocation: "",
	onlyFlag:           false,
	instantVersion:     "latest",
	targetLauncher:     "docker",
	devMode:            false,
}

func stopContainer() {
	commandSlice := []string{"stop", "instant-openhie"}
	suppressErrors := []string{"Error response from daemon: No such container: instant-openhie"}
	_, err := runCommand("docker", suppressErrors, commandSlice...)
	if err != nil {
		log.Fatalf("runCommand() failed: %v", err)
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

func getHelpText(interactive bool) string {
	if interactive {
		return ``
	} else {
		return `Commands: 
		help 		this menu
		install		install fhir npm package on fhir server
					usage: install <ig_url> <fhir_server> <authtype> <user/token> <pass>

					examples:
					install https://intrahealth.github.io/simple-hiv-ig/ http://hapi.fhir.org/baseR4 none
					install <ig_url> <fhir_server> basic smith stuff
					install <ig_url> <fhir_server> token "123"
					install <ig_url> <fhir_server> custom test
		init/up/destroy/down	the deploy command you want to run (brief description below)
					deploy commands:
						init:
						up:
						destroy:
						down:
					custom flags:
						only:
						-t:
						-c:
						--dev:
						-e:
						--env-file:
					usage:
						<deploy command> <custom flags> <package names>
					examples:
						init -t=swarm --dev interoperability-layer-openhim
		`
	}
}

func main() {
	loadConfig()

	//Need to set the default here as we declare the struct before the config is loaded in.
	customOptions.targetLauncher = cfg.DefaultTargetLauncher

	data, err := f.ReadFile("banner.txt")
	if err != nil {
		log.Println(err)
	}
	color.Green(string(data))

	color.Cyan("Version: 1.02b")
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

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
