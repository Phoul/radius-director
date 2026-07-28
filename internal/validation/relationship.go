package validation

import (
	"fmt"

	"github.com/gobcn/radius-director/internal/model"
)

func validateRelationships(configuration model.Configuration) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(configuration.Tenants) {
		validationErrors = append(validationErrors, validateTenantRelationships(identifier, configuration.Tenants[identifier], configuration.GlobalObjects)...)
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

	return validationErrors
}
