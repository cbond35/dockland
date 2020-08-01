package main

import (
	"context"
	"fmt"

	"github.com/cbbond/dockland/daemon"
)

func main() {
	ctx := context.Background()

	di, err := daemon.NewInterface(ctx)
	if err != nil {
		panic(err)
	}

	opts := map[string]string{
		"image":    "nginx",
		"name":     "nginx_box",
		"port":     "80",
		"hostPort": "80",
	}

	config := di.NewContainerConfig(opts)

	id, err := di.NewContainer(ctx, config)
	if err != nil {
		panic(err)
	}

	if err := di.StartContainer(ctx, id); err != nil {
		panic(err)
	}

	for i := 0; i < di.NumContainers(); i++ {
		fmt.Printf("%s\n", di.ContainerList()[i].Names[0])
	}
}
