package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

var (
	binaryFilePath           string
	output                   string
	serviceInitialisedResult string
	serviceBroughtDownResult string
	serviceBroughtUpResult   string
	serviceDestroyedResult   string
	configOptionsResult      string
)

func theServiceIsInitialised() error {
	result, err := runTestCommand(binaryFilePath, "init", "-t=docker", "core")
	if err == nil {
		serviceInitialisedResult = result
	}
	return err
}

func checkTheServiceIsInitialised() error {
	return theLoggedStringsMatch(serviceInitialisedResult, "init -t docker core")
}

func theServiceIsBroughtUp() error {
	result, err := runTestCommand(binaryFilePath, "up", "-t=k8s", "core")
	if err == nil {
		serviceBroughtUpResult = result
	}
	return err
}

func checkTheServiceIsBroughtUp() error {
	return theLoggedStringsMatch(serviceBroughtUpResult, "up -t k8s core")
}

func theServiceIsBroughtDown() error {
	result, err := runTestCommand(binaryFilePath, "down", "-t=swarm", "opencr")
	if err == nil {
		serviceBroughtDownResult = result
	}
	return err
}

func checkTheServiceIsBroughtDown() error {
	return theLoggedStringsMatch(serviceBroughtDownResult, "down -t swarm opencr")
}

func theServiceIsDestroyed() error {
	result, err := runTestCommand(binaryFilePath, "destroy", "-t=k8s", "core")
	if err == nil {
		serviceDestroyedResult = result
	}
	return err
}

func checkTheServiceIsDestroyed() error {
	return theLoggedStringsMatch(serviceDestroyedResult, "destroy -t k8s core")
}

func theServiceConfigOptionsArePassed() error {
	result, err := runTestCommand(binaryFilePath, "init", "-t=docker", "-c=./features", "custom_package", "-e=NODE_ENV=DEV", "--onlyFlag", "--dev")
	if err == nil {
		configOptionsResult = result
	}
	return err
}

func theLoggedStringsMatch(str, strToMatch string) error {
	if !strings.Contains(str, strToMatch) {
		return errors.New("String '" + strToMatch + "' not matched")
	}
	return nil
}

func checkTheServiceConfigOptionsArePassed() error {
	err := theLoggedStringsMatch(configOptionsResult, "init --onlyFlag --dev -t docker custom_package")
	if err != nil {
		return err
	}
	return theLoggedStringsMatch(configOptionsResult, "NODE_ENV=DEV")
}

func InitializeScenario(sc *godog.ScenarioContext) {
	suite := &godog.TestSuite{
		TestSuiteInitializer: func(s *godog.TestSuiteContext) {
			s.AfterSuite(clean)
		},
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			if binaryFilePath == "" {
				binaryFilePath = buildBinary()
			}

			sc.Step(`^the service is initialised$`, theServiceIsInitialised)
			sc.Step(`^the service is brought up$`, theServiceIsBroughtUp)
			sc.Step(`^the service is brought down$`, theServiceIsBroughtDown)
			sc.Step(`^the service is destroyed$`, theServiceIsDestroyed)
			sc.Step(`^the service config options are passed$`, theServiceConfigOptionsArePassed)

			sc.Step(`^check the service is initialised$`, checkTheServiceIsInitialised)
			sc.Step(`^check the service is brought up$`, checkTheServiceIsBroughtUp)
			sc.Step(`^check the service is brought down$`, checkTheServiceIsBroughtDown)
			sc.Step(`^check the service is destroyed$`, checkTheServiceIsDestroyed)
			sc.Step(`^check the service config options are passed$`, checkTheServiceConfigOptionsArePassed)
		},
	}

	if suite.Run() != 0 {
		fmt.Println("Tests failed")
		os.Exit(1)
	}
	os.Exit(0)
}

func buildBinary() string {
	_, err := runTestCommand("/bin/sh", filepath.Join(".", "features", "build-cli.sh"))
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(filepath.Join(".", "features", "test-platform-linux"))
	if err != nil {
		panic(err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(".", "features", "test-platform.exe")
	case "ios":
		return filepath.Join(".", "features", "test-platform-macos")
	case "linux":
		return filepath.Join(".", "features", "test-platform-linux")
	default:
		panic(errors.New("Operating system not supported"))
	}
}

func runTestCommand(commandName string, commandSlice ...string) (string, error) {
	cmd := exec.Command(commandName, commandSlice...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer cmdReader.Close()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var wg sync.WaitGroup
	wg.Add(1)

	var loggedResults string
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			loggedResults += scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Error starting Cmd. "+stderr.String())

	}
	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "Error waiting for Cmd. "+stderr.String())
	}

	wg.Wait()
	return loggedResults, nil
}

func clean() {
	fileList := []string{"test-platform.exe", "test-platform-linux", "test-platform-macos"}
	for _, f := range fileList {
		err := os.Remove(filepath.Join(".", "features", f))
		if err != nil {
			panic(err)
		}
	}

	_, err := runTestCommand("docker", "volume", "rm", "instant")
	if err != nil {
		fmt.Println("Volume not deleted")
	}
}
