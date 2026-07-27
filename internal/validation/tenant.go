package validation

import "github.com/gobcn/radius-director/internal/model"

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
	return nil
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
