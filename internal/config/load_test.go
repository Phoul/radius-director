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
  authentication_profiles:
    default:
      simultaneous_use: 1
  accounting_profiles:
    default:
      stale_session_timeout: 20m
  nas_devices:
    core:
      ip_address: 10.10.10.1
      vendor: mikrotik
  trusted_radius_clients:
    monitoring:
      ip_address: 10.10.10.2
tenants:
  customer-a:
    authentication_profile: default
    database:
      engine: mysql
      host: db.example.com
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
      core:
        nas_device: core
        credential_profile: default
        accounting_profile: default
        monitoring_profile: default
    trusted_radius_client_assignments:
      monitoring:
        trusted_radius_client: monitoring
        credential_profile: default
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
	got := configuration.GlobalObjects.AuthenticationProfiles["default"].SimultaneousUse
	if got == nil || *got != 1 {
		t.Fatalf("authentication profile simultaneous use = %d, want 1", got)
	}
	if got := configuration.Tenants["customer-a"].AuthenticationProfile; got != "default" {
		t.Fatalf("tenant authentication profile = %q, want %q", got, "default")
	}
	if got := configuration.GlobalObjects.AccountingProfiles["default"].StaleSessionTimeout; got != "20m" {
		t.Fatalf("accounting profile stale session timeout = %q, want %q", got, "20m")
	}
	if got := configuration.Tenants["customer-a"].NASAssignments["core"].NASDevice; got != "core" {
		t.Fatalf("NAS assignment device = %q, want %q", got, "core")
	}
	if got := configuration.GlobalObjects.TrustedRADIUSClients["monitoring"].IPAddress; got != "10.10.10.2" {
		t.Fatalf("trusted RADIUS client IP address = %q, want %q", got, "10.10.10.2")
	}
	if got := configuration.Tenants["customer-a"].TrustedRADIUSClientAssignments["monitoring"].TrustedRADIUSClient; got != "monitoring" {
		t.Fatalf("trusted RADIUS client assignment client = %q, want %q", got, "monitoring")
	}
	if got := configuration.Tenants["customer-a"].RADIUSServer.AuthenticationPort; got != 1812 {
		t.Fatalf("RADIUS Server authentication port = %d, want 1812", got)
	}
	if got := configuration.Tenants["customer-a"].RADIUSServer.Version; got != "3.2.10" {
		t.Fatalf("RADIUS Server version = %q, want %q", got, "3.2.10")
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
