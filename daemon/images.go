package daemon

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
)

// MaxImageResults represents the maximum number of results from a search.
var MaxImageResults int = 10

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
	return di.RefreshImages(ctx)
}

// RemoveImage removes an image. id can be the ID or the image name.
func (di *DockerInterface) RemoveImage(ctx context.Context, id string) error {
	if _, err := di.Client.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove image: %s", err)
	}
	return di.RefreshImages(ctx)
}

// SearchImage searches the registry for the given image.
func (di *DockerInterface) SearchImage(ctx context.Context, image string) ([]registry.SearchResult, error) {
	results, err := di.Client.ImageSearch(ctx, image, types.ImageSearchOptions{Limit: MaxImageResults})

	if err != nil {
		return nil, fmt.Errorf("failed to search image: %s", err)
	}
	return results, nil
}

// NumImages return the number of images.
func (di *DockerInterface) NumImages() int {
	return len(di.Images)
}
