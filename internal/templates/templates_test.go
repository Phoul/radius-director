package templates

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func externalTemplateLoader(t *testing.T) Loader {
	t.Helper()

	directory, err := filepath.Abs(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatalf("resolve template directory: %v", err)
	}

	return NewLoader(os.DirFS(directory))
}

type symlinkMapFS struct {
	fstest.MapFS
	links map[string]string
}

func (f symlinkMapFS) ReadLink(name string) (string, error) {
	target, ok := f.links[name]
	if !ok {
		return "", &fs.PathError{
			Op:   "readlink",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}

	return target, nil
}

func (f symlinkMapFS) Lstat(name string) (fs.FileInfo, error) {
	if _, ok := f.links[name]; ok {
		return symlinkFileInfo{name: path.Base(name)}, nil
	}

	return fs.Lstat(f.MapFS, name)
}

func (f symlinkMapFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(f.MapFS, name)
	if err != nil {
		return nil, err
	}

	for linkPath := range f.links {
		if path.Dir(linkPath) != name {
			continue
		}

		entries = append(entries, symlinkDirEntry{
			name: path.Base(linkPath),
		})
	}

	return entries, nil
}

func (f symlinkMapFS) Sub(dir string) (fs.FS, error) {
	return symlinkMapFSSub{
		parent: f,
		dir:    dir,
	}, nil
}

type symlinkMapFSSub struct {
	parent symlinkMapFS
	dir    string
}

func (f symlinkMapFSSub) Open(name string) (fs.File, error) {
	return f.parent.Open(path.Join(f.dir, name))
}

func (f symlinkMapFSSub) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.parent.ReadDir(path.Join(f.dir, name))
}

type symlinkDirEntry struct {
	name string
}

func (e symlinkDirEntry) Name() string {
	return e.name
}

func (e symlinkDirEntry) IsDir() bool {
	return false
}

func (e symlinkDirEntry) Type() fs.FileMode {
	return fs.ModeSymlink
}

func (e symlinkDirEntry) Info() (fs.FileInfo, error) {
	return symlinkFileInfo{name: e.name}, nil
}

type symlinkFileInfo struct {
	name string
}

func (i symlinkFileInfo) Name() string {
	return i.name
}

func (i symlinkFileInfo) Size() int64 {
	return 0
}

func (i symlinkFileInfo) Mode() fs.FileMode {
	return fs.ModeSymlink
}

func (i symlinkFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (i symlinkFileInfo) IsDir() bool {
	return false
}

func (i symlinkFileInfo) Sys() any {
	return nil
}

func TestLoaderSelectsTemplateForSupportedVersion(t *testing.T) {
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/mods-available/sql": &fstest.MapFile{Data: []byte("version=3.2.10 host={{.Host}}\n")},
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
	loader := externalTemplateLoader(t)
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
			if got := loader.SupportsVersion(test.version); got != test.want {
				t.Fatalf("loader.SupportsVersion(%q) = %t, want %t", test.version, got, test.want)
			}
		})
	}
}

func TestSupportsTemplateSet(t *testing.T) {
	loader := externalTemplateLoader(t)
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
			if got := loader.SupportsTemplateSet(test.version, test.templateSet); got != test.want {
				t.Fatalf("loader.SupportsTemplateSet(%q, %q) = %t, want %t", test.version, test.templateSet, got, test.want)
			}
		})
	}
}

func TestSupportsOverlay(t *testing.T) {
	loader := externalTemplateLoader(t)
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
			if got := loader.SupportsOverlay("3.2.10", test.overlay); got != test.want {
				t.Fatalf("loader.SupportsOverlay(%q, %q) = %t, want %t", "3.2.10", test.overlay, got, test.want)
			}
		})
	}
}

