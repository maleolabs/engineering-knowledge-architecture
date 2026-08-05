package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/bootstrap"
	"github.com/spf13/cobra"
)

// newInitCommand builds the `eka init` command: bootstrap a new EKA
// repository. All bootstrap logic lives in the bootstrap engine; this
// command validates arguments, renders the outcome and maps the result
// to the exit code contract.
//
// Exit codes:
//
//	0  init completed and the generated repository validates; dry-run
//	1  init completed but validation found blocking violations
//	2  usage or internal error (unknown flag, invalid project name,
//	   generation failure)
func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Bootstrap a new EKA repository",
		Long: `Bootstrap a new EKA repository at the current directory or at the
directory name (relative to the current directory). If name already
exists as a directory it is adopted; existing files with identical
content are reused, conflicting files are never overwritten silently.

Prompts (project name, namespace, description, README, git init) are
asked only when stdin is a terminal; otherwise deterministic defaults
are used and git is never initialized.

Exit codes:
  0  init completed and the generated repository validates
  1  init completed but validation found blocking violations
  2  usage or internal error (unknown flag, invalid project name,
     generation failure)`,
		Example: `  eka init              bootstrap the current directory
  eka init myproject    create and bootstrap ./myproject
  eka init --dry-run    preview the plan without writing anything`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return fmt.Errorf("init failed: %w", err)
			}
			target := "."
			if len(args) == 1 {
				name := args[0]
				if strings.ContainsAny(name, `/\`) {
					return fmt.Errorf("project name %q must not contain path separators", name)
				}
				target = name
			}
			// Refuse targets that exist as plain files; directories are
			// adopted.
			if info, err := os.Stat(target); err == nil && !info.IsDir() {
				return fmt.Errorf("target %q exists and is not a directory", target)
			}

			outcome, err := bootstrap.Run(bootstrap.Options{
				Target: target,
				DryRun: dryRun,
				Stdin:  cmd.InOrStdin(),
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			})
			if err != nil {
				return fmt.Errorf("init failed: %w", err)
			}

			printInitOutcome(cmd.OutOrStdout(), outcome)
			if outcome.DryRun {
				return nil
			}
			if outcome.Report != nil && !outcome.Report.Pass() {
				fmt.Fprintf(cmd.ErrOrStderr(),
					"eka: init generated a non-conformant repository (see validation results above)\n")
				return &exitError{code: exitFail}
			}
			return nil
		},
	}
	cmd.Flags().Bool("dry-run", false,
		"preview the plan (directories, files, git, validation); writes nothing")
	return cmd
}

// printInitOutcome renders the init result deterministically: summary
// header, plan (dry-run) or generation counts, validation report, next
// steps. The exact bytes are part of the CLI contract (preserved from
// the pre-Cobra implementation).
func printInitOutcome(out io.Writer, o *bootstrap.Outcome) {
	if o.DryRun {
		fmt.Fprintf(out, "EKA Bootstrap Plan (dry-run)\n")
		fmt.Fprintf(out, "============================\n")
		for _, a := range o.Plan {
			fmt.Fprintf(out, "  %s\n", a.String())
		}
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Dry-run: no changes were written.")
		return
	}

	fmt.Fprintf(out, "EKA Repository Init\n")
	fmt.Fprintf(out, "===================\n")
	fmt.Fprintf(out, "%-23s %s\n", "Target:", o.Target)
	fmt.Fprintf(out, "%-23s %s\n", "Project Name:", o.ProjectName)
	fmt.Fprintf(out, "%-23s %s\n", "Namespace:", o.Namespace)
	fmt.Fprintf(out, "%-23s %s\n", "Repository Type:", o.RepoType)
	fmt.Fprintf(out, "%-23s %s\n", "Git Status:", o.GitStatus)
	fmt.Fprintf(out, "%-23s %s\n", "Knowledge Standard:", "EKA v1.0")
	if !o.AlreadyInitialized {
		fmt.Fprintf(out, "%-23s %d\n", "Dirs created:", len(o.CreatedDirs))
		fmt.Fprintf(out, "%-23s %d\n", "Files created:", len(o.CreatedFiles))
		fmt.Fprintf(out, "%-23s %d\n", "Files reused:", len(o.ReusedFiles))
		fmt.Fprintf(out, "%-23s %d\n", "Files overwritten:", len(o.OverwrittenFiles))
		fmt.Fprintf(out, "%-23s %d\n", "Files skipped:", len(o.SkippedFiles))
	}

	fmt.Fprintln(out)
	if o.Report != nil {
		printReport(out, o.Report)
		fmt.Fprintln(out)
		verdict := "PASS"
		if !o.Report.Pass() {
			verdict = "FAIL"
		}
		fmt.Fprintf(out, "Validation Result: %s (%d errors, %d warnings)\n",
			verdict, o.Report.ErrorCount(), o.Report.WarningCount())
	}

	fmt.Fprintln(out)
	if o.AlreadyInitialized {
		fmt.Fprintln(out, "Already initialized: no changes were made; repository validated.")
		return
	}
	fmt.Fprintln(out, "Next steps:")
	fmt.Fprintln(out, "  - create your first artifact under docs/")
	fmt.Fprintf(out, "  - set namespace on future artifacts: %s\n", o.Namespace)
	fmt.Fprintln(out, "  - read docs/README.md for the serialization conventions")
}
