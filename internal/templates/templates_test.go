package templates

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoaderSelectsTemplateForSupportedVersion(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/mods-available/sql": &fstest.MapFile{Data: []byte("version=3.2.10 host={{.Host}}\n")},
	}}

	tests := []struct {
		version string
		want    string
	}{
		{version: "3.2.10", want: "version=3.2.10 host=db.example.com\n"},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			executor, err := loader.Load(test.version, "default", nil, "mods-available/sql")
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			var rendered bytes.Buffer
			if err := executor.Execute(&rendered, struct{ Host string }{Host: "db.example.com"}); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if got := rendered.String(); got != test.want {
				t.Fatalf("rendered template = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSupportsVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{version: "3.2.10", want: true},
		{version: "3.2.9", want: false},
		{version: "3.4.0", want: false},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			if got := SupportsVersion(test.version); got != test.want {
				t.Fatalf("SupportsVersion(%q) = %t, want %t", test.version, got, test.want)
			}
		})
	}
}

func TestSupportsTemplateSet(t *testing.T) {
	tests := []struct {
		version     string
		templateSet string
		want        bool
	}{
		{version: "3.2.10", templateSet: "default", want: true},
		{version: "3.2.10", templateSet: "missing", want: false},
		{version: "3.2.10", templateSet: ".", want: false},
	}

	for _, test := range tests {
		t.Run(test.templateSet, func(t *testing.T) {
			if got := SupportsTemplateSet(test.version, test.templateSet); got != test.want {
				t.Fatalf("SupportsTemplateSet(%q, %q) = %t, want %t", test.version, test.templateSet, got, test.want)
			}
		})
	}
}

func TestSupportsOverlay(t *testing.T) {
	tests := []struct {
		overlay string
		want    bool
	}{
		{overlay: "test-overlay", want: true},
		{overlay: "missing", want: false},
		{overlay: ".", want: false},
	}

	for _, test := range tests {
		t.Run(test.overlay, func(t *testing.T) {
			if got := SupportsOverlay("3.2.10", test.overlay); got != test.want {
				t.Fatalf("SupportsOverlay(%q, %q) = %t, want %t", "3.2.10", test.overlay, got, test.want)
			}
		})
	}
}

func TestLoaderReturnsErrors(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/invalid": &fstest.MapFile{Data: []byte("{{if}}")},
	}}

	tests := []struct {
		name        string
		version     string
		templateSet string
		template    string
		wantErr     string
	}{
		{
			name:        "unsupported version",
			version:     "3.3.0",
			templateSet: "default",
			template:    "mods-available/sql",
			wantErr:     "FreeRADIUS version \"3.3.0\" is not supported",
		},
		{
			name:        "invalid template name",
			version:     "3.2.10",
			templateSet: "default",
			template:    "../sql",
			wantErr:     "template name \"../sql\" is invalid",
		},
		{
			name:        "missing template",
			version:     "3.2.10",
			templateSet: "default",
			template:    "mods-available/sql",
			wantErr:     "load template \"mods-available/sql\" for FreeRADIUS version \"3.2.10\"",
		},
		{
			name:        "invalid template",
			version:     "3.2.10",
			templateSet: "default",
			template:    "invalid",
			wantErr:     "load template \"invalid\" for FreeRADIUS version \"3.2.10\"",
		},
		{
			name:        "invalid template set",
			version:     "3.2.10",
			templateSet: "../default",
			template:    "mods-available/sql",
			wantErr:     "template set \"../default\" is invalid",
		},
		{
			name:        "current directory template set",
			version:     "3.2.10",
			templateSet: ".",
			template:    "mods-available/sql",
			wantErr:     "template set \".\" is invalid",
		},
		{
			name:        "missing template set",
			version:     "3.2.10",
			templateSet: "alternate",
			template:    "mods-available/sql",
			wantErr:     `template set "alternate" is not available for FreeRADIUS version "3.2.10"`,
		},
		{
			name:        "template set containing path separator",
			version:     "3.2.10",
			templateSet: "custom/default",
			template:    "mods-available/sql",
			wantErr:     `template set "custom/default" is invalid`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := loader.Load(test.version, test.templateSet, nil, test.template)
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Load() error = %q, want containing %q", err, test.wantErr)
			}
		})
	}
}

func TestEmbeddedLoaderLoadsCurrentTemplateSet(t *testing.T) {
	if _, err := EmbeddedLoader().Load("3.2.10", "default", nil, "mods-available/sql"); err != nil {
		t.Fatalf("EmbeddedLoader().Load() error = %v", err)
	}
}

func TestManagedTemplates(t *testing.T) {
	got, err := ManagedTemplates("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("ManagedTemplates() error = %v", err)
	}

	want := []string{
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

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ManagedTemplates() = %#v, want %#v", got, want)
	}
}

func TestManagedTemplatesReturnsTemplateSetErrors(t *testing.T) {
	tests := []struct {
		name        string
		templateSet string
		wantErr     string
	}{
		{
			name:        "invalid template set",
			templateSet: "../default",
			wantErr:     `template set "../default" is invalid`,
		},
		{
			name:        "current directory template set",
			templateSet: ".",
			wantErr:     `template set "." is invalid`,
		},
		{
			name:        "missing template set",
			templateSet: "alternate",
			wantErr:     `template set "alternate" is not available for FreeRADIUS version "3.2.10"`,
		},
		{
			name:        "template set containing path separator",
			templateSet: "custom/default",
			wantErr:     `template set "custom/default" is invalid`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ManagedTemplates("3.2.10", test.templateSet, nil)
			if err == nil {
				t.Fatal("ManagedTemplates() error = nil, want error")
			}

			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf(
					"ManagedTemplates() error = %q, want containing %q",
					err,
					test.wantErr,
				)
			}
		})
	}
}

