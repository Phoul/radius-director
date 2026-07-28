// Package renderer renders intermediate FreeRADIUS configuration models.
package renderer

import (
	"strings"

	"github.com/gobcn/radius-director/internal/generator"
)

// RenderClients renders the clients.conf content for a generated configuration.
func RenderClients(configuration generator.Configuration) (string, error) {
	var rendered strings.Builder

	for tenantIndex, tenant := range configuration.Tenants {
		if tenantIndex > 0 {
			rendered.WriteByte('\n')
		}

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
	}

	return rendered.String(), nil
}
