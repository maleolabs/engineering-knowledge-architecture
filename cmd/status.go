package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/workspace"
	"github.com/spf13/cobra"
)

// newStatusCommand builds the `eka status` command: a deterministic
// overview of the EKA workspace — path, schema version, registered
// projects/repositories, canonical store counts and the last sync per
// repository. An uninitialized workspace (no workspace.json) is
// reported informatively and exits 0: status is a read-only probe and
// never creates the workspace.
func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the EKA workspace status",
		Long: `Show the EKA workspace status: the workspace path, the schema
version, the registered projects and repositories, the canonical
store totals (objects, immutable payloads, attachments) and the last
pull/push per repository (from the sync log).

When no workspace exists yet an informational message is printed and
the command exits 0 — status never initializes the workspace.

Exit codes:
  0  status shown (or no workspace yet)
  2  usage or internal error`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := workspace.HomeDir()
			if err != nil {
				return err
			}
			s := styleFor(cmd)
			if _, err := os.Stat(filepath.Join(home, "workspace.json")); err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintf(s.W, "%s\n", s.Accent("Runtime"))
					fmt.Fprintf(s.W, "  %s\n", s.Info("No EKA workspace at "+home+" yet. Run 'eka project register' to create it."))
					return nil
				}
				return fmt.Errorf("status failed: cannot access %s: %w", home, err)
			}

			ws, err := workspace.Ensure()
			if err != nil {
				return err
			}
			defer ws.Close()
			if err := renderStatus(s, ws); err != nil {
				return err
			}
			return nil
		},
	}
}

// renderStatus renders the workspace overview deterministically. Any
// store failure is propagated: status is a monitoring command and must
// never report a false healthy state (exit code 2 on internal error).
func renderStatus(s *ui.Style, ws *workspace.Workspace) error {
	_, id, created := ws.Meta()
	// The store schema version (eka.db) — the meaningful one; the
	// workspace.json file format version is an internal detail.
	sv, err := ws.DB.SchemaVersion()
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}
	projects, err := ws.Projects()
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}
	objects, payloads, attachments, err := ws.Counts()
	if err != nil {
		return fmt.Errorf("status failed: %w", err)
	}

	ui.NewHeader(s, "Runtime").
		Add("Workspace", ws.Path()).
		Add("Schema", "v"+strconv.Itoa(sv)).
		Add("ID", id).
		Add("Created", created).
		Render()

	ui.NewSummary(s).
		Add("Projects", strconv.Itoa(len(projects))).
		Add("Objects", strconv.Itoa(objects)).
		Add("Payloads", strconv.Itoa(payloads)).
		Add("Attachments", strconv.Itoa(attachments)).
		Render()

	if len(projects) == 0 {
		fmt.Fprintf(s.W, "\n%s\n", s.Info("No projects registered. Run 'eka project register' to add one."))
		return nil
	}
	for _, p := range projects {
		repos, err := ws.Repos(p.ID)
		if err != nil {
			return fmt.Errorf("status failed: %w", err)
		}
		fmt.Fprintf(s.W, "\n%s %s\n", ui.IconBullet, s.Accent(p.ID))
		for _, r := range repos {
			last, err := lastSyncDetail(ws, p.ID, r.Name)
			if err != nil {
				return fmt.Errorf("status failed: %w", err)
			}
			fmt.Fprintf(s.W, "  %s %s  (%s)%s\n", ui.IconBullet, s.Info(r.Name), displayPath(r.Path), last)
		}
	}
	return nil
}

// lastSyncDetail renders the most recent sync-log entry of one
// repository ("" when none).
func lastSyncDetail(ws *workspace.Workspace, projectID, repo string) (string, error) {
	entries, err := ws.DB.RecentSyncs(projectID, repo, 1)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", nil
	}
	e := entries[0]
	return fmt.Sprintf("  [%s %s at %s]", e.Direction, shortDigest(e.SnapshotDigest), e.At), nil
}
