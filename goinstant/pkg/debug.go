package pkg

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gookit/color"
	"golang.org/x/net/context"
)

// Debug contains common issue checking for Docker, etc
func Debug() {
	cwd, err := os.Getwd()
	if err != nil {
		color.Red.Println("Can't get current working directory... this is not a great error.")
		panic(err)
	} else {
		color.Green.Println("Running goinstant [v-alpha] from:", cwd)
	}

	// ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		color.Red.Println("Unable to list Docker containers. Please ensure that Docker is downloaded and running!")
		panic(err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		color.Red.Println("Unable to get Docker context. Please ensure that Docker is downloaded and running!")
		panic(err)
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		fmt.Printf("%d bytes memory is allocated.\n", info.MemTotal)
	}

	// fmt.Println(reflect.TypeOf(containers).String())
	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}
