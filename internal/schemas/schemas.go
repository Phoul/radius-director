// Package schemas loads versioned FreeRADIUS database schemas.
package schemas

import (
	"fmt"
	"io/fs"
	"path"
)

// Loader loads database schemas from a schema library filesystem.
type Loader struct {
	files fs.FS
}

// NewLoader creates a Loader for a schema library filesystem.
func NewLoader(files fs.FS) Loader {
	return Loader{files: files}
}

// LoadMySQLSchema loads the MySQL schema for a FreeRADIUS version.
func (l Loader) LoadMySQLSchema(version string) (string, error) {
	if !fs.ValidPath(version) || version == "." {
		return "", fmt.Errorf("FreeRADIUS version %q is invalid", version)
	}

	schemaPath := path.Join(
		"freeradius",
		version,
		"mysql",
		"schema.sql",
	)

	contents, err := fs.ReadFile(l.files, schemaPath)
	if err != nil {
		return "", fmt.Errorf(
			"load MySQL schema for FreeRADIUS version %q: %w",
			version,
			err,
		)
	}

	return string(contents), nil
}
