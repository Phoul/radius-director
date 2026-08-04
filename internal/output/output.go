// Package output assembles rendered FreeRADIUS configuration files.
package output

import (
	"path/filepath"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

const (
	// ClientsFile is the relative path of a tenant's FreeRADIUS clients file.
	ClientsFile = "clients.conf"
	// ProxyFile is the relative path of a tenant's FreeRADIUS proxy file.
	ProxyFile = "proxy.conf"
	// COASiteFile is the relative path of a tenant's FreeRADIUS CoA site.
	COASiteFile = "sites-available/coa"
)

var renderClients = renderer.RenderClients
var renderCOA = renderer.RenderCOA
var renderProxy = renderer.RenderProxy

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
func Build(configuration generator.Configuration) (Output, error) {
	generated := Output{
		Files: make([]File, 0, len(configuration.Tenants)*3),
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

		coa, err := renderCOA(tenant)
		if err != nil {
			return Output{}, err
		}

		generated.Files = append(generated.Files, File{
			Path:    filepath.Join(tenant.Identifier, COASiteFile),
			Content: coa,
		})

		proxy, err := renderProxy(tenant)
		if err != nil {
			return Output{}, err
		}

		generated.Files = append(generated.Files, File{
			Path:    filepath.Join(tenant.Identifier, ProxyFile),
			Content: proxy,
		})
	}

	return generated, nil
}
