package generator

import (
	"sort"

	"github.com/gobcn/radius-director/internal/model"
)

// Generate builds an intermediate FreeRADIUS configuration from a validated model.
func Generate(configuration model.Configuration) Configuration {
	generated := Configuration{
		Tenants: make([]Tenant, 0, len(configuration.Tenants)),
	}

	for _, tenantIdentifier := range sortedKeys(configuration.Tenants) {
		tenant := configuration.Tenants[tenantIdentifier]
		generatedTenant := Tenant{
			Identifier: tenantIdentifier,
			Clients:    make([]Client, 0, len(tenant.NASAssignments)),
			SQL: SQL{
				Engine:   tenant.Database.Engine,
				Host:     tenant.Database.Host,
				Port:     tenant.Database.Port,
				Database: tenant.Database.Database,
				Username: tenant.Database.Username,
				Password: tenant.Database.Password,
			},
			RADIUSServer: RADIUSServer{
				AuthenticationPort: tenant.RADIUSServer.AuthenticationPort,
				AccountingPort:     tenant.RADIUSServer.AccountingPort,
				COAPort:            tenant.RADIUSServer.COAPort,
			},
		}

		for _, assignmentIdentifier := range sortedKeys(tenant.NASAssignments) {
			assignment := tenant.NASAssignments[assignmentIdentifier]
			nasDevice := configuration.GlobalObjects.NASDevices[assignment.NASDevice]
			credentialProfile := configuration.GlobalObjects.CredentialProfiles[assignment.CredentialProfile]

			generatedTenant.Clients = append(generatedTenant.Clients, Client{
				Identifier:   assignmentIdentifier,
				IPAddress:    nasDevice.IPAddress,
				SharedSecret: credentialProfile.SharedSecret,
			})
		}

		generated.Tenants = append(generated.Tenants, generatedTenant)
	}

	return generated
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
