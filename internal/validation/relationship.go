package validation

import (
	"fmt"

	"github.com/gobcn/radius-director/internal/model"
)

func validateRelationships(configuration model.Configuration) []error {
	var validationErrors []error

	validationErrors = append(
		validationErrors,
		validateProxySQLUsernames(configuration)...,
	)

	for _, identifier := range sortedKeys(configuration.Tenants) {
		validationErrors = append(
			validationErrors,
			validateTenantRelationships(
				identifier,
				configuration.Tenants[identifier],
				configuration.GlobalObjects,
			)...,
		)
	}

	return validationErrors
}

func validateTenantRelationships(tenantIdentifier string, tenant model.Tenant, globalObjects model.GlobalObjects) []error {
	var validationErrors []error
	assignmentsByNASDevice := make(map[string]string)

	for _, identifier := range sortedKeys(tenant.NASAssignments) {
		assignment := tenant.NASAssignments[identifier]
		if assignment.NASDevice == "" {
			continue
		}
		if _, exists := globalObjects.NASDevices[assignment.NASDevice]; !exists {
			continue
		}

		if firstAssignment, exists := assignmentsByNASDevice[assignment.NASDevice]; exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas device %q is assigned by both nas assignments %q and %q", tenantIdentifier, assignment.NASDevice, firstAssignment, identifier))
			continue
		}

		assignmentsByNASDevice[assignment.NASDevice] = identifier
	}
	assignmentsByTrustedRADIUSClient := make(map[string]string)
	for _, identifier := range sortedKeys(tenant.TrustedRADIUSClientAssignments) {
		assignment := tenant.TrustedRADIUSClientAssignments[identifier]
		if assignment.TrustedRADIUSClient == "" {
			continue
		}
		if _, exists := globalObjects.TrustedRADIUSClients[assignment.TrustedRADIUSClient]; !exists {
			continue
		}

		if firstAssignment, exists := assignmentsByTrustedRADIUSClient[assignment.TrustedRADIUSClient]; exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: trusted radius client %q is assigned by both trusted radius client assignments %q and %q", tenantIdentifier, assignment.TrustedRADIUSClient, firstAssignment, identifier))
			continue
		}

		assignmentsByTrustedRADIUSClient[assignment.TrustedRADIUSClient] = identifier
	}

	return validationErrors
}

func validateProxySQLUsernames(configuration model.Configuration) []error {
	var validationErrors []error
	usernames := make(map[string]string)

	for _, identifier := range sortedKeys(configuration.Tenants) {
		tenant := configuration.Tenants[identifier]

		deployment := tenant.Database.Deployment
		if deployment == "" {
			deployment = "container"
		}

		if deployment != "proxysql" {
			continue
		}

		username := tenant.Database.Username
		if username == "" {
			continue
		}

		if firstTenant, exists := usernames[username]; exists {
			validationErrors = append(
				validationErrors,
				fmt.Errorf(
					"tenants %q and %q both use ProxySQL database username %q",
					firstTenant,
					identifier,
					username,
				),
			)
			continue
		}

		usernames[username] = identifier
	}

	return validationErrors
}
