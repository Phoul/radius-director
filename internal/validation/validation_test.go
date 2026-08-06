package validation

import (
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/model"
)

func TestValidate(t *testing.T) {
	one := 1
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"default": {SharedSecret: "secret"},
			},
			AuthenticationProfiles: map[string]model.AuthenticationProfile{
				"default": {SimultaneousUse: &one},
			},
			AccountingProfiles: map[string]model.AccountingProfile{
				"default": {},
			},
			MonitoringProfiles: map[string]model.MonitoringProfile{
				"default": {},
			},
			DeploymentProfiles: map[string]model.DeploymentProfile{
				"default": {},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1", Vendor: "mikrotik"},
			},
			TrustedRADIUSClients: map[string]model.TrustedRADIUSClient{
				"monitoring": {IPAddress: "10.10.10.2"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				AuthenticationProfile: "default",
				DeploymentProfile:     "default",
				Database: model.Database{
					Engine:   "mysql",
					Host:     "db.example.com",
					Port:     3306,
					Database: "radius",
					Username: "radius",
					Password: "secret",
				},
				RADIUSServer: model.RADIUSServer{
					Version:            "3.2.10",
					AuthenticationPort: 1812,
					AccountingPort:     1813,
					COAPort:            3799,
				},
				NASAssignments: map[string]model.NASAssignment{
					"core": {
						NASDevice:         "core",
						CredentialProfile: "default",
						AccountingProfile: "default",
						MonitoringProfile: "default",
					},
				},
			},
		},
	}

	if err := Validate(configuration); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidateNASAssignment(t *testing.T) {
	validAssignment := model.NASAssignment{
		NASDevice:         "core",
		CredentialProfile: "default",
		AccountingProfile: "default",
		MonitoringProfile: "default",
	}

	tests := []struct {
		name       string
		assignment model.NASAssignment
		wantErrs   []string
	}{
		{
			name:       "valid NAS Assignment",
			assignment: validAssignment,
		},
		{
			name: "NAS Device missing",
			assignment: model.NASAssignment{
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas_device must be specified`,
			},
		},
		{
			name: "credential profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": credential_profile must be specified`,
			},
		},
		{
			name: "accounting profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: validAssignment.CredentialProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": accounting_profile must be specified`,
			},
		},
		{
			name: "monitoring profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: validAssignment.AccountingProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": monitoring_profile must be specified`,
			},
		},
		{
			name: "multiple properties missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": credential_profile must be specified`,
				`tenant "customer-a": nas assignment "core": accounting_profile must be specified`,
			},
		},
		{
			name: "all properties missing",
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas_device must be specified`,
				`tenant "customer-a": nas assignment "core": credential_profile must be specified`,
				`tenant "customer-a": nas assignment "core": accounting_profile must be specified`,
				`tenant "customer-a": nas assignment "core": monitoring_profile must be specified`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateNASAssignment("customer-a", "core", test.assignment)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateNASAssignment() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateNASAssignment() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateNASAssignmentReferences(t *testing.T) {
	one := 1
	globalObjects := model.GlobalObjects{
		CredentialProfiles: map[string]model.CredentialProfile{
			"default": {},
		},
		AuthenticationProfiles: map[string]model.AuthenticationProfile{
			"default": {SimultaneousUse: &one},
		},
		AccountingProfiles: map[string]model.AccountingProfile{
			"default": {},
		},
		MonitoringProfiles: map[string]model.MonitoringProfile{
			"default": {},
		},
		NASDevices: map[string]model.NASDevice{
			"core": {},
		},
	}
	validAssignment := model.NASAssignment{
		NASDevice:         "core",
		CredentialProfile: "default",
		AccountingProfile: "default",
		MonitoringProfile: "default",
	}

	tests := []struct {
		name       string
		assignment model.NASAssignment
		wantErrs   []string
	}{
		{
			name:       "valid references",
			assignment: validAssignment,
		},
		{
			name: "NAS Device missing",
			assignment: model.NASAssignment{
				NASDevice:         "missing-device",
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas device "missing-device" does not exist`,
			},
		},
		{
			name: "credential profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: "missing-credential-profile",
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": credential profile "missing-credential-profile" does not exist`,
			},
		},
		{
			name: "accounting profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: "missing-accounting-profile",
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": accounting profile "missing-accounting-profile" does not exist`,
			},
		},
		{
			name: "monitoring profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: "missing-monitoring-profile",
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": monitoring profile "missing-monitoring-profile" does not exist`,
			},
		},
		{
			name: "multiple references missing",
			assignment: model.NASAssignment{
				NASDevice:         "missing-device",
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: "missing-accounting-profile",
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas device "missing-device" does not exist`,
				`tenant "customer-a": nas assignment "core": accounting profile "missing-accounting-profile" does not exist`,
			},
		},
		{
			name: "all references missing",
			assignment: model.NASAssignment{
				NASDevice:         "missing-device",
				CredentialProfile: "missing-credential-profile",
				AccountingProfile: "missing-accounting-profile",
				MonitoringProfile: "missing-monitoring-profile",
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas device "missing-device" does not exist`,
				`tenant "customer-a": nas assignment "core": credential profile "missing-credential-profile" does not exist`,
				`tenant "customer-a": nas assignment "core": accounting profile "missing-accounting-profile" does not exist`,
				`tenant "customer-a": nas assignment "core": monitoring profile "missing-monitoring-profile" does not exist`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateNASAssignmentReferences("customer-a", "core", test.assignment, globalObjects)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateNASAssignmentReferences() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateNASAssignmentReferences() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTenantRelationships(t *testing.T) {
	globalObjects := model.GlobalObjects{
		NASDevices: map[string]model.NASDevice{
			"core":    {},
			"edge":    {},
			"gateway": {},
		},
	}

	tests := []struct {
		name     string
		tenant   model.Tenant
		wantErrs []string
	}{
		{
			name: "no duplicates",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "core"},
					"assignment-b": {NASDevice: "edge"},
				},
			},
		},
		{
			name: "one duplicate pair",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "core"},
					"assignment-b": {NASDevice: "core"},
				},
			},
			wantErrs: []string{
				`tenant "customer-a": nas device "core" is assigned by both nas assignments "assignment-a" and "assignment-b"`,
			},
		},
		{
			name: "multiple independent duplicates",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "core"},
					"assignment-b": {NASDevice: "core"},
					"assignment-c": {NASDevice: "edge"},
					"assignment-d": {NASDevice: "edge"},
				},
			},
			wantErrs: []string{
				`tenant "customer-a": nas device "core" is assigned by both nas assignments "assignment-a" and "assignment-b"`,
				`tenant "customer-a": nas device "edge" is assigned by both nas assignments "assignment-c" and "assignment-d"`,
			},
		},
		{
			name: "all duplicate assignments reported",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "core"},
					"assignment-b": {NASDevice: "core"},
					"assignment-c": {NASDevice: "core"},
				},
			},
			wantErrs: []string{
				`tenant "customer-a": nas device "core" is assigned by both nas assignments "assignment-a" and "assignment-b"`,
				`tenant "customer-a": nas device "core" is assigned by both nas assignments "assignment-a" and "assignment-c"`,
			},
		},
		{
			name: "duplicate mixed with valid assignments",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "core"},
					"assignment-b": {NASDevice: "core"},
					"assignment-c": {NASDevice: "gateway"},
				},
			},
			wantErrs: []string{
				`tenant "customer-a": nas device "core" is assigned by both nas assignments "assignment-a" and "assignment-b"`,
			},
		},
		{
			name: "missing NAS Device does not produce relationship errors",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {},
					"assignment-b": {},
				},
			},
		},
		{
			name: "nonexistent NAS Device does not produce relationship errors",
			tenant: model.Tenant{
				NASAssignments: map[string]model.NASAssignment{
					"assignment-a": {NASDevice: "missing-device"},
					"assignment-b": {NASDevice: "missing-device"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateTenantRelationships("customer-a", test.tenant, globalObjects)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateTenantRelationships() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateTenantRelationships() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTenant(t *testing.T) {
	validTenant := model.Tenant{
		AuthenticationProfile: "default",
		DeploymentProfile:     "default",
		Database: model.Database{
			Engine:   "mysql",
			Host:     "db.example.com",
			Port:     3306,
			Database: "radius",
			Username: "radius",
			Password: "secret",
		},
		RADIUSServer: model.RADIUSServer{
			Version:            "3.2.10",
			AuthenticationPort: 1812,
			AccountingPort:     1813,
			COAPort:            3799,
		},
		NASAssignments: map[string]model.NASAssignment{
			"core": {
				NASDevice:         "core",
				CredentialProfile: "default",
				AccountingProfile: "default",
				MonitoringProfile: "default",
			},
		},
	}

	tests := []struct {
		name     string
		tenant   model.Tenant
		wantErrs []string
	}{
		{
			name:   "valid tenant",
			tenant: validTenant,
		},
		{
			name: "database missing",
			tenant: model.Tenant{
				AuthenticationProfile: validTenant.AuthenticationProfile,
				DeploymentProfile:     validTenant.DeploymentProfile,
				RADIUSServer:          validTenant.RADIUSServer,
				NASAssignments:        validTenant.NASAssignments,
			},
			wantErrs: []string{
				`tenant "customer-a": exactly one database must be defined`,
			},
		},
		{
			name: "RADIUS Server missing",
			tenant: model.Tenant{
				AuthenticationProfile: validTenant.AuthenticationProfile,
				DeploymentProfile:     validTenant.DeploymentProfile,
				Database:              validTenant.Database,
				NASAssignments:        validTenant.NASAssignments,
			},
			wantErrs: []string{
				`tenant "customer-a": exactly one radius server must be defined`,
			},
		},
		{
			name: "NAS assignments missing",
			tenant: model.Tenant{
				AuthenticationProfile: validTenant.AuthenticationProfile,
				DeploymentProfile:     validTenant.DeploymentProfile,
				Database:              validTenant.Database,
				RADIUSServer:          validTenant.RADIUSServer,
			},
			wantErrs: []string{
				`tenant "customer-a": at least one nas assignment must be defined`,
			},
		},
		{
			name: "deployment profile missing",
			tenant: model.Tenant{
				AuthenticationProfile: validTenant.AuthenticationProfile,
				Database:              validTenant.Database,
				RADIUSServer:          validTenant.RADIUSServer,
				NASAssignments:        validTenant.NASAssignments,
			},
			wantErrs: []string{
				`tenant "customer-a": deployment_profile must be specified`,
			},
		},
		{
			name: "all required tenant objects missing",
			wantErrs: []string{
				`tenant "customer-a": authentication_profile must be specified`,
				`tenant "customer-a": deployment_profile must be specified`,
				`tenant "customer-a": exactly one database must be defined`,
				`tenant "customer-a": exactly one radius server must be defined`,
				`tenant "customer-a": at least one nas assignment must be defined`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateTenant("customer-a", test.tenant)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateTenant() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateTenant() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateRADIUSServer(t *testing.T) {
	validServer := model.RADIUSServer{
		Version:            "3.2.10",
		AuthenticationPort: 1812,
		AccountingPort:     1813,
		COAPort:            3799,
	}

	tests := []struct {
		name     string
		server   model.RADIUSServer
		wantErrs []string
	}{
		{
			name:   "valid RADIUS Server",
			server: validServer,
		},
		{
			name: "version missing",
			server: model.RADIUSServer{
				AuthenticationPort: 1812,
				AccountingPort:     1813,
				COAPort:            3799,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server version must be specified`,
			},
		},
		{
			name: "version unsupported",
			server: model.RADIUSServer{
				Version:            "3.3.0",
				AuthenticationPort: 1812,
				AccountingPort:     1813,
				COAPort:            3799,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server version "3.3.0" is not supported`,
			},
		},
		{
			name: "minimum valid ports",
			server: model.RADIUSServer{
				Version:            "3.2.10",
				AuthenticationPort: 1,
				AccountingPort:     1,
				COAPort:            1,
			},
		},
		{
			name: "maximum valid ports",
			server: model.RADIUSServer{
				Version:            "3.2.10",
				AuthenticationPort: 65535,
				AccountingPort:     65535,
				COAPort:            65535,
			},
		},
		{
			name: "authentication port missing",
			server: model.RADIUSServer{
				Version:        validServer.Version,
				AccountingPort: validServer.AccountingPort,
				COAPort:        validServer.COAPort,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server authentication_port must be between 1 and 65535`,
			},
		},
		{
			name: "authentication port above range",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: 65536,
				AccountingPort:     validServer.AccountingPort,
				COAPort:            validServer.COAPort,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server authentication_port must be between 1 and 65535`,
			},
		},
		{
			name: "accounting port missing",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: validServer.AuthenticationPort,
				COAPort:            validServer.COAPort,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server accounting_port must be between 1 and 65535`,
			},
		},
		{
			name: "accounting port above range",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: validServer.AuthenticationPort,
				AccountingPort:     65536,
				COAPort:            validServer.COAPort,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server accounting_port must be between 1 and 65535`,
			},
		},
		{
			name: "CoA port missing",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: validServer.AuthenticationPort,
				AccountingPort:     validServer.AccountingPort,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server coa_port must be between 1 and 65535`,
			},
		},
		{
			name: "CoA port above range",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: validServer.AuthenticationPort,
				AccountingPort:     validServer.AccountingPort,
				COAPort:            65536,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server coa_port must be between 1 and 65535`,
			},
		},
		{
			name: "multiple invalid ports",
			server: model.RADIUSServer{
				Version:            validServer.Version,
				AuthenticationPort: 0,
				AccountingPort:     65536,
				COAPort:            -1,
			},
			wantErrs: []string{
				`tenant "customer-a": radius server authentication_port must be between 1 and 65535`,
				`tenant "customer-a": radius server accounting_port must be between 1 and 65535`,
				`tenant "customer-a": radius server coa_port must be between 1 and 65535`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateRADIUSServer("customer-a", test.server)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateRADIUSServer() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateRADIUSServer() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateDatabase(t *testing.T) {
	validDatabase := model.Database{
		Engine:   "mysql",
		Host:     "db.example.com",
		Port:     3306,
		Database: "radius",
		Username: "radius",
		Password: "secret",
	}

	tests := []struct {
		name     string
		database model.Database
		wantErrs []string
	}{
		{
			name:     "valid database",
			database: validDatabase,
		},
		{
			name: "minimum valid port",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     1,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
		},
		{
			name: "maximum valid port",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     65535,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
		},
		{
			name: "engine missing",
			database: model.Database{
				Host:     validDatabase.Host,
				Port:     validDatabase.Port,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database engine must be specified`,
			},
		},
		{
			name: "engine unsupported",
			database: model.Database{
				Engine:   "postgresql",
				Host:     validDatabase.Host,
				Port:     validDatabase.Port,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database engine "postgresql" is not supported`,
			},
		},
		{
			name: "host missing",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Port:     validDatabase.Port,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database host must be specified`,
			},
		},
		{
			name: "port below range",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     0,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database port must be between 1 and 65535`,
			},
		},
		{
			name: "port above range",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     65536,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database port must be between 1 and 65535`,
			},
		},
		{
			name: "database name missing",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     validDatabase.Port,
				Username: validDatabase.Username,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database name must be specified`,
			},
		},
		{
			name: "username missing",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     validDatabase.Port,
				Database: validDatabase.Database,
				Password: validDatabase.Password,
			},
			wantErrs: []string{
				`tenant "customer-a": database username must be specified`,
			},
		},
		{
			name: "password missing",
			database: model.Database{
				Engine:   validDatabase.Engine,
				Host:     validDatabase.Host,
				Port:     validDatabase.Port,
				Database: validDatabase.Database,
				Username: validDatabase.Username,
			},
			wantErrs: []string{
				`tenant "customer-a": database password must be specified`,
			},
		},
		{
			name: "multiple invalid properties",
			wantErrs: []string{
				`tenant "customer-a": database engine must be specified`,
				`tenant "customer-a": database host must be specified`,
				`tenant "customer-a": database port must be between 1 and 65535`,
				`tenant "customer-a": database name must be specified`,
				`tenant "customer-a": database username must be specified`,
				`tenant "customer-a": database password must be specified`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateDatabase("customer-a", test.database)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateDatabase() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateDatabase() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateNASDevice(t *testing.T) {
	tests := []struct {
		name     string
		device   model.NASDevice
		wantErrs []string
	}{
		{
			name:   "IPv4 address and vendor specified",
			device: model.NASDevice{IPAddress: "10.10.10.1", Vendor: "mikrotik"},
		},
		{
			name:   "IPv6 address and vendor specified",
			device: model.NASDevice{IPAddress: "2001:db8::1", Vendor: "mikrotik"},
		},
		{
			name:   "IP address missing",
			device: model.NASDevice{Vendor: "mikrotik"},
			wantErrs: []string{
				`nas device "core": ip_address must be a valid IPv4 or IPv6 address`,
			},
		},
		{
			name:   "IP address invalid",
			device: model.NASDevice{IPAddress: "not-an-ip", Vendor: "mikrotik"},
			wantErrs: []string{
				`nas device "core": ip_address must be a valid IPv4 or IPv6 address`,
			},
		},
		{
			name:   "vendor missing",
			device: model.NASDevice{IPAddress: "10.10.10.1"},
			wantErrs: []string{
				`nas device "core": vendor must be specified`,
			},
		},
		{
			name:   "IP address and vendor missing",
			device: model.NASDevice{},
			wantErrs: []string{
				`nas device "core": ip_address must be a valid IPv4 or IPv6 address`,
				`nas device "core": vendor must be specified`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateNASDevice("core", test.device)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateNASDevice() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateNASDevice() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTrustedRADIUSClient(t *testing.T) {
	tests := []struct {
		name     string
		client   model.TrustedRADIUSClient
		wantErrs []string
	}{
		{
			name:   "IPv4 address specified",
			client: model.TrustedRADIUSClient{IPAddress: "10.10.10.1"},
		},
		{
			name:   "IPv6 address specified",
			client: model.TrustedRADIUSClient{IPAddress: "2001:db8::1"},
		},
		{
			name:   "IP address missing",
			client: model.TrustedRADIUSClient{},
			wantErrs: []string{
				`trusted radius client "monitoring": ip_address must be a valid IPv4 or IPv6 address`,
			},
		},
		{
			name:   "IP address invalid",
			client: model.TrustedRADIUSClient{IPAddress: "not-an-ip"},
			wantErrs: []string{
				`trusted radius client "monitoring": ip_address must be a valid IPv4 or IPv6 address`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateTrustedRADIUSClient("monitoring", test.client)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateTrustedRADIUSClient() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateTrustedRADIUSClient() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTrustedRADIUSClientAssignment(t *testing.T) {
	tests := []struct {
		name       string
		assignment model.TrustedRADIUSClientAssignment
		wantErrs   []string
	}{
		{
			name: "valid assignment",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "monitoring",
				CredentialProfile:   "default",
			},
		},
		{
			name: "trusted RADIUS client missing",
			assignment: model.TrustedRADIUSClientAssignment{
				CredentialProfile: "default",
			},
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": trusted_radius_client must be specified`,
			},
		},
		{
			name: "credential profile missing",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "monitoring",
			},
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": credential_profile must be specified`,
			},
		},
		{
			name: "all properties missing",
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": trusted_radius_client must be specified`,
				`tenant "customer-a": trusted radius client assignment "monitoring": credential_profile must be specified`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateTrustedRADIUSClientAssignment("customer-a", "monitoring", test.assignment)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateTrustedRADIUSClientAssignment() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateTrustedRADIUSClientAssignment() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTrustedRADIUSClientAssignmentReferences(t *testing.T) {
	globalObjects := model.GlobalObjects{
		CredentialProfiles: map[string]model.CredentialProfile{"default": {}},
		TrustedRADIUSClients: map[string]model.TrustedRADIUSClient{
			"monitoring": {},
		},
	}

	tests := []struct {
		name       string
		assignment model.TrustedRADIUSClientAssignment
		wantErrs   []string
	}{
		{
			name: "valid references",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "monitoring",
				CredentialProfile:   "default",
			},
		},
		{
			name: "trusted RADIUS client missing",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "missing-client",
				CredentialProfile:   "default",
			},
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": trusted radius client "missing-client" does not exist`,
			},
		},
		{
			name: "credential profile missing",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "monitoring",
				CredentialProfile:   "missing-credentials",
			},
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": credential profile "missing-credentials" does not exist`,
			},
		},
		{
			name: "all references missing",
			assignment: model.TrustedRADIUSClientAssignment{
				TrustedRADIUSClient: "missing-client",
				CredentialProfile:   "missing-credentials",
			},
			wantErrs: []string{
				`tenant "customer-a": trusted radius client assignment "monitoring": trusted radius client "missing-client" does not exist`,
				`tenant "customer-a": trusted radius client assignment "monitoring": credential profile "missing-credentials" does not exist`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors := validateTrustedRADIUSClientAssignmentReferences("customer-a", "monitoring", test.assignment, globalObjects)
			if len(validationErrors) != len(test.wantErrs) {
				t.Fatalf("validateTrustedRADIUSClientAssignmentReferences() returned %d errors, want %d", len(validationErrors), len(test.wantErrs))
			}

			for index, wantErr := range test.wantErrs {
				if got := validationErrors[index].Error(); got != wantErr {
					t.Errorf("validateTrustedRADIUSClientAssignmentReferences() error = %q, want %q", got, wantErr)
				}
			}
		})
	}
}

func TestValidateTrustedRADIUSClientAssignmentRelationships(t *testing.T) {
	globalObjects := model.GlobalObjects{
		TrustedRADIUSClients: map[string]model.TrustedRADIUSClient{
			"monitoring":   {},
			"provisioning": {},
		},
	}
	tenant := model.Tenant{
		TrustedRADIUSClientAssignments: map[string]model.TrustedRADIUSClientAssignment{
			"assignment-a": {TrustedRADIUSClient: "monitoring"},
			"assignment-b": {TrustedRADIUSClient: "monitoring"},
			"assignment-c": {TrustedRADIUSClient: "provisioning"},
			"assignment-d": {TrustedRADIUSClient: "missing-client"},
		},
	}

	validationErrors := validateTenantRelationships("customer-a", tenant, globalObjects)
	if len(validationErrors) != 1 {
		t.Fatalf("validateTenantRelationships() returned %d errors, want 1", len(validationErrors))
	}
	if got, want := validationErrors[0].Error(), `tenant "customer-a": trusted radius client "monitoring" is assigned by both trusted radius client assignments "assignment-a" and "assignment-b"`; got != want {
		t.Errorf("validateTenantRelationships() error = %q, want %q", got, want)
	}
}

func TestValidateCredentialProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile model.CredentialProfile
		wantErr string
	}{
		{
			name:    "shared secret specified",
			profile: model.CredentialProfile{SharedSecret: "secret"},
		},
		{
			name:    "shared secret missing",
			profile: model.CredentialProfile{},
			wantErr: `credential profile "default": shared_secret must be specified`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := validateCredentialProfile("default", test.profile)
			if len(errs) == 0 {
				if test.wantErr != "" {
					t.Fatal("validateCredentialProfile() returned no errors")
				}
				return
			}

			if test.wantErr == "" {
				t.Fatalf("validateCredentialProfile() error = %v, want none", errs[0])
			}
			if got := errs[0].Error(); got != test.wantErr {
				t.Fatalf("validateCredentialProfile() error = %q, want %q", got, test.wantErr)
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	values := map[string]int{
		"zulu":   1,
		"alpha":  2,
		"middle": 3,
	}

	if got, want := sortedKeys(values), []string{"alpha", "middle", "zulu"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("sortedKeys() = %v, want %v", got, want)
	}
}

func TestValidateAccountingProfileStaleSessionTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout string
		wantErr string
	}{
		{name: "omitted"},
		{name: "valid minutes", timeout: "20m"},
		{name: "valid hour", timeout: "1h"},
		{name: "invalid duration", timeout: "twenty minutes", wantErr: `accounting profile "default": stale_session_timeout must be a valid duration`},
		{name: "zero duration", timeout: "0s", wantErr: `accounting profile "default": stale_session_timeout must be greater than zero`},
		{name: "negative duration", timeout: "-5m", wantErr: `accounting profile "default": stale_session_timeout must be greater than zero`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := validateAccountingProfile("default", model.AccountingProfile{StaleSessionTimeout: test.timeout})
			if test.wantErr == "" {
				if len(errs) != 0 {
					t.Fatalf("validateAccountingProfile() errors = %v, want none", errs)
				}
				return
			}
			if len(errs) != 1 {
				t.Fatalf("validateAccountingProfile() returned %d errors, want 1", len(errs))
			}
			if got := errs[0].Error(); got != test.wantErr {
				t.Fatalf("validateAccountingProfile() error = %q, want %q", got, test.wantErr)
			}
		})
	}
}

func TestValidateAuthenticationProfile(t *testing.T) {
	one := 1
	zero := 0
	minusOne := -1
	tests := []struct {
		name    string
		profile model.AuthenticationProfile
		wantErr string
	}{
		{name: "valid", profile: model.AuthenticationProfile{SimultaneousUse: &one}},
		{
			name:    "unspecified",
			profile: model.AuthenticationProfile{},
		},
		{
			name: "zero",
			profile: model.AuthenticationProfile{
				SimultaneousUse: &zero,
			}, wantErr: `authentication profile "default": simultaneous_use must be greater than zero`},
		{name: "negative", profile: model.AuthenticationProfile{SimultaneousUse: &minusOne}, wantErr: `authentication profile "default": simultaneous_use must be greater than zero`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := validateAuthenticationProfile("default", test.profile)
			if test.wantErr == "" {
				if len(errs) != 0 {
					t.Fatalf("errors = %v, want none", errs)
				}
				return
			}
			if len(errs) != 1 || errs[0].Error() != test.wantErr {
				t.Fatalf("errors = %v, want %q", errs, test.wantErr)
			}
		})
	}
}

func TestValidateTenantAuthenticationProfileReference(t *testing.T) {
	one := 1
	globals := model.GlobalObjects{AuthenticationProfiles: map[string]model.AuthenticationProfile{"default": {SimultaneousUse: &one}}}
	tenant := model.Tenant{AuthenticationProfile: "default"}
	if errs := validateTenantReferences("customer-a", tenant, globals); len(errs) != 0 {
		t.Fatalf("errors = %v, want none", errs)
	}
	tenant.AuthenticationProfile = "missing"
	errs := validateTenantReferences("customer-a", tenant, globals)
	if len(errs) != 1 || errs[0].Error() != `tenant "customer-a": authentication profile "missing" does not exist` {
		t.Fatalf("errors = %v", errs)
	}
}

func TestValidateDeploymentProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile model.DeploymentProfile
		wantErr string
	}{
		{
			name:    "unspecified",
			profile: model.DeploymentProfile{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := validateDeploymentProfile("default", test.profile)
			if test.wantErr == "" {
				if len(errs) != 0 {
					t.Fatalf("errors = %v, want none", errs)
				}
				return
			}
			if len(errs) != 1 || errs[0].Error() != test.wantErr {
				t.Fatalf("errors = %v, want %q", errs, test.wantErr)
			}
		})
	}
}

func TestValidateTenantDeploymentProfileReference(t *testing.T) {
	globals := model.GlobalObjects{DeploymentProfiles: map[string]model.DeploymentProfile{"default": {}}}
	tenant := model.Tenant{DeploymentProfile: "default"}
	if errs := validateTenantReferences("customer-a", tenant, globals); len(errs) != 0 {
		t.Fatalf("errors = %v, want none", errs)
	}
	tenant.DeploymentProfile = "missing"
	errs := validateTenantReferences("customer-a", tenant, globals)
	if len(errs) != 1 || errs[0].Error() != `tenant "customer-a": deployment profile "missing" does not exist` {
		t.Fatalf("errors = %v", errs)
	}
}
