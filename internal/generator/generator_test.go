package generator

import (
	"reflect"
	"testing"
	"time"

	"github.com/gobcn/radius-director/internal/model"
)

func TestGenerateOneNASAssignmentCreatesOneClient(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"default": {SharedSecret: "shared-secret"},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1", Vendor: "mikrotik"},
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
	if got, want := generated.Tenants[0].FreeRADIUSClients, []FreeRADIUSClient{{
		Identifier:   "core-router",
		IPAddress:    "10.10.10.1",
		SharedSecret: "shared-secret",
		Vendor:       "mikrotik",
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated FreeRADIUS clients = %#v, want %#v", got, want)
	}
	if got, want := generated.Tenants[0].HomeServers, []HomeServer{{
		Identifier:   "core-router",
		IPAddress:    "10.10.10.1",
		SharedSecret: "shared-secret",
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated home servers = %#v, want %#v", got, want)
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
				"core": {IPAddress: "10.10.10.1", Vendor: "mikrotik"},
				"edge": {IPAddress: "10.10.10.2", Vendor: "generic"},
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
	if got, want := generated.Tenants[0].FreeRADIUSClients, []FreeRADIUSClient{
		{
			Identifier:   "core-router",
			IPAddress:    "10.10.10.1",
			SharedSecret: "core-secret",
			Vendor:       "mikrotik",
		},
		{
			Identifier:   "edge-router",
			IPAddress:    "10.10.10.2",
			SharedSecret: "edge-secret",
			Vendor:       "generic",
		},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated FreeRADIUS clients = %#v, want %#v", got, want)
	}
	if got, want := generated.Tenants[0].HomeServers, []HomeServer{
		{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "core-secret"},
		{Identifier: "edge-router", IPAddress: "10.10.10.2", SharedSecret: "edge-secret"},
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated home servers = %#v, want %#v", got, want)
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
		Host:     "database-customer-a",
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
			Identifier:         "tenant-a",
			FreeRADIUSClients:  []FreeRADIUSClient{},
			HomeServers:        []HomeServer{},
			AccountingPolicies: []NASAccountingPolicy{},
			DatabaseDeployment: "container",
			SQL: SQL{
				Engine:   "mysql",
				Host:     "database-tenant-a",
				Port:     3306,
				Database: "tenant_a",
				Username: "radius-tenant_a",
				Password: "password-tenant_a",
			},
		},
		{
			Identifier:         "tenant-b",
			FreeRADIUSClients:  []FreeRADIUSClient{},
			HomeServers:        []HomeServer{},
			AccountingPolicies: []NASAccountingPolicy{},
			DatabaseDeployment: "container",
			SQL: SQL{
				Engine:   "mysql",
				Host:     "database-tenant-b",
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

func TestGenerateDefaultsDatabaseDeploymentToContainer(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {
				Database: model.Database{
					Engine:   "mysql",
					Host:     "database.internal",
					Port:     3306,
					Database: "customer_a",
					Username: "radius",
					Password: "database-password",
				},
			},
		},
	}

	generated := Generate(configuration)

	if got, want := generated.Tenants[0].DatabaseDeployment, "container"; got != want {
		t.Fatalf("database deployment = %q, want %q", got, want)
	}
}

func TestGeneratePreservesExternalDatabaseDeployment(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {
				Database: model.Database{
					Engine:     "mysql",
					Deployment: "external",
					Host:       "database.example.com",
					Port:       3306,
					Database:   "customer_a",
					Username:   "radius",
					Password:   "database-password",
				},
			},
		},
	}

	generated := Generate(configuration)

	if got, want := generated.Tenants[0].DatabaseDeployment, "external"; got != want {
		t.Fatalf("database deployment = %q, want %q", got, want)
	}
}

func TestGenerateContainerDatabaseEndpoint(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {
				Database: model.Database{
					Engine:   "mysql",
					Database: "radius",
					Username: "radius",
					Password: "secret",
				},
			},
		},
	}

	generated := Generate(configuration)

	sql := generated.Tenants[0].SQL

	if got, want := sql.Host, "database-customer-a"; got != want {
		t.Fatalf("SQL.Host = %q, want %q", got, want)
	}

	if got, want := sql.Port, 3306; got != want {
		t.Fatalf("SQL.Port = %d, want %d", got, want)
	}
}

func TestGenerateExternalDatabaseEndpoint(t *testing.T) {
	configuration := model.Configuration{
		Tenants: map[string]model.Tenant{
			"customer-a": {
				Database: model.Database{
					Engine:     "mysql",
					Deployment: "external",
					Host:       "db.example.com",
					Port:       3307,
					Database:   "radius",
					Username:   "radius",
					Password:   "secret",
				},
			},
		},
	}

	generated := Generate(configuration)

	sql := generated.Tenants[0].SQL

	if got, want := sql.Host, "db.example.com"; got != want {
		t.Fatalf("SQL.Host = %q, want %q", got, want)
	}

	if got, want := sql.Port, 3307; got != want {
		t.Fatalf("SQL.Port = %d, want %d", got, want)
	}
}

