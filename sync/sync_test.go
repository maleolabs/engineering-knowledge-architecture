package sync

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/compile"
	"github.com/maleolabs/engineering-knowledge-architecture/exchange"
	"github.com/maleolabs/engineering-knowledge-architecture/workspace"
)

// fixtureValid is the conformant sync test fixture (3 artifacts + 1
// attachment, warnings only).
const fixtureValid = "testdata/valid"

// copyFixture copies the fixture tree into a fresh temp dir and
// returns its path.
func copyFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	err := filepath.WalkDir(fixtureValid, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(fixtureValid, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}

// newWorkspace sets EKA_HOME to a temp dir and ensures the workspace.
func newWorkspace(t *testing.T) *workspace.Workspace {
	t.Helper()
	t.Setenv("EKA_HOME", t.TempDir())
	w, err := workspace.Ensure()
	if err != nil {
		t.Fatalf("workspace.Ensure: %v", err)
	}
	t.Cleanup(func() { w.Close() })
	return w
}

// both returns the default options (pull + push).
func both() Options { return Options{Pull: true, Push: true} }

// TestEKAHomeHonored: the workspace resolves to EKA_HOME.
func TestEKAHomeHonored(t *testing.T) {
	home := t.TempDir()
	t.Setenv("EKA_HOME", home)
	w, err := workspace.Ensure()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	if w.Path() != filepath.Clean(home) {
		t.Errorf("workspace path = %q, want %q", w.Path(), home)
	}
}

// TestSyncFreshRepo: a first sync registers the repo, seeds the
// canonical store from the snapshot and writes the snapshot package.
func TestSyncFreshRepo(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)

	report, err := Run(w, repoDir, both())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !report.NewRepo {
		t.Error("first sync must register the repository")
	}
	if report.Project != filepath.Base(repoDir) {
		t.Errorf("project = %q, want %q", report.Project, filepath.Base(repoDir))
	}
	if report.PullSource != "docs" {
		t.Errorf("pull source = %q, want docs (no snapshot yet)", report.PullSource)
	}
	if report.PulledUnits != 4 {
		t.Errorf("pulled units = %d, want 4", report.PulledUnits)
	}
	if report.PulledAttachments != 1 {
		t.Errorf("pulled attachments = %d, want 1", report.PulledAttachments)
	}
	if report.PushedUnits != 4 {
		t.Errorf("pushed units = %d, want 4", report.PushedUnits)
	}
	if report.SnapshotLabel == "" || report.SnapshotDigest == "" {
		t.Errorf("snapshot label/digest must be set: %q / %q", report.SnapshotLabel, report.SnapshotDigest)
	}
	if report.SnapshotLabel != "rsf-repo-eka-sync-fixture-1.1" {
		t.Errorf("snapshot label = %q", report.SnapshotLabel)
	}

	// Snapshot directory exists with the expected entries.
	snapshotDir := filepath.Join(repoDir, "exchange", "snapshots")
	for _, want := range []string{"header.json", "manifest.json", "declarations.json", "integrity.json"} {
		if _, err := os.Stat(filepath.Join(snapshotDir, want)); err != nil {
			t.Errorf("snapshot missing %s: %v", want, err)
		}
	}
	if _, err := os.Stat(filepath.Join(snapshotDir, "units", "eka-sync-fixture", "adr-001-runtime-v1", "content")); err != nil {
		t.Errorf("snapshot unit entry missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(snapshotDir, "attachments", "docs", "architecture", "diagram.txt")); err != nil {
		t.Errorf("snapshot attachment entry missing: %v", err)
	}

	// DB seeded (counts through the store directly — the workspace
	// Counts helper moved to the runtime Knowledge service).
	objects, err := w.Store().RefCount()
	if err != nil {
		t.Fatal(err)
	}
	payloads, err := w.Store().PayloadCount()
	if err != nil {
		t.Fatal(err)
	}
	attachments, err := w.Store().AttachmentCount()
	if err != nil {
		t.Fatal(err)
	}
	if objects != 4 || attachments != 1 {
		t.Errorf("counts = %d objects / %d attachments, want 4 / 1", objects, attachments)
	}
	if payloads != 4 {
		t.Errorf("payloads = %d, want 4 (one immutable payload per unit)", payloads)
	}
	// The stored reference carries its source repo, and the payload
	// preserves the content and the relationships (serialized inside
	// the immutable unit.json).
	r, ok, err := w.Store().Ref("eka-sync-fixture/adr:001-runtime:1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("seeded reference missing")
	}
	if r.SourceRepo != filepath.Base(repoDir) {
		t.Errorf("source repo = %q, want %q", r.SourceRepo, filepath.Base(repoDir))
	}
	unitJSON, content, err := w.Store().Payload(r.ObjectHash)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "ADR-001") {
		t.Errorf("content not preserved: %q", content)
	}
	u, err := exchange.DecodeUnit(unitJSON, content)
	if err != nil {
		t.Fatalf("stored payload must decode: %v", err)
	}
	if len(u.Relationships) != 1 || u.Relationships[0].Type != "depends-on" ||
		u.Relationships[0].Target != "eka-sync-fixture/sto:login-email:1" {
		t.Errorf("relationships = %+v", u.Relationships)
	}
}

