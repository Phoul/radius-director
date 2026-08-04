package renderer

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

func TestRenderClientsOneTenantOneClient(t *testing.T) {
	tenant := generator.Tenant{
		Identifier:   "customer-a",
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
		FreeRADIUSClients: []generator.FreeRADIUSClient{
			{
				Identifier:   "core-router",
				IPAddress:    "10.10.10.1",
				SharedSecret: "shared-secret",
			},
		},
	}

	got, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	want := "client core-router {\n" +
		"    ipaddr = 10.10.10.1\n" +
		"    secret = shared-secret\n" +
		"    nas_type = other\n" +
		"    coa_server = concentrators\n" +
		"}"
	if !strings.Contains(normalizeLineEndings(got), "# Tenant: customer-a") {
		t.Fatalf("RenderClients() output does not identify the tenant")
	}
	if !strings.Contains(normalizeLineEndings(got), want) {
		t.Fatalf("RenderClients() output does not contain %q", want)
	}
}

func TestRenderClientsOneTenantMultipleClients(t *testing.T) {
	tenant := generator.Tenant{
		Identifier:   "customer-a",
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
		FreeRADIUSClients: []generator.FreeRADIUSClient{
			{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "core-secret", Vendor: "mikrotik"},
			{Identifier: "edge-router", IPAddress: "10.10.10.2", SharedSecret: "edge-secret", Vendor: "generic"},
		},
	}

	got, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	coreClient := "client core-router {\n" +
		"    ipaddr = 10.10.10.1\n" +
		"    secret = core-secret\n" +
		"    nas_type = mikrotik_snmp\n" +
		"    coa_server = concentrators\n" +
		"}"
	edgeClient := "client edge-router {\n" +
		"    ipaddr = 10.10.10.2\n" +
		"    secret = edge-secret\n" +
		"    nas_type = other\n" +
		"    coa_server = concentrators\n" +
		"}"
	normalized := normalizeLineEndings(got)
	coreIndex := strings.Index(normalized, coreClient)
	edgeIndex := strings.Index(normalized, edgeClient)
	if coreIndex < 0 || edgeIndex < 0 {
		t.Fatalf("RenderClients() output does not contain both expected client definitions")
	}
	if coreIndex > edgeIndex {
		t.Fatalf("RenderClients() rendered clients out of order")
	}
}

func TestRenderClientsIsDeterministic(t *testing.T) {
	tenant := generator.Tenant{
		Identifier:        "tenant-a",
		RADIUSServer:      generator.RADIUSServer{Version: "3.2.10"},
		FreeRADIUSClients: []generator.FreeRADIUSClient{{Identifier: "a", IPAddress: "10.0.0.1", SharedSecret: "a-secret"}},
	}

	first, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("first RenderClients() error = %v", err)
	}
	second, err := RenderClients(tenant)
	if err != nil {
		t.Fatalf("second RenderClients() error = %v", err)
	}
	if first != second {
		t.Fatalf("RenderClients() returned different results: %q and %q", first, second)
	}
}

func TestRenderClientsUnsupportedVersion(t *testing.T) {
	_, err := RenderClients(generator.Tenant{
		RADIUSServer: generator.RADIUSServer{Version: "9.9.9"},
	})
	if err == nil || !strings.Contains(err.Error(), `FreeRADIUS version "9.9.9" is not supported`) {
		t.Fatalf("RenderClients() error = %v, want unsupported version error", err)
	}
}

func TestRenderClientsPropagatesTemplateExecutionError(t *testing.T) {
	executionError := errors.New("template execution failed")
	useClientsTemplateLoader(t, testTemplateLoader{
		executor: testTemplateExecutor{err: executionError},
	})

	_, err := RenderClients(generator.Tenant{
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
	})
	if !errors.Is(err, executionError) {
		t.Fatalf("RenderClients() error = %v, want %v", err, executionError)
	}
}

func TestRenderSQL(t *testing.T) {
	tenant := generator.Tenant{
		SQL: generator.SQL{
			Engine:   "postgresql",
			Host:     "db.example.com",
			Port:     5432,
			Database: "radius",
			Username: "radius-user",
			Password: "radius-password",
		},
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
	}

	got, err := RenderSQL(tenant)
	if err != nil {
		t.Fatalf("RenderSQL() error = %v", err)
	}

	for _, expected := range []string{
		`dialect = "postgresql"`,
		`server = "db.example.com"`,
		"port = 5432",
		`login = "radius-user"`,
		`password = "radius-password"`,
		`radius_db = "radius"`,
	} {
		if !strings.Contains(got, expected) {
			t.Errorf("RenderSQL() output does not contain %q", expected)
		}
	}
}

func TestRenderSQLIsDeterministic(t *testing.T) {
	tenant := generator.Tenant{
		SQL:          generator.SQL{Engine: "mysql", Host: "db", Port: 3306},
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
	}

	first, err := RenderSQL(tenant)
	if err != nil {
		t.Fatalf("first RenderSQL() error = %v", err)
	}
	second, err := RenderSQL(tenant)
	if err != nil {
		t.Fatalf("second RenderSQL() error = %v", err)
	}
	if first != second {
		t.Fatalf("RenderSQL() returned different results")
	}
}

func TestRenderSQLUnsupportedVersion(t *testing.T) {
	_, err := RenderSQL(generator.Tenant{
		RADIUSServer: generator.RADIUSServer{Version: "9.9.9"},
	})
	if err == nil || !strings.Contains(err.Error(), `FreeRADIUS version "9.9.9" is not supported`) {
		t.Fatalf("RenderSQL() error = %v, want unsupported version error", err)
	}
}

func TestRenderSQLPropagatesMissingTemplateError(t *testing.T) {
	missingTemplate := errors.New("template does not exist")
	useSQLTemplateLoader(t, testTemplateLoader{err: missingTemplate})

	_, err := RenderSQL(generator.Tenant{
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
	})
	if !errors.Is(err, missingTemplate) {
		t.Fatalf("RenderSQL() error = %v, want %v", err, missingTemplate)
	}
}

func TestRenderSQLPropagatesTemplateExecutionError(t *testing.T) {
	executionError := errors.New("template execution failed")
	useSQLTemplateLoader(t, testTemplateLoader{
		executor: testTemplateExecutor{err: executionError},
	})

	_, err := RenderSQL(generator.Tenant{
		RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
	})
	if !errors.Is(err, executionError) {
		t.Fatalf("RenderSQL() error = %v, want %v", err, executionError)
	}
}

func useSQLTemplateLoader(t *testing.T, loader templates.Loader) {
	t.Helper()
	original := sqlTemplateLoader
	sqlTemplateLoader = loader
	t.Cleanup(func() { sqlTemplateLoader = original })
}

func normalizeLineEndings(value string) string {
	return strings.ReplaceAll(value, "\r\n", "\n")
}

func useClientsTemplateLoader(t *testing.T, loader templates.Loader) {
	t.Helper()
	original := clientsTemplateLoader
	clientsTemplateLoader = loader
	t.Cleanup(func() { clientsTemplateLoader = original })
}

type testTemplateLoader struct {
	executor templates.Executor
	err      error
}

func (l testTemplateLoader) Load(version, name string) (templates.Executor, error) {
	return l.executor, l.err
}

type testTemplateExecutor struct {
	err error
}

func (e testTemplateExecutor) Execute(io.Writer, any) error {
	return e.err
}
