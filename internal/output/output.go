// Package output assembles rendered FreeRADIUS configuration files.
package output

import (
	"path/filepath"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

const (
	ClientsFile = "clients.conf"
)

var renderClients = renderer.RenderClients

// Output contains the generated files that are ready to be written.
type Output struct {
	Files []File
}

// File is a generated file and its content.
type File struct {
	Path    string
	Content string
}

// Build assembles generated files from an intermediate configuration model.
func Build(configuration generator.Configuration) (Output, error) {
	generated := Output{
		Files: make([]File, 0, len(configuration.Tenants)),
	}

	for _, tenant := range configuration.Tenants {
		clients, err := renderClients(tenant)
		if err != nil {
			return Output{}, err
		}

		generated.Files = append(generated.Files, File{
			Path:    filepath.Join(tenant.Identifier, ClientsFile),
			Content: clients,
		})
	}

	return generated, nil
}
