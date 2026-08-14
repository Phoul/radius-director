// Package output assembles rendered FreeRADIUS configuration files.
package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobcn/radius-director/internal/generator"
	"github.com/gobcn/radius-director/internal/renderer"
)

// Output contains the generated files that are ready to be written.
type Output struct {
	Files   []File
	Remove  []string
	Tenants []string
}

// FileKind identifies the type of filesystem object in the generated output.
type FileKind int

const (
	FileKindRegular FileKind = iota
	FileKindSymlink
)

// File is a generated filesystem object.
//
// Regular files use Content. Symlinks use Target.
type File struct {
	Path        string
	Kind        FileKind
	Content     string
	Target      string
	Permissions os.FileMode
}

// Build renders each tenant's files and assembles them into an Output object.
func Build(configuration generator.Configuration, templateRenderer renderer.Renderer) (Output, error) {
	generated := Output{
		Files: make([]File, 0),
	}

	for _, tenant := range configuration.Tenants {
		generated.Tenants = append(generated.Tenants, tenant.Identifier)
		renderedFiles, err := templateRenderer.Render(tenant)
		if err != nil {
			return Output{}, err
		}

		for _, rendered := range renderedFiles {
			file := File{
				Path: filepath.Join(tenant.Identifier, rendered.RelativePath),
			}

			switch rendered.Kind {
			case renderer.RenderedFileKindRegular:
				file.Kind = FileKindRegular
				file.Content = rendered.Content

			case renderer.RenderedFileKindSymlink:
				file.Kind = FileKindSymlink
				file.Target = rendered.Target

			default:
				return Output{}, fmt.Errorf(
					"unsupported rendered file kind for %q",
					rendered.RelativePath,
				)
			}

			generated.Files = append(generated.Files, file)
		}

		for _, removePath := range tenant.Remove {
			generated.Remove = append(
				generated.Remove,
				filepath.Join(tenant.Identifier, removePath),
			)
		}
	}

	return generated, nil
}
