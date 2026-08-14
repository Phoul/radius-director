package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobcn/radius-director/internal/output"
)

func TestWriteSingleFile(t *testing.T) {
	root := t.TempDir()
	generated := output.Output{Files: []output.File{
		{Path: "clients.conf", Content: "client router {}\n"},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, filepath.Join(root, "clients.conf"), "client router {}\n")
}

func TestWriteSymlink(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join("..", "mods-available", "sql")

	generated := output.Output{Files: []output.File{
		{
			Path:   "mods-enabled/sql",
			Kind:   output.FileKindSymlink,
			Target: target,
		},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	path := filepath.Join(root, "mods-enabled", "sql")

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Lstat() error = %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("path mode = %v, want symbolic link", info.Mode())
	}

	actualTarget, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("Readlink() error = %v", err)
	}

	if actualTarget != target {
		t.Fatalf("symlink target = %q, want %q",
			actualTarget, target)
	}
}

func TestWriteReplacesSymlinkWithFile(t *testing.T) {
	root := t.TempDir()

	path := filepath.Join(root, "config")

	if err := os.Symlink("old.conf", path); err != nil {
		t.Fatalf("Symlink() setup error = %v", err)
	}

	generated := output.Output{Files: []output.File{
		{
			Path:    "config",
			Content: "new content",
		},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, path, "new content")

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Lstat() error = %v", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("config is still a symbolic link")
	}
}

func TestWriteReplacesFileWithSymlink(t *testing.T) {
	root := t.TempDir()

	path := filepath.Join(root, "config")

	if err := os.WriteFile(path, []byte("old content"), 0o644); err != nil {
		t.Fatalf("WriteFile() setup error = %v", err)
	}

	generated := output.Output{Files: []output.File{
		{
			Path:   "config",
			Kind:   output.FileKindSymlink,
			Target: "new.conf",
		},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Lstat() error = %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("config is not a symbolic link")
	}

	target, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("Readlink() error = %v", err)
	}

	if target != "new.conf" {
		t.Fatalf("symlink target = %q, want %q", target, "new.conf")
	}
}

func TestWriteReplacesSymlinkWithSymlink(t *testing.T) {
	root := t.TempDir()

	path := filepath.Join(root, "config")

	if err := os.Symlink("old.conf", path); err != nil {
		t.Fatalf("Symlink() setup error = %v", err)
	}

	generated := output.Output{Files: []output.File{
		{
			Path:   "config",
			Kind:   output.FileKindSymlink,
			Target: "new.conf",
		},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Lstat() error = %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("config is not a symbolic link")
	}

	target, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("Readlink() error = %v", err)
	}

	if target != "new.conf" {
		t.Fatalf("symlink target = %q, want %q", target, "new.conf")
	}
}

func TestWriteWithRelativeRoot(t *testing.T) {
	t.Chdir(t.TempDir())
	generated := output.Output{Files: []output.File{
		{Path: "clients.conf", Content: "client router {}\n"},
	}}

	if err := Write(".", generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, "clients.conf", "client router {}\n")
}

func TestWriteCreatesNestedDirectories(t *testing.T) {
	root := t.TempDir()
	generated := output.Output{Files: []output.File{
		{Path: "mods-enabled/sql/config", Content: "database configuration"},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, filepath.Join(root, "mods-enabled", "sql", "config"), "database configuration")
}

func TestWriteOverwritesExistingFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "clients.conf")
	if err := os.WriteFile(path, []byte("old content"), 0o644); err != nil {
		t.Fatalf("WriteFile() setup error = %v", err)
	}

	generated := output.Output{Files: []output.File{
		{Path: "clients.conf", Content: "new content"},
	}}
	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, path, "new content")
}

