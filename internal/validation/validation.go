// Package validation validates RADIUS Director configuration models.
package validation

import (
	"errors"
	"sort"

	"github.com/gobcn/radius-director/internal/model"
)

// Validate validates a RADIUS Director configuration.
func Validate(configuration model.Configuration) error {
	validationErrors := validateGlobalObjects(configuration.GlobalObjects)
	validationErrors = append(validationErrors, validateTenants(configuration.Tenants)...)
	validationErrors = append(validationErrors, validateReferences(configuration)...)

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
