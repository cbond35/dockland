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
	di.StartContainer(ctx, 0)

	di.RestartContainer(ctx, 0)
	di.RestartContainer(ctx, 0)

	for _, container := range di.Running() {
		fmt.Printf("%s %s\n", container.ID, container.Image)
	}

	di.StopContainer(ctx, 0)
	di.StopContainer(ctx, 0)

	for _, container := range di.Stopped() {
		fmt.Printf("%s %s\n", container.ID, container.Image)
	}
}
