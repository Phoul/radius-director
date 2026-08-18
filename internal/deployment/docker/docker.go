package docker

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/output"
	"github.com/gobcn/radius-director/internal/schemas"
	"github.com/gobcn/radius-director/internal/templates"
)

type ComposeData struct {
	Tenants              []ComposeTenant
	ProxySQL             *ComposeProxySQL
	HasContainerDatabase bool
}

type ComposeTenant struct {
	Name                string
	ServiceName         string
	DatabaseDeployment  string
	DatabaseServiceName string
	DatabaseVolumeName  string
	DatabaseName        string
	DatabaseUsername    string
	DatabasePassword    string
	RadiusVersion       string
	AuthenticationPort  int
	AccountingPort      int
	COAPort             int
}

type ComposeProxySQL struct {
	ServiceName string
	Backends    []ComposeProxySQLBackend
}

type ComposeProxySQLBackend struct {
	TenantName string

	Host string
	Port int

	Database string
	Username string
	Password string

	Hostgroup int
}

type Deployment struct {
	Files []output.File
}

const (
	ComposeOutputPath    = "docker-compose.yml"
	EntrypointOutputPath = "entrypoint.sh"
	ProxySQLOutputPath   = "proxysql.cnf"
)

func NewComposeData(configuration generator.Configuration) ComposeData {
	data := ComposeData{
		Tenants: make([]ComposeTenant, 0, len(configuration.Tenants)),
	}

	for _, tenant := range configuration.Tenants {
		data.Tenants = append(data.Tenants, ComposeTenant{
			Name:                tenant.Identifier,
			ServiceName:         "radius-" + tenant.Identifier,
			DatabaseDeployment:  tenant.DatabaseDeployment,
			DatabaseServiceName: "database-" + tenant.Identifier,
			DatabaseVolumeName:  "database-" + tenant.Identifier + "-data",
			DatabaseName:        tenant.SQL.Database,
			DatabaseUsername:    tenant.SQL.Username,
			DatabasePassword:    tenant.SQL.Password,
			RadiusVersion:       tenant.RADIUSServer.Version,
			AuthenticationPort:  tenant.RADIUSServer.AuthenticationPort,
			AccountingPort:      tenant.RADIUSServer.AccountingPort,
			COAPort:             tenant.RADIUSServer.COAPort,
		})

		if tenant.DatabaseDeployment == "container" {
			data.HasContainerDatabase = true
		}

		if tenant.ProxySQL == nil {
			continue
		}

		if data.ProxySQL == nil {
			data.ProxySQL = &ComposeProxySQL{
				ServiceName: "proxysql",
				Backends:    make([]ComposeProxySQLBackend, 0),
			}
		}

		hostgroup := 10 + len(data.ProxySQL.Backends)

		data.ProxySQL.Backends = append(
			data.ProxySQL.Backends,
			ComposeProxySQLBackend{
				TenantName: tenant.Identifier,
				Host:       tenant.ProxySQL.BackendHost,
				Port:       tenant.ProxySQL.BackendPort,
				Database:   tenant.SQL.Database,
				Username:   tenant.SQL.Username,
				Password:   tenant.SQL.Password,
				Hostgroup:  hostgroup,
			},
		)
	}

	return data
}

func GenerateDeployment(
	loader templates.Loader,
	schemaLoader schemas.Loader,
	configuration generator.Configuration,
) (Deployment, error) {
	data := NewComposeData(configuration)

	composeExecutor, err := loader.LoadDeployment("docker", "docker-compose.yml")
	if err != nil {
		return Deployment{}, fmt.Errorf("load Docker Compose template: %w", err)
	}

	var compose bytes.Buffer

	if err := composeExecutor.Execute(&compose, data); err != nil {
		return Deployment{}, fmt.Errorf("render Docker Compose template: %w", err)
	}

	entrypointExecutor, err := loader.LoadDeployment("docker", "entrypoint.sh")
	if err != nil {
		return Deployment{}, fmt.Errorf("load Docker entrypoint template: %w", err)
	}

	var entrypoint bytes.Buffer

	if err := entrypointExecutor.Execute(&entrypoint, data); err != nil {
		return Deployment{}, fmt.Errorf("render Docker entrypoint template: %w", err)
	}

	var proxySQL bytes.Buffer

	if data.ProxySQL != nil {
		proxySQLExecutor, err := loader.LoadDeployment("docker", "proxysql.cnf")
		if err != nil {
			return Deployment{}, fmt.Errorf(
				"load ProxySQL configuration template: %w",
				err,
			)
		}

		if err := proxySQLExecutor.Execute(&proxySQL, data); err != nil {
			return Deployment{}, fmt.Errorf(
				"render ProxySQL configuration template: %w",
				err,
			)
		}
	}

	var files = []output.File{
		{
			Path:    ComposeOutputPath,
			Kind:    output.FileKindRegular,
			Content: compose.String(),
		},
		{
			Path:        EntrypointOutputPath,
			Kind:        output.FileKindRegular,
			Permissions: 0o755,
			Content:     entrypoint.String(),
		},
	}

	if data.ProxySQL != nil {
		files = append(files, output.File{
			Path:    ProxySQLOutputPath,
			Kind:    output.FileKindRegular,
			Content: proxySQL.String(),
		})
	}

	for _, tenant := range configuration.Tenants {
		if tenant.DatabaseDeployment != "container" {
			continue
		}

		schema, err := schemaLoader.LoadMySQLSchema(tenant.RADIUSServer.Version)
		if err != nil {
			return Deployment{}, fmt.Errorf(
				"load MySQL schema for tenant %q: %w",
				tenant.Identifier,
				err,
			)
		}

		files = append(files, output.File{
			Path:    filepath.Join(tenant.Identifier, "database", "schema.sql"),
			Kind:    output.FileKindRegular,
			Content: schema,
		})
	}

	return Deployment{
		Files: files,
	}, nil
}
