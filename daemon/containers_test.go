package daemon

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

// Used to compare containers.
type containerCompare struct {
	name     string
	image    string
	port     int
	hostPort int
	hostIP   string
}

// Get container resource by id.
func getContainer(id string) (types.Container, error) {
	di, _ := NewInterface(context.TODO())

	for _, container := range di.Containers {
		if container.ID == id {
			return container, nil
		}
	}
	return types.Container{}, fmt.Errorf("no container %s found", id)
}

// TestNewContainer
func TestNewContainer(t *testing.T) {
	tables := []struct {
		opts   map[string]string
		fields containerCompare
	}{
		{
			map[string]string{"name": "test1", "image": "nginx"},
			containerCompare{"/test1", "nginx", 80, 0, ""},
		},
		{
			map[string]string{"name": "test2", "image": "nginx", "port": "80", "hostPort": "8080"},
			containerCompare{"/test2", "nginx", 80, 8080, "::"},
		},
		{
			map[string]string{"name": "test3", "image": "nginx", "env": "IS_TEST=TRUE"},
			containerCompare{"/test3", "nginx", 80, 0, ""},
		},
		{
			map[string]string{"name": "test4", "image": "alpine",
				"entrypoint": "/bin/echo", "cmd": "Hello World!"},
			containerCompare{"/test4", "alpine", 0, 0, "::"},
		},
	}

	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	for _, table := range tables {
		id, err := di.NewContainer(ctx, table.opts)

		if err != nil {
			t.Logf("got error creating container: %s", err)
			t.FailNow()
		}
		defer di.RemoveContainer(ctx, id)

		if err := di.StartContainer(ctx, id); err != nil {
			t.Logf("got error starting container: %s", err)
			t.FailNow()
		}

		container, err := getContainer(id)
		if err != nil {
			t.Logf("got error finding container: %s", err)
			t.FailNow()
		}

		want := table.fields
		got := containerCompare{
			container.Names[0], container.Image, 0, 0, "::"}

		if len(container.Ports) > 0 {
			got.port = int(container.Ports[0].PrivatePort)
			got.hostPort = int(container.Ports[0].PublicPort)
			got.hostIP = container.Ports[0].IP
		}

		if got != want {
			if got.name != want.name {
				t.Errorf("container name does not match for container %s: got %s, want %s",
					table.opts["name"], got.name, want.name)
			} else if got.image != want.image {
				t.Errorf("container image does not match for container %s: got %s, want %s",
					table.opts["name"], got.image, want.image)
			} else if got.port != want.port {
				t.Errorf("container port does not match for container %s: got %d, want %d ",
					table.opts["name"], got.port, want.port)
			} else if got.hostPort != want.hostPort {
				t.Errorf("host port does not match for container %s: got %d, want %d",
					table.opts["name"], got.hostPort, want.hostPort)
			} else if got.hostIP != want.hostIP {
				t.Errorf("host IP does not match for container %s: got %s, want %s",
					table.opts["name"], got.hostIP, want.hostIP)
			}
		}
	}
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
		t.Logf("got error renaming container: %s", err)
		t.FailNow()
	}
	if di.NumContainers() != want {
		t.Errorf("got %d containers, want %d", di.NumContainers(), want)
	}
}