func TestGenerateTrustedRADIUSClientAssignments(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"core-credentials":         {SharedSecret: "core-secret"},
				"monitoring-credentials":   {SharedSecret: "monitoring-secret"},
				"provisioning-credentials": {SharedSecret: "provisioning-secret"},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1", Vendor: "mikrotik"},
			},
			TrustedRADIUSClients: map[string]model.TrustedRADIUSClient{
				"monitoring":   {IPAddress: "10.10.10.10"},
				"provisioning": {IPAddress: "10.10.10.11"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				NASAssignments: map[string]model.NASAssignment{
					"core": {NASDevice: "core", CredentialProfile: "core-credentials"},
				},
				TrustedRADIUSClientAssignments: map[string]model.TrustedRADIUSClientAssignment{
					"provisioning": {
						TrustedRADIUSClient: "provisioning",
						CredentialProfile:   "provisioning-credentials",
					},
					"monitoring": {
						TrustedRADIUSClient: "monitoring",
						CredentialProfile:   "monitoring-credentials",
					},
				},
			},
		},
	}

	generated := Generate(configuration)
	want := []FreeRADIUSClient{
		{Identifier: "core", IPAddress: "10.10.10.1", SharedSecret: "core-secret", Vendor: "mikrotik"},
		{Identifier: "monitoring", IPAddress: "10.10.10.10", SharedSecret: "monitoring-secret"},
		{Identifier: "provisioning", IPAddress: "10.10.10.11", SharedSecret: "provisioning-secret"},
	}
	if got := generated.Tenants[0].FreeRADIUSClients; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated FreeRADIUS clients = %#v, want %#v", got, want)
	}
	if got, want := generated.Tenants[0].HomeServers, []HomeServer{{
		Identifier:   "core",
		IPAddress:    "10.10.10.1",
		SharedSecret: "core-secret",
	}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated home servers = %#v, want %#v", got, want)
	}
}

func TestGenerateOneTenantCreatesOneRADIUSServer(t *testing.T) {
	radiusServer := model.RADIUSServer{
		Version:            "3.2.10",
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
		Version:            "3.2.10",
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
					Version:            "3.2.10",
					AuthenticationPort: 2812,
					AccountingPort:     2813,
					COAPort:            4799,
				},
			},
			"tenant-a": {
				RADIUSServer: model.RADIUSServer{
					Version:            "3.2.10",
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
		Version:            "3.2.10",
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
		Version:            "3.2.10",
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

func TestGenerateNASAccountingPolicies(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			AccountingProfiles: map[string]model.AccountingProfile{
				"disabled": {},
				"long":     {StaleSessionTimeout: "1h"},
				"standard": {StaleSessionTimeout: "20m"},
			},
			NASDevices: map[string]model.NASDevice{
				"alpha": {IPAddress: "10.10.10.1"},
				"beta":  {IPAddress: "10.10.10.2"},
				"gamma": {IPAddress: "10.10.10.3"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				NASAssignments: map[string]model.NASAssignment{
					"gamma-router": {NASDevice: "gamma", AccountingProfile: "disabled"},
					"alpha-router": {NASDevice: "alpha", AccountingProfile: "standard"},
					"beta-router":  {NASDevice: "beta", AccountingProfile: "long"},
				},
			},
		},
	}

	generated := Generate(configuration)
	standard := 20 * time.Minute
	long := time.Hour
	want := []NASAccountingPolicy{
		{NASAssignmentIdentifier: "alpha-router", NASDeviceIdentifier: "alpha", IPAddress: "10.10.10.1", StaleSessionTimeout: &standard},
		{NASAssignmentIdentifier: "beta-router", NASDeviceIdentifier: "beta", IPAddress: "10.10.10.2", StaleSessionTimeout: &long},
		{NASAssignmentIdentifier: "gamma-router", NASDeviceIdentifier: "gamma", IPAddress: "10.10.10.3", StaleSessionTimeout: nil},
	}
	if got := generated.Tenants[0].AccountingPolicies; !reflect.DeepEqual(got, want) {
		t.Fatalf("generated accounting policies = %#v, want %#v", got, want)
	}
}

func TestGenerateTrustedRADIUSClientsDoNotCreateAccountingPolicies(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			TrustedRADIUSClients: map[string]model.TrustedRADIUSClient{
				"sonar": {IPAddress: "20.104.33.4"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				TrustedRADIUSClientAssignments: map[string]model.TrustedRADIUSClientAssignment{
					"sonar": {TrustedRADIUSClient: "sonar"},
				},
			},
		},
	}

	generated := Generate(configuration)
	if got := generated.Tenants[0].AccountingPolicies; len(got) != 0 {
		t.Fatalf("generated accounting policies = %#v, want none", got)
	}
}

func TestGenerateResolvesTenantAuthenticationPolicy(t *testing.T) {
	one := 1
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{AuthenticationProfiles: map[string]model.AuthenticationProfile{"standard": {SimultaneousUse: &one}}},
		Tenants:       map[string]model.Tenant{"customer-a": {AuthenticationProfile: "standard"}},
	}
	generated := Generate(configuration)
	got := generated.Tenants[0].AuthenticationPolicy.SimultaneousUse

	if got == nil {
		t.Fatal("simultaneous use = nil, want 1")
	}

	if *got != 1 {
		t.Fatalf("simultaneous use = %d, want 1", *got)
	}
}

func TestGenerateResolvesDeploymentProfileTemplate(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			DeploymentProfiles: map[string]model.DeploymentProfile{
				"experimental": {
					Template: "alternate",
					Overlays: []string{
						"coa-relay-test",
						"debug-logging",
					},
				},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				DeploymentProfile: "experimental",
			},
		},
	}

	generated := Generate(configuration)

	if got, want := generated.Tenants[0].Template, "alternate"; got != want {
		t.Fatalf("tenant template = %q, want %q", got, want)
	}

	if got := len(generated.Tenants[0].Overlays); got != 2 {
		t.Fatalf("tenant overlay count = %d, want %d", got, 2)
	}

	if got := generated.Tenants[0].Overlays[0]; got != "coa-relay-test" {
		t.Fatalf("first tenant overlay = %q, want %q", got, "coa-relay-test")
	}

	if got := generated.Tenants[0].Overlays[1]; got != "debug-logging" {
		t.Fatalf("second tenant overlay = %q, want %q", got, "debug-logging")
	}
}
