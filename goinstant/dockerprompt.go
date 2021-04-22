package main

// composeUpCore is not used any longer
// func composeUpCore() {

// 	home, _ := os.UserHomeDir()
// 	composefile := path.Join(home, ".instant/instant/core/docker/docker-compose.yml")
// 	color.Yellow.Println("Running on", runtime.GOOS)
// 	switch runtime.GOOS {
// 	case "linux", "darwin":
// 		cmd := exec.Command("docker-compose", "-f", composefile, "up", "-d")
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			log.Fatalf("cmd.Run() failed with %s\n", err)
// 		}
// 	case "windows":
// 		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "up", "-d")
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	default:
// 		consoleSender(server, "What operating system is this?")
// 	}

// }

// composeDownCore is not used any longer
// func composeDownCore() {

// 	home, _ := os.UserHomeDir()
// 	composefile := path.Join(home, ".instant/instant/core/docker/docker-compose.yml")
// 	color.Yellow.Println("Running on", runtime.GOOS)
// 	switch runtime.GOOS {
// 	case "linux", "darwin":
// 		cmd := exec.Command("docker-compose", "-f", composefile, "down")
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			log.Fatalf("cmd.Run() failed with %s\n", err)
// 		}
// 	case "windows":
// 		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", composefile, "down")
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	default:
// 		consoleSender(server, "What operating system is this?")
// 	}

// }

// ComposeUp is not used any longer
// func composeUp(composeFile string) {

// 	color.Yellow.Println("Running on", runtime.GOOS)
// 	switch runtime.GOOS {
// 	case "linux", "darwin":
// 		cmd := exec.Command("docker-compose", "-f", "-", "up", "-d")
// 		cmd.Stdin = strings.NewReader(composeFile)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			log.Fatalf("cmd.Run() failed with %s\n", err)
// 		}
// 	case "windows":
// 		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "up", "-d")
// 		cmd.Stdin = strings.NewReader(composeFile)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	default:
// 		consoleSender(server, "What operating system is this?")
// 	}

// }

// ComposeDown is not used any longer
// func composeDown(composeFile string) {
// 	color.Yellow.Println("Running on", runtime.GOOS)
// 	switch runtime.GOOS {
// 	case "linux", "darwin":
// 		cmd := exec.Command("docker-compose", "-f", "-", "down")
// 		cmd.Stdin = strings.NewReader(composeFile)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Run()
// 		if err != nil {
// 			log.Fatalf("cmd.Run() failed with %s\n", err)
// 		}
// 	case "windows":
// 		cmd := exec.Command("cmd", "/C", "docker-compose", "-f", "-", "down")
// 		cmd.Stdin = strings.NewReader(composeFile)
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println("Error: ", err)
// 		}
// 	default:
// 		consoleSender(server, "What operating system is this?")
// 	}
// }

// ComposeGet is not used any longer
// func composeGet(url string) string {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		color.Red.Println("Are you connected to the Internet? Error:", err)
// 		consoleSender(server, "Are you connected to the Internet? Error")
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		color.Red.Println("Strange error reading the downloaded body. Error:", err)
// 		consoleSender(server, "Strange error reading the downloaded body.")
// 	}
// 	fmt.Println(string(body))
// 	return (string(body))
// }
