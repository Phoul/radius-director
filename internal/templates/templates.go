// Package templates loads embedded, versioned managed configuration templates.
package templates

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/gobcn/radius-director/templates"
)

// Executor executes a parsed managed configuration template.
type Executor interface {
	Execute(io.Writer, any) error
}

// Loader selects and loads a managed configuration template for a FreeRADIUS
// version.
type Loader interface {
	Load(version, templateSet string, overlays []string, name string) (Executor, error)
}

type loader struct {
	files fs.FS
}

type resolvedTemplates struct {
	files map[string]string
}

func validSetName(name string) bool {
	return name != "." && fs.ValidPath(name) && !strings.Contains(name, "/")
}

// EmbeddedLoader returns a Loader backed by templates embedded in the binary.
func EmbeddedLoader() Loader {
	return loader{files: templateassets.Files}
}

// SupportsVersion reports whether an embedded template set supports version.
func SupportsVersion(version string) bool {
	return loader{files: templateassets.Files}.supportsVersion(version)
}

// SupportsTemplateSet reports whether an embedded template set is available
// for a FreeRADIUS version.
func SupportsTemplateSet(version, templateSet string) bool {
	return loader{files: templateassets.Files}.supportsTemplateSet(version, templateSet)
}

// SupportsOverlay reports whether an embedded overlay is available for a
// FreeRADIUS version.
func SupportsOverlay(version, overlay string) bool {
	return loader{files: templateassets.Files}.supportsOverlay(version, overlay)
}

// ManagedTemplates returns every managed template for a supported
// FreeRADIUS version.
func ManagedTemplates(version, templateSet string, overlays []string) ([]string, error) {
	return loader{files: templateassets.Files}.managedTemplates(version, templateSet, overlays)
}

func (l loader) Load(
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

func (l loader) supportsVersion(version string) bool {
	if !fs.ValidPath(version) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("sets", version))
	return err == nil && info.IsDir()
}

func (l loader) supportsTemplateSet(version, templateSet string) bool {
	if !l.supportsVersion(version) || !validSetName(templateSet) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("sets", version, templateSet))
	return err == nil && info.IsDir()
}

func (l loader) supportsOverlay(version, overlay string) bool {
	if !l.supportsVersion(version) || !validSetName(overlay) {
		return false
	}

	info, err := fs.Stat(l.files, path.Join("overlays", version, overlay))
	return err == nil && info.IsDir()
}

func (l loader) mergeTemplateTree(
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

func (l loader) resolve(
	version,
	templateSet string,
	overlays []string,
) (*resolvedTemplates, error) {
	if !l.supportsVersion(version) {
		return nil, fmt.Errorf("FreeRADIUS version %q is not supported", version)
	}

	if !validSetName(templateSet) {
		return nil, fmt.Errorf("template set %q is invalid", templateSet)
	}

	if !l.supportsTemplateSet(version, templateSet) {
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

		if !l.supportsOverlay(version, overlay) {
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

func (l loader) managedTemplates(
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
