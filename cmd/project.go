package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/workspace"
	"github.com/spf13/cobra"
)

// newProjectCommand builds the `eka project` command tree: the project
// and repository registry of the EKA workspace.
//
// Exit codes:
//
//	0  registration succeeded (new or already registered); list printed
//	2  usage or internal error (workspace resolution, registry failure)
func newProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage EKA workspace projects",
		Long: `Manage the projects and repositories registered in the EKA
workspace (default ~/.eka, or $EKA_HOME).

A project groups one or more repositories; the repository name is its
directory basename, the project name is the --name flag value or the
same basename. The canonical store attributes every pulled object to
its repository.

Exit codes:
  0  success
  2  usage or internal error`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newProjectRegisterCommand(), newProjectListCommand())
	return cmd
}

// newProjectRegisterCommand builds `eka project register [path]
// [--name NAME]`.
func newProjectRegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [path]",
		Short: "Register a repository in the workspace",
		Long: `Register the EKA repository at path (default: the current
directory) in the EKA workspace. The repository name is the directory
basename of the absolute path; the project name is the --name flag
value or the same basename. Registering the same repository again is a
no-op (the stored path is refreshed).

Exit codes:
  0  registration succeeded
  2  usage or internal error (workspace resolution, registry failure)`,
		Example: `  eka project register
  eka project register /path/to/repo
  eka project register /path/to/repo --name myproject`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return fmt.Errorf("project register failed: %w", err)
			}
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			if info, err := os.Stat(path); err != nil {
				return fmt.Errorf("project register failed: cannot access %s: %w", path, err)
			} else if !info.IsDir() {
				return fmt.Errorf("project register failed: %s is not a directory", path)
			}
			ws, err := workspace.Ensure()
			if err != nil {
				return err
			}
			defer ws.Close()

			project, repo, created, err := ws.RegisterRepo(path, name)
			if err != nil {
				return err
			}
			status := "already registered"
			if created {
				status = "registered"
			}
			s := styleFor(cmd)
			ui.NewHeader(s, "Project").
				Add("Project", project.ID).
				Add("Repository", repo.Name).
				Add("Path", repo.Path).
				Render()
			ui.NewSummary(s).
				Add("Project", project.ID).
				Add("Repository", repo.Name).
				Add("Path", repo.Path).
				Add("Status", status).
				Render()
			return nil
		},
	}
	cmd.Flags().String("name", "", "project name (default: the repository basename)")
	return cmd
}

// newProjectListCommand builds `eka project list`.
func newProjectListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered projects and repositories",
		Long: `List every project registered in the EKA workspace with its
repositories (name and path), sorted deterministically: projects by
id, repositories by name. A workspace with no registered projects
prints an informational message and exits 0.

Exit codes:
  0  success
  2  usage or internal error`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, err := workspace.Ensure()
			if err != nil {
				return err
			}
			defer ws.Close()

			projects, err := ws.Projects()
			if err != nil {
				return err
			}
			s := styleFor(cmd)
			ui.NewHeader(s, "Projects").
				Add("Workspace", ws.Path()).
				Render()
			if len(projects) == 0 {
				fmt.Fprintf(s.W, "\n%s\n", s.Info("No projects registered yet. Run 'eka project register' to add one."))
				return nil
			}
			for _, p := range projects {
				repos, err := ws.Repos(p.ID)
				if err != nil {
					return err
				}
				fmt.Fprintf(s.W, "\n%s %s\n", ui.IconBullet, s.Accent(p.ID))
				for _, r := range repos {
					fmt.Fprintf(s.W, "  %s %s  (%s)\n", ui.IconBullet, s.Info(r.Name), displayPath(r.Path))
				}
			}
			return nil
		},
	}
	return cmd
}

// displayPath renders a repository path relative to the current
// directory when it is a descendant, else absolute — shorter output,
// still deterministic.
func displayPath(path string) string {
	wd, err := filepath.Abs(".")
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(wd, path)
	if err != nil || rel == ".." || len(rel) > 2 && rel[:3] == ".."+string(filepath.Separator) {
		return path
	}
	return rel
}
