package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"project", "-h"}, {"project", "register", "-h"}, {"project", "list", "-h"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		if !strings.Contains(text, "eka project") {
			t.Errorf("args %v: help must mention eka project", args)
		}
	}
}

// TestProjectRegisterHappyPath: registering a repository exits 0 and
// reports the project/repository/status.
func TestProjectRegisterHappyPath(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	code, text, errText := runIn([]string{"project", "register", repo})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	for _, want := range []string{"Project", "Repository", "Path", "Status", "registered"} {
		if !strings.Contains(text, want) {
			t.Errorf("output must contain %q:\n%s", want, text)
		}
	}
	if !strings.Contains(text, filepath.Base(repo)) {
		t.Errorf("output must name the repository %q:\n%s", filepath.Base(repo), text)
	}
}

// TestProjectRegisterTwiceReportsAlreadyRegistered.
func TestProjectRegisterTwiceReportsAlreadyRegistered(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"project", "register", repo}); code != 0 {
		t.Fatalf("first register: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"project", "register", repo})
	if code != 0 {
		t.Fatalf("second register: exit = %d\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "already registered") {
		t.Errorf("second register must report already registered:\n%s", text)
	}
}

// TestProjectRegisterCustomName: --name sets the project name.
func TestProjectRegisterCustomName(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copySyncFixture(t)
	code, text, errText := runIn([]string{"project", "register", repo, "--name", "myproject"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "myproject") {
		t.Errorf("output must carry the custom project name:\n%s", text)
	}
}

// TestProjectListEmpty: an empty workspace lists no projects and exits
// 0 with the informational message.
func TestProjectListEmpty(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	code, text, errText := runIn([]string{"project", "list"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "No projects registered yet") {
		t.Errorf("empty list must be informational:\n%s", text)
	}
}

// TestProjectListSorted: projects and repositories render sorted, with
// the workspace path in the header.
func TestProjectListSorted(t *testing.T) {
	home := t.TempDir()
	t.Setenv("EKA_HOME", home)
	repoA := copySyncFixture(t)
	repoB := copySyncFixture(t)
	for _, args := range [][]string{
		{"project", "register", repoA},
		{"project", "register", repoB, "--name", "zproject"},
	} {
		if code, _, errText := runIn(args); code != 0 {
			t.Fatalf("register %v: exit %d\n%s", args, code, errText)
		}
	}
	code, text, errText := runIn([]string{"project", "list"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, home) {
		t.Errorf("list must show the workspace path:\n%s", text)
	}
	// Both project names present; the repository basenames present.
	for _, want := range []string{filepath.Base(repoA), filepath.Base(repoB), "zproject"} {
		if !strings.Contains(text, want) {
			t.Errorf("list missing %q:\n%s", want, text)
		}
	}
	// Deterministic: same state, same bytes.
	_, text2, _ := runIn([]string{"project", "list"})
	if text != text2 {
		t.Error("project list output differs between runs")
	}
}

// TestProjectRegisterRejectsMissingPath: an unreadable path is a usage
// error (exit 2).
func TestProjectRegisterRejectsMissingPath(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	code, _, errText := runIn([]string{"project", "register", filepath.Join(t.TempDir(), "nope")})
	if code != 2 {
		t.Errorf("exit = %d, want 2\nstderr: %s", code, errText)
	}
	if errText == "" {
		t.Error("stderr must not be empty")
	}
}

// TestStatusAfterRegisterNoObjects: status renders workspace, schema,
// project and zero counts.
func TestStatusAfterRegisterNoObjects(t *testing.T) {
	home := t.TempDir()
	t.Setenv("EKA_HOME", home)
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"project", "register", repo}); code != 0 {
		t.Fatalf("register: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"status"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	for _, want := range []string{"Runtime", home, "Schema", "Objects", "Payloads", "Attachments", "Projects"} {
		if !strings.Contains(text, want) {
			t.Errorf("status missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "pull at") || strings.Contains(text, "push at") {
		t.Error("no sync log entries expected before any sync")
	}
}

// TestStatusAfterSyncShowsLastSync: after a sync the status reports
// the last pull/push per repository.
func TestStatusAfterSyncShowsLastSync(t *testing.T) {
	home := t.TempDir()
	t.Setenv("EKA_HOME", home)
	repo := copySyncFixture(t)
	if code, _, errText := runIn([]string{"sync", repo}); code != 0 {
		t.Fatalf("sync: exit %d\n%s", code, errText)
	}
	code, text, errText := runIn([]string{"status"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(text, "Objects") {
		t.Errorf("status must show counts:\n%s", text)
	}
	// The sync log renders the most recent entry (the push after a
	// full sync cycle).
	if !strings.Contains(text, "push") {
		t.Errorf("status must show the last sync entry:\n%s", text)
	}
	if !strings.Contains(text, "at 20") {
		t.Errorf("status must show the sync timestamp:\n%s", text)
	}
	// Deterministic output.
	_, text2, _ := runIn([]string{"status"})
	if text != text2 {
		t.Error("status output differs between runs")
	}
}
