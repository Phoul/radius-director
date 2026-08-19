package runtime

import (
	"context"
	"errors"
	"testing"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type fakeNetworkClient struct {
	inspectResults []client.NetworkInspectResult
	inspectErrors  []error

	createResult client.NetworkCreateResult
	createError  error

	removeError error

	inspectCalls int
	createCalls  int
	removeCalls  int

	lastInspectedNetwork string
	lastCreatedName      string
	lastCreateOptions    client.NetworkCreateOptions
	lastRemovedNetwork   string
	lastRemoveOptions    client.NetworkRemoveOptions
}

func (f *fakeNetworkClient) NetworkInspect(
	_ context.Context,
	networkID string,
	_ client.NetworkInspectOptions,
) (client.NetworkInspectResult, error) {
	index := f.inspectCalls
	f.inspectCalls++

	f.lastInspectedNetwork = networkID

	var result client.NetworkInspectResult
	if index < len(f.inspectResults) {
		result = f.inspectResults[index]
	}

	var err error
	if index < len(f.inspectErrors) {
		err = f.inspectErrors[index]
	}

	return result, err
}

func (f *fakeNetworkClient) NetworkCreate(
	_ context.Context,
	name string,
	options client.NetworkCreateOptions,
) (client.NetworkCreateResult, error) {
	f.createCalls++
	f.lastCreatedName = name
	f.lastCreateOptions = options

	return f.createResult, f.createError
}

func (f *fakeNetworkClient) NetworkRemove(
	_ context.Context,
	networkID string,
	options client.NetworkRemoveOptions,
) (client.NetworkRemoveResult, error) {
	f.removeCalls++
	f.lastRemovedNetwork = networkID
	f.lastRemoveOptions = options

	return client.NetworkRemoveResult{}, f.removeError
}

func TestEnsureDockerNetworkCreatesMissingNetwork(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectErrors: []error{
			errdefs.ErrNotFound,
		},
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err != nil {
		t.Fatalf("EnsureDockerNetwork() error = %v", err)
	}

	if !created {
		t.Fatal("created = false, want true")
	}

	if fake.inspectCalls != 1 {
		t.Fatalf("inspect calls = %d, want 1", fake.inspectCalls)
	}

	if fake.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", fake.createCalls)
	}

	if fake.lastCreatedName != runtime.NetworkName {
		t.Fatalf(
			"created network name = %q, want %q",
			fake.lastCreatedName,
			runtime.NetworkName,
		)
	}

	if fake.lastCreateOptions.Driver != "bridge" {
		t.Fatalf(
			"created network driver = %q, want %q",
			fake.lastCreateOptions.Driver,
			"bridge",
		)
	}

	if fake.lastCreateOptions.Labels[NetworkLabelRuntimeID] != runtime.ID {
		t.Fatalf(
			"runtime ID label = %q, want %q",
			fake.lastCreateOptions.Labels[NetworkLabelRuntimeID],
			runtime.ID,
		)
	}
}

func TestEnsureDockerNetworkAcceptsOwnedNetwork(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectResults: []client.NetworkInspectResult{
			{
				Network: network.Inspect{
					Network: network.Network{
						Name: runtime.NetworkName,
						Labels: map[string]string{
							NetworkLabelRuntimeID: runtime.ID,
						},
					},
				},
			},
		},
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err != nil {
		t.Fatalf("EnsureDockerNetwork() error = %v", err)
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if fake.createCalls != 0 {
		t.Fatalf("create calls = %d, want 0", fake.createCalls)
	}
}

func TestEnsureDockerNetworkRejectsUnownedNetwork(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectResults: []client.NetworkInspectResult{
			{
				Network: network.Inspect{
					Network: network.Network{
						Name:   runtime.NetworkName,
						Labels: map[string]string{},
					},
				},
			},
		},
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err == nil {
		t.Fatal("EnsureDockerNetwork() expected an error")
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if fake.createCalls != 0 {
		t.Fatalf("create calls = %d, want 0", fake.createCalls)
	}
}

func TestEnsureDockerNetworkRejectsDifferentRuntime(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectResults: []client.NetworkInspectResult{
			{
				Network: network.Inspect{
					Network: network.Network{
						Name: runtime.NetworkName,
						Labels: map[string]string{
							NetworkLabelRuntimeID: "different-runtime",
						},
					},
				},
			},
		},
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err == nil {
		t.Fatal("EnsureDockerNetwork() expected an error")
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if fake.createCalls != 0 {
		t.Fatalf("create calls = %d, want 0", fake.createCalls)
	}
}

func TestEnsureDockerNetworkHandlesCreationRace(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectErrors: []error{
			errdefs.ErrNotFound,
			nil,
		},
		inspectResults: []client.NetworkInspectResult{
			{},
			{
				Network: network.Inspect{
					Network: network.Network{
						Name: runtime.NetworkName,
						Labels: map[string]string{
							NetworkLabelRuntimeID: runtime.ID,
						},
					},
				},
			},
		},
		createError: errdefs.ErrAlreadyExists,
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err != nil {
		t.Fatalf("EnsureDockerNetwork() error = %v", err)
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if fake.inspectCalls != 2 {
		t.Fatalf("inspect calls = %d, want 2", fake.inspectCalls)
	}

	if fake.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", fake.createCalls)
	}
}

func TestEnsureDockerNetworkRejectsCreationRaceWithDifferentOwner(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	fake := &fakeNetworkClient{
		inspectErrors: []error{
			errdefs.ErrNotFound,
			nil,
		},
		inspectResults: []client.NetworkInspectResult{
			{},
			{
				Network: network.Inspect{
					Network: network.Network{
						Name: runtime.NetworkName,
						Labels: map[string]string{
							NetworkLabelRuntimeID: "different-runtime",
						},
					},
				},
			},
		},
		createError: errdefs.ErrAlreadyExists,
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err == nil {
		t.Fatal("EnsureDockerNetwork() expected an error")
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if fake.inspectCalls != 2 {
		t.Fatalf("inspect calls = %d, want 2", fake.inspectCalls)
	}

	if fake.createCalls != 1 {
		t.Fatalf("create calls = %d, want 1", fake.createCalls)
	}
}

func TestEnsureDockerNetworkPropagatesInspectError(t *testing.T) {
	runtime := Runtime{
		ID:          "runtime-123",
		NetworkName: "radius-director-test",
	}

	expected := errors.New("docker unavailable")

	fake := &fakeNetworkClient{
		inspectErrors: []error{expected},
	}

	created, err := EnsureDockerNetwork(context.Background(), fake, runtime)
	if err == nil {
		t.Fatal("EnsureDockerNetwork() expected an error")
	}

	if created {
		t.Fatal("created = true, want false")
	}

	if !errors.Is(err, expected) {
		t.Fatalf("error = %v, want wrapped %v", err, expected)
	}

	if fake.createCalls != 0 {
		t.Fatalf("create calls = %d, want 0", fake.createCalls)
	}
}

func TestFakeNetworkClientNetworkRemove(t *testing.T) {
	fake := &fakeNetworkClient{}

	_, err := fake.NetworkRemove(
		context.Background(),
		"radius-director-test",
		client.NetworkRemoveOptions{},
	)
	if err != nil {
		t.Fatalf("NetworkRemove() error = %v", err)
	}

	if fake.removeCalls != 1 {
		t.Fatalf("remove calls = %d, want 1", fake.removeCalls)
	}

	if fake.lastRemovedNetwork != "radius-director-test" {
		t.Fatalf(
			"removed network = %q, want %q",
			fake.lastRemovedNetwork,
			"radius-director-test",
		)
	}
}
