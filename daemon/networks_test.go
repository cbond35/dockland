package daemon

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

// Used to compare networks.
type networkCompare struct {
	name     string
	scope    string
	driver   string
	ipv6     bool
	ingress  bool
	internal bool
}

// Get network resource by id.
func getNetwork(id string) (types.NetworkResource, error) {
	di, _ := NewInterface(context.TODO())

	for _, network := range di.Networks {
		if network.ID == id {
			return network, nil
		}
	}

	return types.NetworkResource{}, fmt.Errorf("network %s not found", id[:idLen])
}

// TestNewNetwork
func TestNewNetwork(t *testing.T) {
	ctx := context.TODO()
	di, err := NewInterface(ctx)

	if err != nil {
		t.Errorf("%s", err)
	}

	tables := []struct {
		opts   map[string]string
		fields networkCompare
	}{
		{
			map[string]string{"name": "test1"},
			networkCompare{"test1", "local", "bridge", false, false, false},
		},
		{
			map[string]string{"name": "test2", "internal": "y"},
			networkCompare{"test2", "local", "bridge", false, false, true},
		},
		{
			map[string]string{"name": "test3", "driver": "macvlan", "internal": "y"},
			networkCompare{"test3", "local", "macvlan", false, false, true},
		},
	}

	for _, table := range tables {
		id, err := di.NewNetwork(ctx, table.opts)
		defer di.RemoveNetwork(ctx, id)

		if err != nil {
			t.Errorf("%s", err)
		}

		network, _ := getNetwork(id)

		want := table.fields
		got := networkCompare{
			network.Name, network.Scope, network.Driver,
			network.EnableIPv6, network.Ingress, network.Internal}

		if got != want {
			t.Fail()
		}
	}
}

// TestRemoveNetwork
func TestRemoveNetwork(t *testing.T) {
	ctx := context.TODO()
	di, err := NewInterface(ctx)

	if err != nil {
		t.Errorf("%s", err)
	}

	testNetwork := map[string]string{"name": "test"}
	want := di.NumNetworks()

	id, err := di.NewNetwork(ctx, testNetwork)

	if err != nil {
		t.Fail()
	}

	err = di.RemoveNetwork(ctx, id)

	if err != nil || di.NumNetworks() != want {
		t.Fail()
	}
}
