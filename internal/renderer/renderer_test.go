package renderer

import (
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
)

func TestRenderClientsOneTenantOneClient(t *testing.T) {
	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier: "customer-a",
				Clients: []generator.Client{
					{
						Identifier:   "core-router",
						IPAddress:    "10.10.10.1",
						SharedSecret: "shared-secret",
					},
				},
			},
		},
	}

	got, err := RenderClients(configuration)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	want := "# Tenant: customer-a\n" +
		"client core-router {\n" +
		"    ipaddr = 10.10.10.1\n" +
		"    secret = shared-secret\n" +
		"}\n"
	if got != want {
		t.Fatalf("RenderClients() = %q, want %q", got, want)
	}
}

func TestRenderClientsOneTenantMultipleClients(t *testing.T) {
	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier: "customer-a",
				Clients: []generator.Client{
					{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "core-secret"},
					{Identifier: "edge-router", IPAddress: "10.10.10.2", SharedSecret: "edge-secret"},
				},
			},
		},
	}

	got, err := RenderClients(configuration)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	want := "# Tenant: customer-a\n" +
		"client core-router {\n" +
		"    ipaddr = 10.10.10.1\n" +
		"    secret = core-secret\n" +
		"}\n" +
		"client edge-router {\n" +
		"    ipaddr = 10.10.10.2\n" +
		"    secret = edge-secret\n" +
		"}\n"
	if got != want {
		t.Fatalf("RenderClients() = %q, want %q", got, want)
	}
}

func TestRenderClientsMultipleTenants(t *testing.T) {
	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier: "customer-a",
				Clients: []generator.Client{
					{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "core-secret"},
				},
			},
			{
				Identifier: "customer-b",
				Clients: []generator.Client{
					{Identifier: "edge-router", IPAddress: "10.20.20.1", SharedSecret: "edge-secret"},
				},
			},
		},
	}

	got, err := RenderClients(configuration)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	want := "# Tenant: customer-a\n" +
		"client core-router {\n" +
		"    ipaddr = 10.10.10.1\n" +
		"    secret = core-secret\n" +
		"}\n\n" +
		"# Tenant: customer-b\n" +
		"client edge-router {\n" +
		"    ipaddr = 10.20.20.1\n" +
		"    secret = edge-secret\n" +
		"}\n"
	if got != want {
		t.Fatalf("RenderClients() = %q, want %q", got, want)
	}
}

func TestRenderClientsIsDeterministic(t *testing.T) {
	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{Identifier: "tenant-a", Clients: []generator.Client{{Identifier: "a", IPAddress: "10.0.0.1", SharedSecret: "a-secret"}}},
			{Identifier: "tenant-b", Clients: []generator.Client{{Identifier: "b", IPAddress: "10.0.0.2", SharedSecret: "b-secret"}}},
		},
	}

	first, err := RenderClients(configuration)
	if err != nil {
		t.Fatalf("first RenderClients() error = %v", err)
	}
	second, err := RenderClients(configuration)
	if err != nil {
		t.Fatalf("second RenderClients() error = %v", err)
	}
	if first != second {
		t.Fatalf("RenderClients() returned different results: %q and %q", first, second)
	}
}
