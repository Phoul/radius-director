// Package generator builds intermediate FreeRADIUS configuration models.
package generator

// Configuration is an intermediate representation of generated FreeRADIUS
// configuration, grouped by tenant.
type Configuration struct {
	Tenants []Tenant
}

// Tenant contains generated configuration for one RADIUS Director tenant.
type Tenant struct {
	Identifier string
	Clients    []Client
}

// Client contains the information needed to render a FreeRADIUS client definition.
type Client struct {
	Identifier   string
	IPAddress    string
	SharedSecret string
}
