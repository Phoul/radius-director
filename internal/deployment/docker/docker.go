package docker

import (
	"bytes"
	"fmt"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/output"
	"github.com/gobcn/radius-director/internal/templates"
)

type ComposeData struct {
	Tenants []ComposeTenant
}

type ComposeTenant struct {
	Name               string
	ServiceName        string
	RadiusVersion      string
	AuthenticationPort int
	AccountingPort     int
	COAPort            int
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
			Name:               tenant.Identifier,
			ServiceName:        "radius-" + tenant.Identifier,
			RadiusVersion:      tenant.RADIUSServer.Version,
			AuthenticationPort: tenant.RADIUSServer.AuthenticationPort,
			AccountingPort:     tenant.RADIUSServer.AccountingPort,
			COAPort:            tenant.RADIUSServer.COAPort,
		})
	}

	return data
}

func GenerateDeployment(
	loader templates.Loader,
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

	return Deployment{
		Files: []output.File{
			{
				Path:    ComposeOutputPath,
				Kind:    output.FileKindRegular,
				Content: compose.String(),
			},
			{
				Path:    EntrypointOutputPath,
				Kind:    output.FileKindRegular,
				Content: entrypoint.String(),
			},
		},
	}, nil
}
