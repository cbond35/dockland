package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	cli "github.com/docker/docker/client"
)

// DockerInterface is our primary source of information about
// the Docker daemon and associated containers, images, networks,
// and volumes.
type DockerInterface struct {
	Client     *cli.Client
	Containers []types.Container
	Images     []types.ImageSummary
	Info       types.Info
	Networks   []types.NetworkResource
	Volumes    volume.VolumeListOKBody
}

// Refresh updates the given DockerInterface's fields with the latest
// information from the Docker API.
func (di *DockerInterface) Refresh(ctx context.Context) error {
	var err error

	di.Containers, err = di.Client.ContainerList(
		ctx, types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to fetch containers: %s", err)
	}

	di.Images, err = di.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to fetch images: %s", err)
	}

	di.Info, err = di.Client.Info(ctx)
	if err != nil {
		return fmt.Errorf("Failed to fetch client information: %s", err)
	}

	di.Networks, err = di.Client.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to fetch networks: %s", err)
	}

	di.Volumes, err = di.Client.VolumeList(ctx, filters.Args{})
	if err != nil {
		return fmt.Errorf("Failed to fetch volumes: %s", err)
	}

	return err
}

// NewInterface returns a DockerInterface with information about the
// Docker daemon and all containers, images, networks, and volumes.
func NewInterface(ctx context.Context) (*DockerInterface, error) {
	di := &DockerInterface{}
	var err error

	di.Client, err = cli.NewClientWithOpts(
		cli.FromEnv, cli.WithAPIVersionNegotiation())

	if err != nil {
		return nil, fmt.Errorf("Failed to initialize client: %s", err)
	}

	err = di.Refresh(ctx)
	if err != nil {
		return nil, err
	}

	return di, nil
}
