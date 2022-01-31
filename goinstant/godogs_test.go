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
	binaryFilePath string
	cmdResult      string
)

func theCommandIsRun(command string) error {
	res, err := runTestCommand(binaryFilePath, strings.Split(command, " ")...)
	if err == nil {
		cmdResult = res
	}
	return nil
}

func checkTheCLIOutputIs(command string) error {
	err := theLoggedStringsMatch(cmdResult, command)
	if err != nil {
		return err
	}
	return nil
}

func theLoggedStringsMatch(str, strToMatch string) error {
	if !strings.Contains(str, strToMatch) {
		return errors.New("String '" + strToMatch + "' not matched")
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	suite := &godog.TestSuite{
		TestSuiteInitializer: func(s *godog.TestSuiteContext) {
			s.AfterSuite(cleanBinaries)
		},
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			if binaryFilePath == "" {
				binaryFilePath = buildBinary()
			}
			sc.Step(`^check the CLI output is "([^"]*)"$`, checkTheCLIOutputIs)
			sc.Step(`^the command "([^"]*)" is run$`, theCommandIsRun)
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

func cleanBinaries() {
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
