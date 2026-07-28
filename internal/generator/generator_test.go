package generator

import (
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/model"
)

func TestGenerateOneNASAssignmentCreatesOneClient(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"default": {SharedSecret: "shared-secret"},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				NASAssignments: map[string]model.NASAssignment{
					"core-router": {
						NASDevice:         "core",
						CredentialProfile: "default",
					},
				},
			},
		},
	}

	generated := Generate(configuration)
	if len(generated.Tenants) != 1 {
		t.Fatalf("generated tenants = %d, want 1", len(generated.Tenants))
	}
	if got, want := generated.Tenants[0].Identifier, "customer-a"; got != want {
		t.Fatalf("tenant identifier = %q, want %q", got, want)
	}
	if got, want := generated.Tenants[0].Clients, []Client{{
		Identifier:   "core-router",
		IPAddress:    "10.10.10.1",
		SharedSecret: "shared-secret",
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated clients = %#v, want %#v", got, want)
	}
}

func TestGenerateMultipleNASAssignmentsCreatesMultipleClients(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"core-credentials": {SharedSecret: "core-secret"},
				"edge-credentials": {SharedSecret: "edge-secret"},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1"},
				"edge": {IPAddress: "10.10.10.2"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				NASAssignments: map[string]model.NASAssignment{
					"edge-router": {
						NASDevice:         "edge",
						CredentialProfile: "edge-credentials",
					},
					"core-router": {
						NASDevice:         "core",
						CredentialProfile: "core-credentials",
					},
				},
			},
		},
	}

	generated := Generate(configuration)
	if len(generated.Tenants) != 1 {
		t.Fatalf("generated tenants = %d, want 1", len(generated.Tenants))
	}
	if got, want := generated.Tenants[0].Clients, []Client{
		{
			Identifier:   "core-router",
			IPAddress:    "10.10.10.1",
			SharedSecret: "core-secret",
		},
		{
			Identifier:   "edge-router",
			IPAddress:    "10.10.10.2",
			SharedSecret: "edge-secret",
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated clients = %#v, want %#v", got, want)
	}
}
