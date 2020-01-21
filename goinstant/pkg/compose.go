package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

// ComposeGet gets a docker-cmpose from github
func ComposeGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		color.Red.Println("Are you connected to the Internet? Error:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red.Println("Strange error reading the downloaded body. Error:", err)
	}
	fmt.Println(string(body))
	return (string(body))
}

// ComposeUp brings up based on docker-compose
func ComposeUp(composeFile string) {

	fmt.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", "-", "up", "-d")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "up", "-d")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		fmt.Println("What operating system is this?", runtime.GOOS)
	}

}

// ComposeDown stops containers based on docker-compose
func ComposeDown(composeFile string) {
	fmt.Println("Running on", runtime.GOOS)
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("docker-compose", "-f", "-", "down")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "down")
		cmd.Stdin = strings.NewReader(composeFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	default:
		fmt.Println("What operating system is this?", runtime.GOOS)

	}
}
