// Command eka is the official EKA conformance validator CLI.
//
// It wraps the reusable conformance package
// (github.com/maleolabs/engineering-knowledge-architecture/conformance) and
// formats its Report for human consumption.
//
// Usage:
//
//	eka validate [path]     validate the EKA repository rooted at path
//	                        (default: current directory)
//
// Exit codes (deterministic contract):
//
//	0  fully compliant (warnings allowed)
//	1  blocking violations present
//	2  usage or internal error (unknown command, invalid path,
//	   unreadable scan root)
//
// Warnings never affect the exit code.
//
// This implementation provides ONLY the validate command. doctor, import,
// export and graph are explicitly out of scope; any other command is
// rejected with an error mentioning the validate-only scope.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

const usageText = `eka — official EKA conformance validator

Usage:
  eka validate [path]    validate the EKA repository rooted at path
                         (default: current directory)

Exit codes:
  0  fully compliant (warnings allowed)
  1  blocking violations present
  2  usage or internal error (invalid path, unreadable root)

This implementation provides only the 'validate' command.
`

// run executes the CLI with the given arguments and returns the exit code.
// It is separated from main so tests can exercise it without a subprocess.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usageText)
		return 2
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usageText)
		return 0
	case "validate":
		// fallthrough to validation below
	default:
		fmt.Fprintf(stderr, "eka: unknown command %q — this implementation provides only the 'validate' command\n", args[0])
		return 2
	}

	if len(args) > 2 {
		fmt.Fprintf(stderr, "eka: too many arguments for 'validate'\n%s", usageText)
		return 2
	}
	path := "."
	if len(args) == 2 {
		switch args[1] {
		case "-h", "--help":
			fmt.Fprint(stdout, usageText)
			return 0
		default:
			path = args[1]
		}
	}

	report, err := conformance.Validate(path)
	if err != nil {
		fmt.Fprintf(stderr, "eka: validate failed: %v\n", err)
		return 2
	}

	printReport(stdout, report)
	if report.Pass() {
		return 0
	}
	return 1
}

// printReport renders the report to out in a deterministic format:
// scan summary, sorted results, execution summary.
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

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
