package daemon

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/volume"
)

// newVolumeConfig takes a map of options and creates the necessary
// configuration struct to create a new volume.
func newVolumeConfig(opts map[string]string) *volume.VolumeCreateBody {
	config := &volume.VolumeCreateBody{
		Name:       opts["name"],
		Driver:     opts["driver"],
		Labels:     make(map[string]string),
		DriverOpts: make(map[string]string),
	}

	for _, label := range strings.Split(opts["labels"], ",") {
		key_value := strings.SplitN(label, "=", 2)

		if len(key_value) > 1 {
			config.Labels[key_value[0]] = key_value[1]
		}
	}

	for _, driver_opt := range strings.Split(opts["options"], ",") {
		key_value := strings.SplitN(driver_opt, "=", 2)

		if len(key_value) > 1 {
			config.DriverOpts[key_value[0]] = key_value[1]
		}
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
		return "", err
	}

	return response.Name, di.RefreshVolumes(ctx)
}

// RemoveVolume removes a volume.
func (di *DockerInterface) RemoveVolume(ctx context.Context, id string) error {
	if err := di.Client.VolumeRemove(ctx, id, true); err != nil {
		return err
	}

	return di.RefreshVolumes(ctx)
}

// NumVolumes returns the current number of volumes.
func (di *DockerInterface) NumVolumes() int {
	return len(di.Volumes)
}
