package cmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
	"github.com/spf13/cobra"
)

// newImportCommand builds the `eka import` command: import an RSF package
// (.ekapkg ZIP container or directory layout) into the EKA repository
// rooted at the current directory. All import logic lives in the exchange
// engine; this command validates arguments, renders the outcome and maps
// the result to the exit code contract.
//
// Exit codes:
//
//	0  import succeeded (artifacts imported, or every artifact was a no-op
//	   duplicate)
//	1  repository-side failure: target not an EKA repository, repository
//	   failed validation (pre-import gate or post-commit revalidation with
//	   rollback), conflicts, unresolved non-draft relationships
//	2  usage or internal error (missing/extra arguments, unreadable or
//	   malformed package, integrity failure, unsupported versions,
//	   filesystem failure)
func newImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <package-path>",
		Short: "Import an RSF package into the current repository",
		Long: `Import an EKA package in the Reference Serialization Format (RSF)
v1.0 into the EKA repository rooted at the current directory.

The package (a single-file .ekapkg ZIP container or the equivalent
directory layout) is verified first: entry structure, strict JSON
decoding (unknown fields are rejected), SHA-256 integrity (package,
per-unit, per-attachment digests) and manifest self-consistency.
The import then runs the Exchange Contract phases: contract
validation, identity, state, structure, referential resolution
(local -> global -> declared external, with draft tolerance),
conflict detection (reject-by-default) and duplicate detection
(identical artifacts are skipped as no-ops).

The integration strategy is a conservative merge: only NEW artifacts
are written; identical duplicates are skipped; any difference on an
existing identity is a conflict and aborts the import with a summary
before anything is written. Nothing is ever overwritten or deleted.
After the commit the repository is revalidated; on any failure every
written file is rolled back.

Exit codes:
  0  import succeeded (or all no-op)
  1  repository invalid, conflicts, unresolved relationships, or
     post-commit revalidation failure (rolled back)
  2  package invalid/malformed/unsupported, usage or fs failure`,
		Example: `  eka import rsf-repo-acme-1.ekapkg
  eka import /tmp/package-layout/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := exchange.Import(args[0], exchange.ImportOptions{Root: "."})
			if err != nil {
				var ive *exchange.ImportValidationError
				if errors.As(err, &ive) {
					if ive.Report != nil {
						printReport(cmd.OutOrStdout(), ive.Report)
					}
					fmt.Fprintf(cmd.ErrOrStderr(), "eka: %s\n", ive.Error())
					return &exitError{code: exitFail}
				}
				var ce *exchange.ConflictError
				if errors.As(err, &ce) {
					fmt.Fprintf(cmd.ErrOrStderr(), "eka: %s\n", ce.Error())
					for _, c := range ce.Conflicts {
						fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", c.String())
					}
					return &exitError{code: exitFail}
				}
				var re *exchange.RelationshipError
				if errors.As(err, &re) {
					fmt.Fprintf(cmd.ErrOrStderr(), "eka: %s\n", re.Error())
					for _, d := range re.Details {
						fmt.Fprintf(cmd.ErrOrStderr(), "  - %s\n", d)
					}
					return &exitError{code: exitFail}
				}
				return err // Mapped to exit 2 by Execute.
			}
			printImportOutcome(cmd.OutOrStdout(), res)
			return nil
		},
	}
	return cmd
}

// printImportOutcome renders the import result deterministically: summary
// header, repository, package identity, versions verified, verdict counts,
// validation status (pre + post), next steps. The exact bytes are part of
// the CLI contract.
func printImportOutcome(out io.Writer, r *exchange.ImportResult) {
	fmt.Fprintf(out, "EKA Import\n")
	fmt.Fprintf(out, "==========\n")
	fmt.Fprintf(out, "%-20s %s\n", "Repository Root:", r.Root)
	fmt.Fprintf(out, "%-20s %s\n", "Package Label:", r.PackageLabel)
	fmt.Fprintf(out, "%-20s %s\n", "Versions:", "RSF v1.0 (serialization 1, exchange format 1, specification 1.0)")
	fmt.Fprintf(out, "%-20s %d\n", "Imported:", len(r.ImportedArtifacts))
	fmt.Fprintf(out, "%-20s %d\n", "Skipped (no-op):", len(r.SkippedArtifacts))
	fmt.Fprintf(out, "%-20s %d\n", "Attachments:", len(r.AttachmentsImported))
	fmt.Fprintf(out, "%-20s %d\n", "Attachments skipped:", len(r.AttachmentsSkipped))
	fmt.Fprintf(out, "%-20s %d\n", "Conflicts:", len(r.Conflicts))
	fmt.Fprintf(out, "%-20s %s\n", "Validation (pre):", "PASS ("+strconv.Itoa(r.PreValidation.ErrorCount())+" errors, "+strconv.Itoa(r.PreValidation.WarningCount())+" warnings)")
	fmt.Fprintf(out, "%-20s %s\n", "Validation (post):", "PASS ("+strconv.Itoa(r.Validation.ErrorCount())+" errors, "+strconv.Itoa(r.Validation.WarningCount())+" warnings)")
	if len(r.Warnings) > 0 {
		fmt.Fprintf(out, "%-20s %d\n", "Warnings:", len(r.Warnings))
		for _, w := range r.Warnings {
			fmt.Fprintf(out, "  - %s\n", w)
		}
	}
	fmt.Fprintln(out, "Next steps: run `eka validate` and review the imported artifacts under docs/")
}
