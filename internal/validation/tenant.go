package validation

import (
	"fmt"

	"github.com/gobcn/radius-director/internal/model"
)

func validateTenants(tenants map[string]model.Tenant) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(tenants) {
		validationErrors = append(validationErrors, validateTenant(identifier, tenants[identifier])...)
	}

	return validationErrors
}

func validateTenant(identifier string, tenant model.Tenant) []error {
	validationErrors := validateDatabase(identifier, tenant.Database)
	validationErrors = append(validationErrors, validateRADIUSServers(identifier, tenant.RADIUSServers)...)
	validationErrors = append(validationErrors, validateNASAssignments(identifier, tenant.NASAssignments)...)

	return validationErrors
}

func validateDatabase(tenantIdentifier string, database model.Database) []error {
	var validationErrors []error
	if database.Engine == "" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database engine must be specified", tenantIdentifier))
	} else if database.Engine != "mysql" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database engine %q is not supported", tenantIdentifier, database.Engine))
	}
	if database.Host == "" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database host must be specified", tenantIdentifier))
	}
	if database.Port < 1 || database.Port > 65535 {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database port must be between 1 and 65535", tenantIdentifier))
	}
	if database.Database == "" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database name must be specified", tenantIdentifier))
	}
	if database.Username == "" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database username must be specified", tenantIdentifier))
	}
	if database.Password == "" {
		validationErrors = append(validationErrors, fmt.Errorf("tenant %q: database password must be specified", tenantIdentifier))
	}

	return validationErrors
}

func validateRADIUSServers(tenantIdentifier string, servers map[string]model.RADIUSServer) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(servers) {
		validationErrors = append(validationErrors, validateRADIUSServer(tenantIdentifier, identifier, servers[identifier])...)
	}

	return validationErrors
}

func validateRADIUSServer(tenantIdentifier, identifier string, server model.RADIUSServer) []error {
	return nil
}

func validateNASAssignments(tenantIdentifier string, assignments map[string]model.NASAssignment) []error {
	var validationErrors []error
	for _, identifier := range sortedKeys(assignments) {
		validationErrors = append(validationErrors, validateNASAssignment(tenantIdentifier, identifier, assignments[identifier])...)
	}

	return validationErrors
}

func validateNASAssignment(tenantIdentifier, identifier string, assignment model.NASAssignment) []error {
	return nil
}
