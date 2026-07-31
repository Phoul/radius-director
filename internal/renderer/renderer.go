// Package renderer renders intermediate FreeRADIUS configuration models.
package renderer

import (
	"bytes"
	"strings"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

var sqlTemplateLoader = templates.EmbeddedLoader()

// RenderClients renders the clients.conf content for a generated tenant.
func RenderClients(tenant generator.Tenant) (string, error) {
	var rendered strings.Builder

	rendered.WriteString("# Tenant: ")
	rendered.WriteString(tenant.Identifier)
	rendered.WriteByte('\n')

	for _, client := range tenant.Clients {
		rendered.WriteString("client ")
		rendered.WriteString(client.Identifier)
		rendered.WriteString(" {\n")
		rendered.WriteString("    ipaddr = ")
		rendered.WriteString(client.IPAddress)
		rendered.WriteString("\n")
		rendered.WriteString("    secret = ")
		rendered.WriteString(client.SharedSecret)
		rendered.WriteString("\n}\n")
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
