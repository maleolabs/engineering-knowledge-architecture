// Package cmd implements the EKA CLI as a thin Cobra command layer.
//
// The command tree (root, validate, init) is the only part of the
// codebase that knows about argument parsing, flags, help text, output
// rendering and exit codes. It contains no domain logic: validate
// delegates to the conformance package, init delegates to the bootstrap
// engine, both of which are public, reusable application packages.
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
//	0  fully compliant (warnings allowed); init completed and validates
//	1  blocking violations present (validate; or init produced a
//	   non-conformant repository)
//	2  usage or internal error (unknown command, invalid path,
//	   unreadable scan root)
package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

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

eka validate checks a repository against the EKA conformance rules and
eka init bootstraps a new EKA repository from the embedded skeleton,
validating the result afterwards.

Exit codes:
  0  fully compliant (warnings allowed)
  1  blocking violations present
  2  usage or internal error (unknown command, invalid path,
     unreadable root)`,
		// `eka` without a subcommand is a usage error: print the usage to
		// stderr and exit 2 (preserved pre-Cobra behavior). Cobra's
		// default for a non-runnable root would print help and exit 0,
		// so the root is runnable on purpose.
		RunE: func(cmd *cobra.Command, args []string) error {
			// UsageString + ErrOrStderr: cobra's Usage() writes through
			// OutOrStderr(), which actually resolves to the stdout
			// writer; usage must go to stderr.
			fmt.Fprint(cmd.ErrOrStderr(), cmd.UsageString())
			return &exitError{code: exitUsage}
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
	root.AddCommand(newValidateCommand(), newInitCommand())
	return root
}
