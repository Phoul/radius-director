// Package generator builds intermediate FreeRADIUS configuration models.
package generator

// Configuration is an intermediate representation of generated FreeRADIUS
// configuration, grouped by tenant.
type Configuration struct {
	Tenants []Tenant
}

// Tenant contains generated configuration for one RADIUS Director tenant.
type Tenant struct {
	Identifier   string
	Clients      []Client
	SQL          SQL
	RADIUSServer RADIUSServer
}

// Client contains the information needed to render a FreeRADIUS client definition.
type Client struct {
	Identifier   string
	IPAddress    string
	SharedSecret string
}

// SQL contains the information needed to render a FreeRADIUS sql module.
type SQL struct {
	Engine   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// RADIUSServer contains the information needed to render a FreeRADIUS server.
type RADIUSServer struct {
	AuthenticationPort int
	AccountingPort     int
	COAPort            int
}
