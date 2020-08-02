package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// newVolumeConfig takes a map of options and creates the necessary
// configuration struct to create a new volume.
func newVolumeConfig(opts map[string]string) *volume.VolumeCreateBody {
	config := &volume.VolumeCreateBody{
		Name:   opts["name"],
		Driver: opts["driver"],
	}

	return config
}

// NewVolume creates a new volume with the provided options and returns the
// volume's name.
func (di *DockerInterface) NewVolume(ctx context.Context, opts map[string]string) (string, error) {
	config := *newVolumeConfig(opts)

	response, err := di.Client.VolumeCreate(
		ctx,
		config)

	if err != nil {
		return "", fmt.Errorf("failed to create volume: %s", err)
	}

	return response.Name, di.RefreshVolumes(ctx)
}

// RemoveVolume removes a volume.
func (di *DockerInterface) RemoveVolume(ctx context.Context, id string) error {
	if err := di.Client.VolumeRemove(ctx, id, true); err != nil {
		return fmt.Errorf("failed to remove volume %s: %s", id[:idLen], err)
	}

	return di.RefreshVolumes(ctx)
}

// PruneVolumes removes any unused volumes.
func (di *DockerInterface) PruneVolumes(ctx context.Context) error {
	if _, err := di.Client.VolumesPrune(ctx, filters.Args{}); err != nil {
		return fmt.Errorf("failed to prune volumes: %s", err)
	}

	return di.RefreshVolumes(ctx)
}

// NumVolumes returns the current number of volumes.
func (di *DockerInterface) NumVolumes() int {
	return len(di.Volumes)
}
