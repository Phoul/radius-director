// Package renderer renders intermediate FreeRADIUS configuration models.
package renderer

import (
	"bytes"
	"fmt"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

// Render renders every managed configuration file for a tenant.
func Render(tenant generator.Tenant, templateSet string) ([]RenderedFile, error) {
	version := tenant.RADIUSServer.Version

	paths, err := templates.ManagedTemplates(version, templateSet, tenant.Overlays)
	if err != nil {
		return nil, err
	}

	loader := templates.EmbeddedLoader()

	files := make([]RenderedFile, 0, len(paths))

	for _, path := range paths {
		tmpl, err := loader.Load(version, templateSet, tenant.Overlays, path)
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
