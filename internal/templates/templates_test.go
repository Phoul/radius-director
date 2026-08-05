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
		"3.2.10/mods-available/sql": &fstest.MapFile{Data: []byte("version=3.2.10 host={{.Host}}\n")},
	}}

	tests := []struct {
		version string
		want    string
	}{
		{version: "3.2.10", want: "version=3.2.10 host=db.example.com\n"},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			executor, err := loader.Load(test.version, "mods-available/sql")
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

func TestLoaderReturnsErrors(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.10/invalid": &fstest.MapFile{Data: []byte("{{if}}")},
	}}

	tests := []struct {
		name     string
		version  string
		template string
		wantErr  string
	}{
		{
			name:     "unsupported version",
			version:  "3.3.0",
			template: "mods-available/sql",
			wantErr:  "FreeRADIUS version \"3.3.0\" is not supported",
		},
		{
			name:     "invalid template name",
			version:  "3.2.10",
			template: "../sql",
			wantErr:  "template name \"../sql\" is invalid",
		},
		{
			name:     "missing template",
			version:  "3.2.10",
			template: "mods-available/sql",
			wantErr:  "load template \"mods-available/sql\" for FreeRADIUS version \"3.2.10\"",
		},
		{
			name:     "invalid template",
			version:  "3.2.10",
			template: "invalid",
			wantErr:  "load template \"invalid\" for FreeRADIUS version \"3.2.10\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := loader.Load(test.version, test.template)
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
	if _, err := EmbeddedLoader().Load("3.2.10", "mods-available/sql"); err != nil {
		t.Fatalf("EmbeddedLoader().Load() error = %v", err)
	}
}

func TestManagedTemplates(t *testing.T) {
	got, err := ManagedTemplates("3.2.10")
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
