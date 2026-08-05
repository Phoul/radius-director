package renderer

import (
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
)

func TestRender(t *testing.T) {
	one := 1

	tenant := generator.Tenant{
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

	if len(files) != 5 {
		t.Fatalf("Render() returned %d files, want 5", len(files))
	}

	expected := []string{
		"clients.conf",
		"mods-available/sql",
		"mods-config/files/authorize",
		"proxy.conf",
		"sites-available/coa",
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
