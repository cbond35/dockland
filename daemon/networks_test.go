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
	di, err := NewInterface(ctx)

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	for _, table := range tables {
		id, err := di.NewNetwork(ctx, table.opts)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		defer di.RemoveNetwork(ctx, id)

		network, err := getNetwork(id)
		if err != nil {
			t.Errorf("couldn't find network %s", id)
		}

		want := table.fields
		got := networkCompare{
			network.Name, network.Scope, network.Driver,
			network.EnableIPv6, network.Ingress, network.Internal}

		if got != want {
			t.Errorf("networks do not match")
		}
	}
}

// TestRemoveNetwork
func TestRemoveNetwork(t *testing.T) {
	ctx := context.TODO()
	di, err := NewInterface(ctx)

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	testNetwork := map[string]string{"name": "test_network"}
	want := di.NumNetworks() + 1

	id, err := di.NewNetwork(ctx, testNetwork)

	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if di.NumNetworks() != want {
		t.Fail()
	}

	if err = di.RemoveNetwork(ctx, id); err != nil {
		t.Errorf("got error: %s", err)
	}

	want--
	if di.NumNetworks() != want {
		t.Errorf("got %d for NumNetworks, got %d", di.NumNetworks(), want)
	}
}

// TestConnectNetwork + TestDisconnectNetwork
func TestConnectDisconnectNetwork(t *testing.T) {
	ctx := context.TODO()
	di, err := NewInterface(ctx)

	if err != nil {
		t.Errorf("got error: %s", err)
	}

	testNetwork := map[string]string{"name": "test_network"}
	testContainer := map[string]string{"name": "test_container", "cmd": "bash"}

	net_id, err := di.NewNetwork(ctx, testNetwork)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	defer di.RemoveNetwork(ctx, net_id)

	con_id, err := di.NewContainer(ctx, testContainer)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	defer di.RemoveContainer(ctx, con_id)

	if err := di.ConnectNetwork(ctx, net_id, con_id); err != nil {
		t.Errorf("got error: %s", err)
	}

	if err := di.DisconnectNetwork(ctx, net_id, con_id); err != nil {
		t.Errorf("got error: %s", err)
	}
}
