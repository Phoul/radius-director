// Package output assembles rendered FreeRADIUS configuration files.
package output

import (
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
	clients, err := renderClients(configuration)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Files: []File{
			{
				Path:    ClientsFile,
				Content: clients,
			},
		},
	}, nil
}
