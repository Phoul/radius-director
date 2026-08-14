package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/output"
	"github.com/gobcn/radius-director/internal/schemas"
	"github.com/gobcn/radius-director/internal/templates"
)

func TestGenerateDeployment(t *testing.T) {
	directory, err := filepath.Abs(filepath.Join("..", "..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	loader := templates.NewLoader(os.DirFS(directory))

	schemaDirectory, err := filepath.Abs(filepath.Join("..", "..", "..", "schemas"))
	if err != nil {
		t.Fatalf("resolve schema directory: %v", err)
	}

	schemaLoader := schemas.NewLoader(os.DirFS(schemaDirectory))

	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier:         "customer-a",
				DatabaseDeployment: "container",
				SQL: generator.SQL{
					Engine:   "mysql",
					Database: "customer_a",
					Username: "radius",
					Password: "database-password",
				},
				RADIUSServer: generator.RADIUSServer{
					Version: "3.2.10",
				},
			},
			{
				Identifier:         "customer-b",
				DatabaseDeployment: "external",
				RADIUSServer: generator.RADIUSServer{
					Version: "3.3.0",
				},
			},
		},
	}

	deployment, err := GenerateDeployment(loader, schemaLoader, configuration)
	if err != nil {
		t.Fatalf("GenerateDeployment() error = %v", err)
	}

	if len(deployment.Files) != 3 {
		t.Fatalf("GenerateDeployment() returned %d files, want 3", len(deployment.Files))
	}

	files := make(map[string]output.File)

	for _, file := range deployment.Files {
		files[file.Path] = file
	}

	compose, ok := files[ComposeOutputPath]
	if !ok {
		t.Fatalf("generated deployment does not contain %q", ComposeOutputPath)
	}

	expectedComposeContent := []string{
		"radius-customer-a:",
		"image: freeradius/freeradius-server:3.2.10",
		"radius-customer-b:",
		"image: freeradius/freeradius-server:3.3.0",

		"database-customer-a:",
		"image: mysql:8.4",
		"MYSQL_RANDOM_ROOT_PASSWORD: \"yes\"",
		"MYSQL_DATABASE: customer_a",
		"MYSQL_USER: radius",
		"MYSQL_PASSWORD: database-password",
		"database-customer-a-data:/var/lib/mysql",
		"./customer-a/database/schema.sql:/docker-entrypoint-initdb.d/01-radius-schema.sql:ro",

		"condition: service_healthy",
	}

	for _, value := range expectedComposeContent {
		if !strings.Contains(compose.Content, value) {
			t.Errorf("generated Docker Compose file does not contain %q:\n%s", value, compose.Content)
		}
	}

	if strings.Contains(compose.Content, "database-customer-b:") {
		t.Error("generated Docker Compose file unexpectedly contains database-customer-b")
	}

	if strings.Contains(compose.Content, "database-customer-b-data:") {
		t.Error("generated Docker Compose file unexpectedly contains database-customer-b-data")
	}

	entrypoint, ok := files[EntrypointOutputPath]
	if !ok {
		t.Fatalf("generated deployment does not contain %q", EntrypointOutputPath)
	}

	if entrypoint.Permissions != 0o755 {
		t.Errorf("entrypoint permissions = %o, want 755", entrypoint.Permissions)
	}

	if !strings.Contains(entrypoint.Content, "#!/bin/sh") {
		t.Errorf("generated entrypoint does not contain shell shebang:\n%s", entrypoint.Content)
	}

	if !strings.Contains(entrypoint.Content, "freeradius -f") {
		t.Errorf("generated entrypoint does not contain FreeRADIUS startup command:\n%s", entrypoint.Content)
	}

	schemaPath := filepath.Join("customer-a", "database", "schema.sql")

	schema, ok := files[schemaPath]
	if !ok {
		t.Fatalf("generated deployment does not contain %q", schemaPath)
	}

	if !strings.Contains(schema.Content, "CREATE TABLE IF NOT EXISTS radacct") {
		t.Errorf(
			"generated MySQL schema does not contain radacct table definition:\n%s",
			schema.Content,
		)
	}

	externalSchemaPath := filepath.Join("customer-b", "database", "schema.sql")

	if _, ok := files[externalSchemaPath]; ok {
		t.Errorf("generated deployment unexpectedly contains %q", externalSchemaPath)
	}
}
