// Package templates loads versioned managed configuration templates.
package templates

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
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

// DeploymentLoader loads a deployment template.
type DeploymentLoader interface {
	LoadDeployment(deployment, name string) (Executor, error)
}

type resolvedEntryKind int

const (
	resolvedFile resolvedEntryKind = iota
	resolvedSymlink
)

type resolvedEntry struct {
	kind   resolvedEntryKind
	source string
	target string
}

type resolvedTemplates struct {
	entries map[string]resolvedEntry
}

// ManagedTemplateKind identifies the type of a managed template object.
type ManagedTemplateKind int

const (
	ManagedTemplateKindRegular ManagedTemplateKind = iota
	ManagedTemplateKindSymlink
)

// ManagedTemplate describes one resolved template filesystem object.
type ManagedTemplate struct {
	Path   string
	Kind   ManagedTemplateKind
	Target string
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

	entry, ok := resolved.entries[name]
	if !ok {
		return nil, fmt.Errorf(
			"load template %q for FreeRADIUS version %q: template does not exist",
			name,
			version,
		)
	}

	if entry.kind != resolvedFile {
		return nil, fmt.Errorf(
			"load template %q for FreeRADIUS version %q: template is not a regular file",
			name,
			version,
		)
	}

	contents, err := fs.ReadFile(l.files, entry.source)
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

func (l Loader) LoadDeployment(deployment, name string) (Executor, error) {
	if !fs.ValidPath(deployment) || deployment == "." {
		return nil, fmt.Errorf("deployment template %q is invalid", deployment)
	}

	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("deployment template name %q is invalid", name)
	}

	templatePath := path.Join("deployment", deployment, name)

	parsed, err := template.ParseFS(l.files, templatePath)
	if err != nil {
		return nil, fmt.Errorf(
			"load deployment template %q for %q: %w",
			name,
			deployment,
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

func (l Loader) validateSymlink(
	layerRoot string,
	symlinkPath string,
) (string, error) {
	visited := make(map[string]bool)
	return l.resolveSymlink(layerRoot, symlinkPath, visited)
}

func (l Loader) resolveSymlink(
	layerRoot string,
	symlinkPath string,
	visited map[string]bool,
) (string, error) {
	if visited[symlinkPath] {
		return "", fmt.Errorf("symlink cycle detected at %q", symlinkPath)
	}
	visited[symlinkPath] = true

	target, err := fs.ReadLink(l.files, symlinkPath)
	if err != nil {
		return "", fmt.Errorf("read symlink target: %w", err)
	}

	if target == "" {
		return "", fmt.Errorf("symlink target is empty")
	}

	if filepath.VolumeName(target) != "" {
		return "", fmt.Errorf("symlink target %q is absolute", target)
	}

	if filepath.IsAbs(target) {
		return "", fmt.Errorf("symlink target %q is absolute", target)
	}

	target = filepath.ToSlash(target)

	if strings.Contains(target, `\`) {
		return "", fmt.Errorf(
			"symlink target %q contains a backslash",
			target,
		)
	}

	if path.IsAbs(target) {
		return "", fmt.Errorf("symlink target %q is absolute", target)
	}

	if len(target) >= 2 &&
		((target[0] >= 'A' && target[0] <= 'Z') ||
			(target[0] >= 'a' && target[0] <= 'z')) &&
		target[1] == ':' {
		return "", fmt.Errorf(
			"symlink target %q is absolute",
			target,
		)
	}

	linkDir := path.Dir(symlinkPath)
	targetPath := path.Clean(path.Join(linkDir, target))
	rootPath := path.Clean(layerRoot)

	if targetPath != rootPath && !strings.HasPrefix(targetPath, rootPath+"/") {
		return "", fmt.Errorf(
			"symlink target %q escapes template layer %q",
			target,
			layerRoot,
		)
	}

	info, err := fs.Lstat(l.files, targetPath)
	if err != nil {
		return "", fmt.Errorf(
			"symlink target %q does not exist: %w",
			target,
			err,
		)
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		if _, err := l.resolveSymlink(layerRoot, targetPath, visited); err != nil {
			return "", err
		}
	} else if !info.Mode().IsRegular() && !info.IsDir() {
		return "", fmt.Errorf(
			"symlink target %q is not a regular file or directory",
			target,
		)
	}

	return target, nil
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

		fullPath := path.Join(rootPath, relativePath)

		if d.Type()&fs.ModeSymlink != 0 {
			target, err := l.validateSymlink(rootPath, fullPath)
			if err != nil {
				return fmt.Errorf("validate symlink %q: %w", fullPath, err)
			}

			resolved.entries[relativePath] = resolvedEntry{
				kind:   resolvedSymlink,
				target: target,
			}

			return nil
		}

		if d.Type().IsRegular() {
			resolved.entries[relativePath] = resolvedEntry{
				kind:   resolvedFile,
				source: fullPath,
			}

			return nil
		}

		return fmt.Errorf("unsupported template filesystem entry %q", fullPath)
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
		entries: make(map[string]resolvedEntry),
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

// ManagedTemplateEntries returns every managed template object for a supported
// FreeRADIUS version.
func (l Loader) ManagedTemplateEntries(
	version,
	templateSet string,
	overlays []string,
) ([]ManagedTemplate, error) {
	resolved, err := l.resolve(version, templateSet, overlays)
	if err != nil {
		return nil, err
	}

	entries := make([]ManagedTemplate, 0, len(resolved.entries))

	for name, entry := range resolved.entries {
		managed := ManagedTemplate{
			Path: name,
		}

		switch entry.kind {
		case resolvedFile:
			managed.Kind = ManagedTemplateKindRegular

		case resolvedSymlink:
			managed.Kind = ManagedTemplateKindSymlink
			managed.Target = entry.target

		default:
			return nil, fmt.Errorf(
				"unsupported resolved template entry %q",
				name,
			)
		}

		entries = append(entries, managed)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return entries, nil
}

// ManagedTemplates returns every managed template for a supported FreeRADIUS
// version.
func (l Loader) ManagedTemplates(
	version,
	templateSet string,
	overlays []string,
) ([]string, error) {
	entries, err := l.ManagedTemplateEntries(version, templateSet, overlays)
	if err != nil {
		return nil, err
	}

	templates := make([]string, 0, len(entries))

	for _, entry := range entries {
		templates = append(templates, entry.Path)
	}

	return templates, nil
}
