package daemon

import (
	"context"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

// Used to compare volumes.
type volumeCompare struct {
	name    string
	driver  string
	labels  map[string]string
	options map[string]string
}

// Get volume resource by name.
func getVolume(name string) *types.Volume {
	di, _ := NewInterface(context.TODO())

	for _, volume := range di.Volumes {
		if volume.Name == name {
			return volume
		}
	}

	return &types.Volume{}
}

// TestNewVolume
func TestNewVolume(t *testing.T) {
	tables := []struct {
		opts   map[string]string
		fields volumeCompare
	}{
		{
			map[string]string{"name": "test1"},
			volumeCompare{"test1", "local", map[string]string{}, map[string]string{}},
		},
		{
			map[string]string{"name": "test2", "labels": "case=2"},
			volumeCompare{"test2", "local", map[string]string{"case": "2"}, map[string]string{}},
		},
		{
			map[string]string{"name": "test3", "options": "type=tmpfs,device=tmpfs"},
			volumeCompare{"test3", "local", map[string]string{}, map[string]string{
				"type": "tmpfs", "device": "tmpfs"}},
		},
	}

	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	for _, table := range tables {
		name, err := di.NewVolume(ctx, table.opts)

		if err != nil {
			t.Errorf("got error creating volume: %s", err)
		}
		defer di.RemoveVolume(ctx, name)

		volume := getVolume(name)
		want := table.fields
		got := volumeCompare{
			volume.Name, volume.Driver, volume.Labels, volume.Options}

		if got.name != want.name || got.driver != want.driver || !reflect.DeepEqual(
			got.labels, want.labels) || !reflect.DeepEqual(got.options, want.options) {
			t.Errorf("volumes do not match for volume %s", table.opts["name"])
		}
	}
}

// TestRemoveVolume
func TestRemoveVolume(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := di.NumVolumes()

	testVolume := map[string]string{"name": "remove_volume"}
	name, _ := di.NewVolume(ctx, testVolume)

	if err := di.RemoveVolume(ctx, name); err != nil {
		t.Errorf("got error removing volume: %s", name)
	}

	if di.NumVolumes() != want {
		t.Errorf("got %d volumes, want %d", di.NumVolumes(), want)
	}
}
