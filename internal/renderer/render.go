package renderer

import "github.com/gobcn/radius-director/internal/generator"

// Render renders every managed configuration file for a tenant.
//
// The initial implementation delegates to the existing per-file renderers.
// Future implementations will discover and render templates automatically.
func Render(tenant generator.Tenant) ([]RenderedFile, error) {
	files := make([]RenderedFile, 0, 4)

	clients, err := RenderClients(tenant)
	if err != nil {
		return nil, err
	}

	files = append(files, RenderedFile{
		RelativePath: "clients.conf",
		Content:      clients,
	})

	proxy, err := RenderProxy(tenant)
	if err != nil {
		return nil, err
	}

	files = append(files, RenderedFile{
		RelativePath: "proxy.conf",
		Content:      proxy,
	})

	coa, err := RenderCOA(tenant)
	if err != nil {
		return nil, err
	}

	files = append(files, RenderedFile{
		RelativePath: "sites-available/coa",
		Content:      coa,
	})

	authorize, err := RenderAuthorize(tenant)
	if err != nil {
		return nil, err
	}

	files = append(files, RenderedFile{
		RelativePath: "mods-config/files/authorize",
		Content:      authorize,
	})

	return files, nil
}
