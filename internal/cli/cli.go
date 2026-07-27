// Package cli provides the RADIUS Director command-line interface.
package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/gobcn/radius-director/internal/config"
	"github.com/gobcn/radius-director/internal/validation"
)

// Run executes the command-line interface and returns its exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "validate" {
		return runValidate(args[1:], stdout, stderr)
	}

	flags := flag.NewFlagSet("radius-director", flag.ContinueOnError)
	flags.SetOutput(stderr)
	help := flags.Bool("help", false, "Show this help message.")
	flags.BoolVar(help, "h", false, "Show this help message.")
	flags.Usage = func() {
		fmt.Fprintln(stdout, "RADIUS Director manages declarative FreeRADIUS configuration.")
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  radius-director validate <config.yaml>")
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Options:")
		fmt.Fprintln(stdout, "  -h, --help")
		fmt.Fprintln(stdout, "        Show this help message.")
	}

	if err := flags.Parse(args); err != nil {
		return 2
	}

	if *help {
		flags.Usage()
		return 0
	}

	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "unexpected argument: %s\n", flags.Arg(0))
		flags.Usage()
		return 2
	}

	flags.Usage()
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  radius-director validate <config.yaml>")
		return 0
	}

	if len(args) != 1 {
		fmt.Fprintln(stderr, "validate requires exactly one configuration file")
		return 2
	}

	configuration, err := config.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := validation.Validate(configuration); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintln(stdout, "Configuration parsed and validated successfully.")
	return 0
}
