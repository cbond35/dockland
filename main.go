package main

import (
	"context"

	"github.com/cbbond/dockland/daemon"
)

func main() {
	ctx := context.Background()

	di, err := daemon.NewInterface(ctx)
	if err != nil {
		panic(err)
	}

	if err = di.PruneContainers(ctx); err != nil {
		panic(err)
	}

	if err = di.PruneImages(ctx); err != nil {
		panic(err)
	}

	results, searchErr := di.SearchImage(ctx, "debian")

	if searchErr != nil {
		panic(err)
	}

	err = di.PullImage(ctx, results[0].Name)

	if err != nil {
		panic(err)
	}
}
