package daemon

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
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

	return di.RefreshImages(ctx)
}

// RemoveImage removes an image.
func (di *DockerInterface) RemoveImage(ctx context.Context, id string) error {
	if _, err := di.Client.ImageRemove(ctx, id, types.ImageRemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove image %s: %s", id[:idLen], err)
	}

	return di.RefreshImages(ctx)
}

// SearchImages searches the registry for term.
func (di *DockerInterface) SearchImage(ctx context.Context, term string) ([]registry.SearchResult, error) {
	results, err := di.Client.ImageSearch(ctx, term, types.ImageSearchOptions{Limit: 10})

	if err != nil {
		return nil, fmt.Errorf("failed to search for image %s: %s", term, err)
	}

	return results, nil
}

// PruneImages removes unused image data.
func (di *DockerInterface) PruneImages(ctx context.Context) error {
	if _, err := di.Client.ImagesPrune(ctx, filters.Args{}); err != nil {
		return fmt.Errorf("failed to prune images: %s", err)
	}

	return di.RefreshImages(ctx)
}

// NumImages return the number of images.
func (di *DockerInterface) NumImages() int {
	return len(di.Images)
}
