package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (di *DockerInterface) StopContainer(ctx context.Context, idx int) error {
	id := di.RunningContainers[idx].ID

	if err := di.Client.ContainerStop(ctx, id, nil); err != nil {
		return fmt.Errorf("Failed to stop container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)
}

func (di *DockerInterface) StartContainer(ctx context.Context, idx int) error {
	id := di.StoppedContainers[idx].ID

	if err := di.Client.ContainerStart(
		ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("Failed to start container %s: %s", id, err)
	}

	return di.RefreshContainers(ctx)

}
