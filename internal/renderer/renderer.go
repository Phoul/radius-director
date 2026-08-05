// Package renderer renders intermediate FreeRADIUS configuration models.
package renderer

import (
	"bytes"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

var clientsTemplateLoader = templates.EmbeddedLoader()
var coaTemplateLoader = templates.EmbeddedLoader()
var proxyTemplateLoader = templates.EmbeddedLoader()
var sqlTemplateLoader = templates.EmbeddedLoader()
var authorizeTemplateLoader = templates.EmbeddedLoader()

// RenderClients renders the clients.conf content for a generated tenant.
func RenderClients(tenant generator.Tenant) (string, error) {
	template, err := clientsTemplateLoader.Load(tenant.RADIUSServer.Version, "clients.conf")
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := template.Execute(&rendered, tenant); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

// RenderCOA renders the sites-available/coa content for a generated tenant.
func RenderCOA(tenant generator.Tenant) (string, error) {
	template, err := coaTemplateLoader.Load(tenant.RADIUSServer.Version, "sites-available/coa")
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := template.Execute(&rendered, nil); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

// RenderProxy renders the proxy.conf content for a generated tenant.
func RenderProxy(tenant generator.Tenant) (string, error) {
	template, err := proxyTemplateLoader.Load(tenant.RADIUSServer.Version, "proxy.conf")
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := template.Execute(&rendered, tenant); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

// RenderSQL renders the mods-available/sql content for a generated tenant.
func RenderSQL(tenant generator.Tenant) (string, error) {
	template, err := sqlTemplateLoader.Load(tenant.RADIUSServer.Version, "mods-available/sql")
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := template.Execute(&rendered, tenant.SQL); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

// RenderAuthorize renders the mods-config/files/authorize content for a generated tenant.
func RenderAuthorize(tenant generator.Tenant) (string, error) {
	template, err := authorizeTemplateLoader.Load(
		tenant.RADIUSServer.Version,
		"mods-config/files/authorize",
	)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := template.Execute(&rendered, tenant); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
