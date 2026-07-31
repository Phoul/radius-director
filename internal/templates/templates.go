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

//go:embed all:*
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

// SupportsVersion reports whether an embedded template set supports version.
func SupportsVersion(version string) bool {
	return loader{files: embeddedFiles}.supportsVersion(version)
}

func (l loader) Load(version, name string) (Executor, error) {
	if !l.supportsVersion(version) {
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

func (l loader) supportsVersion(version string) bool {
	if !fs.ValidPath(version) {
		return false
	}

	info, err := fs.Stat(l.files, version)
	return err == nil && info.IsDir()
}
