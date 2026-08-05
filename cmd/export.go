package cmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
	"github.com/spf13/cobra"
)

// newExportCommand builds the `eka export` command: export an EKA
// repository (rooted at the current directory) to the Reference
// Serialization Format (RSF) v1.0 package. All export logic lives in the
// exchange engine; this command only validates flags/arguments, renders
// the outcome and maps the result to the exit code contract.
//
// Exit codes:
//
//	0  export produced a package (validation gate passed, warnings allowed)
//	1  repository validation failed: the export is refused and no package
//	   is produced (the full report is printed)
//	2  usage or internal error (malformed/unknown export target, flag
//	   error, serialization or filesystem failure)
func newExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export [target ...]",
		Short: "Export an EKA repository to an RSF package",
		Long: `Export the EKA repository rooted at the current directory to a
package in the EKA Reference Serialization Format (RSF) v1.0.

The repository is validated against the conformance rules first
(R0-R9); a repository with blocking violations is refused and no
package is produced. Warnings never block an export.

Export targets are canonical reference forms:

  <type>:<id>[:<instance-version>]        same namespace
  <namespace>/<type>:<id>[:<version>]     cross namespace

With no target the whole repository is exported (Repository scope).
One target without an instance-version exports its Artifact Line
(all instances); one versioned target exports exactly that instance;
several targets export their union (Collection scope).

The package is written as a single .ekapkg ZIP container named
rsf-<scope>-<namespace>-1.ekapkg, or as a directory layout when
--output names an existing directory or ends with a path separator.
Exports are deterministic: identical repository state yields
byte-identical packages.

Exit codes:
  0  export succeeded
  1  repository validation failed (no package produced)
  2  usage or internal error (bad target, unwritable output)`,
		Example: `  eka export
  eka export sto:login-email
  eka export adr:001-login-serialization:1
  eka export plan:rilis-1 sto:login-email
  eka export -o /tmp/package
  eka export -o my-package.ekapkg`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
			res, err := exchange.Export(".", exchange.Options{Targets: args, Output: output})
			if err != nil {
				var ve *exchange.ValidationError
				if errors.As(err, &ve) {
					// Blocking violations: print the full report and exit 1.
					printReport(cmd.OutOrStdout(), ve.Report)
					fmt.Fprintf(cmd.ErrOrStderr(), "eka: %s\n", ve.Error())
					return &exitError{code: exitFail}
				}
				return err // Mapped to exit 2 by Execute.
			}
			printExportOutcome(cmd.OutOrStdout(), res)
			return nil
		},
	}
	cmd.Flags().StringP("output", "o", "",
		"output path: a .ekapkg file, an existing directory, or a path ending in / (default: <label>.ekapkg in the current directory)")
	return cmd
}

// printExportOutcome renders the export result deterministically: summary
// header, package identity, scope, counts, destination. The exact bytes
// are part of the CLI contract.
func printExportOutcome(out io.Writer, r *exchange.Result) {
	fmt.Fprintf(out, "EKA Export\n")
	fmt.Fprintf(out, "==========\n")
	fmt.Fprintf(out, "%-20s %s\n", "Package Label:", r.Label)
	fmt.Fprintf(out, "%-20s %s\n", "Scope:", r.Package.Header.ExportScope)
	fmt.Fprintf(out, "%-20s %s\n", "Namespace:", r.Package.Header.Namespace)
	fmt.Fprintf(out, "%-20s %d\n", "Units:", r.Units)
	fmt.Fprintf(out, "%-20s %d\n", "Attachments:", r.Attachments)
	fmt.Fprintf(out, "%-20s %d\n", "External refs:", r.ExternalReferences)
	kind := "ZIP container"
	if r.Directory {
		kind = "directory"
	}
	fmt.Fprintf(out, "%-20s %s (%s)\n", "Output:", r.Output, kind)
	fmt.Fprintf(out, "%-20s %s\n", "Validation:", "PASS (0 errors, "+strconv.Itoa(r.Validation.WarningCount())+" warnings)")
	fmt.Fprintf(out, "%-20s %s\n", "Serialization:", "RSF v1.0 (serialization version 1)")
}
