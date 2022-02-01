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
	logs           string
	directoryNames []string
)

func theCommandIsRun(command string) error {
	res, err := runTestCommand(binaryFilePath, strings.Split(command, " ")...)
	if err == nil {
		logs = res
	}
	return nil
}

func checkTheCLIOutputIs(expectedOutput string) error {
	err := compareStrings(logs, expectedOutput)
	if err != nil {
		return err
	}
	return nil
}

func checkCustomPackages(packages *godog.Table) error {
	head := packages.Rows[0].Cells

	for i := 1; i < len(packages.Rows); i++ {

		for n, cell := range packages.Rows[i].Cells {
			switch head[n].Value {
			case "directory":
				directoryNames = append(directoryNames, cell.Value)
			case "location":
				err := compareStrings(logs, cell.Value)
				if err != nil {
					return err
				}
			default:
				return errors.New("Unexpected column name: " + head[n].Value)
			}
		}
	}

	// Cleanup downloaded packages post test
	deleteContentAtFilePath([]string{".", "features"}, directoryNames)
	return nil
}

func compareStrings(inputLogs, expectedOutput string) error {
	if !strings.Contains(inputLogs, expectedOutput) {
		return errors.New("Logs received: '" + inputLogs + "\nSubstring expected: " + expectedOutput)
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	suite := &godog.TestSuite{
		TestSuiteInitializer: func(s *godog.TestSuiteContext) {
			s.AfterSuite(cleanUp)
		},
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			if binaryFilePath == "" {
				binaryFilePath = buildBinary()
			}
			sc.Step(`^check the CLI output is "([^"]*)"$`, checkTheCLIOutputIs)
			sc.Step(`^the command "([^"]*)" is run$`, theCommandIsRun)
			sc.Step(`^check that the CLI added custom packages$`, checkCustomPackages)
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

func deleteContentAtFilePath(filePath []string, content []string) {
	for _, c := range content {
		err := os.RemoveAll(filepath.Join(filepath.Join(filePath...), c))
		if err != nil {
			panic(err)
		}
	}
}

func cleanUp() {
	deleteContentAtFilePath([]string{".", "features"}, []string{"test-platform.exe", "test-platform-linux", "test-platform-macos"})
	deleteContentAtFilePath([]string{"."}, directoryNames)

	_, errContainer := runTestCommand("docker", "rm", "instant-openhie")
	if errContainer == nil {
		fmt.Println("Deleted Instant OpenHIE container")
	}

	_, errVolume := runTestCommand("docker", "volume", "rm", "instant")
	if errVolume != nil {
		fmt.Println("Instant Docker volume not deleted")
	}
}
