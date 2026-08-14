package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/deployment/docker"
	"github.com/gobcn/radius-director/internal/templates"
)

func testTemplateLoader(t *testing.T) templates.Loader {
	t.Helper()

	directory, err := filepath.Abs(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	return templates.NewLoader(os.DirFS(directory))
}

func exampleConfigPath(t *testing.T) string {
	t.Helper()

	path, err := filepath.Abs(filepath.Join("..", "..", "examples", "example.yaml"))
	if err != nil {
		t.Fatalf("resolve example configuration: %v", err)
	}

	return path
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"--help"}, &stdout, &stderr, testTemplateLoader(t)); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}

	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("help output = %q, want usage information", stdout.String())
	}
	if count := strings.Count(stdout.String(), "Usage:"); count != 1 {
		t.Fatalf("help output contains %d usage sections, want 1", count)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunValidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("global_objects: {}\ntenants: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"validate", path}, &stdout, &stderr, testTemplateLoader(t)); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}
	if got := stdout.String(); got != "Configuration parsed and validated successfully.\n" {
		t.Fatalf("stdout = %q, want success message", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunGenerate(t *testing.T) {
	configPath := exampleConfigPath(t)
	outputDir := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run(
		[]string{"generate", configPath, outputDir},
		&stdout,
		&stderr,
		testTemplateLoader(t),
	); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0; stderr = %q", exitCode, stderr.String())
	}

	if stdout.String() != "Configuration generated successfully.\n" {
		t.Fatalf("stdout = %q, want success message", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	tenantRoot := filepath.Join(outputDir, "customer-a")

	expectedFiles := []string{
		"clients.conf",
		"clients.d/radius-director.conf",
		"mods-available/sql",
		"mods-config/files/authorize",
		"mods-config/files/authorize.d/radius-director",
		"proxy.conf",
		"proxy.d/radius-director.conf",
		"sites-available/coa",
		"sites-available/default",
	}

	for _, relativePath := range expectedFiles {
		path := filepath.Join(tenantRoot, filepath.FromSlash(relativePath))

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected generated file %q: %v", relativePath, err)
			continue
		}

		if !info.Mode().IsRegular() {
			t.Errorf("generated path %q is not a regular file", relativePath)
		}
	}

	expectedSymlinks := map[string]string{
		"mods-enabled/sql":      filepath.Join("..", "mods-available", "sql"),
		"sites-enabled/coa":     filepath.Join("..", "sites-available", "coa"),
		"sites-enabled/default": filepath.Join("..", "sites-available", "default"),
		"users":                 filepath.Join("mods-config", "files", "authorize"),
	}

	for relativePath, wantTarget := range expectedSymlinks {
		path := filepath.Join(tenantRoot, filepath.FromSlash(relativePath))

		target, err := os.Readlink(path)
		if err != nil {
			t.Errorf("expected generated symlink %q: %v", relativePath, err)
			continue
		}

		if target != wantTarget {
			t.Errorf("symlink %q = %q, want %q", relativePath, target, wantTarget)
		}
	}

	composePath := filepath.Join(outputDir, docker.ComposeOutputPath)

	composeContent, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("expected generated Docker Compose file: %v", err)
	}

	compose := string(composeContent)

	expectedComposeContent := []string{
		"radius-customer-a:",
		"image: freeradius/freeradius-server:3.2.10",
	}

	for _, value := range expectedComposeContent {
		if !strings.Contains(compose, value) {
			t.Errorf("generated Docker Compose file does not contain %q:\n%s", value, compose)
		}
	}

	entrypointPath := filepath.Join(outputDir, docker.EntrypointOutputPath)

	entrypointContent, err := os.ReadFile(entrypointPath)
	if err != nil {
		t.Fatalf("expected generated Docker entrypoint: %v", err)
	}

	entrypoint := string(entrypointContent)

	if !strings.Contains(entrypoint, "#!/bin/sh") {
		t.Errorf("generated entrypoint does not contain shell shebang:\n%s", entrypoint)
	}

	if !strings.Contains(entrypoint, "freeradius -f") {
		t.Errorf("generated entrypoint does not contain FreeRADIUS startup command:\n%s", entrypoint)
	}
}

func TestRunMaintenanceAccountingHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"maintenance", "accounting", "--help"}, &stdout, &stderr, testTemplateLoader(t)); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "maintenance accounting <config.yaml> <tenant>") {
		t.Fatalf("stdout = %q, want accounting maintenance usage", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunMaintenanceAccountingUnknownTenant(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("global_objects: {}\ntenants: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"maintenance", "accounting", path, "missing"}, &stdout, &stderr, testTemplateLoader(t)); exitCode != 1 {
		t.Fatalf("Run() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), `tenant "missing" does not exist`) {
		t.Fatalf("stderr = %q, want missing tenant error", stderr.String())
	}
}

func TestRunMaintenanceAccountingNoEnabledPolicies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	configuration := `global_objects:
  credential_profiles:
    default:
      shared_secret: secret
  authentication_profiles:
    default: {}
  accounting_profiles:
    default: {}
  monitoring_profiles:
    default: {}
  deployment_profiles:
    default:
      template: default
  nas_devices:
    router:
      ip_address: 192.0.2.1
      vendor: mikrotik
  trusted_radius_clients: {}
tenants:
  customer-a:
    authentication_profile: default
    deployment_profile: default
    database:
      engine: mysql
      host: db
      port: 3306
      database: radius
      username: radius
      password: secret
    radius_server:
      version: 3.2.10
      authentication_port: 1812
      accounting_port: 1813
      coa_port: 3799
    nas_assignments:
      router:
        nas_device: router
        credential_profile: default
        accounting_profile: default
        monitoring_profile: default
    trusted_radius_client_assignments: {}
`
	if err := os.WriteFile(path, []byte(configuration), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"maintenance", "accounting", path, "customer-a"}, &stdout, &stderr, testTemplateLoader(t)); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0; stderr = %q", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "no stale-session maintenance policies are enabled") {
		t.Fatalf("stdout = %q, want no enabled policies message", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