func TestWriteMultipleFiles(t *testing.T) {
	root := t.TempDir()
	generated := output.Output{Files: []output.File{
		{Path: "clients.conf", Content: "clients"},
		{Path: "mods-enabled/sql", Content: "sql"},
		{Path: "sites-enabled/default", Content: "server default"},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, filepath.Join(root, "clients.conf"), "clients")
	assertFileContent(t, filepath.Join(root, "mods-enabled", "sql"), "sql")
	assertFileContent(t, filepath.Join(root, "sites-enabled", "default"), "server default")
}

func TestWritePreservesFileContent(t *testing.T) {
	root := t.TempDir()
	content := "line one\n\n  indented line\n\x00binary byte\n"
	generated := output.Output{Files: []output.File{
		{Path: "exact.conf", Content: content},
	}}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(t, filepath.Join(root, "exact.conf"), content)
}

func TestWritePerTenantManifests(t *testing.T) {
	root := t.TempDir()

	generated := output.Output{
		Tenants: []string{"customer-a", "customer-b"},
		Files: []output.File{
			{Path: "customer-a/clients.conf", Content: "client router {}\n"},
			{Path: "customer-b/clients.conf", Content: "client other {}\n"},
		},
		Remove: []string{
			"customer-a/sites-enabled/inner-tunnel",
			"customer-a/mods-enabled/foo",
			"customer-b/sites-enabled/inner-tunnel",
		},
	}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(
		t,
		filepath.Join(root, "customer-a", ".radius-director", "manifest.yaml"),
		"remove:\n    - sites-enabled/inner-tunnel\n    - mods-enabled/foo\n",
	)

	assertFileContent(
		t,
		filepath.Join(root, "customer-b", ".radius-director", "manifest.yaml"),
		"remove:\n    - sites-enabled/inner-tunnel\n",
	)

	assertFileContent(
		t,
		filepath.Join(root, "customer-a", "clients.conf"),
		"client router {}\n",
	)

	assertFileContent(
		t,
		filepath.Join(root, "customer-b", "clients.conf"),
		"client other {}\n",
	)
}

func TestWritePerTenantManifestWithoutRemovals(t *testing.T) {
	root := t.TempDir()

	generated := output.Output{
		Tenants: []string{"customer-a"},
		Files: []output.File{
			{Path: "customer-a/clients.conf", Content: "client router {}\n"},
		},
	}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	assertFileContent(
		t,
		filepath.Join(root, "customer-a", ".radius-director", "manifest.yaml"),
		"{}\n",
	)
}

func TestWriteRejectsPathTraversal(t *testing.T) {
	root := filepath.Join(t.TempDir(), "output")
	outsidePath := filepath.Join(filepath.Dir(root), "outside.txt")
	generated := output.Output{Files: []output.File{
		{Path: "../outside.txt", Content: "must not be written"},
	}}

	if err := Write(root, generated); err == nil {
		t.Fatal("Write() error = nil, want path traversal error")
	}
	if _, err := os.Stat(outsidePath); !os.IsNotExist(err) {
		t.Fatalf("outside file state = %v, want not exist", err)
	}
}

func TestWritePropagatesFilesystemError(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "output-file")
	if err := os.WriteFile(root, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("WriteFile() setup error = %v", err)
	}

	generated := output.Output{Files: []output.File{
		{Path: "clients.conf", Content: "clients"},
	}}
	if err := Write(root, generated); err == nil {
		t.Fatal("Write() error = nil, want filesystem error")
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	if string(content) != want {
		t.Fatalf("file content = %q, want %q", content, want)
	}
}

func TestWriteUsesFilePermissions(t *testing.T) {
	root := t.TempDir()

	generated := output.Output{
		Files: []output.File{
			{
				Path:        "entrypoint.sh",
				Kind:        output.FileKindRegular,
				Content:     "#!/bin/sh\n",
				Permissions: 0o755,
			},
		},
	}

	if err := Write(root, generated); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	path := filepath.Join(root, "entrypoint.sh")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if info.Mode().Perm() != 0o755 {
		t.Errorf("file permissions = %o, want 755", info.Mode().Perm())
	}
}
