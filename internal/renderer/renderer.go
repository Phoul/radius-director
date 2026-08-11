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

	paths, err := r.templateLoader.ManagedTemplates(version, tenant.Template, tenant.Overlays)
	if err != nil {
		return nil, err
	}

	files := make([]RenderedFile, 0, len(paths))

	for _, path := range paths {
		tmpl, err := r.templateLoader.Load(version, tenant.Template, tenant.Overlays, path)
		if err != nil {
			return nil, err
		}

		var rendered bytes.Buffer

		if err := tmpl.Execute(&rendered, tenant); err != nil {
			return nil, fmt.Errorf("render %q: %w", path, err)
		}

		files = append(files, RenderedFile{
			RelativePath: path,
			Content:      rendered.String(),
		})
	}

	return files, nil
}
