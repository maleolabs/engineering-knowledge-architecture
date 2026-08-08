package workspace

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

// This file implements the project + repository registry over the
// workspace database (store). A Project groups one or more
// repositories; a Repo is one registered repository path with its
// registry name.
//
// Registry semantics:
//   - project id = name (both are the flag value or the basename of the
//     absolute repository path);
//   - the repository path is normalized to an absolute clean path;
//   - a repository path is owned by the project that registered it
//     first: re-registering the path under the same project is a no-op
//     (idempotent), re-registering it under a different project is
//     refused (deterministic ownership — the path column is unique, so
//     provenance and sync resolution can never be ambiguous);
//   - the provenance pair of a repository is (project_id, name) — the
//     same composite key used by the canonical store, so objects and
//     attachments are attributed workspace-uniquely.
//
// All SQL is parameterized.

// Project is one registered project.
type Project struct {
	ID      string
	Name    string
	Created string
}

// Repo is one registered repository.
type Repo struct {
	ProjectID string
	Name      string
	Path      string
	Created   string
}

// RegisterRepo registers the repository at path under a project. The
// repository name is always the basename of the cleaned absolute path.
// The project name is the name flag value when non-empty, else the same
// basename; the project is created when missing (project id = name).
// The repo row (project_id, name) is upserted with the normalized
// absolute path; re-registering an already owned path under the same
// project updates nothing (idempotent). created reports whether the
// repo row was newly created. Two repositories with different basenames
// can therefore share one project (the same --name), which is how a
// project groups multiple repositories.
func (w *Workspace) RegisterRepo(path, name string) (Project, Repo, bool, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot resolve repository path %q: %w", path, err)
	}
	abs = filepath.Clean(abs)
	repoName := filepath.Base(abs)
	if repoName == "" || repoName == "." || repoName == string(filepath.Separator) {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot derive a repository name from path %q", abs)
	}
	projectName := name
	if projectName == "" {
		projectName = repoName
	}

	// Path ownership: a path registered under a different project is
	// refused — the first project owns it (registry determinism).
	if existing, found, err := w.FindRepo(abs); err != nil {
		return Project{}, Repo{}, false, err
	} else if found && existing.ProjectID != projectName {
		return Project{}, Repo{}, false, fmt.Errorf(
			"workspace: repository path %s is already registered under project %q; register it under that project or choose another path",
			abs, existing.ProjectID)
	}

	now := time.Now().Format("2006-01-02")

	// Upsert the project; read its record back.
	if _, err := w.Store().DB().Exec(`INSERT INTO projects (id, name, created) VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET name = excluded.name`, projectName, projectName, now); err != nil {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot register project %q: %w", projectName, err)
	}
	project := Project{ID: projectName, Name: projectName, Created: now}
	if err := w.Store().DB().QueryRow(`SELECT created FROM projects WHERE id = ?`, projectName).Scan(&project.Created); err != nil {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot read project %q: %w", projectName, err)
	}

	// Upsert the repo; created = whether a fresh row was inserted
	// (checked before the upsert, since the row exists afterwards).
	exists, err := repoExists(w, projectName, repoName)
	if err != nil {
		return Project{}, Repo{}, false, err
	}
	created := !exists
	if _, err := w.Store().DB().Exec(`INSERT INTO repos (project_id, name, path, created) VALUES (?, ?, ?, ?)
		ON CONFLICT(project_id, name) DO UPDATE SET path = excluded.path`, projectName, repoName, abs, now); err != nil {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot register repository %q: %w", abs, err)
	}
	repo := Repo{ProjectID: projectName, Name: repoName, Path: abs, Created: now}
	if err := w.Store().DB().QueryRow(`SELECT created FROM repos WHERE project_id = ? AND name = ?`, projectName, repoName).Scan(&repo.Created); err != nil {
		return Project{}, Repo{}, false, fmt.Errorf("workspace: cannot read repository record: %w", err)
	}
	return project, repo, created, nil
}

// repoExists reports whether a repo row already exists (checked before
// the upsert, so a pre-existing row means the upsert was an update).
func repoExists(w *Workspace, projectID, name string) (bool, error) {
	var n int
	err := w.Store().DB().QueryRow(`SELECT COUNT(*) FROM repos WHERE project_id = ? AND name = ?`, projectID, name).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("workspace: cannot check repository existence: %w", err)
	}
	return n > 0, nil
}

// FindRepo returns the repository registered at absPath (by normalized
// absolute path), if any. The path column is unique, so at most one
// repository row can match.
func (w *Workspace) FindRepo(absPath string) (Repo, bool, error) {
	abs, err := filepath.Abs(absPath)
	if err != nil {
		return Repo{}, false, fmt.Errorf("workspace: cannot resolve path %q: %w", absPath, err)
	}
	abs = filepath.Clean(abs)
	var r Repo
	err = w.Store().DB().QueryRow(`SELECT project_id, name, path, created FROM repos WHERE path = ?`, abs).
		Scan(&r.ProjectID, &r.Name, &r.Path, &r.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Repo{}, false, nil
		}
		return Repo{}, false, fmt.Errorf("workspace: cannot find repository %s: %w", abs, err)
	}
	return r, true, nil
}

// Repos returns every repository of one project, sorted by name.
func (w *Workspace) Repos(projectID string) ([]Repo, error) {
	rows, err := w.Store().DB().Query(`SELECT project_id, name, path, created FROM repos
		WHERE project_id = ? ORDER BY name`, projectID)
	if err != nil {
		return nil, fmt.Errorf("workspace: cannot list repositories: %w", err)
	}
	defer rows.Close()
	var out []Repo
	for rows.Next() {
		var r Repo
		if err := rows.Scan(&r.ProjectID, &r.Name, &r.Path, &r.Created); err != nil {
			return nil, fmt.Errorf("workspace: cannot scan repository row: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("workspace: cannot list repositories: %w", err)
	}
	return out, nil
}

// Projects returns every registered project, sorted by id.
func (w *Workspace) Projects() ([]Project, error) {
	rows, err := w.Store().DB().Query(`SELECT id, name, created FROM projects ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("workspace: cannot list projects: %w", err)
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Created); err != nil {
			return nil, fmt.Errorf("workspace: cannot scan project row: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("workspace: cannot list projects: %w", err)
	}
	return out, nil
}

// SortedProjectIDs returns the project ids sorted (deterministic
// iteration for consumers that aggregate across projects).
func SortedProjectIDs(projects []Project) []string {
	ids := make([]string, 0, len(projects))
	for _, p := range projects {
		ids = append(ids, p.ID)
	}
	sort.Strings(ids)
	return ids
}
