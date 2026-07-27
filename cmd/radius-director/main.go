// Command radius-director manages declarative FreeRADIUS configuration.
package main

import (
	"os"

	"github.com/gobcn/radius-director/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
