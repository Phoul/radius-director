package assets

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExport(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	createTestAssetTree(t, source)

	if err := Export(source, destination); err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	expectedFiles := map[string]string{
		filepath.Join("templates", "default", "radiusd.conf"):     "radiusd configuration",
		filepath.Join("templates", "default", "sites", "default"): "default site",
		filepath.Join("schemas", "default.sql"):                   "CREATE TABLE radcheck;",
	}

	for relativePath, expectedContent := range expectedFiles {
		path := filepath.Join(destination, relativePath)

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read exported file %q: %v", relativePath, err)
		}

		if string(content) != expectedContent {
			t.Errorf(
				"exported file %q = %q, want %q",
				relativePath,
				string(content),
				expectedContent,
			)
		}
	}
}

func TestExportPreservesSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	source := t.TempDir()
	destination := t.TempDir()

	createTestAssetTree(t, source)

	exportedLink := filepath.Join(
		destination,
		"templates",
		"default",
		"sites-enabled",
		"default",
	)

	info, err := os.Lstat(exportedLink)
	if err != nil {
		t.Fatalf("Lstat() exported symlink: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf(
			"exported path %q is not a symlink",
			exportedLink,
		)
	}

	target, err := os.Readlink(exportedLink)
	if err != nil {
		t.Fatalf("Readlink() exported symlink: %v", err)
	}

	if target != "../../sites/default" {
		t.Errorf(
			"exported symlink target = %q, want %q",
			target,
			"../../sites/default",
		)
	}
}

func TestExportRefusesExistingTemplatesDirectory(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	createTestAssetTree(t, source)

	existingTemplates := filepath.Join(destination, "templates")
	if err := os.MkdirAll(existingTemplates, 0o755); err != nil {
		t.Fatalf("create existing templates directory: %v", err)
	}

	customFile := filepath.Join(existingTemplates, "custom.conf")
	const customContent = "administrator customization"

	if err := os.WriteFile(customFile, []byte(customContent), 0o644); err != nil {
		t.Fatalf("write custom template: %v", err)
	}

	err := Export(source, destination)
	if err == nil {
		t.Fatal("Export() succeeded, want error for existing templates directory")
	}

	content, err := os.ReadFile(customFile)
	if err != nil {
		t.Fatalf("read existing custom file: %v", err)
	}

	if string(content) != customContent {
		t.Errorf(
			"existing custom file was modified: got %q, want %q",
			string(content),
			customContent,
		)
	}
}

func TestExportRefusesExistingSchemasDirectory(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	createTestAssetTree(t, source)

	existingSchemas := filepath.Join(destination, "schemas")
	if err := os.MkdirAll(existingSchemas, 0o755); err != nil {
		t.Fatalf("create existing schemas directory: %v", err)
	}

	customFile := filepath.Join(existingSchemas, "custom.sql")
	const customContent = "administrator schema customization"

	if err := os.WriteFile(customFile, []byte(customContent), 0o644); err != nil {
		t.Fatalf("write custom schema: %v", err)
	}

	err := Export(source, destination)
	if err == nil {
		t.Fatal("Export() succeeded, want error for existing schemas directory")
	}

	content, err := os.ReadFile(customFile)
	if err != nil {
		t.Fatalf("read existing custom schema: %v", err)
	}

	if string(content) != customContent {
		t.Errorf(
			"existing custom schema was modified: got %q, want %q",
			string(content),
			customContent,
		)
	}
}

func TestExportFailsWhenTemplatesSourceDoesNotExist(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	if err := os.MkdirAll(filepath.Join(source, "schemas"), 0o755); err != nil {
		t.Fatalf("create schemas directory: %v", err)
	}

	err := Export(source, destination)
	if err == nil {
		t.Fatal("Export() succeeded, want error when templates source is missing")
	}
}

func TestExportFailsWhenSchemasSourceDoesNotExist(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()

	if err := os.MkdirAll(filepath.Join(source, "templates"), 0o755); err != nil {
		t.Fatalf("create templates directory: %v", err)
	}

	err := Export(source, destination)
	if err == nil {
		t.Fatal("Export() succeeded, want error when schemas source is missing")
	}
}

func createTestAssetTree(t *testing.T, source string) {
	t.Helper()

	files := map[string]string{
		filepath.Join("templates", "default", "radiusd.conf"):     "radiusd configuration",
		filepath.Join("templates", "default", "sites", "default"): "default site",
		filepath.Join("schemas", "default.sql"):                   "CREATE TABLE radcheck;",
	}

	for relativePath, content := range files {
		path := filepath.Join(source, relativePath)

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create directory for %q: %v", relativePath, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write test asset %q: %v", relativePath, err)
		}
	}

	if runtime.GOOS != "windows" {
		linkPath := filepath.Join(
			source,
			"templates",
			"default",
			"sites-enabled",
			"default",
		)

		if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
			t.Fatalf("create symlink directory: %v", err)
		}

		if err := os.Symlink("../../sites/default", linkPath); err != nil {
			t.Fatalf("create test symlink: %v", err)
		}
	}
}
