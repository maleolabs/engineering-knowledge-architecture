package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// copySyncFixture copies the sync test fixture tree into a fresh temp
// dir and returns its path.
func copySyncFixture(t *testing.T) string {
	t.Helper()
	src := filepath.Join("..", "sync", "testdata", "valid")
	dst := t.TempDir()
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
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

func TestSyncHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"sync", "-h"}, {"sync", "pull", "-h"}, {"sync", "push", "-h"}, {"sync", "--help"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		if !strings.Contains(text, "eka sync") {
			t.Errorf("args %v: help must mention eka sync", args)
		}
	}
}

func TestSyncHelpDocumentsModel(t *testing.T) {
	_, text, _ := runIn([]string{"sync", "-h"})
	for _, want := range []string{"workspace", "$EKA_HOME", "exchange/snapshots", "idempotent", "Deletions"} {
		if !strings.Contains(text, want) {
			t.Errorf("sync help must document %q", want)
		}
	}
}

// TestSyncHappyPath: syncing a fresh fixture repository exits 0 and
// renders the deterministic report; the snapshot directory is created.
func TestSyncHappyPath(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	code, text, errText := runIn([]string{"sync", repo})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	for _, want := range []string{"Runtime", "Workspace", "Project", "Repository", "Pull", "Push", "Snapshot", "rsf-repo-eka-sync-fixture-1.1", "registered (new)", "docs: 4 units, 1 attachment", "4 units, 1 attachment"} {
		if !strings.Contains(text, want) {
			t.Errorf("output must contain %q:\n%s", want, text)
		}
	}
	if strings.Contains(errText, "\x1b") {
		t.Errorf("stderr must not contain ANSI escapes: %q", errText)
	}
	for _, want := range []string{"header.json", "manifest.json", "declarations.json", "integrity.json"} {
		if _, err := os.Stat(filepath.Join(repo, "exchange", "snapshots", want)); err != nil {
			t.Errorf("snapshot missing %s: %v", want, err)
		}
	}
}

// TestSyncSecondRunUnchanged: the second sync reports unchanged.
func TestSyncSecondRunUnchanged(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"sync", repo}); code != 0 {
		t.Fatalf("first sync: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"sync", repo})
	if code != 0 {
		t.Fatalf("second sync: exit = %d\nstdout: %s\nstderr: %s", code, text, errText)
	}
	for _, want := range []string{"unchanged", "snapshot up to date"} {
		if !strings.Contains(text, want) {
			t.Errorf("second sync must report unchanged, missing %q:\n%s", want, text)
		}
	}
}

// TestSyncPullFromDocs: pull --from-docs re-seeds from the docs tree.
func TestSyncPullFromDocs(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"sync", repo}); code != 0 {
		t.Fatalf("first sync: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"sync", "pull", repo, "--from-docs"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "docs: 4 units, 1 attachment") {
		t.Errorf("pull --from-docs must report the docs source:\n%s", text)
	}
}

// TestSyncPushOnly: push alone assembles the snapshot without pulling.
func TestSyncPushOnly(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	code, _, errText := runIn([]string{"sync", "push", repo})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errText)
	}
	// Nothing pulled, so the store is empty and the push is a no-op.
	if _, err := os.Stat(filepath.Join(repo, "exchange", "snapshots")); !os.IsNotExist(err) {
		t.Error("push of an empty store must not write a snapshot")
	}
	// Pull first, then push.
	if code, _, errText := runIn([]string{"sync", "pull", repo}); code != 0 {
		t.Fatalf("pull: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"sync", "push", repo})
	if code != 0 {
		t.Fatalf("push after pull: exit = %d\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "rsf-repo-eka-sync-fixture-1.1") {
		t.Errorf("push output must carry the snapshot label:\n%s", text)
	}
	if _, err := os.Stat(filepath.Join(repo, "exchange", "snapshots", "header.json")); err != nil {
		t.Error("push must write the snapshot")
	}
}

// TestSyncValidationFailureExitsOne: a non-conformant repository (no
// snapshot yet) is refused by the docs gate with the report — exit 1.
func TestSyncValidationFailureExitsOne(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := t.TempDir()
	dir := filepath.Join(repo, "docs", "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	bad := "---\nnamespace: x\nid: 1\n---\n# bad\n" // type missing: R0 error
	if err := os.WriteFile(filepath.Join(dir, "bad.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	code, text, errText := runIn([]string{"sync", repo})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "Verdict: FAIL") {
		t.Errorf("stdout must contain the validation report:\n%s", text)
	}
	if !strings.Contains(errText, "knowledge compile refused") {
		t.Errorf("stderr must explain the refusal, got %q", errText)
	}
}

// TestSyncCorruptSnapshotExitsOne: a tampered snapshot is an integrity
// failure — exit 1.
func TestSyncCorruptSnapshotExitsOne(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"sync", repo}); code != 0 {
		t.Fatalf("first sync: exit %d\n%s", code, errText)
	}
	content := filepath.Join(repo, "exchange", "snapshots", "units", "eka-sync-fixture", "adr-001-runtime-v1", "content")
	data, err := os.ReadFile(content)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(content, append([]byte("X"), data[1:]...), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, errText := runIn([]string{"sync", repo})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstderr: %s", code, errText)
	}
	if !strings.Contains(errText, "snapshot package refused") {
		t.Errorf("stderr must explain the integrity refusal, got %q", errText)
	}
}

// TestSyncOutputDeterministic: two sync runs of identical state
// produce byte-identical output. The first run settles the repository
// (docs pull + registration); the second and third runs both report
// "unchanged".
func TestSyncOutputDeterministic(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"sync", repo}); code != 0 {
		t.Fatalf("settling sync: exit %d\n%s", code, errText)
	}
	run := func() string {
		_, text, _ := runIn([]string{"sync", repo})
		return text
	}
	first := run()
	second := run()
	if first != second {
		t.Error("sync output differs between identical runs")
	}
}
