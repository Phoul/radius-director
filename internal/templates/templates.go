// Package templates loads versioned managed configuration templates.
package templates

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"text/template"
)

// Executor executes a parsed managed configuration template.
type Executor interface {
	Execute(io.Writer, any) error
}

// Loader selects and loads managed configuration templates from a template
// library filesystem.
type Loader struct {
	files fs.FS
}

type resolvedTemplates struct {
	files map[string]string
}

// NewLoader creates a Loader for a template library filesystem.
func NewLoader(files fs.FS) Loader {
	return Loader{files: files}
}

func validSetName(name string) bool {
	return name != "." && fs.ValidPath(name) && !strings.Contains(name, "/")
}

func (l Loader) Load(
	version,
	templateSet string,
	overlays []string,
	name string,
) (Executor, error) {
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("template name %q is invalid", name)
	}

	resolved, err := l.resolve(version, templateSet, overlays)
	if err != nil {
		return nil, err
	}

	templatePath, ok := resolved.files[name]
	if !ok {
		return nil, fmt.Errorf(
			"load template %q for FreeRADIUS version %q: template does not exist",
			name,
			version,
		)
	}

	contents, err := fs.ReadFile(l.files, templatePath)
	if err != nil {
		return nil, fmt.Errorf(
			"load template %q for FreeRADIUS version %q: %w",
			name,
			version,
			err,
		)
	}

	parsed, err := template.New(path.Base(name)).Parse(string(contents))
	if err != nil {
		return nil, fmt.Errorf(
			"load template %q for FreeRADIUS version %q: %w",
			name,
			version,
			err,
		)
	}

	return parsed, nil
}

// SupportsVersion reports whether the template library supports version.
func (l Loader) SupportsVersion(version string) bool {
	if !fs.ValidPath(version) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("sets", version))
	return err == nil && info.IsDir()
}

// SupportsTemplateSet reports whether a template set is available for a
// FreeRADIUS version.
func (l Loader) SupportsTemplateSet(version, templateSet string) bool {
	if !l.SupportsVersion(version) || !validSetName(templateSet) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("sets", version, templateSet))
	return err == nil && info.IsDir()
}

// SupportsOverlay reports whether an overlay is available for a FreeRADIUS
// version.
func (l Loader) SupportsOverlay(version, overlay string) bool {
	if !l.SupportsVersion(version) || !validSetName(overlay) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("overlays", version, overlay))
	return err == nil && info.IsDir()
}

func (l Loader) mergeTemplateTree(
	resolved *resolvedTemplates,
	rootPath string,
) error {
	root, err := fs.Sub(l.files, rootPath)
	if err != nil {
		return err
	}

	return fs.WalkDir(root, ".", func(relativePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		resolved.files[relativePath] = path.Join(rootPath, relativePath)

		return nil
	})
}

func (l Loader) resolve(
	version,
	templateSet string,
	overlays []string,
) (*resolvedTemplates, error) {
	if !l.SupportsVersion(version) {
		return nil, fmt.Errorf("FreeRADIUS version %q is not supported", version)
	}

	if !validSetName(templateSet) {
		return nil, fmt.Errorf("template set %q is invalid", templateSet)
	}

	if !l.SupportsTemplateSet(version, templateSet) {
		return nil, fmt.Errorf(
			"template set %q is not available for FreeRADIUS version %q",
			templateSet,
			version,
		)
	}

	resolved := &resolvedTemplates{
		files: make(map[string]string),
	}

	if err := l.mergeTemplateTree(
		resolved,
		path.Join("sets", version, templateSet),
	); err != nil {
		return nil, err
	}

	for _, overlay := range overlays {
		if !validSetName(overlay) {
			return nil, fmt.Errorf("overlay %q is invalid", overlay)
		}

		if !l.SupportsOverlay(version, overlay) {
			return nil, fmt.Errorf(
				"overlay %q is not available for FreeRADIUS version %q",
				overlay,
				version,
			)
		}

		overlayPath := path.Join("overlays", version, overlay)

		if err := l.mergeTemplateTree(resolved, overlayPath); err != nil {
			return nil, err
		}
	}

	return resolved, nil
}

// ManagedTemplates returns every managed template for a supported FreeRADIUS
// version.
func (l Loader) ManagedTemplates(
	version,
	templateSet string,
	overlays []string,
) ([]string, error) {
	resolved, err := l.resolve(version, templateSet, overlays)
	if err != nil {
		return nil, err
	}

	templates := make([]string, 0, len(resolved.files))

	for name := range resolved.files {
		templates = append(templates, name)
	}

	sort.Strings(templates)

	return templates, nil
}
