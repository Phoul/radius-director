package output

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
	"github.com/gobcn/radius-director/internal/templates"
)

func testRenderer(t *testing.T) renderer.Renderer {
	t.Helper()

	directory, err := filepath.Abs(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	return renderer.New(templates.NewLoader(os.DirFS(directory)))
}

func TestBuildCreatesFilesForEachTenant(t *testing.T) {
	configuration := testConfiguration()
	templateRenderer := testRenderer(t)

	generated, err := Build(configuration, templateRenderer)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	want := Output{}

	for _, tenant := range configuration.Tenants {
		rendered, err := templateRenderer.Render(tenant)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		for _, file := range rendered {
			want.Files = append(want.Files, File{
				Path:    filepath.Join(tenant.Identifier, file.RelativePath),
				Content: file.Content,
			})
		}
	}

	if !reflect.DeepEqual(generated, want) {
		t.Fatalf("Build() = %#v, want %#v", generated, want)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	configuration := testConfiguration()
	templateRenderer := testRenderer(t)

	first, err := Build(configuration, templateRenderer)
	if err != nil {
		t.Fatalf("first Build() error = %v", err)
	}

	second, err := Build(configuration, templateRenderer)
	if err != nil {
		t.Fatalf("second Build() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("Build() returned different results: %#v and %#v", first, second)
	}
}

func TestBuildPropagatesRendererError(t *testing.T) {
	generated, err := Build(
		testConfiguration(),
		renderer.New(templates.NewLoader(fstest.MapFS{})),
	)

	if err == nil {
		t.Fatal("Build() error = nil, want error")
	}

	if !reflect.DeepEqual(generated, Output{}) {
		t.Fatalf("Build() output = %#v, want empty output", generated)
	}
}

func testConfiguration() generator.Configuration {
	return generator.Configuration{
		Tenants: []generator.Tenant{
			{
				Identifier:   "customer-a",
				Template:     "default",
				RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
				FreeRADIUSClients: []generator.FreeRADIUSClient{
					{
						Identifier:   "core-router",
						IPAddress:    "10.10.10.1",
						SharedSecret: "shared-secret",
					},
				},
				HomeServers: []generator.HomeServer{
					{Identifier: "core-router", IPAddress: "10.10.10.1", SharedSecret: "shared-secret"},
				},
			},
			{
				Identifier:   "customer-b",
				Template:     "default",
				RADIUSServer: generator.RADIUSServer{Version: "3.2.10"},
				FreeRADIUSClients: []generator.FreeRADIUSClient{
					{
						Identifier:   "edge-router",
						IPAddress:    "10.20.20.1",
						SharedSecret: "other-secret",
					},
				},
				HomeServers: []generator.HomeServer{
					{Identifier: "edge-router", IPAddress: "10.20.20.1", SharedSecret: "other-secret"},
				},
			},
		},
	}
}
