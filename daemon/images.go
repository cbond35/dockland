package daemon

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// PullImage pulls the image with the given img name.
func (di *DockerInterface) PullImage(ctx context.Context, img string) error {
	response, err := di.Client.ImagePull(ctx, img, types.ImagePullOptions{})

	if err != nil {
		return fmt.Errorf("failed to pull image: %s", err)
	}

	scan := bufio.NewScanner(response)
	for scan.Scan() {
	}

	if err := scan.Err(); err != nil {
		return fmt.Errorf("failed to pull image: %s", err)
	}

	return nil
}

// RemoveImage removes the image at idx in the image list.
func (di *DockerInterface) RemoveImage(ctx context.Context, idx int) error {
	if idx < 0 || idx >= di.NumImages() {
		return fmt.Errorf("invalid image index %d", idx)
	}

	id := di.ImageList()[idx].ID

	if _, err := di.Client.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove image: %s", err)
	}

	return nil
}

// ImageList returns a list of all images in the Docker host.
func (di *DockerInterface) ImageList() []types.ImageSummary {
	return di.Images
}

func (di *DockerInterface) NumImages() int {
	return len(di.ImageList())
}
