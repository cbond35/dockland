package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// RestartContainer restarts the container at idx in the DockerInterface's
// slice of running containers.
func (di *DockerInterface) RestartContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= len(di.RunningContainers) {
		return fmt.Errorf("Invalid index %s", idx)
	}

	id := di.RunningContainers[idx].ID

	if err := di.Client.ContainerRestart(ctx, id, nil); err != nil {
		return fmt.Errorf("Failed to stop container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// StopContainer stops the container at idx in the DockerInterface's
// slice of running containers.
func (di *DockerInterface) StopContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= len(di.RunningContainers) {
		return fmt.Errorf("Invalid index %s", idx)
	}

	id := di.RunningContainers[idx].ID

	if err := di.Client.ContainerStop(ctx, id, nil); err != nil {
		return fmt.Errorf("Failed to stop container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// StartContainer starts the container at idx in the DockerInterface's
// slice of stopped containers.
func (di *DockerInterface) StartContainer(ctx context.Context, idx int) error {
	if idx < 0 || idx >= len(di.StoppedContainers) {
		return fmt.Errorf("Invalid index %s", idx)
	}

	id := di.StoppedContainers[idx].ID

	if err := di.Client.ContainerStart(
		ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("Failed to start container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

// Running returns a slice of current running containers.
func (di *DockerInterface) Running() []types.Container {
	return di.RunningContainers
}

// Stopped returns a slice of current stopped containers.
func (di *DockerInterface) Stopped() []types.Container {
	return di.StoppedContainers
}
