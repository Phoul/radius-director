package validation

import (
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/model"
)

func TestValidate(t *testing.T) {
	configuration := model.Configuration{
		GlobalObjects: model.GlobalObjects{
			CredentialProfiles: map[string]model.CredentialProfile{
				"default": {SharedSecret: "secret"},
			},
			AuthenticationProfiles: map[string]model.AuthenticationProfile{
				"default": {},
			},
			AccountingProfiles: map[string]model.AccountingProfile{
				"default": {},
			},
			MonitoringProfiles: map[string]model.MonitoringProfile{
				"default": {},
			},
			NASDevices: map[string]model.NASDevice{
				"core": {IPAddress: "10.10.10.1", Vendor: "mikrotik"},
			},
		},
		Tenants: map[string]model.Tenant{
			"customer-a": {
				Database: model.Database{
					Engine:   "mysql",
					Host:     "db.example.com",
					Port:     3306,
					Database: "radius",
					Username: "radius",
					Password: "secret",
				},
				RADIUSServer: model.RADIUSServer{
					AuthenticationPort: 1812,
					AccountingPort:     1813,
					COAPort:            3799,
				},
				NASAssignments: map[string]model.NASAssignment{
					"core": {
						NASDevice:             "core",
						CredentialProfile:     "default",
						AuthenticationProfile: "default",
						AccountingProfile:     "default",
						MonitoringProfile:     "default",
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
		NASDevice:             "core",
		CredentialProfile:     "default",
		AuthenticationProfile: "default",
		AccountingProfile:     "default",
		MonitoringProfile:     "default",
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
				CredentialProfile:     validAssignment.CredentialProfile,
				AuthenticationProfile: validAssignment.AuthenticationProfile,
				AccountingProfile:     validAssignment.AccountingProfile,
				MonitoringProfile:     validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": nas_device must be specified`,
			},
		},
		{
			name: "credential profile missing",
			assignment: model.NASAssignment{
				NASDevice:             validAssignment.NASDevice,
				AuthenticationProfile: validAssignment.AuthenticationProfile,
				AccountingProfile:     validAssignment.AccountingProfile,
				MonitoringProfile:     validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": credential_profile must be specified`,
			},
		},
		{
			name: "authentication profile missing",
			assignment: model.NASAssignment{
				NASDevice:         validAssignment.NASDevice,
				CredentialProfile: validAssignment.CredentialProfile,
				AccountingProfile: validAssignment.AccountingProfile,
				MonitoringProfile: validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": authentication_profile must be specified`,
			},
		},
		{
			name: "accounting profile missing",
			assignment: model.NASAssignment{
				NASDevice:             validAssignment.NASDevice,
				CredentialProfile:     validAssignment.CredentialProfile,
				AuthenticationProfile: validAssignment.AuthenticationProfile,
				MonitoringProfile:     validAssignment.MonitoringProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": accounting_profile must be specified`,
			},
		},
		{
			name: "monitoring profile missing",
			assignment: model.NASAssignment{
				NASDevice:             validAssignment.NASDevice,
				CredentialProfile:     validAssignment.CredentialProfile,
				AuthenticationProfile: validAssignment.AuthenticationProfile,
				AccountingProfile:     validAssignment.AccountingProfile,
			},
			wantErrs: []string{
				`tenant "customer-a": nas assignment "core": monitoring_profile must be specified`,
			},
		},
		{
			name: "multiple properties missing",
			assignment: model.NASAssignment{
				NASDevice:             validAssignment.NASDevice,
				AuthenticationProfile: validAssignment.AuthenticationProfile,
				MonitoringProfile:     validAssignment.MonitoringProfile,
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
				`tenant "customer-a": nas assignment "core": authentication_profile must be specified`,
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

func TestValidateTenant(t *testing.T) {
	validTenant := model.Tenant{
		Database: model.Database{
			Engine:   "mysql",
			Host:     "db.example.com",
			Port:     3306,
			Database: "radius",
			Username: "radius",
			Password: "secret",
		},
		RADIUSServer: model.RADIUSServer{
			AuthenticationPort: 1812,
			AccountingPort:     1813,
			COAPort:            3799,
		},
		NASAssignments: map[string]model.NASAssignment{
			"core": {
				NASDevice:             "core",
				CredentialProfile:     "default",
				AuthenticationProfile: "default",
				AccountingProfile:     "default",
				MonitoringProfile:     "default",
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
				RADIUSServer:   validTenant.RADIUSServer,
				NASAssignments: validTenant.NASAssignments,
			},
			wantErrs: []string{
				`tenant "customer-a": exactly one database must be defined`,
			},
		},
		{
			name: "RADIUS Server missing",
			tenant: model.Tenant{
				Database:       validTenant.Database,
				NASAssignments: validTenant.NASAssignments,
			},
			wantErrs: []string{
				`tenant "customer-a": exactly one radius server must be defined`,
			},
		},
		{
			name: "NAS assignments missing",
			tenant: model.Tenant{
				Database:     validTenant.Database,
				RADIUSServer: validTenant.RADIUSServer,
			},
			wantErrs: []string{
				`tenant "customer-a": at least one nas assignment must be defined`,
			},
		},
		{
			name: "all required tenant objects missing",
			wantErrs: []string{
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
			name: "minimum valid ports",
			server: model.RADIUSServer{
				AuthenticationPort: 1,
				AccountingPort:     1,
				COAPort:            1,
			},
		},
		{
			name: "maximum valid ports",
			server: model.RADIUSServer{
				AuthenticationPort: 65535,
				AccountingPort:     65535,
				COAPort:            65535,
			},
		},
		{
			name: "authentication port missing",
			server: model.RADIUSServer{
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
