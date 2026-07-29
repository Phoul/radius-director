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
				Database: testDatabase("customer_a"),
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
				Database: testDatabase("customer_a"),
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

func TestGenerateOneTenantCreatesOneSQL(t *testing.T) {
	database := model.Database{
		Engine:   "mysql",
		Host:     "database.internal",
		Port:     3306,
		Database: "customer_a",
		Username: "radius",
		Password: "database-password",
	}
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {Database: database},
		},
	}

	generated := Generate(configuration)
	if len(generated.Tenants) != 1 {
		t.Fatalf("generated tenants = %d, want 1", len(generated.Tenants))
	}
	want := SQL{
		Engine:   "mysql",
		Host:     "database.internal",
		Port:     3306,
		Database: "customer_a",
		Username: "radius",
		Password: "database-password",
	}
	if got := generated.Tenants[0].SQL; got != want {
		t.Fatalf("generated SQL = %#v, want %#v", got, want)
	}
}

func TestGenerateMultipleTenantsCreatesSQLInDeterministicOrder(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"tenant-b": {Database: testDatabase("tenant_b")},
			"tenant-a": {Database: testDatabase("tenant_a")},
		},
	}

	generated := Generate(configuration)
	want := []Tenant{
		{
			Identifier: "tenant-a",
			Clients:    []Client{},
			SQL: SQL{
				Engine:   "mysql",
				Host:     "database.tenant_a.internal",
				Port:     3306,
				Database: "tenant_a",
				Username: "radius-tenant_a",
				Password: "password-tenant_a",
			},
		},
		{
			Identifier: "tenant-b",
			Clients:    []Client{},
			SQL: SQL{
				Engine:   "mysql",
				Host:     "database.tenant_b.internal",
				Port:     3306,
				Database: "tenant_b",
				Username: "radius-tenant_b",
				Password: "password-tenant_b",
			},
		},
	}
	if got := generated.Tenants; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated tenants = %#v, want %#v", got, want)
	}
}

func TestGenerateOneTenantCreatesOneRADIUSServer(t *testing.T) {
	radiusServer := model.RADIUSServer{
		Version:            "3.2.9",
		AuthenticationPort: 1812,
		AccountingPort:     1813,
		COAPort:            3799,
	}
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {RADIUSServer: radiusServer},
		},
	}

	generated := Generate(configuration)
	if len(generated.Tenants) != 1 {
		t.Fatalf("generated tenants = %d, want 1", len(generated.Tenants))
	}
	want := RADIUSServer{
		Version:            "3.2.9",
		AuthenticationPort: 1812,
		AccountingPort:     1813,
		COAPort:            3799,
	}
	if got := generated.Tenants[0].RADIUSServer; got != want {
		t.Fatalf("generated RADIUS server = %#v, want %#v", got, want)
	}
}

func TestGenerateMultipleTenantsCreatesRADIUSServersInDeterministicOrder(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"tenant-b": {
				RADIUSServer: model.RADIUSServer{
					Version:            "3.4.0",
					AuthenticationPort: 2812,
					AccountingPort:     2813,
					COAPort:            4799,
				},
			},
			"tenant-a": {
				RADIUSServer: model.RADIUSServer{
					Version:            "3.2.9",
					AuthenticationPort: 1812,
					AccountingPort:     1813,
					COAPort:            3799,
				},
			},
		},
	}

	generated := Generate(configuration)
	if len(generated.Tenants) != 2 {
		t.Fatalf("generated tenants = %d, want 2", len(generated.Tenants))
	}
	if got, want := generated.Tenants[0].Identifier, "tenant-a"; got != want {
		t.Fatalf("first tenant identifier = %q, want %q", got, want)
	}
	if got, want := generated.Tenants[0].RADIUSServer, (RADIUSServer{
		Version:            "3.2.9",
		AuthenticationPort: 1812,
		AccountingPort:     1813,
		COAPort:            3799,
	}); got != want {
		t.Fatalf("first generated RADIUS server = %#v, want %#v", got, want)
	}
	if got, want := generated.Tenants[1].Identifier, "tenant-b"; got != want {
		t.Fatalf("second tenant identifier = %q, want %q", got, want)
	}
	if got, want := generated.Tenants[1].RADIUSServer, (RADIUSServer{
		Version:            "3.4.0",
		AuthenticationPort: 2812,
		AccountingPort:     2813,
		COAPort:            4799,
	}); got != want {
		t.Fatalf("second generated RADIUS server = %#v, want %#v", got, want)
	}
}

func testDatabase(name string) model.Database {
	return model.Database{
		Engine:   "mysql",
		Host:     "database." + name + ".internal",
		Port:     3306,
		Database: name,
		Username: "radius-" + name,
		Password: "password-" + name,
	}
}