// TestSyncSecondRunIdempotent: the second sync pulls from the snapshot
// (unchanged), re-pushes byte-identical snapshot files, and leaves the
// store untouched.
func TestSyncSecondRunIdempotent(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	if _, err := Run(w, repoDir, both()); err != nil {
		t.Fatal(err)
	}
	snapshotDir := filepath.Join(repoDir, "exchange", "snapshots")
	before := snapshotBytes(t, snapshotDir)

	report, err := Run(w, repoDir, both())
	if err != nil {
		t.Fatalf("second Run: %v", err)
	}
	if !report.Unchanged {
		t.Error("second sync must report unchanged")
	}
	if report.PullSource != "snapshot" {
		t.Errorf("pull source = %q, want snapshot", report.PullSource)
	}
	if report.PulledUnits != 0 {
		t.Errorf("second sync must pull 0 units, got %d", report.PulledUnits)
	}
	if len(report.Warnings) == 0 {
		t.Error("unchanged pull must carry a warning note")
	}

	after := snapshotBytes(t, snapshotDir)
	if !bytesEqualMaps(before, after) {
		t.Error("re-push must produce byte-identical snapshot files")
	}
	objects, err := w.Store().RefCount()
	if err != nil {
		t.Fatal(err)
	}
	attachments, err := w.Store().AttachmentCount()
	if err != nil {
		t.Fatal(err)
	}
	if objects != 4 || attachments != 1 {
		t.Errorf("store changed by second sync: %d objects / %d attachments", objects, attachments)
	}
}