func TestLoaderReturnsErrors(t *testing.T) {
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/invalid": &fstest.MapFile{Data: []byte("{{if}}")},
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

func TestExternalLoaderLoadsCurrentTemplateSet(t *testing.T) {
	if _, err := externalTemplateLoader(t).Load("3.2.10", "default", nil, "mods-available/sql"); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestManagedTemplates(t *testing.T) {
	got, err := externalTemplateLoader(t).ManagedTemplates("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("ManagedTemplates() error = %v", err)
	}

	want := []string{
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
			_, err := externalTemplateLoader(t).ManagedTemplates("3.2.10", test.templateSet, nil)
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
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/base-only": &fstest.MapFile{
			Data: []byte("base only"),
		},
		"sets/3.2.10/default/replaced": &fstest.MapFile{
			Data: []byte("base"),
		},
		"overlays/3.2.10/first/replaced": &fstest.MapFile{
			Data: []byte("first"),
		},
		"overlays/3.2.10/first/added": &fstest.MapFile{
			Data: []byte("added"),
		},
		"overlays/3.2.10/second/replaced": &fstest.MapFile{
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

	want := map[string]resolvedEntry{
		"base-only": {
			kind:   resolvedFile,
			source: "sets/3.2.10/default/base-only",
		},
		"added": {
			kind:   resolvedFile,
			source: "overlays/3.2.10/first/added",
		},
		"replaced": {
			kind:   resolvedFile,
			source: "overlays/3.2.10/second/replaced",
		},
	}

	if !reflect.DeepEqual(resolved.entries, want) {
		t.Fatalf("resolve() entries = %#v, want %#v", resolved.entries, want)
	}
}

func TestResolveTemplatesWithSymlink(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/real.conf": &fstest.MapFile{
					Data: []byte("real"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "./real.conf",
			},
		},
	}

	resolved, err := loader.resolve("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	want := resolvedEntry{
		kind:   resolvedSymlink,
		target: "./real.conf",
	}

	if got := resolved.entries["link.conf"]; !reflect.DeepEqual(got, want) {
		t.Fatalf("resolved.entries[link.conf] = %#v, want %#v", got, want)
	}
}

func TestResolveTemplatesNormalizesNativeSymlinkTarget(t *testing.T) {
	nativeTarget := "." + string(filepath.Separator) + "real.conf"

	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/real.conf": &fstest.MapFile{
					Data: []byte("real"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": nativeTarget,
			},
		},
	}

	resolved, err := loader.resolve("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	want := resolvedEntry{
		kind:   resolvedSymlink,
		target: "./real.conf",
	}

	if got := resolved.entries["link.conf"]; !reflect.DeepEqual(got, want) {
		t.Fatalf("resolved.entries[link.conf] = %#v, want %#v", got, want)
	}
}

func TestResolveTemplatesWithContainedParentSymlink(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/shared/real.conf": &fstest.MapFile{
					Data: []byte("real"),
				},
				"sets/3.2.10/default/sub": &fstest.MapFile{
					Mode: fs.ModeDir,
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/sub/link.conf": "../shared/real.conf",
			},
		},
	}

	resolved, err := loader.resolve("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	if got := resolved.entries["sub/link.conf"]; got.target != "../shared/real.conf" {
		t.Fatalf("symlink target = %q, want %q", got.target, "../shared/real.conf")
	}
}

func TestResolveTemplatesRejectsAbsoluteSymlink(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "/etc/freeradius/real.conf",
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `symlink target "/etc/freeradius/real.conf" is absolute`) {
		t.Fatalf("resolve() error = %q, want absolute-target error", err)
	}
}

func TestResolveTemplatesRejectsEscapingSymlink(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "../../outside",
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `escapes template layer`) {
		t.Fatalf("resolve() error = %q, want containment error", err)
	}
}

func TestResolveTemplatesRejectsMissingSymlinkTarget(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "./missing.conf",
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `does not exist`) {
		t.Fatalf("resolve() error = %q, want missing-target error", err)
	}
}

func TestResolveTemplatesWithSymlinkChain(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/real.conf": &fstest.MapFile{
					Data: []byte("real"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/a": "./b",
				"sets/3.2.10/default/b": "./real.conf",
			},
		},
	}

	resolved, err := loader.resolve("3.2.10", "default", nil)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	if got := resolved.entries["a"]; got.target != "./b" {
		t.Fatalf("a target = %q, want %q", got.target, "./b")
	}
}

func TestResolveTemplatesRejectsSymlinkCycle(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/a": "./b",
				"sets/3.2.10/default/b": "./a",
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "symlink cycle detected") {
		t.Fatalf("resolve() error = %q, want cycle error", err)
	}
}

func TestResolveTemplatesRejectsEmptySymlinkTarget(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "",
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "symlink target is empty") {
		t.Fatalf("resolve() error = %q, want empty-target error", err)
	}
}

