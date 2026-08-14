// Package writer writes generated files to the filesystem.
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobcn/radius-director/internal/output"
	"gopkg.in/yaml.v3"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

type manifest struct {
	Remove []string `yaml:"remove,omitempty"`
}

// Write writes every generated file beneath root.
func Write(root string, generated output.Output) error {
	normalizedRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	if err := writeManifests(normalizedRoot, generated); err != nil {
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

		if file.Kind != output.FileKindRegular && file.Kind != output.FileKindSymlink {
			return fmt.Errorf("unsupported output file kind for %q", file.Path)
		}

		if err := removeExisting(destination); err != nil {
			return err
		}

		switch file.Kind {
		case output.FileKindRegular:
			permissions := file.Permissions
			if permissions == 0 {
				permissions = filePermissions
			}

			if err := os.WriteFile(destination, []byte(file.Content), permissions); err != nil {
				return err
			}

		case output.FileKindSymlink:
			if err := os.Symlink(file.Target, destination); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported output file kind for %q", file.Path)
		}
	}

	return nil
}

func writeManifests(root string, generated output.Output) error {
	removalsByTenant := make(map[string][]string)

	for _, removePath := range generated.Remove {
		parts := strings.SplitN(filepath.ToSlash(removePath), "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid generated removal path %q", removePath)
		}

		tenant := parts[0]
		relativePath := parts[1]

		removalsByTenant[tenant] = append(
			removalsByTenant[tenant],
			relativePath,
		)
	}

	for _, tenant := range generated.Tenants {
		manifestPath := filepath.Join(
			root,
			tenant,
			".radius-director",
			"manifest.yaml",
		)

		if err := os.MkdirAll(filepath.Dir(manifestPath), directoryPermissions); err != nil {
			return err
		}

		content, err := yaml.Marshal(manifest{
			Remove: removalsByTenant[tenant],
		})
		if err != nil {
			return err
		}

		if err := os.WriteFile(manifestPath, content, filePermissions); err != nil {
			return err
		}
	}

	return nil
}

func removeExisting(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("cannot replace directory %q", path)
	}

	return os.Remove(path)
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