// TestPullCorruptSnapshot: a tampered snapshot is refused with an
// error (integrity class), never silently skipped.
func TestPullCorruptSnapshot(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	if _, err := Run(w, repoDir, both()); err != nil {
		t.Fatal(err)
	}
	content := filepath.Join(repoDir, "exchange", "snapshots", "units", "eka-sync-fixture", "adr-001-runtime-v1", "content")
	data, err := os.ReadFile(content)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(content, append([]byte("X"), data[1:]...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(w, repoDir, both()); err == nil {
		t.Fatal("corrupt snapshot must error")
	}
}

// TestMultiRepoOneProject: two repositories registered under one
// project sync into the union in the canonical store; each push
// carries exactly its own objects.
func TestMultiRepoOneProject(t *testing.T) {
	w := newWorkspace(t)
	repoA := copyFixture(t)
	// A second repository with a disjoint namespace (the fixtures must
	// not collide in the canonical store, which is keyed by identity).
	repoB := copyFixture(t)
	rewriteFiles(t, repoB, "eka-sync-fixture", "eka-sync-fixture-b")

	// Register both under one project.
	_, _, _, err := w.RegisterRepo(repoA, "myproject")
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, err = w.RegisterRepo(repoB, "myproject")
	if err != nil {
		t.Fatal(err)
	}

	// Sync both.
	for _, repo := range []string{repoA, repoB} {
		report, err := Run(w, repo, both())
		if err != nil {
			t.Fatalf("sync %s: %v", repo, err)
		}
		if report.Project != "myproject" {
			t.Errorf("project = %q, want myproject", report.Project)
		}
		if report.PulledUnits != 4 {
			t.Errorf("pulled units = %d, want 4", report.PulledUnits)
		}
	}

	// Union in the DB: both repos' refs present, one payload each.
	objects, err := w.Store().RefCount()
	if err != nil {
		t.Fatal(err)
	}
	if objects != 8 {
		t.Errorf("union objects = %d, want 8", objects)
	}
	refsA, err := w.Store().Refs("myproject", filepath.Base(repoA))
	if err != nil {
		t.Fatal(err)
	}
	refsB, err := w.Store().Refs("myproject", filepath.Base(repoB))
	if err != nil {
		t.Fatal(err)
	}
	if len(refsA) != 4 || len(refsB) != 4 {
		t.Errorf("per-repo slices = %d/%d, want 4/4", len(refsA), len(refsB))
	}
	for _, r := range refsA {
		if r.SourceRepo != filepath.Base(repoA) {
			t.Errorf("reference %s attributed to %q, want %q", r.Form, r.SourceRepo, filepath.Base(repoA))
		}
	}
	// Each snapshot carries only its own namespace.
	for _, repo := range []string{repoA, repoB} {
		pkg, err := exchange.LoadPackage(filepath.Join(repo, "exchange", "snapshots"))
		if err != nil {
			t.Fatalf("LoadPackage(%s): %v", repo, err)
		}
		if len(pkg.Units) != 4 {
			t.Errorf("snapshot of %s carries %d units, want 4", repo, len(pkg.Units))
		}
	}
}

// rewriteFiles rewrites one string to another in every file under dir.
func rewriteFiles(t *testing.T, dir, from, to string) {
	t.Helper()
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.Contains(string(data), from) {
			return nil
		}
		return os.WriteFile(path, []byte(strings.ReplaceAll(string(data), from, to)), 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestMigrationPullSeedsFromDocs: a repository without a snapshot is
// seeded from its docs tree (source "docs").
func TestMigrationPullSeedsFromDocs(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	report, err := Run(w, repoDir, Options{Pull: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if report.PullSource != "docs" {
		t.Errorf("pull source = %q, want docs", report.PullSource)
	}
	if report.PulledUnits != 4 {
		t.Errorf("pulled units = %d, want 4", report.PulledUnits)
	}
	if report.PushedUnits != 0 {
		t.Errorf("pull-only run must push 0 units, got %d", report.PushedUnits)
	}
	// No snapshot written by a pull-only run.
	if _, err := os.Stat(filepath.Join(repoDir, "exchange", "snapshots")); !os.IsNotExist(err) {
		t.Error("pull-only run must not write a snapshot")
	}
}

// TestFromDocsReseeds: --from-docs re-seeds the canonical store from
// the docs tree even when a snapshot exists, reported as source
// "docs".
func TestFromDocsReseeds(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	if _, err := Run(w, repoDir, both()); err != nil {
		t.Fatal(err)
	}
	report, err := Run(w, repoDir, Options{Pull: true, FromDocs: true})
	if err != nil {
		t.Fatalf("Run with FromDocs: %v", err)
	}
	if report.PullSource != "docs" {
		t.Errorf("pull source = %q, want docs", report.PullSource)
	}
	if report.PulledUnits != 4 {
		t.Errorf("pulled units = %d, want 4", report.PulledUnits)
	}
	// The docs digest differs from the snapshot digest (fresh pull).
	if report.Unchanged {
		t.Error("from-docs pull must not report unchanged")
	}
}

// TestDocsModeValidationGate: a non-conformant repository is refused
// by the docs-mode gate with a compile.ValidationError (the Knowledge
// Compiler's typed validation failure).
func TestDocsModeValidationGate(t *testing.T) {
	w := newWorkspace(t)
	repoDir := t.TempDir()
	dir := filepath.Join(repoDir, "docs", "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	bad := "---\nnamespace: x\nid: 1\n---\n# bad\n" // type missing: R0 error
	if err := os.WriteFile(filepath.Join(dir, "bad.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Run(w, repoDir, both())
	if err == nil {
		t.Fatal("non-conformant repository must be refused")
	}
	var ve *compile.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("error type = %T, want *compile.ValidationError", err)
	}
}

// TestPushOnlyNoObjects: pushing a registered repository with no
// stored objects is a no-op (no files written).
func TestPushOnlyNoObjects(t *testing.T) {
	w := newWorkspace(t)
	repoDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoDir, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := w.RegisterRepo(repoDir, ""); err != nil {
		t.Fatal(err)
	}
	report, err := Run(w, repoDir, Options{Push: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if report.PushedUnits != 0 || report.SnapshotLabel != "" || report.SnapshotDigest != "" {
		t.Errorf("no-op push result = %+v", report)
	}
	if _, err := os.Stat(filepath.Join(repoDir, "exchange", "snapshots")); !os.IsNotExist(err) {
		t.Error("no-op push must not write a snapshot")
	}
}

// snapshotBytes walks the snapshot dir and returns rel path -> bytes.
func snapshotBytes(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// bytesEqualMaps compares two name->bytes maps for byte equality.
func bytesEqualMaps(a, b map[string][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for name, data := range a {
		other, ok := b[name]
		if !ok || !bytes.Equal(data, other) {
			return false
		}
	}
	return true
}

// TestBasenameCollisionAcrossProjects (M1 regression): two
// repositories with the SAME basename registered under DIFFERENT
// projects must never leak objects into each other's snapshots —
// provenance is the (project_id, source_repo) pair, not the basename.
func TestBasenameCollisionAcrossProjects(t *testing.T) {
	w := newWorkspace(t)
	parent := t.TempDir()
	repoA := filepath.Join(parent, "samename", "a")
	repoB := filepath.Join(parent, "samename", "b")
	copyFixtureTo(t, repoA)
	copyFixtureTo(t, repoB)
	// Disjoint namespaces: without provenance isolation the second
	// push would package the first repository's units too.
	rewriteFiles(t, repoB, "eka-sync-fixture", "eka-sync-fixture-b")

	if _, _, _, err := w.RegisterRepo(repoA, "proj-a"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := w.RegisterRepo(repoB, "proj-b"); err != nil {
		t.Fatal(err)
	}
	for _, repo := range []string{repoA, repoB} {
		if _, err := Run(w, repo, both()); err != nil {
			t.Fatalf("sync %s: %v", repo, err)
		}
	}

	for _, tc := range []struct{ repo, wantNS string }{
		{repoA, "eka-sync-fixture"},
		{repoB, "eka-sync-fixture-b"},
	} {
		pkg, err := exchange.LoadPackage(filepath.Join(tc.repo, "exchange", "snapshots"))
		if err != nil {
			t.Fatalf("LoadPackage(%s): %v", tc.repo, err)
		}
		if len(pkg.Units) != 4 {
			t.Errorf("snapshot of %s carries %d units, want 4 (its own only)", tc.repo, len(pkg.Units))
		}
		for _, u := range pkg.Units {
			if u.Identity.Namespace != tc.wantNS {
				t.Errorf("snapshot of %s carries foreign namespace %q", tc.repo, u.Identity.Namespace)
			}
		}
	}
}

// TestAttachmentIsolationAcrossRepos (M5 regression): attachments are
// attributed to their repository; a push never packages another
// repository's attachments into a snapshot.
func TestAttachmentIsolationAcrossRepos(t *testing.T) {
	w := newWorkspace(t)
	parent := t.TempDir()
	repoA := filepath.Join(parent, "repo-a")
	repoB := filepath.Join(parent, "repo-b")
	copyFixtureTo(t, repoA)
	copyFixtureTo(t, repoB)
	rewriteFiles(t, repoB, "eka-sync-fixture", "eka-sync-fixture-b")
	// Remove repo B's own attachment: with provenance isolation its
	// snapshot must carry zero attachments, never repo A's.
	if err := os.Remove(filepath.Join(repoB, "docs", "architecture", "diagram.txt")); err != nil {
		t.Fatal(err)
	}

	if _, _, _, err := w.RegisterRepo(repoA, "myproject"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := w.RegisterRepo(repoB, "myproject"); err != nil {
		t.Fatal(err)
	}
	for _, repo := range []string{repoA, repoB} {
		if _, err := Run(w, repo, both()); err != nil {
			t.Fatalf("sync %s: %v", repo, err)
		}
	}

	attPath := filepath.Join("attachments", "docs", "architecture", "diagram.txt")
	pkgA, err := exchange.LoadPackage(filepath.Join(repoA, "exchange", "snapshots"))
	if err != nil {
		t.Fatal(err)
	}
	pkgB, err := exchange.LoadPackage(filepath.Join(repoB, "exchange", "snapshots"))
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgA.Attachments) != 1 {
		t.Errorf("repo A snapshot carries %d attachments, want 1", len(pkgA.Attachments))
	}
	if len(pkgB.Attachments) != 0 {
		t.Errorf("repo B snapshot carries %d attachments, want 0 (isolation)", len(pkgB.Attachments))
	}
	_ = attPath
}

// TestCrossRepoOverwriteRecorded (D4 contract): pulling an identity
// already owned by a different repository with different content
// applies deterministic last-wins AND records the overwrite in the
// report.
func TestCrossRepoOverwriteRecorded(t *testing.T) {
	w := newWorkspace(t)
	parent := t.TempDir()
	repoA := filepath.Join(parent, "repo-a")
	repoB := filepath.Join(parent, "repo-b")
	copyFixtureTo(t, repoA)
	copyFixtureTo(t, repoB)
	// Same namespaces, same identities, different content: change the
	// body of the ADR in repo B so the digests differ.
	adrB := filepath.Join(repoB, "docs", "decisions", "adr-001-runtime.md")
	data, err := os.ReadFile(adrB)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(adrB, append(data, []byte("\n# divergent revision\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, _, _, err := w.RegisterRepo(repoA, "myproject"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := w.RegisterRepo(repoB, "myproject"); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(w, repoA, both()); err != nil {
		t.Fatal(err)
	}
	reportB, err := Run(w, repoB, both())
	if err != nil {
		t.Fatal(err)
	}
	if len(reportB.Overwrites) == 0 {
		t.Error("repo B pull must record cross-repository overwrites")
	}
	if !strings.Contains(reportB.Overwrites[0], "001-runtime") {
		t.Errorf("overwrite record = %q, want the colliding identity named", reportB.Overwrites[0])
	}
	if !strings.Contains(reportB.Overwrites[0], "repo-a") {
		t.Errorf("overwrite record = %q, want the previous owner named", reportB.Overwrites[0])
	}

	// Last-wins is symmetric: a forced re-seed of repo A (docs mode
	// never skips) overwrites repo B's divergent copy again — the
	// store keeps the last pull.
	reportA, err := Run(w, repoA, Options{Pull: true, FromDocs: true})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, o := range reportA.Overwrites {
		if strings.Contains(o, "repo-b") {
			found = true
		}
	}
	if !found {
		t.Errorf("repo A re-pull must overwrite repo B's divergent copy, got %v", reportA.Overwrites)
	}
}

// TestPushFailsOnStoreReadError (M2 regression): a store read failure
// during push must abort the push — a snapshot must never be emitted
// with silently dropped units.
func TestPushFailsOnStoreReadError(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	if _, err := Run(w, repoDir, both()); err != nil {
		t.Fatal(err)
	}
	// Break the store behind the push's back.
	if _, err := w.Store().DB().Exec(`DROP TABLE object_refs`); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(w, repoDir, Options{Push: true}); err == nil {
		t.Fatal("push must fail when reference reads fail")
	}
}

// TestRegisterPathOwnedByFirstProject (m6 regression): a repository
// path is owned by the project that registered it first; registering
// the same path under a different project is refused.
func TestRegisterPathOwnedByFirstProject(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	if _, _, _, err := w.RegisterRepo(repoDir, "proj-one"); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := w.RegisterRepo(repoDir, "proj-two"); err == nil {
		t.Fatal("re-registering an owned path under another project must be refused")
	}
	// Re-registering under the same project stays a no-op.
	_, _, created, err := w.RegisterRepo(repoDir, "proj-one")
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Error("re-registering the same path under the same project must be a no-op")
	}
}

// copyFixtureTo copies the fixture tree into the target directory.
func copyFixtureTo(t *testing.T, dst string) {
	t.Helper()
	err := filepath.WalkDir(fixtureValid, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(fixtureValid, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestPushReportsSnapshotChanged (M1 regression): when the canonical
// store is tampered behind the runtime's back, the next sync rewrites
// the snapshot with a DIFFERENT digest — the report must say so
// instead of claiming "unchanged", so a user never commits a
// corruption-derived snapshot believing nothing changed.
func TestPushReportsSnapshotChanged(t *testing.T) {
	w := newWorkspace(t)
	repoDir := copyFixture(t)
	first, err := Run(w, repoDir, both())
	if err != nil {
		t.Fatal(err)
	}
	if first.PushChanged {
		t.Error("first sync must not report a snapshot change")
	}

	// Tamper one payload behind the store's back.
	rows, err := w.Store().AllPayloads()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("no payloads seeded")
	}
	if _, err := w.Store().DB().Exec(`UPDATE object_payloads SET content = ? WHERE object_hash = ?`,
		[]byte("tampered"), rows[0].ObjectHash); err != nil {
		t.Fatal(err)
	}

	second, err := Run(w, repoDir, both())
	if err != nil {
		t.Fatal(err)
	}
	if !second.PushChanged {
		t.Error("push after tampering must report the snapshot digest change")
	}
	found := false
	for _, warn := range second.Warnings {
		if strings.Contains(warn, "snapshot rewritten") {
			found = true
		}
	}
	if !found {
		t.Errorf("warnings must carry the snapshot-rewritten note, got %v", second.Warnings)
	}
}
