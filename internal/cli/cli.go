// Package cli provides the RADIUS Director command-line interface.
package cli

import (
	"flag"
	"fmt"
	"io"
)

// Run executes the command-line interface and returns its exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("radius-director", flag.ContinueOnError)
	flags.SetOutput(stderr)
	help := flags.Bool("help", false, "Show this help message.")
	flags.BoolVar(help, "h", false, "Show this help message.")
	flags.Usage = func() {
		fmt.Fprintln(stdout, "RADIUS Director manages declarative FreeRADIUS configuration.")
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  radius-director [options]")
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
