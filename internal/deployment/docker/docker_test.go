package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/output"
	"github.com/gobcn/radius-director/internal/templates"
)

func TestGenerateDeployment(t *testing.T) {
	directory, err := filepath.Abs(filepath.Join("..", "..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	loader := templates.NewLoader(os.DirFS(directory))

	configuration := generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier: "customer-a",
				RADIUSServer: generator.RADIUSServer{
					Version: "3.2.10",
				},
			},
			{
				Identifier: "customer-b",
				RADIUSServer: generator.RADIUSServer{
					Version: "3.3.0",
				},
			},
		},
	}

	deployment, err := GenerateDeployment(loader, configuration)
	if err != nil {
		t.Fatalf("GenerateDeployment() error = %v", err)
	}

	if len(deployment.Files) != 2 {
		t.Fatalf("GenerateDeployment() returned %d files, want 2", len(deployment.Files))
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
	}

	for _, value := range expectedComposeContent {
		if !strings.Contains(compose.Content, value) {
			t.Errorf("generated Docker Compose file does not contain %q:\n%s", value, compose.Content)
		}
	}

	entrypoint, ok := files[EntrypointOutputPath]
	if !ok {
		t.Fatalf("generated deployment does not contain %q", EntrypointOutputPath)
	}

	if !strings.Contains(entrypoint.Content, "#!/bin/sh") {
		t.Errorf("generated entrypoint does not contain shell shebang:\n%s", entrypoint.Content)
	}

	if !strings.Contains(entrypoint.Content, "freeradius -f") {
		t.Errorf("generated entrypoint does not contain FreeRADIUS startup command:\n%s", entrypoint.Content)
	}
}
