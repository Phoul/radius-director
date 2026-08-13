package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/templates"
)

func testRenderer(t *testing.T) Renderer {
	t.Helper()

	directory, err := filepath.Abs(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	return New(templates.NewLoader(os.DirFS(directory)))
}

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

	files, err := testRenderer(t).Render(tenant)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := []string{
		"clients.conf",
		"clients.d/radius-director.conf",
		"mods-available/sql",
		"mods-config/files/authorize",
		"mods-config/files/authorize.d/radius-director",
		"mods-enabled/sql",
		"proxy.conf",
		"proxy.d/radius-director.conf",
		"sites-available/coa",
		"sites-available/default",
		"sites-enabled/coa",
		"sites-enabled/default",
		"users",
	}

	if len(files) != len(expected) {
		t.Fatalf("Render() returned %d files, want %d", len(files), len(expected))
	}

	for i, want := range expected {
		if files[i].RelativePath != want {
			t.Errorf("files[%d].RelativePath = %q, want %q",
				i, files[i].RelativePath, want)
		}

		switch files[i].Kind {
		case RenderedFileKindRegular:
			if files[i].Content == "" {
				t.Errorf("files[%d].Content is empty", i)
			}

		case RenderedFileKindSymlink:
			if files[i].Target == "" {
				t.Errorf("files[%d].Target is empty", i)
			}

		default:
			t.Errorf("files[%d].Kind = %v, want a known kind", i, files[i].Kind)
		}
	}
}

func TestRenderWithSymlink(t *testing.T) {
	directory := t.TempDir()

	templateDirectory := filepath.Join(directory, "sets", "3.2.10", "default")
	if err := os.MkdirAll(templateDirectory, 0o755); err != nil {
		t.Fatalf("create template directory: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(templateDirectory, "real.conf"),
		[]byte("real = true\n"),
		0o644,
	); err != nil {
		t.Fatalf("write real.conf: %v", err)
	}

	if err := os.Symlink(
		"real.conf",
		filepath.Join(templateDirectory, "link.conf"),
	); err != nil {
		t.Fatalf("create link.conf: %v", err)
	}

	loader := templates.NewLoader(os.DirFS(directory))
	renderer := New(loader)

	tenant := generator.Tenant{
		Template: "default",
		RADIUSServer: generator.RADIUSServer{
			Version: "3.2.10",
		},
	}

	files, err := renderer.Render(tenant)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Render() returned %d files, want 2", len(files))
	}

	if files[0].RelativePath != "link.conf" {
		t.Errorf(
			"files[0].RelativePath = %q, want %q",
			files[0].RelativePath,
			"link.conf",
		)
	}

	if files[0].Kind != RenderedFileKindSymlink {
		t.Errorf(
			"files[0].Kind = %v, want %v",
			files[0].Kind,
			RenderedFileKindSymlink,
		)
	}

	if files[0].Target != "real.conf" {
		t.Errorf(
			"files[0].Target = %q, want %q",
			files[0].Target,
			"real.conf",
		)
	}

	if files[0].Content != "" {
		t.Errorf(
			"files[0].Content = %q, want empty",
			files[0].Content,
		)
	}

	if files[1].RelativePath != "real.conf" {
		t.Errorf(
			"files[1].RelativePath = %q, want %q",
			files[1].RelativePath,
			"real.conf",
		)
	}

	if files[1].Kind != RenderedFileKindRegular {
		t.Errorf(
			"files[1].Kind = %v, want %v",
			files[1].Kind,
			RenderedFileKindRegular,
		)
	}

	if files[1].Content != "real = true\n" {
		t.Errorf(
			"files[1].Content = %q, want %q",
			files[1].Content,
			"real = true\n",
		)
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

	files, err := testRenderer(t).Render(tenant)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var found bool

	for _, file := range files {
		if file.RelativePath != "overlay-test.conf" {
			continue
		}

		found = true

		if file.Kind != RenderedFileKindRegular {
			t.Errorf("file.Kind = %v, want %v",
				file.Kind, RenderedFileKindRegular)
		}

		got := strings.ReplaceAll(file.Content, "\r\n", "\n")
		want := "# This file exists only to verify external overlay resolution.\noverlay_test = \"customer-a\""

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
