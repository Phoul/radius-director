package renderer

import (
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
)

func TestRender(t *testing.T) {
	one := 1

	tenant := generator.Tenant{
		Template: "default",
		RADIUSServer: generator.RADIUSServer{
			Version: "3.2.10",
		},
		AuthenticationPolicy: generator.AuthenticationPolicy{
			SimultaneousUse: &one,
		},
	}

	files, err := Render(tenant)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if len(files) != 9 {
		t.Fatalf("Render() returned %d files, want 9", len(files))
	}

	expected := []string{
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

	for i, want := range expected {
		if files[i].RelativePath != want {
			t.Errorf("files[%d].RelativePath = %q, want %q",
				i, files[i].RelativePath, want)
		}

		if files[i].Content == "" {
			t.Errorf("files[%d].Content is empty", i)
		}
	}
}

func TestRenderWithOverlay(t *testing.T) {
	tenant := generator.Tenant{
		Identifier: "customer-a",
		RADIUSServer: generator.RADIUSServer{
			Version: "3.2.10",
		},
		Template: "default",
		Overlays: []string{"test-overlay"},
	}

	files, err := Render(tenant)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var found bool

	for _, file := range files {
		if file.RelativePath != "overlay-test.conf" {
			continue
		}

		found = true

		got := strings.ReplaceAll(file.Content, "\r\n", "\n")
		want := "# This file exists only to verify embedded overlay resolution.\noverlay_test = \"customer-a\""

		if got != want {
			t.Fatalf(
				"overlay-test.conf content = %q, want %q",
				file.Content,
				want,
			)
		}
	}

	if !found {
		t.Fatal("Render() did not produce overlay-test.conf")
	}
}
