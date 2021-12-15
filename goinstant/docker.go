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
	"golang.org/x/net/context"
)

func debugDocker() {

	fmt.Printf("...checking your Docker setup")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Can't get current working directory... this is not a great error.")
		// panic(err)
	} else {
		fmt.Println(cwd)
	}

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Println("Unable to get Docker context. Please ensure that Docker is downloaded and running")
		panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		str1 := "bytes memory is allocated\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		fmt.Println(result)
		fmt.Println("Docker setup looks good")
	}

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
		case sliceContains([]string{"docker", "kubernetes", "k8s"}, option):
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

func RunDirectDockerCommand(startupCommands []string) {
	fmt.Println("Note: Initial setup takes 1-5 minutes.\nWait for the DONE message.\n--------------------------")

	environmentVariables, deployCommand, otherFlags, deployEnvironment, packages, customPackagePaths, instantVersion := extractCommands(startupCommands)

	fmt.Println("Environment:", deployEnvironment)
	fmt.Println("Action:", deployCommand)
	fmt.Println("Package IDs:", packages)
	fmt.Println("Custom package paths:", customPackagePaths)
	fmt.Println("Environment Variables:", environmentVariables)
	fmt.Println("Other Flags:", otherFlags)
	fmt.Println("InstantVersion:", instantVersion)

	instantImage := "openhie/instant:" + instantVersion

	if deployCommand == "init" {
		fmt.Println("\n\nDelete a pre-existing instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		runCommand(deployEnvironment, commandSlice...)
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
	runCommand(deployEnvironment, commandSlice...)

	fmt.Println("Adding 3rd party packages to instant volume:")

	for _, c := range customPackagePaths {
		fmt.Print("- " + c)
		mountCustomPackage(deployEnvironment, c)
	}

	fmt.Println("\nRun Instant OpenHIE Installer Container")
	commandSlice = []string{"start", "-a", "instant-openhie"}
	runCommand(deployEnvironment, commandSlice...)

	if deployCommand == "destroy" {
		fmt.Println("Delete instant volume...")
		commandSlice := []string{"volume", "rm", "instant"}
		runCommand(deployEnvironment, commandSlice...)
	}
}

func runCommand(commandName string, commandSlice ...string) (pathToPackage string) {
	cmd := exec.Command(commandName, commandSlice...)
	cmdReader, err := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd.", stderr.String(), err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd.", stderr.String(), err)
		return
	}

	if commandName == "git" {
		if len(commandSlice) < 2 {
			fmt.Fprintln(os.Stderr, "Not enough arguments for git command", stderr.String(), err)
			return
		}
		pathToPackage = commandSlice[1]
		// Get name of repo
		urlSplit := strings.Split(pathToPackage, ".")
		urlPathSplit := strings.Split(urlSplit[len(urlSplit)-2], "/")
		repoName := urlPathSplit[len(urlPathSplit)-1]

		pathToPackage = filepath.Join(".", repoName)
	}
	return
}

func mountCustomPackage(deployEnvironment string, pathToPackage string) {
	gitRegex := regexp.MustCompile(`\.git`)
	httpRegex := regexp.MustCompile("http")
	zipRegex := regexp.MustCompile(`\.zip`)
	tarRegex := regexp.MustCompile(`\.tar`)

	if gitRegex.MatchString(pathToPackage) {
		pathToPackage = runCommand("git", []string{"clone", pathToPackage}...)
	} else if httpRegex.MatchString(pathToPackage) {
		resp, err := http.Get(pathToPackage)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error in dowloading custom package", err)
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			panic("Error in dowloading custom package - HTTP status code: " + strconv.Itoa(resp.StatusCode))
		}

		if zipRegex.MatchString(pathToPackage) {
			pathToPackage = unzipPackage(resp.Body)
		} else if tarRegex.MatchString(pathToPackage) {
			pathToPackage = untarPackage(resp.Body)
		}
	}

	commandSlice := []string{"cp", pathToPackage, "instant-openhie:instant/"}
	runCommand(deployEnvironment, commandSlice...)
}

func createZipFile(file string, content io.Reader) {
	output, err := os.Create(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in creating zip file:")
		panic(err)
	}
	defer output.Close()

	_, err = io.Copy(output, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in copying zip file content:")
		panic(err)
	}
}

func unzipPackage(zipContent io.ReadCloser) (pathToPackage string) {
	tempZipFile := "temp.zip"
	createZipFile(tempZipFile, zipContent)

	// Unzip file
	archive, err := zip.OpenReader(tempZipFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in unzipping file:")
		panic(err)
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
			fmt.Fprintln(os.Stderr, "Error in unzipping file:")
			panic(err)
		}

		dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error in unzipping file:")
			panic(err)
		}
		_, err = io.Copy(dest, content)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error in copying unzipped files:")
			panic(err)
		}
		content.Close()
	}

	// Remove temp zip file
	tempFilePath := filepath.Join(".", tempZipFile)
	archive.Close()
	err = os.Remove(tempFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in deleting temp.zip file:")
		panic(err)
	}

	pathToPackage = filepath.Join(".", packageName)
	return
}

func untarPackage(tarContent io.ReadCloser) (pathToPackage string) {
	packageName := ""
	gzipReader, err := gzip.NewReader(tarContent)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in extracting tar file:")
		panic(err)
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
			fmt.Fprintln(os.Stderr, "Error in extracting tar file:")
			panic(err)
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
			fmt.Fprintln(os.Stderr, "Error in untaring file:")
			panic(err)
		}
		if _, err := io.Copy(dest, tarReader); err != nil {
			fmt.Fprintln(os.Stderr, "Error in extracting tar file:")
			panic(err)
		}
	}
	pathToPackage = filepath.Join(".", packageName)
	return
}
