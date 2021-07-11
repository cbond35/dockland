package daemon

import (
	"context"
	"testing"
)

// TestNewContainer
func TestNewContainer(t *testing.T) {
}

// TestRestartContainer
func TestRestartContainer(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	conID, _ := di.NewContainer(ctx, testContainer)
	defer di.RemoveContainer(ctx, conID)

	if err := di.RestartContainer(ctx, conID); err != nil {
		t.Errorf("got error restarting container: %s", err)
	}
}

// TestStopContainer
func TestStopContainer(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	conID, _ := di.NewContainer(ctx, testContainer)
	defer di.RemoveContainer(ctx, conID)

	if err := di.StopContainer(ctx, conID); err != nil {
		t.Errorf("got error stopping container: %s", err)
	}
}

// TestStartContainer
func TestStartContainer(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	conID, _ := di.NewContainer(ctx, testContainer)
	defer di.RemoveContainer(ctx, conID)

	di.StopContainer(ctx, conID)
	if err := di.StartContainer(ctx, conID); err != nil {
		t.Errorf("got error starting container: %s", err)
	}
}

// TestRenameContainer
func TestRenameContainer(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}
	newName := "renamed_container"

	conID, _ := di.NewContainer(ctx, testContainer)
	defer di.RemoveContainer(ctx, conID)

	if err := di.RenameContainer(ctx, conID, newName); err != nil {
		t.Errorf("got error renaming container: %s", err)
	}

	// TODO: Confirm the container's name has actually been changed.
	// Will it be useful to implement a GetContainer(name), Network, etc.
	// for the UI? Will make testing easier as well.
}

// TestRemoveContainer
func TestRemoveContainer(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := di.NumContainers()
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	conID, _ := di.NewContainer(ctx, testContainer)

	if err := di.RemoveContainer(ctx, conID); err != nil {
		t.Errorf("got error renaming container: %s", err)
	}
	if di.NumContainers() != want {
		t.Errorf("got %d containers, want %d", di.NumContainers(), want)
	}
}
