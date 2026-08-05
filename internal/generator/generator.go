package generator

import (
	"sort"
	"time"

	"github.com/gobcn/radius-director/internal/model"
)

// Generate builds an intermediate FreeRADIUS configuration from a validated model.
func Generate(configuration model.Configuration) Configuration {
	generated := Configuration{
		Tenants: make([]Tenant, 0, len(configuration.Tenants)),
	}

	for _, tenantIdentifier := range sortedKeys(configuration.Tenants) {
		tenant := configuration.Tenants[tenantIdentifier]
		authenticationProfile := configuration.GlobalObjects.AuthenticationProfiles[tenant.AuthenticationProfile]
		generatedTenant := Tenant{
			Identifier: tenantIdentifier,
			AuthenticationPolicy: AuthenticationPolicy{
				SimultaneousUse: authenticationProfile.SimultaneousUse,
			},
			FreeRADIUSClients:  make([]FreeRADIUSClient, 0, len(tenant.NASAssignments)+len(tenant.TrustedRADIUSClientAssignments)),
			HomeServers:        make([]HomeServer, 0, len(tenant.NASAssignments)),
			AccountingPolicies: make([]NASAccountingPolicy, 0, len(tenant.NASAssignments)),
			SQL: SQL{
				Engine:   tenant.Database.Engine,
				Host:     tenant.Database.Host,
				Port:     tenant.Database.Port,
				Database: tenant.Database.Database,
				Username: tenant.Database.Username,
				Password: tenant.Database.Password,
			},
			RADIUSServer: RADIUSServer{
				Version:            tenant.RADIUSServer.Version,
				AuthenticationPort: tenant.RADIUSServer.AuthenticationPort,
				AccountingPort:     tenant.RADIUSServer.AccountingPort,
				COAPort:            tenant.RADIUSServer.COAPort,
			},
		}

		for _, assignmentIdentifier := range sortedKeys(tenant.NASAssignments) {
			assignment := tenant.NASAssignments[assignmentIdentifier]
			nasDevice := configuration.GlobalObjects.NASDevices[assignment.NASDevice]
			credentialProfile := configuration.GlobalObjects.CredentialProfiles[assignment.CredentialProfile]
			accountingProfile := configuration.GlobalObjects.AccountingProfiles[assignment.AccountingProfile]

			var staleSessionTimeout *time.Duration
			if accountingProfile.StaleSessionTimeout != "" {
				parsedTimeout, _ := time.ParseDuration(accountingProfile.StaleSessionTimeout)
				staleSessionTimeout = &parsedTimeout
			}

			generatedTenant.FreeRADIUSClients = append(generatedTenant.FreeRADIUSClients, FreeRADIUSClient{
				Identifier:   assignmentIdentifier,
				IPAddress:    nasDevice.IPAddress,
				SharedSecret: credentialProfile.SharedSecret,
				Vendor:       nasDevice.Vendor,
			})
			generatedTenant.HomeServers = append(generatedTenant.HomeServers, HomeServer{
				Identifier:   assignmentIdentifier,
				IPAddress:    nasDevice.IPAddress,
				SharedSecret: credentialProfile.SharedSecret,
			})
			generatedTenant.AccountingPolicies = append(generatedTenant.AccountingPolicies, NASAccountingPolicy{
				NASAssignmentIdentifier: assignmentIdentifier,
				NASDeviceIdentifier:     assignment.NASDevice,
				IPAddress:               nasDevice.IPAddress,
				StaleSessionTimeout:     staleSessionTimeout,
			})
		}

		for _, assignmentIdentifier := range sortedKeys(tenant.TrustedRADIUSClientAssignments) {
			assignment := tenant.TrustedRADIUSClientAssignments[assignmentIdentifier]
			trustedRADIUSClient := configuration.GlobalObjects.TrustedRADIUSClients[assignment.TrustedRADIUSClient]
			credentialProfile := configuration.GlobalObjects.CredentialProfiles[assignment.CredentialProfile]
			generatedTenant.FreeRADIUSClients = append(generatedTenant.FreeRADIUSClients, FreeRADIUSClient{
				Identifier:   assignmentIdentifier,
				IPAddress:    trustedRADIUSClient.IPAddress,
				SharedSecret: credentialProfile.SharedSecret,
			})
		}
		sort.SliceStable(generatedTenant.FreeRADIUSClients, func(left, right int) bool {
			return generatedTenant.FreeRADIUSClients[left].Identifier < generatedTenant.FreeRADIUSClients[right].Identifier
		})

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
