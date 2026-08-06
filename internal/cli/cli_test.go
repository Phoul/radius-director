package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"--help"}, &stdout, &stderr); exitCode != 0 {
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

	if exitCode := Run([]string{"validate", path}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}
	if got := stdout.String(); got != "Configuration parsed and validated successfully.\n" {
		t.Fatalf("stdout = %q, want success message", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunMaintenanceAccountingHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"maintenance", "accounting", "--help"}, &stdout, &stderr); exitCode != 0 {
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

	if exitCode := Run([]string{"maintenance", "accounting", path, "missing"}, &stdout, &stderr); exitCode != 1 {
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

	if exitCode := Run([]string{"maintenance", "accounting", path, "customer-a"}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0; stderr = %q", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "no stale-session maintenance policies are enabled") {
		t.Fatalf("stdout = %q, want no enabled policies message", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
