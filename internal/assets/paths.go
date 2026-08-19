package assets

import "os"

const (
	// FactoryRoot contains the immutable assets bundled into the image.
	// It is used only by "export assets".
	FactoryRoot = "/app/factory"

	// AdminRootEnv is the environment variable used to specify the
	// administrator asset library.
	AdminRootEnv = "RADIUS_DIRECTOR_ASSETS"
)

// AdminRoot returns the administrator asset library root.
//
// When RADIUS_DIRECTOR_ASSETS is not set, the current working directory
// is used. This supports local development from the repository root.
//
// The factory asset library is never used as a runtime fallback.
func AdminRoot() string {
	if root := os.Getenv(AdminRootEnv); root != "" {
		return root
	}

	return "."
}
