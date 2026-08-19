// Command radius-director manages declarative FreeRADIUS configuration.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobcn/radius-director/internal/assets"
	"github.com/gobcn/radius-director/internal/cli"
	"github.com/gobcn/radius-director/internal/runtime"
	"github.com/gobcn/radius-director/internal/schemas"
	"github.com/gobcn/radius-director/internal/templates"
)

func main() {
	var templateLoader templates.Loader
	var schemaLoader schemas.Loader

	var runtimeInit cli.RuntimeInitializer

	if commandNeedsDocker(os.Args[1:]) {
		networkClient, err := runtime.NewDockerNetworkClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer networkClient.Close()

		runtimeInit = runtimeInitializer{
			networkClient: networkClient,
		}
	}

	if commandNeedsLibraries(os.Args[1:]) {
		adminRoot := assets.AdminRoot()

		adminTemplates := filepath.Join(adminRoot, "templates")
		adminSchemas := filepath.Join(adminRoot, "schemas")

		if err := validateAdminAssetLibrary(adminTemplates, adminSchemas); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var err error

		templateLoader, err = loadTemplateLibrary(adminTemplates)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		schemaLoader, err = loadSchemaLibrary(adminSchemas)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Exit(cli.Run(
		os.Args[1:],
		os.Stdout,
		os.Stderr,
		templateLoader,
		schemaLoader,
		runtimeInit,
	))
}

func commandNeedsLibraries(args []string) bool {
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "validate", "generate", "maintenance":
		return true
	default:
		return false
	}
}

func commandNeedsDocker(args []string) bool {
	if len(args) < 2 || args[0] != "init" {
		return false
	}

	return args[1] != "-h" && args[1] != "--help"
}

func validateAdminAssetLibrary(templatesPath, schemasPath string) error {
	required := []string{
		templatesPath,
		schemasPath,
	}

	for _, directory := range required {
		info, err := os.Stat(directory)
		if err != nil {
			return fmt.Errorf(
				"administrator asset library is incomplete: required directory %q: %w",
				directory,
				err,
			)
		}

		if !info.IsDir() {
			return fmt.Errorf(
				"administrator asset library is incomplete: %q is not a directory",
				directory,
			)
		}
	}

	return nil
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

func loadSchemaLibrary(directory string) (schemas.Loader, error) {
	info, err := os.Stat(directory)
	if err != nil {
		return schemas.Loader{}, fmt.Errorf("load schema directory %q: %w", directory, err)
	}
	if !info.IsDir() {
		return schemas.Loader{}, fmt.Errorf("load schema directory %q: not a directory", directory)
	}
	if _, err := os.ReadDir(directory); err != nil {
		return schemas.Loader{}, fmt.Errorf("load schema directory %q: %w", directory, err)
	}

	return schemas.NewLoader(os.DirFS(directory)), nil
}

type runtimeInitializer struct {
	networkClient runtime.NetworkClient
}

func (r runtimeInitializer) Init(
	ctx context.Context,
	root string,
	networkName string,
) error {
	return runtime.Init(ctx, root, networkName, r.networkClient)
}
