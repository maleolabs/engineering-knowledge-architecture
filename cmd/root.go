// Package cmd implements the EKA CLI as a thin Cobra command layer.
//
// The command tree (root, validate, init, export) is the only part of the
// codebase that knows about argument parsing, flags, help text, output
// rendering and exit codes. It contains no domain logic: validate
// delegates to the conformance package, init delegates to the bootstrap
// engine, export delegates to the exchange engine — all public, reusable
// application packages.
//
// Layout rationale: the reusable engines stay where they are
// (bootstrap/, conformance/, skeletonembed.go at the module root). There
// is deliberately no internal/ or pkg/ directory — those add indirection
// without serving any immediate consumer, and the task rules out
// speculative abstraction. cmd/ is a leaf: nothing imports it except
// cmd/eka/main.go.
//
// Exit codes (deterministic contract, preserved from the pre-Cobra CLI):
//
//	0  fully compliant (warnings allowed); init completed and validates;
//	   export produced a package
//	1  blocking violations present (validate; init produced a
//	   non-conformant repository; export refused the repository)
//	2  usage or internal error (unknown command, invalid path,
//	   unreadable scan root, bad export target, export failure)
package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/spf13/cobra"
)

// Exit codes of the deterministic contract documented above.
const (
	exitOK    = 0
	exitFail  = 1
	exitUsage = 2
)

// exitError carries an explicit exit code out of a command's RunE. Commands
// that must exit non-zero after printing their own diagnostics return it;
// Execute maps the code to the process exit code without printing anything.
type exitError struct{ code int }

func (e *exitError) Error() string { return fmt.Sprintf("exit code %d", e.code) }

// Execute runs the EKA CLI with the given arguments and streams and
// returns the process exit code. It is the single testable entry point:
// main() delegates to it, tests call it directly.
func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if args == nil {
		args = []string{}
	}
	root := newRootCommand()
	root.SetArgs(args)
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)

	err := root.Execute()
	if err == nil {
		return exitOK
	}
	var ee *exitError
	if errors.As(err, &ee) {
		return ee.code
	}
	fmt.Fprintf(stderr, "eka: %s\n", renderError(err))
	return exitUsage
}

// renderError formats an execution error deterministically. Cobra's
// "unknown command" errors are augmented with the list of available
// commands, mirroring the pre-Cobra CLI contract. Any other error is
// passed through verbatim ("eka: <error>").
func renderError(err error) string {
	msg := err.Error()
	if strings.HasPrefix(msg, "unknown command") {
		msg += " — available commands: " + availableCommands()
	}
	return msg
}

// availableCommands lists the registered subcommands in registration
// order, excluding the built-in help/completion commands (which a fresh
// tree does not contain yet).
func availableCommands() string {
	names := make([]string, 0, 2)
	for _, c := range newRootCommand().Commands() {
		names = append(names, c.Name())
	}
	return strings.Join(names, ", ")
}

// newRootCommand builds the root command with all subcommands registered.
// A fresh tree is built per Execute call so that SetArgs/SetIn/SetOut/
// SetErr never leak between invocations (and concurrent Executes stay
// safe).
func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "eka",
		Short: "Official EKA CLI: conformance validation and repository bootstrapping",
		Long: `The official EKA CLI.

eka validate checks a repository against the EKA conformance rules,
eka init bootstraps a new EKA repository from the embedded skeleton
(validating the result afterwards), eka export projects a repository
to a deterministic package in the EKA Reference Serialization Format
(RSF) v1.0, and eka import consumes such a package.

Command output is deterministic: the same input always produces the
same bytes. On a terminal the output is colored and progress is shown
in place; when piped or redirected it is plain text with no control
sequences. Use --verbose/-v for additional detail lines (per-unit
lists, plan actions).

Exit codes:
  0  fully compliant (warnings allowed)
  1  blocking violations present
  2  usage or internal error (unknown command, invalid path,
     unreadable root)`,
		// `eka` without a subcommand shows the product landing: a calm
		// orientation (what the CLI is, its commands, help and version
		// pointers) instead of the raw usage dump. Landing is
		// informational output — it exits 0. Unknown subcommands remain
		// usage errors (exit 2).
		RunE: func(cmd *cobra.Command, args []string) error {
			printLanding(styleFor(cmd))
			return nil
		},
		// The CLI owns all error output: SilenceErrors + SilenceUsage on
		// the root suppress cobra's "Error: …" prefix and its automatic
		// usage dumps (children inherit both flags). Execute renders
		// every error as a single deterministic "eka: …" line and maps
		// it to the exit code contract. Help output is unaffected: the
		// flag.ErrHelp path prints help to stdout and exits 0.
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.PersistentFlags().BoolP(flagVerbose, "v", false,
		"verbose output: additional detail lines (per-unit lists, plan actions)")
	root.AddCommand(newValidateCommand(), newInitCommand(), newExportCommand(), newImportCommand(), newVersionCommand())
	return root
}

// printLanding renders the root landing page: a calm product orientation
// without banners or decoration — heading, one-line description, compact
// command overview, and pointers to help and version. Deterministic on
// non-TTY output; the heading is accent-colored on a color TTY.
func printLanding(s *ui.Style) {
	fmt.Fprintln(s.W, s.Accent("Engineering Knowledge Architecture (EKA)"))
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "  The official command-line interface for the EKA engineering")
	fmt.Fprintln(s.W, "  knowledge standard: bootstrap, validate, and exchange")
	fmt.Fprintln(s.W, "  engineering knowledge repositories.")
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Commands")
	for _, c := range newRootCommand().Commands() {
		fmt.Fprintf(s.W, "  %-12s %s\n", c.Name(), c.Short)
	}
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Help")
	fmt.Fprintln(s.W, "  Run 'eka help <command>' for command details,")
	fmt.Fprintln(s.W, "  or 'eka <command> --help' for usage.")
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Version")
	fmt.Fprintf(s.W, "  %s (EKA standard %s)\n", version, standardVersion)
}