func TestResolveTemplatesRejectsWindowsAbsoluteSymlinkTarget(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/base": &fstest.MapFile{
					Data: []byte("base"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": `C:/Windows/System32/something`,
			},
		},
	}

	_, err := loader.resolve("3.2.10", "default", nil)
	if err == nil {
		t.Fatal("resolve() error = nil, want error")
	}

	if !strings.Contains(err.Error(), `is absolute`) {
		t.Fatalf("resolve() error = %q, want Windows absolute-path error", err)
	}
}

func TestResolveTemplatesOverlayReplacesSymlinkTargetFile(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/real.conf": &fstest.MapFile{
					Data: []byte("base"),
				},
				"overlays/3.2.10/customer-a/real.conf": &fstest.MapFile{
					Data: []byte("overlay"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link.conf": "./real.conf",
			},
		},
	}

	resolved, err := loader.resolve(
		"3.2.10",
		"default",
		[]string{"customer-a"},
	)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	if got := resolved.entries["link.conf"]; got.kind != resolvedSymlink ||
		got.target != "./real.conf" {
		t.Fatalf("link.conf = %#v, want symlink to ./real.conf", got)
	}

	if got := resolved.entries["real.conf"]; got.kind != resolvedFile ||
		got.source != "overlays/3.2.10/customer-a/real.conf" {
		t.Fatalf("real.conf = %#v, want overlay file", got)
	}
}

func TestResolveTemplatesOverlayReplacesSymlink(t *testing.T) {
	loader := Loader{
		files: symlinkMapFS{
			MapFS: fstest.MapFS{
				"sets/3.2.10/default/one": &fstest.MapFile{
					Data: []byte("one"),
				},
				"sets/3.2.10/default/two": &fstest.MapFile{
					Data: []byte("two"),
				},
				"overlays/3.2.10/customer-a/two": &fstest.MapFile{
					Data: []byte("overlay two"),
				},
			},
			links: map[string]string{
				"sets/3.2.10/default/link":        "./one",
				"overlays/3.2.10/customer-a/link": "./two",
			},
		},
	}

	resolved, err := loader.resolve(
		"3.2.10",
		"default",
		[]string{"customer-a"},
	)
	if err != nil {
		t.Fatalf("resolve() error = %v", err)
	}

	got := resolved.entries["link"]

	if got.kind != resolvedSymlink || got.target != "./two" {
		t.Fatalf("link = %#v, want symlink to ./two", got)
	}
}

func TestResolveTemplatesReturnsOverlayErrors(t *testing.T) {
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/base": &fstest.MapFile{
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
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/example": &fstest.MapFile{
			Data: []byte("base {{.Value}}"),
		},
		"overlays/3.2.10/first/example": &fstest.MapFile{
			Data: []byte("first {{.Value}}"),
		},
		"overlays/3.2.10/second/example": &fstest.MapFile{
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
	loader := Loader{files: fstest.MapFS{
		"sets/3.2.10/default/base-only": &fstest.MapFile{
			Data: []byte("base"),
		},
		"sets/3.2.10/default/replaced": &fstest.MapFile{
			Data: []byte("base"),
		},
		"overlays/3.2.10/test/replaced": &fstest.MapFile{
			Data: []byte("overlay"),
		},
		"overlays/3.2.10/test/added": &fstest.MapFile{
			Data: []byte("added"),
		},
	}}

	got, err := loader.ManagedTemplates(
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

func TestExternalLoaderLoadsOverlay(t *testing.T) {
	executor, err := externalTemplateLoader(t).Load(
		"3.2.10",
		"default",
		[]string{"test-overlay"},
		"overlay-test.conf",
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
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
	want := "# This file exists only to verify external overlay resolution.\noverlay_test = \"customer-a\""

	if got != want {
		t.Fatalf("rendered template = %q, want %q", got, want)
	}
}

func TestManagedTemplatesIncludesExternalOverlay(t *testing.T) {
	got, err := externalTemplateLoader(t).ManagedTemplates(
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
		"mods-enabled/sql",
		"overlay-test.conf",
		"proxy.conf",
		"proxy.d/radius-director.conf",
		"sites-available/coa",
		"sites-available/default",
		"sites-enabled/coa",
		"sites-enabled/default",
		"users",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ManagedTemplates() = %#v, want %#v", got, want)
	}
}
