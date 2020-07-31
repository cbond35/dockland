package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// Options for creating a new container.
type NewContainerOpts struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	Name             string
}

// NewContainer creates a new container with the provided config options.
func (di *DockerInterface) NewContainer(
	ctx context.Context, config *NewContainerOpts) error {
	_, err := di.Client.ContainerCreate(
		ctx,
		config.Config,
		config.HostConfig,
		config.NetworkingConfig,
		config.Name)

	if err != nil {
		return fmt.Errorf("failed to create container: %s", err)
	}

	di.RefreshContainers(ctx)
	return nil
}

// RestartContainer restarts the container at idx in the DockerInterface's
// slice of running containers.
func (di *DockerInterface) RestartContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= di.NumRunning() {
		return fmt.Errorf("invalid container index %d", idx)
	}

	id := di.RunningList()[idx].ID

	if err := di.Client.ContainerRestart(ctx, id, nil); err != nil {
		return fmt.Errorf("failed to restart container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// StopContainer stops the running container at idx.
func (di *DockerInterface) StopContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= di.NumRunning() {
		return fmt.Errorf("invalid  container index %d", idx)
	}

	id := di.RunningList()[idx].ID

	if err := di.Client.ContainerStop(ctx, id, nil); err != nil {
		return fmt.Errorf("failed to stop container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// StartContainer starts the stopped container at idx.
func (di *DockerInterface) StartContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= di.NumStopped() {
		return fmt.Errorf("invalid container index %d", idx)
	}

	id := di.StoppedList()[idx].ID

	if err := di.Client.ContainerStart(
		ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// RunningList returns a slice of current running containers.
func (di *DockerInterface) RunningList() []types.Container {
	return di.RunningContainers
}

// StoppedList returns a slice of current stopped containers.
func (di *DockerInterface) StoppedList() []types.Container {
	return di.StoppedContainers
}

// NumRunning returns the current number of running containers.
func (di *DockerInterface) NumRunning() int {
	return len(di.RunningList())
}

// NumStopped returns the current number of stopped containers.
func (di *DockerInterface) NumStopped() int {
	return len(di.StoppedList())
}
