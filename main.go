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
		"name": "test_network",
	}

	_, err = di.NewNetwork(ctx, opts)
	if err != nil {
		panic(err)
	}

	for i := 0; i < di.NumNetworks(); i++ {
		fmt.Printf("%s\n", di.Networks[i].Name)
	}
}
