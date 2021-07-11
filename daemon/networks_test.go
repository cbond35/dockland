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
	return types.NetworkResource{}, fmt.Errorf("no network %s found", id)
}

// TestNewNetwork
func TestNewNetwork(t *testing.T) {
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

	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	for _, table := range tables {
		id, err := di.NewNetwork(ctx, table.opts)

		if err != nil {
			t.Logf("got error creating network: %s", err)
			t.FailNow()
		}
		defer di.RemoveNetwork(ctx, id)

		network, err := getNetwork(id)
		if err != nil {
			t.Logf("got error finding network: %s", err)
			t.FailNow()
		}

		want := table.fields
		got := networkCompare{
			network.Name, network.Scope, network.Driver,
			network.EnableIPv6, network.Ingress, network.Internal}

		if got != want {
			t.Errorf("networks do not match for network %s", table.opts["name"])
		}
	}
}

// TestRemoveNetwork
func TestRemoveNetwork(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)
	want := di.NumNetworks()

	testNetwork := map[string]string{"name": "test_network"}
	id, _ := di.NewNetwork(ctx, testNetwork)

	if err := di.RemoveNetwork(ctx, id); err != nil {
		t.Logf("got error removing network: %s", err)
		t.FailNow()
	}
	if di.NumNetworks() != want {
		t.Errorf("got %d networks, want %d", di.NumNetworks(), want)
	}
}

// TestConnectNetwork
func TestConnectNetwork(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	testNetwork := map[string]string{"name": "test_network"}
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	netID, _ := di.NewNetwork(ctx, testNetwork)
	defer di.RemoveNetwork(ctx, netID)

	conID, _ := di.NewContainer(ctx, testContainer)
	defer di.RemoveContainer(ctx, conID)

	if err := di.ConnectNetwork(ctx, netID, conID); err != nil {
		t.Errorf("got error connecting network: %s", err)
	}
}

// TestDisconnectNetwork
func TestDisconnectNetwork(t *testing.T) {
	ctx := context.TODO()
	di, _ := NewInterface(ctx)

	testNetwork := map[string]string{"name": "test_network"}
	testContainer := map[string]string{"name": "test_container", "image": "nginx"}

	netID, _ := di.NewNetwork(ctx, testNetwork)
	defer di.RemoveNetwork(ctx, netID)

	conID, _ := di.NewContainer(ctx, testContainer)
	di.ConnectNetwork(ctx, netID, conID)
	defer di.RemoveContainer(ctx, conID)

	if err := di.DisconnectNetwork(ctx, netID, conID); err != nil {
		t.Errorf("got error disconnecting network: %s", err)
	}
}