func TestResolveTemplates(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/base-only": &fstest.MapFile{
			Data: []byte("base only"),
		},
		"3.2.10/default/replaced": &fstest.MapFile{
			Data: []byte("base"),
		},
		"3.2.10/overlays/first/replaced": &fstest.MapFile{
			Data: []byte("first"),
		},
		"3.2.10/overlays/first/added": &fstest.MapFile{
			Data: []byte("added"),
		},
		"3.2.10/overlays/second/replaced": &fstest.MapFile{
			Data: []byte("second"),
		},
	}}

	resolved, err := loader.resolve(
		"3.2.10",
		"default",
		[]string{"first", "second"},
	)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	want := map[string]string{
		"base-only": "3.2.10/default/base-only",
		"added":     "3.2.10/overlays/first/added",
		"replaced":  "3.2.10/overlays/second/replaced",
	}

	if !reflect.DeepEqual(resolved.files, want) {
		t.Fatalf("resolve() files = %#v, want %#v", resolved.files, want)
	}
}

func TestResolveTemplatesReturnsOverlayErrors(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/base": &fstest.MapFile{
			Data: []byte("base"),
		},
	}}

	tests := []struct {
		name     string
		overlays []string
		wantErr  string
	}{
		{
			name:     "invalid overlay",
			overlays: []string{"../test"},
			wantErr:  `overlay "../test" is invalid`,
		},
		{
			name:     "current directory overlay",
			overlays: []string{"."},
			wantErr:  `overlay "." is invalid`,
		},
		{
			name:     "overlay containing path separator",
			overlays: []string{"experiments/test"},
			wantErr:  `overlay "experiments/test" is invalid`,
		},
		{
			name:     "missing overlay",
			overlays: []string{"missing"},
			wantErr:  `overlay "missing" is not available for FreeRADIUS version "3.2.10"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := loader.resolve(
				"3.2.10",
				"default",
				test.overlays,
			)
			if err == nil {
				t.Fatal("resolve() error = nil, want error")
			}

			if !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf(
					"resolve() error = %q, want containing %q",
					err,
					test.wantErr,
				)
			}
		})
	}
}

func TestLoaderLoadsWinningOverlay(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/example": &fstest.MapFile{
			Data: []byte("base {{.Value}}"),
		},
		"3.2.10/overlays/first/example": &fstest.MapFile{
			Data: []byte("first {{.Value}}"),
		},
		"3.2.10/overlays/second/example": &fstest.MapFile{
			Data: []byte("second {{.Value}}"),
		},
	}}

	executor, err := loader.Load(
		"3.2.10",
		"default",
		[]string{"first", "second"},
		"example",
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	var rendered bytes.Buffer

	if err := executor.Execute(
		&rendered,
		struct{ Value string }{Value: "value"},
	); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := rendered.String(), "second value"; got != want {
		t.Fatalf("rendered template = %q, want %q", got, want)
	}
}

func TestManagedTemplatesIncludesOverlayFiles(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/default/base-only": &fstest.MapFile{
			Data: []byte("base"),
		},
		"3.2.10/default/replaced": &fstest.MapFile{
			Data: []byte("base"),
		},
		"3.2.10/overlays/test/replaced": &fstest.MapFile{
			Data: []byte("overlay"),
		},
		"3.2.10/overlays/test/added": &fstest.MapFile{
			Data: []byte("added"),
		},
	}}

	got, err := loader.managedTemplates(
		"3.2.10",
		"default",
		[]string{"test"},
	)
	if err != nil {
		t.Fatalf("managedTemplates() error = %v", err)
	}

	want := []string{
		"added",
		"base-only",
		"replaced",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("managedTemplates() = %#v, want %#v", got, want)
	}
}

func TestEmbeddedLoaderLoadsOverlay(t *testing.T) {
	executor, err := EmbeddedLoader().Load(
		"3.2.10",
		"default",
		[]string{"test-overlay"},
		"overlay-test.conf",
	)
	if err != nil {
		t.Fatalf("EmbeddedLoader().Load() error = %v", err)
	}

	var rendered bytes.Buffer

	if err := executor.Execute(
		&rendered,
		struct {
			Identifier string
		}{
			Identifier: "customer-a",
		},
	); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := strings.ReplaceAll(rendered.String(), "\r\n", "\n")
	want := "# This file exists only to verify embedded overlay resolution.\noverlay_test = \"customer-a\""

	if got != want {
		t.Fatalf("rendered template = %q, want %q", got, want)
	}
}

func TestManagedTemplatesIncludesEmbeddedOverlay(t *testing.T) {
	got, err := ManagedTemplates(
		"3.2.10",
		"default",
		[]string{"test-overlay"},
	)
	if err != nil {
		t.Fatalf("ManagedTemplates() error = %v", err)
	}

	want := []string{
		"clients.conf",
		"clients.d/radius-director.conf",
		"mods-available/sql",
		"mods-config/files/authorize",
		"mods-config/files/authorize.d/radius-director",
		"overlay-test.conf",
		"proxy.conf",
		"proxy.d/radius-director.conf",
		"sites-available/coa",
		"sites-available/default",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ManagedTemplates() = %#v, want %#v", got, want)
	}
}
