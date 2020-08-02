package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// Options for creating a new container.
type containerConfig struct {
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	Name             string
}

// newContainerConfig takes a map of options and creates the necessary
// configuration structs to create a new container.
func newContainerConfig(opts map[string]string) *containerConfig {
	config := &containerConfig{
		Config:           &container.Config{},
		HostConfig:       &container.HostConfig{},
		NetworkingConfig: &network.NetworkingConfig{},
		Name:             opts["name"],
	}

	config.Config.Image = opts["image"]

	config.Config.AttachStdin = true
	config.Config.AttachStdout = true
	config.Config.AttachStderr = true

	port := opts["port"]
	hostPort := opts["hostPort"]
	hostIP := opts["hostIP"]

	if port != "" && hostPort != "" {
		if hostIP == "" {
			hostIP = "0.0.0.0"
		}

		bindings := []nat.PortBinding{
			nat.PortBinding{HostIP: hostIP, HostPort: hostPort},
		}

		portMap := nat.PortMap{nat.Port(port + "/tcp"): bindings}
		config.HostConfig.PortBindings = portMap
	}

	return config
}

// NewContainer creates a new container with the provided options and
// returns the container's ID.
func (di *DockerInterface) NewContainer(ctx context.Context, opts map[string]string) (string, error) {
	config := newContainerConfig(opts)

	response, err := di.Client.ContainerCreate(
		ctx,
		config.Config,
		config.HostConfig,
		config.NetworkingConfig,
		config.Name)

	if err != nil {
		return "", fmt.Errorf("failed to create container: %s", err)
	}

	di.RefreshContainers(ctx)
	return response.ID, nil
}

// RestartContainer restarts a running container.
func (di *DockerInterface) RestartContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerRestart(ctx, id, nil); err != nil {
		return fmt.Errorf("failed to restart container %s: %s", id[:idLen], err)
	}

	return di.RefreshContainers(ctx)
}

// StopContainer stops a running container.
func (di *DockerInterface) StopContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerStop(ctx, id, nil); err != nil {
		return fmt.Errorf("failed to stop container %s: %s", id[:idLen], err)
	}

	return di.RefreshContainers(ctx)
}

// StartContainer starts a stopped container.
func (di *DockerInterface) StartContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerStart(
		ctx, id, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %s", id[:idLen], err)
	}

	return di.RefreshContainers(ctx)
}

// RenameContainer renames a container to name.
func (di *DockerInterface) RenameContainer(ctx context.Context, id string, name string) error {
	if err := di.Client.ContainerRename(ctx, id, name); err != nil {
		return fmt.Errorf("failed to rename container %s to %s: %s", id[:idLen], name, err)
	}

	return di.RefreshContainers(ctx)
}

// RemoveContainer removes a container.
func (di *DockerInterface) RemoveContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerRemove(
		ctx, id, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container %s: %s", id[:idLen], err)
	}

	return nil
}

// PruneContainers deletes any unused container data.
func (di *DockerInterface) PruneContainers(ctx context.Context) error {
	if _, err := di.Client.ContainersPrune(ctx, filters.Args{}); err != nil {
		return fmt.Errorf("failed to prune containers: %s", err)
	}

	return di.RefreshContainers(ctx)
}

// NumStopped returns the current number of containers.
func (di *DockerInterface) NumContainers() int {
	return len(di.Containers)
}
