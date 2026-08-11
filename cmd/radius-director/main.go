// Command radius-director manages declarative FreeRADIUS configuration.
package main

import (
	"fmt"
	"os"

	"github.com/gobcn/radius-director/internal/cli"
	"github.com/gobcn/radius-director/internal/templates"
)

func main() {
	templateLoader, err := loadTemplateLibrary("./templates")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, templateLoader))
}

func loadTemplateLibrary(directory string) (templates.Loader, error) {
	info, err := os.Stat(directory)
	if err != nil {
		return templates.Loader{}, fmt.Errorf("load template directory %q: %w", directory, err)
	}
	if !info.IsDir() {
		return templates.Loader{}, fmt.Errorf("load template directory %q: not a directory", directory)
	}
	if _, err := os.ReadDir(directory); err != nil {
		return templates.Loader{}, fmt.Errorf("load template directory %q: %w", directory, err)
	}

	return templates.NewLoader(os.DirFS(directory)), nil
}
