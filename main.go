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

	di.StartContainer(ctx, 0)
	di.StartContainer(ctx, 1)

	for _, container := range di.RunningContainers {
		fmt.Printf("%s %s\n", container.ID, container.Image)
	}

	di.StopContainer(ctx, 0)
	di.StopContainer(ctx, 1)

	for _, container := range di.StoppedContainers {
		fmt.Printf("%s %s\n", container.ID, container.Image)
	}
}
