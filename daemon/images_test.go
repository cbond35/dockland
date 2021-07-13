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

	if err := di.PullImage(ctx, "debian"); err != nil {
		t.Logf("got error pulling image: %s", err)
		t.FailNow()
	}
	defer di.RemoveImage(ctx, "debian")

	if di.NumImages() != want {
		t.Errorf("got %d images, want %d", di.NumImages(), want)
	}
}

// TestRemoveImage
func TestRemoveImage(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	di.PullImage(ctx, "ubuntu")

	if err := di.RemoveImage(ctx, "ubuntu"); err != nil {
		t.Logf("got error removing image: %s", err)
		t.FailNow()
	}
	if err := di.RemoveImage(ctx, "no_such_image"); err == nil {
		t.Error("expected error removing image")
	}
}

// TestSearchImages
func TestSearchImages(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := MaxImageResults

	results, err := di.SearchImage(ctx, "busybox")

	if err != nil {
		t.Logf("got error searching images: %s", err)
		t.FailNow()
	}

	if len(results) != want {
		t.Errorf("got %d images in search, want %d", len(results), want)
	}
}
