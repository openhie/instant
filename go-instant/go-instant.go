package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}

	// ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// fmt.Println(reflect.TypeOf(containers).String())

	info, err := cli.Info(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d bytes memory is allocated.\n", info.MemTotal)

	for _, container := range containers {
		fmt.Printf("ContainerID: %s Status: %s Image: %s\n", container.ID[:10], container.State, container.Image)
	}

}
