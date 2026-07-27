package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"--help"}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}

	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("help output = %q, want usage information", stdout.String())
	}
	if count := strings.Count(stdout.String(), "Usage:"); count != 1 {
		t.Fatalf("help output contains %d usage sections, want 1", count)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunValidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("global_objects: {}\ntenants: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if exitCode := Run([]string{"validate", path}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("Run() exit code = %d, want 0", exitCode)
	}
	if got := stdout.String(); got != "Configuration parsed successfully.\n" {
		t.Fatalf("stdout = %q, want success message", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
