package daemon

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// newContainerCreateConfig takes a map of options and creates the necessary
// configuration structs to create a new container.
func newContainerCreateConfig(opts map[string]string) *types.ContainerCreateConfig {
	config := &types.ContainerCreateConfig{
		Config:     &container.Config{},
		HostConfig: &container.HostConfig{},
		Name:       opts["name"],
	}

	config.Config.Image = opts["image"]
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
			{HostIP: hostIP, HostPort: hostPort},
			{HostIP: "::", HostPort: hostPort},
		}
		config.HostConfig.PortBindings = nat.PortMap{nat.Port(port + "/tcp"): bindings}
	}

	if env := opts["env"]; env != "" {
		for _, arg := range strings.Split(env, ",") {
			config.Config.Env = append(config.Config.Env, strings.TrimSpace(arg))
		}
	}
	if cmd := opts["cmd"]; cmd != "" {
		for _, arg := range strings.Split(cmd, ",") {
			config.Config.Cmd = append(config.Config.Cmd, strings.TrimSpace(arg))
		}
	}
	if entryPoint := opts["entrypoint"]; entryPoint != "" {
		for _, arg := range strings.Split(entryPoint, ",") {
			config.Config.Entrypoint = append(config.Config.Entrypoint, strings.TrimSpace(arg))
		}

	}
	return config
}

// NewContainer creates a new container with the provided options and
// returns the container's ID.
func (di *DockerInterface) NewContainer(ctx context.Context,
	opts map[string]string) (string, error) {
	config := newContainerCreateConfig(opts)

	response, err := di.Client.ContainerCreate(
		ctx,
		config.Config,
		config.HostConfig,
		nil,
		nil,
		config.Name)

	if err != nil {
		return "", err
	}

	return response.ID, di.RefreshContainers(ctx)
}

// RestartContainer restarts a running container.
func (di *DockerInterface) RestartContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerRestart(ctx, id, nil); err != nil {
		return err
	}
	return di.RefreshContainers(ctx)
}

// StopContainer stops a running container.
func (di *DockerInterface) StopContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerStop(ctx, id, nil); err != nil {
		return err
	}
	return di.RefreshContainers(ctx)
}

// StartContainer starts a stopped container.
func (di *DockerInterface) StartContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerStart(
		ctx, id, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return di.RefreshContainers(ctx)
}

// RenameContainer renames a container to name.
func (di *DockerInterface) RenameContainer(ctx context.Context,
	id string, name string) error {
	if err := di.Client.ContainerRename(ctx, id, name); err != nil {
		return err
	}
	return di.RefreshContainers(ctx)
}

// RemoveContainer removes a container.
func (di *DockerInterface) RemoveContainer(ctx context.Context, id string) error {
	if err := di.Client.ContainerRemove(
		ctx, id, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}
	return di.RefreshContainers(ctx)
}

// NumContainers returns the current number of containers.
func (di *DockerInterface) NumContainers() int {
	return len(di.Containers)
}
