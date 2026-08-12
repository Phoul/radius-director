// Package renderer renders intermediate FreeRADIUS configuration models.
package renderer

import (
	"bytes"
	"fmt"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

// Renderer renders managed configuration using a template library.
type Renderer struct {
	templateLoader templates.Loader
}

// New creates a Renderer that uses templateLoader.
func New(templateLoader templates.Loader) Renderer {
	return Renderer{templateLoader: templateLoader}
}

// Render renders every managed configuration file for a tenant.
func (r Renderer) Render(tenant generator.Tenant) ([]RenderedFile, error) {
	version := tenant.RADIUSServer.Version

	entries, err := r.templateLoader.ManagedTemplateEntries(
		version,
		tenant.Template,
		tenant.Overlays,
	)
	if err != nil {
		return nil, err
	}

	files := make([]RenderedFile, 0, len(entries))

	for _, entry := range entries {
		switch entry.Kind {
		case templates.ManagedTemplateKindRegular:
			tmpl, err := r.templateLoader.Load(
				version,
				tenant.Template,
				tenant.Overlays,
				entry.Path,
			)
			if err != nil {
				return nil, err
			}

			var rendered bytes.Buffer

			if err := tmpl.Execute(&rendered, tenant); err != nil {
				return nil, fmt.Errorf("render %q: %w", entry.Path, err)
			}

			files = append(files, RenderedFile{
				RelativePath: entry.Path,
				Kind:         RenderedFileKindRegular,
				Content:      rendered.String(),
			})

		case templates.ManagedTemplateKindSymlink:
			files = append(files, RenderedFile{
				RelativePath: entry.Path,
				Kind:         RenderedFileKindSymlink,
				Target:       entry.Target,
			})

		default:
			return nil, fmt.Errorf(
				"unsupported managed template kind for %q",
				entry.Path,
			)
		}
	}

	return files, nil
}
