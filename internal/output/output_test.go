package output

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

func TestBuildCreatesClientsFileForEachTenant(t *testing.T) {
	configuration := testConfiguration()
	wantCustomerA, err := renderer.RenderClients(configuration.Tenants[0])
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	wantCustomerB, err := renderer.RenderClients(configuration.Tenants[1])
	if err != nil {
		t.Fatalf("RenderClients() error = %v", err)
	}
	wantCustomerAProxy, err := renderer.RenderProxy(configuration.Tenants[0])
	if err != nil {
		t.Fatalf("RenderProxy() error = %v", err)
	}
	wantCustomerBProxy, err := renderer.RenderProxy(configuration.Tenants[1])
	if err != nil {
		t.Fatalf("RenderProxy() error = %v", err)
	}
	wantCustomerACOA, err := renderer.RenderCOA(configuration.Tenants[0])
	if err != nil {
		t.Fatalf("RenderCOA() error = %v", err)
	}
	wantCustomerBCOA, err := renderer.RenderCOA(configuration.Tenants[1])
	if err != nil {
		t.Fatalf("RenderCOA() error = %v", err)
	}
	wantCustomerAAuthorize, err := renderer.RenderAuthorize(configuration.Tenants[0])
	if err != nil {
		t.Fatalf("RenderAuthorize() error = %v", err)
	}

	wantCustomerBAuthorize, err := renderer.RenderAuthorize(configuration.Tenants[1])
	if err != nil {
		t.Fatalf("RenderAuthorize() error = %v", err)
	}

	generated, err := Build(configuration)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	want := Output{
		Files: []File{
			{Path: filepath.Join("customer-a", ClientsFile), Content: wantCustomerA},
			{Path: filepath.Join("customer-a", COASiteFile), Content: wantCustomerACOA},
			{Path: filepath.Join("customer-a", ProxyFile), Content: wantCustomerAProxy},
			{Path: filepath.Join("customer-a", AuthorizeFile), Content: wantCustomerAAuthorize},

			{Path: filepath.Join("customer-b", ClientsFile), Content: wantCustomerB},
			{Path: filepath.Join("customer-b", COASiteFile), Content: wantCustomerBCOA},
			{Path: filepath.Join("customer-b", ProxyFile), Content: wantCustomerBProxy},
			{Path: filepath.Join("customer-b", AuthorizeFile), Content: wantCustomerBAuthorize},
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
	renderClients = func(generator.Tenant) (string, error) {
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

func TestBuildPropagatesCOARendererError(t *testing.T) {
	want := errors.New("render coa failed")
	originalRenderCOA := renderCOA
	renderCOA = func(generator.Tenant) (string, error) {
		return "", want
	}
	t.Cleanup(func() {
		renderCOA = originalRenderCOA
	})

	generated, err := Build(testConfiguration())
	if !errors.Is(err, want) {
		t.Fatalf("Build() error = %v, want %v", err, want)
	}
	if !reflect.DeepEqual(generated, Output{}) {
		t.Fatalf("Build() output = %#v, want empty output", generated)
	}
}

func TestBuildPropagatesProxyRendererError(t *testing.T) {
	want := errors.New("render proxy failed")
	originalRenderProxy := renderProxy
	renderProxy = func(generator.Tenant) (string, error) {
		return "", want
	}
	t.Cleanup(func() {
		renderProxy = originalRenderProxy
	})

	generated, err := Build(testConfiguration())
	if !errors.Is(err, want) {
		t.Fatalf("Build() error = %v, want %v", err, want)
	}
	if !reflect.DeepEqual(generated, Output{}) {
		t.Fatalf("Build() output = %#v, want empty output", generated)
	}
}

func TestBuildPropagatesAuthorizeRendererError(t *testing.T) {
	want := errors.New("render authorize failed")

	originalRenderAuthorize := renderAuthorize

	renderAuthorize = func(generator.Tenant) (string, error) {
		return "", want
	}

	t.Cleanup(func() {
		renderAuthorize = originalRenderAuthorize
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
