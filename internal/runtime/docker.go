package runtime

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const NetworkLabelRuntimeID = "com.gobcn.radius-director.runtime-id"

type NetworkClient interface {
	NetworkInspect(
		ctx context.Context,
		networkID string,
		options client.NetworkInspectOptions,
	) (client.NetworkInspectResult, error)

	NetworkCreate(
		ctx context.Context,
		name string,
		options client.NetworkCreateOptions,
	) (client.NetworkCreateResult, error)

	NetworkRemove(
		ctx context.Context,
		networkID string,
		options client.NetworkRemoveOptions,
	) (client.NetworkRemoveResult, error)
}

type DockerNetworkClient struct {
	client *client.Client
}

func NewDockerNetworkClient() (*DockerNetworkClient, error) {
	dockerClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("create Docker client: %w", err)
	}

	return &DockerNetworkClient{
		client: dockerClient,
	}, nil
}

func (c *DockerNetworkClient) NetworkInspect(
	ctx context.Context,
	networkID string,
	options client.NetworkInspectOptions,
) (client.NetworkInspectResult, error) {
	return c.client.NetworkInspect(ctx, networkID, options)
}

func (c *DockerNetworkClient) NetworkCreate(
	ctx context.Context,
	name string,
	options client.NetworkCreateOptions,
) (client.NetworkCreateResult, error) {
	return c.client.NetworkCreate(ctx, name, options)
}

func (c *DockerNetworkClient) NetworkRemove(
	ctx context.Context,
	networkID string,
	options client.NetworkRemoveOptions,
) (client.NetworkRemoveResult, error) {
	return c.client.NetworkRemove(ctx, networkID, options)
}

func (c *DockerNetworkClient) Close() error {
	return c.client.Close()
}

func EnsureDockerNetwork(
	ctx context.Context,
	networkClient NetworkClient,
	runtime Runtime,
) (bool, error) {
	inspect, err := networkClient.NetworkInspect(
		ctx,
		runtime.NetworkName,
		client.NetworkInspectOptions{},
	)

	if err == nil {
		if err := verifyRuntimeNetwork(inspect.Network, runtime); err != nil {
			return false, err
		}

		return false, nil
	}

	if !errdefs.IsNotFound(err) {
		return false, fmt.Errorf(
			"inspect Docker network %q: %w",
			runtime.NetworkName,
			err,
		)
	}

	_, err = networkClient.NetworkCreate(
		ctx,
		runtime.NetworkName,
		client.NetworkCreateOptions{
			Driver: "bridge",
			Labels: map[string]string{
				NetworkLabelRuntimeID: runtime.ID,
			},
		},
	)

	if err == nil {
		return true, nil
	}

	if !errdefs.IsAlreadyExists(err) {
		return false, fmt.Errorf(
			"create Docker network %q: %w",
			runtime.NetworkName,
			err,
		)
	}

	// Another process may have created the network after our
	// initial inspection. Inspect it again and verify ownership.
	inspect, err = networkClient.NetworkInspect(
		ctx,
		runtime.NetworkName,
		client.NetworkInspectOptions{},
	)
	if err != nil {
		return false, fmt.Errorf(
			"inspect Docker network %q after creation conflict: %w",
			runtime.NetworkName,
			err,
		)
	}

	if err := verifyRuntimeNetwork(inspect.Network, runtime); err != nil {
		return false, err
	}

	return false, nil
}

func verifyRuntimeNetwork(
	network network.Inspect,
	runtime Runtime,
) error {
	runtimeID := network.Labels[NetworkLabelRuntimeID]

	if runtimeID == "" {
		return fmt.Errorf(
			"Docker network %q already exists but is not owned by RADIUS Director",
			runtime.NetworkName,
		)
	}

	if runtimeID != runtime.ID {
		return fmt.Errorf(
			"Docker network %q belongs to a different RADIUS Director runtime",
			runtime.NetworkName,
		)
	}

	return nil
}
