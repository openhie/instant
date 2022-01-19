package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func debugDocker() error {
	fmt.Printf("...checking your Docker setup")

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Can't get current working directory")
	}

	fmt.Println(cwd)

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to get Docker context. Please ensure that Docker is downloaded and running")
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		str1 := "bytes memory is allocated\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		fmt.Println(result)
		fmt.Println("Docker setup looks good")
	}

	return nil
}

func getPackagePaths(inputArr []string, flags []string) (packagePaths []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				packagePath := strings.Replace(i, flag, "", 1)
				packagePath = strings.Trim(packagePath, "\"")
				packagePaths = append(packagePaths, packagePath)
			}
		}
	}
	return
}

func getEnvironmentVariables(inputArr []string, flags []string) (environmentVariables []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				environmentVariables = append(environmentVariables, strings.SplitN(i, "=", 2)...)
			}
		}
	}
	return
}

func sliceContains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func extractCommands(startupCommands []string) (environmentVariables []string, deployCommand string, otherFlags []string, deployEnvironment string, packages []string, customPackagePaths []string, instantVersion string) {
	instantVersion = "latest"

	for _, option := range startupCommands {
		switch {
		case sliceContains([]string{"init", "up", "down", "destroy"}, option):
			deployCommand = option
		case sliceContains([]string{"docker", "kubernetes", "k8s", "swarm"}, option):
			deployEnvironment = option
		case strings.HasPrefix(option, "-c=") || strings.HasPrefix(option, "--custom-package="):
			customPackagePaths = append(customPackagePaths, option)
		case strings.HasPrefix(option, "-e=") || strings.HasPrefix(option, "--env-file="):
			environmentVariables = append(environmentVariables, option)
		case strings.HasPrefix(option, "--instant-version="):
			instantVersion = strings.Split(option, "--instant-version=")[1]
		case strings.HasPrefix(option, "-") || strings.HasPrefix(option, "--"):
			otherFlags = append(otherFlags, option)
		default:
			packages = append(packages, option)
		}
	}

	if len(customPackagePaths) > 0 {
		customPackagePaths = getPackagePaths(customPackagePaths, []string{"-c=", "--custom-package="})
	}

	if len(environmentVariables) > 0 {
		environmentVariables = getEnvironmentVariables(environmentVariables, []string{"-e=", "--env-file="})
	}
	return
}

