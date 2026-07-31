package templates

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoaderSelectsTemplateForSupportedVersion(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.9/mods-available/sql": &fstest.MapFile{Data: []byte("version=3.2.9 host={{.Host}}\n")},
		"3.4.0/mods-available/sql": &fstest.MapFile{Data: []byte("version=3.4.0 host={{.Host}}\n")},
	}}

	tests := []struct {
		version string
		want    string
	}{
		{version: "3.2.9", want: "version=3.2.9 host=db.example.com\n"},
		{version: "3.4.0", want: "version=3.4.0 host=db.example.com\n"},
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

func TestLoaderReturnsErrors(t *testing.T) {
	loader := loader{files: fstest.MapFS{
		"3.2.9/invalid": &fstest.MapFile{Data: []byte("{{if}}")},
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
			version:  "3.2.9",
			template: "../sql",
			wantErr:  "template name \"../sql\" is invalid",
		},
		{
			name:     "missing template",
			version:  "3.2.9",
			template: "mods-available/sql",
			wantErr:  "load template \"mods-available/sql\" for FreeRADIUS version \"3.2.9\"",
		},
		{
			name:     "invalid template",
			version:  "3.2.9",
			template: "invalid",
			wantErr:  "load template \"invalid\" for FreeRADIUS version \"3.2.9\"",
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

func TestEmbeddedLoaderReturnsMissingTemplateError(t *testing.T) {
	_, err := EmbeddedLoader().Load("3.2.9", "mods-available/sql")
	if err == nil {
		t.Fatal("EmbeddedLoader().Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "load template \"mods-available/sql\" for FreeRADIUS version \"3.2.9\"") {
		t.Fatalf("EmbeddedLoader().Load() error = %q, want missing template error", err)
	}
}
