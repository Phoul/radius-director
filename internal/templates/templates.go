// Package templates loads embedded, versioned managed configuration templates.
package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
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
	Load(version, templateSet, name string) (Executor, error)
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

// ManagedTemplates returns every managed template for a supported
// FreeRADIUS version.
func ManagedTemplates(version, templateSet string) ([]string, error) {
	return loader{files: embeddedFiles}.managedTemplates(version, templateSet)
}

func (l loader) Load(version, templateSet, name string) (Executor, error) {
	if !l.supportsVersion(version) {
		return nil, fmt.Errorf("FreeRADIUS version %q is not supported", version)
	}
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("template name %q is invalid", name)
	}
	if !fs.ValidPath(templateSet) {
		return nil, fmt.Errorf("template set %q is invalid", templateSet)
	}

	templatePath := path.Join(version, templateSet, name)
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

func (l loader) managedTemplates(version, templateSet string) ([]string, error) {
	if !l.supportsVersion(version) {
		return nil, fmt.Errorf("FreeRADIUS version %q is not supported", version)
	}

	root, err := fs.Sub(l.files, path.Join(version, templateSet))
	if err != nil {
		return nil, err
	}

	var templates []string

	err = fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		templates = append(templates, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(templates)

	return templates, nil
}
