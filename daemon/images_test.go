package daemon

import (
	"context"
	"testing"
)

// TestPullImage
func TestPullImage(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := di.NumImages() + 1

	if err := di.PullImage(ctx, "alpine"); err != nil {
		t.Errorf("got error pulling image: %s", err)
	}
	defer di.RemoveImage(ctx, "alpine")

	if di.NumImages() != want {
		t.Errorf("got %d images, want %d", di.NumImages(), want)
	}
}

// TestRemoveImage
func TestRemoveImage(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := di.NumImages()

	di.PullImage(ctx, "alpine")

	if err := di.RemoveImage(ctx, "alpine"); err != nil {
		t.Errorf("got error removing image: %s", err)
	}

	if di.NumImages() != want {
		t.Errorf("got %d images, want %d", di.NumImages(), want)
	}
}

// TestSearchImages
func TestSearchImages(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := MaxResults // Number of images in the list.

	results, err := di.SearchImage(ctx, "alpine")

	if err != nil {
		t.Errorf("got error searching images: %s", err)
	}

	if len(results) != want {
		t.Errorf("got %d images in search, want %d", len(results), want)
	}
}

// TestPruneImages
func TestPruneImages(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	if err := di.PruneImages(ctx); err != nil {
		t.Errorf("got error pruning images: %s", err)
	}
}
