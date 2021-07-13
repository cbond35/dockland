// The daemon package controls all interaction with the Docker daemon. Any
// options for containers, images, networks, and volumes are all routed
// through this package.

package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// ResourceRefreshError is returned whenever a container, image,
// info, network, or volume refresh fails.
type ResourceRefreshError struct {
	Resource string
	Err      error
}

// DockerInterface is our primary source of information about
// the Docker daemon and associated containers, images, networks,
// and volumes.
type DockerInterface struct {
	Client     *client.Client
	Containers []types.Container
	Images     []types.ImageSummary
	Info       types.Info
	Networks   []types.NetworkResource
	Volumes    []*types.Volume
}

// Error is called whenever the daemon fails to send us an updated resource list.
func (r *ResourceRefreshError) Error() string {
	return fmt.Sprintf("failed to refresh %s: %s", r.Resource, r.Err)
}

// RefreshContainers updates the DockerInterface's Containers fields
// with the latest information from the Docker API.
func (di *DockerInterface) RefreshContainers(ctx context.Context) error {
	var err error

	if di.Containers, err = di.Client.ContainerList(
		ctx, types.ContainerListOptions{All: true}); err != nil {
		return &ResourceRefreshError{"container list", err}
	}
	return nil
}

// RefreshImages updates the DockerInterface's Images field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshImages(ctx context.Context) error {
	var err error

	if di.Images, err = di.Client.ImageList(
		ctx, types.ImageListOptions{All: true}); err != nil {
		return &ResourceRefreshError{"image list", err}
	}
	return nil
}

// RefreshInfo updates the DockerInterface's Info field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshInfo(ctx context.Context) error {
	var err error

	if di.Info, err = di.Client.Info(ctx); err != nil {
		return &ResourceRefreshError{"docker info", err}
	}
	return nil
}

// RefreshNetworks updates the DockerInterface's Networks field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshNetworks(ctx context.Context) error {
	var err error

	if di.Networks, err = di.Client.NetworkList(
		ctx, types.NetworkListOptions{}); err != nil {
		return &ResourceRefreshError{"network list", err}
	}
	return nil
}

// RefreshVolumes updates the DockerInterface's Volumes field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshVolumes(ctx context.Context) error {
	var err error
	var volumeBody volume.VolumeListOKBody

	if volumeBody, err = di.Client.VolumeList(ctx, filters.Args{}); err != nil {
		return &ResourceRefreshError{"volume list", err}
	}

	di.Volumes = volumeBody.Volumes
	return nil
}

// NewInterface returns a DockerInterface with information about the
// Docker daemon and all containers, images, networks, and volumes.
func NewInterface(ctx context.Context) (*DockerInterface, error) {
	var err error
	di := &DockerInterface{}

	di.Client, err = client.NewClientWithOpts(
		client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(fmt.Errorf("failed to connect to docker daemon: %s", err))
	}

	if err = di.RefreshContainers(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshImages(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshInfo(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshNetworks(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshVolumes(ctx); err != nil {
		return nil, err
	}
	return di, nil
}
