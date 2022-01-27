package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

var binaryFilePath string

func theOpenHIMServiceIsInitialised() error {
	_, err := runCommand(binaryFilePath, nil, "init", "-t=docker", "core")
	return err
}

func theOpenHIMServiceIsBroughtDown() error {
	_, err := runCommand(binaryFilePath, nil, "down", "-t=docker", "core")
	return err
}

func theOpenHIMServiceIsBroughtUp() error {
	_, err := runCommand(binaryFilePath, nil, "up", "-t=docker", "core")
	return err
}

func theOpenHIMServiceIsDestroyed() error {
	_, err := runCommand(binaryFilePath, nil, "destroy", "-t=docker", "core")
	if err != nil {
		return err
	}

	fileList := []string{"test-platform.exe", "test-platform-linux", "test-platform-macos"}
	for _, f := range fileList {
		err = os.Remove(filepath.Join(".", "bin", f))
		if err != nil {
			return err
		}
	}

	return nil
}

func theServiceShouldBeReachable() error {
	if !CheckOpenHimStatus() {
		return errors.New("The service is not running")
	}
	return nil
}

func theServiceShouldNotBeReachable() error {
	if CheckOpenHimStatus() {
		return errors.New("The service is running")
	}
	return nil
}

func CheckOpenHimStatus() bool {
	resp, err := http.Get("http://localhost:9000")
	if resp == nil || resp.StatusCode != 200 {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return true
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	if binaryFilePath == "" {
		binaryFilePath = buildBinary()
	}

	ctx.Step(`^the OpenHIM service is brought down$`, theOpenHIMServiceIsBroughtDown)
	ctx.Step(`^the OpenHIM service is brought up$`, theOpenHIMServiceIsBroughtUp)
	ctx.Step(`^the OpenHIM service is destroyed$`, theOpenHIMServiceIsDestroyed)
	ctx.Step(`^the OpenHIM service is initialised$`, theOpenHIMServiceIsInitialised)
	ctx.Step(`^the OpenHIM service is initialised and running$`, theServiceShouldBeReachable)
	ctx.Step(`^the OpenHIM service is initialised but not running$`, theServiceShouldNotBeReachable)
	ctx.Step(`^the OpenHIM service is not instantiated$`, theServiceShouldNotBeReachable)
	ctx.Step(`^the service should be reachable$`, theServiceShouldBeReachable)
	ctx.Step(`^the service should not be reachable$`, theServiceShouldNotBeReachable)
}

func buildBinary() string {
	loggedResults, err := runTestCommand("/bin/sh", filepath.Join(".", "features", "build-cli.sh"))
	if err != nil {
		panic(err)
	}
	fmt.Println(loggedResults)
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var loggedResults string
	scanner := bufio.NewScanner(cmdReader)
	go func() {
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

	return loggedResults, nil
}
