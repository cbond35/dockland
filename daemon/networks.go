package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// Options for creating a new network.
type networkConfig struct {
	Config *types.NetworkCreate
	Name   string
}

// newNetworkConfig takes a map of options and creates the necessary
// configuration struct to create a new network.
func newNetworkConfig(opts map[string]string) *networkConfig {
	config := &networkConfig{
		Config: &types.NetworkCreate{},
		Name:   opts["name"],
	}

	config.Config.Driver = opts["driver"]
	config.Config.Scope = opts["scope"]
	config.Config.Attachable = true

	if opts["ipv6"] != "" {
		config.Config.EnableIPv6 = true
	}

	if opts["ingress"] != "" {
		config.Config.Ingress = true
	}

	if opts["internal"] != "" {
		config.Config.Internal = true
	}
	return config
}

// NewNetwork creates a new network with the provided options and returns
// the network's ID.
func (di *DockerInterface) NewNetwork(ctx context.Context, opts map[string]string) (string, error) {
	config := newNetworkConfig(opts)

	response, err := di.Client.NetworkCreate(
		ctx,
		config.Name,
		*config.Config)

	if err != nil {
		return "", fmt.Errorf("failed to create network: %s", err)
	}
	return response.ID, di.RefreshNetworks(ctx)
}

// RemoveNetwork removes a network.
func (di *DockerInterface) RemoveNetwork(ctx context.Context, id string) error {
	if err := di.Client.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network: %s", err)
	}
	return di.RefreshNetworks(ctx)
}

// ConnectNetwork connects a container to a network.
func (di *DockerInterface) ConnectNetwork(ctx context.Context, net, container string) error {
	if err := di.Client.NetworkConnect(
		ctx, net, container, &network.EndpointSettings{}); err != nil {
		return fmt.Errorf("failed to connect network: %s", err)
	}
	return nil
}

// DisconnectNetwork removes a container from a network.
func (di *DockerInterface) DisconnectNetwork(ctx context.Context, net, container string) error {
	if err := di.Client.NetworkDisconnect(ctx, net, container, true); err != nil {
		return fmt.Errorf("failed to disconnect network: %s", err)
	}
	return nil
}

// NumNetworks returns the current number of networks.
func (di *DockerInterface) NumNetworks() int {
	return len(di.Networks)
}
