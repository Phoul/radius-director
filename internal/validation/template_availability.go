package validation

import (
	"fmt"

	"github.com/gobcn/radius-director/internal/model"
	"github.com/gobcn/radius-director/internal/templates"
)

// validateTemplateAvailability verifies that each tenant's deployment profile
// can be resolved for its selected FreeRADIUS version.
func validateTemplateAvailability(configuration model.Configuration) []error {
	var validationErrors []error

	for _, tenantIdentifier := range sortedKeys(configuration.Tenants) {
		tenant := configuration.Tenants[tenantIdentifier]
		if tenant.DeploymentProfile == "" || tenant.RADIUSServer.Version == "" || !templates.SupportsVersion(tenant.RADIUSServer.Version) {
			continue
		}

		profile, exists := configuration.GlobalObjects.DeploymentProfiles[tenant.DeploymentProfile]
		if !exists {
			continue
		}

		if validDeploymentProfileIdentifier(profile.Template) && !templates.SupportsTemplateSet(tenant.RADIUSServer.Version, profile.Template) {
			validationErrors = append(validationErrors, fmt.Errorf(
				"tenant %q: deployment profile %q: template set %q is not available for FreeRADIUS version %q",
				tenantIdentifier,
				tenant.DeploymentProfile,
				profile.Template,
				tenant.RADIUSServer.Version,
			))
		}

		for _, overlay := range profile.Overlays {
			if validDeploymentProfileIdentifier(overlay) && !templates.SupportsOverlay(tenant.RADIUSServer.Version, overlay) {
				validationErrors = append(validationErrors, fmt.Errorf(
					"tenant %q: deployment profile %q: overlay %q is not available for FreeRADIUS version %q",
					tenantIdentifier,
					tenant.DeploymentProfile,
					overlay,
					tenant.RADIUSServer.Version,
				))
			}
		}
	}

	return validationErrors
}
