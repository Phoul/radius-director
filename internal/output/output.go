// Package output assembles rendered FreeRADIUS configuration files.
package output

import (
	"path/filepath"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

// Output contains the generated files that are ready to be written.
type Output struct {
	Files []File
}

// File is a generated file and its content.
type File struct {
	Path    string
	Content string
}

// Build renders each tenant's files and assembles them into an Output object.
func Build(configuration generator.Configuration, templateRenderer renderer.Renderer) (Output, error) {
	generated := Output{
		Files: make([]File, 0),
	}

	for _, tenant := range configuration.Tenants {
		renderedFiles, err := templateRenderer.Render(tenant)
		if err != nil {
			return Output{}, err
		}

		for _, rendered := range renderedFiles {
			generated.Files = append(generated.Files, File{
				Path:    filepath.Join(tenant.Identifier, rendered.RelativePath),
				Content: rendered.Content,
			})
		}
	}

	return generated, nil
}
