package output

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

func TestBuildCreatesFilesForEachTenant(t *testing.T) {
	configuration := testConfiguration()

	generated, err := Build(configuration)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	want := Output{}

	for _, tenant := range configuration.Tenants {
		rendered, err := renderer.Render(tenant)
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

	first, err := Build(configuration)
	if err != nil {
		t.Fatalf("first Build() error = %v", err)
	}

	second, err := Build(configuration)
	if err != nil {
		t.Fatalf("second Build() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("Build() returned different results: %#v and %#v", first, second)
	}
}

func TestBuildPropagatesRendererError(t *testing.T) {
	want := errors.New("render failed")

	originalRender := render

	render = func(generator.Tenant) ([]renderer.RenderedFile, error) {
		return nil, want
	}

	t.Cleanup(func() {
		render = originalRender
	})

	generated, err := Build(testConfiguration())

	if !errors.Is(err, want) {
		t.Fatalf("Build() error = %v, want %v", err, want)
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
