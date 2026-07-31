// Package templates loads embedded, versioned managed configuration templates.
package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
	"text/template"
)

//go:embed all:3.2.9 all:3.4.0
var embeddedFiles embed.FS

// Executor executes a parsed managed configuration template.
type Executor interface {
	Execute(io.Writer, any) error
}

// Loader selects and loads a managed configuration template for a FreeRADIUS
// version.
type Loader interface {
	Load(version, name string) (Executor, error)
}

type loader struct {
	files fs.FS
}

// EmbeddedLoader returns a Loader backed by templates embedded in the binary.
func EmbeddedLoader() Loader {
	return loader{files: embeddedFiles}
}

func (l loader) Load(version, name string) (Executor, error) {
	if !isSupportedVersion(version) {
		return nil, fmt.Errorf("FreeRADIUS version %q is not supported", version)
	}
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("template name %q is invalid", name)
	}

	templatePath := path.Join(version, name)
	parsed, err := template.ParseFS(l.files, templatePath)
	if err != nil {
		return nil, fmt.Errorf("load template %q for FreeRADIUS version %q: %w", name, version, err)
	}

	return parsed, nil
}

func isSupportedVersion(version string) bool {
	switch version {
	case "3.2.9", "3.4.0":
		return true
	default:
		return false
	}
}
