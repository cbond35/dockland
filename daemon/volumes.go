package daemon

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/volume"
)

// newVolumeCreateBody takes a map of options and creates the necessary
// configuration struct to create a new volume.
func newVolumeCreateBody(opts map[string]string) *volume.VolumeCreateBody {
	config := &volume.VolumeCreateBody{
		Name:       opts["name"],
		Driver:     opts["driver"],
		Labels:     make(map[string]string),
		DriverOpts: make(map[string]string),
	}

	for _, label := range strings.Split(opts["labels"], ",") {
		optValue := strings.SplitN(label, "=", 2)

		if len(optValue) > 1 {
			config.Labels[optValue[0]] = optValue[1]
		}
	}

	for _, driverOpt := range strings.Split(opts["options"], ",") {
		optValue := strings.SplitN(driverOpt, "=", 2)

		if len(optValue) > 1 {
			config.DriverOpts[optValue[0]] = optValue[1]
		}
	}

	return config
}

// NewVolume creates a new volume with the provided options and returns the
// volume's name.
func (di *DockerInterface) NewVolume(ctx context.Context, opts map[string]string) (string, error) {
	config := *newVolumeCreateBody(opts)

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
