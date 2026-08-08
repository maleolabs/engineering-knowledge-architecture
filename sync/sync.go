// Package sync implements the EKA Knowledge Runtime synchronization
// engine (milestone EKA v0.2.0): the bidirectional transport between a
// registered repository and the EKA workspace canonical store.
//
// Model (documented):
//
//   - The workspace is canonical: immutable Engineering Knowledge
//     Objects (content-addressed payloads), their mutable references,
//     attachments and the sync log live in the workspace database
//     (eka.db). Relationships and change logs are serialized inside
//     the immutable payload (unit.json), never stored separately.
//   - The transport between a repository and the workspace is the
//     snapshot directory <repo>/exchange/snapshots: an RSF package in
//     directory layout, verified byte-exact on every read
//     (exchange.LoadPackageWithEntries) and emitted deterministically
//     (exchange.Emit).
//   - Pull (snapshot mode): the snapshot package is verified and its
//     units are seeded as immutable payloads (store.PutUnit: unit.json
//     entry bytes verbatim + content, content-addressed), attributed to
//     the repository via the reference (project_id + source_repo).
//     Idempotent: an unchanged package digest skips the work; re-seeding
//     the same package is a no-op. Deletions are never applied.
//   - Pull (docs mode): the repository's docs tree is compiled through
//     the Knowledge Compiler (compile.Compile: the authoring
//     conformance gate, then the package assembled exactly as
//     `eka export` would via exchange.RepositoryPackage), then seeded
//     the same way with unit.json bytes serialized via
//     exchange.MarshalUnit. This is the migration path for repositories
//     without a snapshot, and the --from-docs re-seed path.
//   - Push: the repository's references in the canonical store are
//     resolved to their immutable payloads, the units are reconstructed
//     (exchange.DecodeUnit) and assembled into an RSF package
//     (namespace resolution: existing-snapshot header, else most common
//     namespace, else error) and emitted into
//     <repo>/exchange/snapshots atomically (write to .snapshots-tmp,
//     then swap).
//
// The sync log (store.RecordSync) records every pull/push run and
// backs the idempotent-pull check and `eka status`.
//
// Error classes: validation failures (docs gate) and integrity
// failures (corrupt snapshot) are returned as typed errors mapped to
// exit code 1 by the CLI; workspace/registry/usage failures map to
// exit code 2.
package sync

import (
	"fmt"
	"path/filepath"

	"github.com/maleolabs/engineering-knowledge-architecture/workspace"
)

// Options configures one sync run. Pull and Push both default to true
// (`eka sync`); a pull-only run sets Push=false, a push-only run sets
// Pull=false.
type Options struct {
	Pull bool
	Push bool
	// FromDocs seeds the canonical store from the repository's docs
	// tree (migration/re-seed) instead of the snapshot directory.
	FromDocs bool
}

// Report is the deterministic outcome of one sync run.
type Report struct {
	// Workspace is the workspace root.
	Workspace string
	// Project and Repo identify the synced repository.
	Project string
	Repo    string
	// PullSource is "snapshot" or "docs" ("" when no pull ran).
	PullSource string
	// PulledUnits/PulledAttachments count the pull work.
	PulledUnits       int
	PulledAttachments int
	// PushedUnits/PushedAttachments count the push work (0 in a
	// no-op push with no stored objects).
	PushedUnits       int
	PushedAttachments int
	// SnapshotLabel is the pushed package label ("" when nothing was
	// pushed); SnapshotDigest is the package digest ("" when neither
	// pull nor push produced a digest).
	SnapshotLabel  string
	SnapshotDigest string
	// NewRepo reports whether the repository was registered by this
	// run.
	NewRepo bool
	// Unchanged reports an idempotent pull (snapshot digest already
	// current — no pull work done).
	Unchanged bool
	// PushChanged reports that the push rewrote the snapshot with a
	// DIFFERENT digest than the one it replaced: the canonical store
	// and the repository snapshot were out of sync (a tampered store
	// is the typical cause). A clean re-push of identical state
	// reports false.
	PushChanged bool
	// Overwrites lists the identities this run replaced from another
	// repository (deterministic order, empty when none).
	Overwrites []string
	// Warnings are informational notes, in deterministic order.
	Warnings []string
}

// Run executes one sync of the repository at repoPath: resolve and
// (auto-)register the repository, then pull and/or push per opts.
// Errors are wrapped with context; validation and integrity failures
// carry their typed classes (see pull.go).
func Run(ws *workspace.Workspace, repoPath string, opts Options) (*Report, error) {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("sync: cannot resolve repository path %q: %w", repoPath, err)
	}
	abs = filepath.Clean(abs)

	// Resolve the registry: an already-registered repository keeps its
	// project; otherwise it is registered (project name = basename).
	repo, found, err := ws.FindRepo(abs)
	if err != nil {
		return nil, err
	}
	report := &Report{Workspace: ws.Path(), Repo: filepath.Base(abs)}
	if !found {
		var project workspace.Project
		project, repo, _, err = ws.RegisterRepo(abs, "")
		if err != nil {
			return nil, err
		}
		report.Project = project.ID
		report.NewRepo = true
	} else {
		report.Project = repo.ProjectID
	}

	if opts.Pull {
		result, err := Pull(ws, repo, opts.FromDocs)
		if err != nil {
			return nil, err
		}
		report.PullSource = result.Source
		report.PulledUnits = result.Units
		report.PulledAttachments = result.Attachments
		report.Unchanged = result.Unchanged
		report.Overwrites = result.Overwrites
		if result.Digest != "" {
			report.SnapshotDigest = result.Digest
		}
	}

	if opts.Push {
		result, err := Push(ws, repo)
		if err != nil {
			return nil, err
		}
		report.PushedUnits = result.Units
		report.PushedAttachments = result.Attachments
		report.PushChanged = result.Changed
		if result.Label != "" {
			report.SnapshotLabel = result.Label
		}
		if result.Digest != "" {
			report.SnapshotDigest = result.Digest
		}
	}

	if report.Unchanged && !report.PushChanged {
		report.Warnings = append(report.Warnings, "no changes: snapshot already up to date")
	}
	if report.PushChanged {
		report.Warnings = append(report.Warnings,
			"snapshot rewritten: digest "+report.SnapshotDigest+" differs from the replaced snapshot; the canonical store and the repository were out of sync (run 'eka integrity check' to verify the store)")
	}
	for _, o := range report.Overwrites {
		report.Warnings = append(report.Warnings, o)
	}
	return report, nil
}
