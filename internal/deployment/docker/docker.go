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
	Tenants []ComposeTenant
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

type Deployment struct {
	Files []output.File
}

const (
	ComposeOutputPath    = "docker-compose.yml"
	EntrypointOutputPath = "entrypoint.sh"
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