func RunDirectDockerCommand(startupCommands []string) error {
	fmt.Println("Note: Initial setup takes 1-5 minutes.\nWait for the DONE message.\n--------------------------")

	environmentVariables, deployCommand, otherFlags, deployEnvironment, packages, customPackagePaths, instantVersion := extractCommands(startupCommands)

	fmt.Println("Environment:", deployEnvironment)
	fmt.Println("Action:", deployCommand)
	fmt.Println("Package IDs:", packages)
	fmt.Println("Custom package paths:", customPackagePaths)
	fmt.Println("Environment Variables:", environmentVariables)
	fmt.Println("Other Flags:", otherFlags)
	fmt.Println("InstantVersion:", instantVersion)

	instantImage := cfg.Image + ":" + instantVersion

	var err error
	if deployCommand == "init" {
		fmt.Println("\n\nDelete a pre-existing instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		_, err = runCommand("docker", []string{"Error: No such volume: instant"}, commandSlice...)
		if err != nil {
			return err
		}
	}

	fmt.Println("Creating fresh instant container with volumes...")
	commandSlice := []string{
		"create",
		"--rm",
		"--mount=type=volume,src=instant,dst=/instant",
		"--name", "instant-openhie",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"--network", "host",
	}
	commandSlice = append(commandSlice, environmentVariables...)
	commandSlice = append(commandSlice, []string{instantImage, deployCommand}...)
	commandSlice = append(commandSlice, otherFlags...)
	commandSlice = append(commandSlice, []string{"-t", deployEnvironment}...)
	commandSlice = append(commandSlice, packages...)
	_, err = runCommand("docker", nil, commandSlice...)
	if err != nil {
		return err
	}

	fmt.Println("Adding 3rd party packages to instant volume:")

	for _, c := range customPackagePaths {
		fmt.Print("- " + c)
		err = mountCustomPackage(c)
		if err != nil {
			return err
		}
	}

	fmt.Println("\nRun Instant OpenHIE Installer Container")
	commandSlice = []string{"start", "-a", "instant-openhie"}
	_, err = runCommand("docker", nil, commandSlice...)
	if err != nil {
		return err
	}

	if deployCommand == "destroy" {
		fmt.Println("Delete instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		_, err = runCommand("docker", nil, commandSlice...)
	}

	return err
}

func runCommand(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
	cmd := exec.Command(commandName, commandSlice...)
	cmdReader, err := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err != nil {
		if suppressErrors != nil && sliceContains(suppressErrors, strings.TrimSpace(stderr.String())) {
			return pathToPackage, nil
		}

		return pathToPackage, errors.Wrap(err, "Error creating StdoutPipe for Cmd.")
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		if suppressErrors != nil && sliceContains(suppressErrors, strings.TrimSpace(stderr.String())) {

		} else {
			return pathToPackage, errors.Wrap(err, "Error starting Cmd. "+stderr.String())
		}
	}

	err = cmd.Wait()
	if err != nil {
		if suppressErrors != nil && sliceContains(suppressErrors, strings.TrimSpace(stderr.String())) {

		} else {
			return pathToPackage, errors.Wrap(err, "Error waiting for Cmd. "+stderr.String())
		}
	}

	if commandName == "git" {
		if len(commandSlice) < 2 {
			return pathToPackage, errors.New("Not enough arguments for git command")
		}
		pathToPackage = commandSlice[1]
		// Get name of repo
		urlSplit := strings.Split(pathToPackage, ".")
		urlPathSplit := strings.Split(urlSplit[len(urlSplit)-2], "/")
		repoName := urlPathSplit[len(urlPathSplit)-1]

		pathToPackage = filepath.Join(".", repoName)
	}

	return pathToPackage, nil
}

func mountCustomPackage(pathToPackage string) error {
	gitRegex := regexp.MustCompile(`\.git`)
	httpRegex := regexp.MustCompile("http")
	zipRegex := regexp.MustCompile(`\.zip`)
	tarRegex := regexp.MustCompile(`\.tar`)

	var err error
	if gitRegex.MatchString(pathToPackage) {
		pathToPackage, err = runCommand("git", nil, []string{"clone", pathToPackage}...)
		if err != nil {
			return err
		}
	} else if httpRegex.MatchString(pathToPackage) {
		resp, err := http.Get(pathToPackage)
		if err != nil {
			return errors.Wrap(err, "Error in downloading custom package")
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.Wrapf(err, "Error in downloading custom package - HTTP status code: %v", strconv.Itoa(resp.StatusCode))
		}

		if zipRegex.MatchString(pathToPackage) {
			pathToPackage, err = unzipPackage(resp.Body)
			if err != nil {
				return err
			}
		} else if tarRegex.MatchString(pathToPackage) {
			pathToPackage, err = untarPackage(resp.Body)
			if err != nil {
				return err
			}
		}
	}

	commandSlice := []string{"cp", pathToPackage, "instant-openhie:instant/"}
	_, err = runCommand("docker", nil, commandSlice...)
	return err
}

func createZipFile(file string, content io.Reader) error {
	output, err := os.Create(file)
	if err != nil {
		return errors.Wrap(err, "Error in creating zip file:")
	}
	defer output.Close()

	_, err = io.Copy(output, content)
	if err != nil {
		return errors.Wrap(err, "Error in copying zip file content:")
	}

	return nil
}

func unzipPackage(zipContent io.ReadCloser) (pathToPackage string, err error) {
	tempZipFile := "temp.zip"
	err = createZipFile(tempZipFile, zipContent)
	if err != nil {
		return "", err
	}

	// Unzip file
	archive, err := zip.OpenReader(tempZipFile)
	if err != nil {
		return "", errors.Wrap(err, "Error in unzipping file:")
	}

	packageName := ""
	for _, file := range archive.File {
		filePath := filepath.Join(".", file.Name)

		if file.FileInfo().IsDir() {
			if packageName == "" {
				packageName = file.Name
			}
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		content, err := file.Open()
		if err != nil {
			return "", errors.Wrap(err, "Error in unzipping file:")
		}

		dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", errors.Wrap(err, "Error in unzipping file:")
		}
		_, err = io.Copy(dest, content)
		if err != nil {
			return "", errors.Wrap(err, "Error in copying unzipping file:")
		}
		content.Close()
	}

	// Remove temp zip file
	tempFilePath := filepath.Join(".", tempZipFile)
	archive.Close()
	err = os.Remove(tempFilePath)
	if err != nil {
		return "", errors.Wrap(err, "Error in deleting temp.zip file:")
	}

	return filepath.Join(".", packageName), nil
}

func untarPackage(tarContent io.ReadCloser) (pathToPackage string, err error) {
	packageName := ""
	gzipReader, err := gzip.NewReader(tarContent)
	if err != nil {
		return "", errors.Wrap(err, "Error in extracting tar file")
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if file == nil {
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "Error in extracting tar file")
		}

		filePath := filepath.Join(".", file.Name)
		if file.Typeflag == tar.TypeDir {
			if packageName == "" {
				packageName = filePath
			}
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return "", errors.Wrap(err, "Error in untaring file")
		}

		_, err = io.Copy(dest, tarReader)
		if err != nil {
			return "", errors.Wrap(err, "Error in extracting tar file")
		}
	}
	pathToPackage = filepath.Join(".", packageName)

	return pathToPackage, nil
}
