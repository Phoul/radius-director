// Package writer writes generated files to the filesystem.
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobcn/radius-director/internal/output"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

// Write writes every generated file beneath root.
func Write(root string, generated output.Output) error {
	normalizedRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	for _, file := range generated.Files {
		destination, err := filepath.Abs(filepath.Join(normalizedRoot, file.Path))
		if err != nil {
			return err
		}
		if !isWithinRoot(normalizedRoot, destination) {
			return fmt.Errorf("generated file path %q escapes output root %q", file.Path, root)
		}

		if err := os.MkdirAll(filepath.Dir(destination), directoryPermissions); err != nil {
			return err
		}
		if err := os.WriteFile(destination, []byte(file.Content), filePermissions); err != nil {
			return err
		}
	}

	return nil
}

func isWithinRoot(root, path string) bool {
	relativePath, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	return relativePath != ".." &&
		!strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) &&
		!filepath.IsAbs(relativePath)
}
