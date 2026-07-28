package renderer

import (
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
)

func TestRenderClientsOneTenantOneClient(t *testing.T) {
	tenant := generator.Tenant{
		Identifier: "customer-a",
		Clients: []generator.Client{
			{
				Identifier:   "core-router",
				IPAddress:    "10.10.10.1",
				SharedSecret: "shared-secret",
			},
		},
	}

	got, err := RenderClients(tenant)
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
	tenant := generator.Tenant{
		Identifier: "customer-a",
		Clients: []generator.Client{
			{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "core-secret"},
			{Identifier: "edge-router", IPAddress: "10.10.10.2", SharedSecret: "edge-secret"},
		},
	}

	got, err := RenderClients(tenant)
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

func TestRenderClientsIsDeterministic(t *testing.T) {
	tenant := generator.Tenant{
		Identifier: "tenant-a",
		Clients:    []generator.Client{{Identifier: "a", IPAddress: "10.0.0.1", SharedSecret: "a-secret"}},
	}

	first, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("first RenderClients() error = %v", err)
	}
	second, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("second RenderClients() error = %v", err)
	}
	if first != second {
		t.Fatalf("RenderClients() returned different results: %q and %q", first, second)
	}
}
