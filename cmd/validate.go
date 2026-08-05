package cmd

import (
	"fmt"
	"io"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
	"github.com/spf13/cobra"
)

// newValidateCommand builds the `eka validate` command: conformance
// validation of the repository rooted at an optional path (default: the
// current directory). All validation logic lives in the conformance
// package; this command only validates arguments, renders the report and
// maps the result to the exit code contract.
func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate an EKA repository against the conformance rules",
		Long: `Validate the EKA repository rooted at path against the conformance
rules.

With no path, the current directory is validated. Warnings never affect
the exit code; blocking violations exit 1. Usage or internal errors
(unknown flag, too many arguments, unreadable root) exit 2.`,
		Example: `  eka validate              validate the current directory
  eka validate docs          validate the repository rooted at docs`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			report, err := conformance.Validate(path)
			if err != nil {
				return fmt.Errorf("validate failed: %w", err)
			}
			printReport(cmd.OutOrStdout(), report)
			if !report.Pass() {
				return &exitError{code: exitFail}
			}
			return nil
		},
	}
}

// printReport renders the report to out in a deterministic format:
// scan summary, sorted results, execution summary. The exact bytes are
// part of the CLI contract (preserved from the pre-Cobra implementation).
func printReport(out io.Writer, r *conformance.Report) {
	fmt.Fprintf(out, "EKA Conformance Validation\n")
	fmt.Fprintf(out, "==========================\n")
	fmt.Fprintf(out, "Root:      %s\n", r.Root)
	fmt.Fprintf(out, ".md files: %d\n", r.FilesScanned)
	fmt.Fprintf(out, "Artifacts: %d\n", r.Artifacts)
	fmt.Fprintf(out, "Errors:    %d\n", r.ErrorCount())
	fmt.Fprintf(out, "Warnings:  %d\n", r.WarningCount())

	fmt.Fprintf(out, "\nResults (sorted by file, then rule):\n")
	results := r.SortedResults()
	if len(results) == 0 {
		fmt.Fprintf(out, "  (no violations found)\n")
	} else {
		for _, res := range results {
			fmt.Fprintf(out, "  [%s] %s %s: %s\n", res.Severity, res.Rule, res.File, res.Message)
		}
	}

	verdict := "PASS"
	if !r.Pass() {
		verdict = "FAIL"
	}
	fmt.Fprintf(out, "\nExecution: %s (%d errors, %d warnings)\n",
		verdict, r.ErrorCount(), r.WarningCount())
}
