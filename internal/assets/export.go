// Package assets exports the RADIUS Director assets bundled with a release.
package assets

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

// Export copies the bundled templates and schemas into destination.
//
// Existing templates or schemas directories are never overwritten.
func Export(sourceRoot, destination string) error {
	destination, err := filepath.Abs(destination)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destination, directoryPermissions); err != nil {
		return err
	}

	for _, name := range []string{"templates", "schemas"} {
		source := filepath.Join(sourceRoot, name)
		target := filepath.Join(destination, name)

		if _, err := os.Lstat(target); err == nil {
			return fmt.Errorf(
				"destination already contains %q; refusing to overwrite",
				name,
			)
		} else if !os.IsNotExist(err) {
			return err
		}

		if err := copyTree(source, target); err != nil {
			return fmt.Errorf("export %s: %w", name, err)
		}
	}

	return nil
}

func copyTree(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}

	return copyEntry(source, destination, info)
}

func copyEntry(source, destination string, info os.FileInfo) error {
	switch {
	case info.IsDir():
		if err := os.MkdirAll(destination, directoryPermissions); err != nil {
			return err
		}

		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			sourcePath := filepath.Join(source, entry.Name())
			destinationPath := filepath.Join(destination, entry.Name())

			entryInfo, err := os.Lstat(sourcePath)
			if err != nil {
				return err
			}

			if err := copyEntry(sourcePath, destinationPath, entryInfo); err != nil {
				return err
			}
		}

		return os.Chmod(destination, info.Mode().Perm())

	case info.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}

		return os.Symlink(target, destination)

	case info.Mode().IsRegular():
		sourceFile, err := os.Open(source)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destinationFile, err := os.OpenFile(
			destination,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL,
			filePermissions,
		)
		if err != nil {
			return err
		}

		if _, err := io.Copy(destinationFile, sourceFile); err != nil {
			destinationFile.Close()
			return err
		}

		if err := destinationFile.Close(); err != nil {
			return err
		}

		return os.Chmod(destination, info.Mode().Perm())

	default:
		return fmt.Errorf(
			"unsupported filesystem entry %q",
			source,
		)
	}
}
