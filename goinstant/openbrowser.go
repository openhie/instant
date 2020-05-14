package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// OpenBrowser for all platforms,
// also see: https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
