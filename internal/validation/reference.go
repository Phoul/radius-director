package validation

import (
	"fmt"

	"github.com/gobcn/radius-director/internal/model"
)

func validateReferences(configuration model.Configuration) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(configuration.Tenants) {
		validationErrors = append(validationErrors, validateTenantReferences(identifier, configuration.Tenants[identifier], configuration.GlobalObjects)...)
	}

	return validationErrors
}

func validateTenantReferences(tenantIdentifier string, tenant model.Tenant, globalObjects model.GlobalObjects) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(tenant.NASAssignments) {
		validationErrors = append(validationErrors, validateNASAssignmentReferences(tenantIdentifier, identifier, tenant.NASAssignments[identifier], globalObjects)...)
	}

	return validationErrors
}

func validateNASAssignmentReferences(tenantIdentifier, identifier string, assignment model.NASAssignment, globalObjects model.GlobalObjects) []error {
	var validationErrors []error
	if assignment.NASDevice != "" {
		if _, exists := globalObjects.NASDevices[assignment.NASDevice]; !exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas assignment %q: nas device %q does not exist", tenantIdentifier, identifier, assignment.NASDevice))
		}
	}
	if assignment.CredentialProfile != "" {
		if _, exists := globalObjects.CredentialProfiles[assignment.CredentialProfile]; !exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas assignment %q: credential profile %q does not exist", tenantIdentifier, identifier, assignment.CredentialProfile))
		}
	}
	if assignment.AuthenticationProfile != "" {
		if _, exists := globalObjects.AuthenticationProfiles[assignment.AuthenticationProfile]; !exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas assignment %q: authentication profile %q does not exist", tenantIdentifier, identifier, assignment.AuthenticationProfile))
		}
	}
	if assignment.AccountingProfile != "" {
		if _, exists := globalObjects.AccountingProfiles[assignment.AccountingProfile]; !exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas assignment %q: accounting profile %q does not exist", tenantIdentifier, identifier, assignment.AccountingProfile))
		}
	}
	if assignment.MonitoringProfile != "" {
		if _, exists := globalObjects.MonitoringProfiles[assignment.MonitoringProfile]; !exists {
			validationErrors = append(validationErrors, fmt.Errorf("tenant %q: nas assignment %q: monitoring profile %q does not exist", tenantIdentifier, identifier, assignment.MonitoringProfile))
		}
	}

	return validationErrors
}
