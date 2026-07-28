package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	contents := []byte(`global_objects:
  credential_profiles:
    default:
      shared_secret: secret
  nas_devices:
    core:
      ip_address: 10.10.10.1
      vendor: mikrotik
tenants:
  customer-a:
    database:
      engine: mysql
      host: db.example.com
      port: 3306
      database: radius
      username: radius
      password: secret
    radius_server:
      authentication_port: 1812
      accounting_port: 1813
      coa_port: 3799
    nas_assignments:
      core:
        nas_device: core
        credential_profile: default
        authentication_profile: default
        accounting_profile: default
        monitoring_profile: default
`)
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}

	configuration, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := configuration.GlobalObjects.CredentialProfiles["default"].SharedSecret; got != "secret" {
		t.Fatalf("credential profile shared secret = %q, want %q", got, "secret")
	}
	if got := configuration.Tenants["customer-a"].NASAssignments["core"].NASDevice; got != "core" {
		t.Fatalf("NAS assignment device = %q, want %q", got, "core")
	}
	if got := configuration.Tenants["customer-a"].RADIUSServer.AuthenticationPort; got != 1812 {
		t.Fatalf("RADIUS Server authentication port = %d, want 1812", got)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(path, []byte("global_objects: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load() error = nil, want YAML parsing error")
	}
}
