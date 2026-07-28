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
