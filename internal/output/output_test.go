package output

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

func TestBuildCreatesClientsFile(t *testing.T) {
	configuration := testConfiguration()
	wantContent, err := renderer.RenderClients(configuration)
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}

	generated, err := Build(configuration)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	want := Output{
		Files: []File{
			{Path: ClientsFile, Content: wantContent},
		},
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
	want := errors.New("render clients failed")
	originalRenderClients := renderClients
	renderClients = func(generator.Configuration) (string, error) {
		return "", want
	}
	t.Cleanup(func() {
		renderClients = originalRenderClients
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
				Identifier: "customer-a",
				Clients: []generator.Client{
					{
						Identifier:   "core-router",
						IPAddress:    "10.10.10.1",
						SharedSecret: "shared-secret",
					},
				},
			},
		},
	}
}
