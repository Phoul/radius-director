// Package validation validates RADIUS Director configuration models.
package validation

import (
	"errors"
	"sort"

	"github.com/gobcn/radius-director/internal/model"
	"github.com/gobcn/radius-director/internal/templates"
)

// Validate validates a RADIUS Director configuration.
func Validate(configuration model.Configuration, templateLoader templates.Loader) error {
	validationErrors := validateGlobalObjects(configuration.GlobalObjects)
	validationErrors = append(validationErrors, validateTenants(configuration.Tenants, templateLoader)...)
	validationErrors = append(validationErrors, validateTemplateAvailability(configuration, templateLoader)...)
	validationErrors = append(validationErrors, validateReferences(configuration)...)
	validationErrors = append(validationErrors, validateRelationships(configuration)...)

	return errors.Join(validationErrors...)
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
