package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Execute() unit tests: argument handling, output contract and exit
// codes. Ported from the pre-Cobra cmd/eka/main_test.go (which tested
// run()) and extended with help/completion coverage.
//
// Exit code contract (documented in root.go):
//
//	0  fully compliant (warnings allowed)
//	1  blocking violations present
//	2  usage or internal error (unknown command, invalid path,
//	   unreadable scan root)

// runIn calls Execute with a non-interactive stdin (deterministic
// defaults for init) and returns the collected stdout/stderr.
func runIn(args []string) (int, string, string) {
	var out, errb bytes.Buffer
	code := Execute(args, strings.NewReader(""), &out, &errb)
	return code, out.String(), errb.String()
}

func TestExitCodeUsage(t *testing.T) {
	code, _, errText := runIn(nil)
	if code != 2 {
		t.Errorf("no args: exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "Usage:") {
		t.Errorf("no args: stderr must contain the usage, got %q", errText)
	}
	code, _, errText = runIn([]string{"doctor"})
	if code != 2 {
		t.Errorf("unknown command: exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "validate") {
		t.Errorf("unknown command message must mention the validate command, got %q", errText)
	}
	if !strings.Contains(errText, "init") {
		t.Errorf("unknown command message must mention the init command, got %q", errText)
	}
	// Cobra sorts command lists alphabetically (EnableCommandSorting);
	// the message must list exactly the registered commands.
	if !strings.Contains(errText, "available commands: init, validate") {
		t.Errorf("unknown command message must list the available commands, got %q", errText)
	}
	code, _, _ = runIn([]string{"validate", "a", "b"})
	if code != 2 {
		t.Errorf("too many args: exit = %d, want 2", code)
	}
	// Unknown root-level flag is a usage error, not an unknown command.
	code, _, errText = runIn([]string{"--bogus"})
	if code != 2 {
		t.Errorf("unknown root flag: exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "--bogus") {
		t.Errorf("unknown root flag: stderr must name the flag, got %q", errText)
	}
}

func TestExitCodeBadPath(t *testing.T) {
	code, _, _ := runIn([]string{"validate", filepath.Join(t.TempDir(), "nope")})
	if code != 2 {
		t.Errorf("invalid path: exit = %d, want 2", code)
	}
}

// TestHelpExitsZero covers the root help entry points (-h, --help, help)
// and the validate help flags. All must exit 0 and print to stdout.
// Subcommand help is `eka help <command>` (Cobra); `eka validate help`
// would be treated as a path argument, matching pre-Cobra behavior.
func TestHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"-h"}, {"--help"}, {"help"}} {
		var out, errb bytes.Buffer
		if code := Execute(args, strings.NewReader(""), &out, &errb); code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		text := out.String()
		if !strings.Contains(text, "Usage:") {
			t.Errorf("args %v: root help must contain Usage:", args)
		}
		for _, cmdName := range []string{"validate", "init"} {
			if !strings.Contains(text, cmdName) {
				t.Errorf("args %v: root help must mention the %s command", args, cmdName)
			}
		}
	}
	for _, args := range [][]string{{"validate", "-h"}, {"validate", "--help"}} {
		var out, errb bytes.Buffer
		if code := Execute(args, strings.NewReader(""), &out, &errb); code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		if !strings.Contains(out.String(), "eka validate") {
			t.Errorf("args %v: validate help text missing usage", args)
		}
	}
}

func TestValidateValidRepoExitsZero(t *testing.T) {
	code, text, errText := runIn([]string{"validate", filepath.Join("..", "conformance", "testdata", "valid")})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errText)
	}
	if !strings.Contains(text, "PASS") {
		t.Errorf("output must contain PASS:\n%s", text)
	}
	if !strings.Contains(text, "Artifacts: 6") {
		t.Errorf("output must contain the artifact count:\n%s", text)
	}
}

func TestValidateInvalidRepoExitsOne(t *testing.T) {
	code, text, _ := runIn([]string{"validate", filepath.Join("..", "conformance", "testdata", "invalid-dimension")})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\noutput: %s", code, text)
	}
	if !strings.Contains(text, "FAIL") {
		t.Errorf("output must contain FAIL:\n%s", text)
	}
	if !strings.Contains(text, "R6") {
		t.Errorf("output must list the R6 violations:\n%s", text)
	}
}

// TestWarningsDoNotAffectExitCode builds a repo whose only finding is a
// warning (a draft artifact with an unresolved reference) and verifies
// exit code 0.
func TestWarningsDoNotAffectExitCode(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "specifications")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `---
namespace: eka-cli
type: spec
id: 001-draft
instance-version: 1
revision: 1
content-state: draft
existence-state: active
dimension: specifications
created: 2026-08-05
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on:
  - sto:ghost
change-log:
  - date: 2026-08-05
    domain: existence-state
    from: "-"
    to: active
    by: Engineering
  - date: 2026-08-05
    domain: content-state
    from: "-"
    to: draft
    by: Engineering
---
# Spec

## Purpose

p

## Content

c
`
	if err := os.WriteFile(filepath.Join(dir, "spec-001-draft.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	code, text, _ := runIn([]string{"validate", root})
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (warnings must not block)\noutput:\n%s", code, text)
	}
	if !strings.Contains(text, "Warnings:  1") {
		t.Errorf("output must count the warning:\n%s", text)
	}
}

// TestDefaultPathIsCurrentDirectory verifies `eka validate` uses ".".
func TestDefaultPathIsCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Empty directory is trivially compliant.
	code, text, _ := runIn([]string{"validate"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, text)
	}
	if !strings.Contains(text, "Root:      .") {
		t.Errorf("default root must be '.':\n%s", text)
	}
}

// TestOutputIsDeterministic verifies stable ordering in printed results.
func TestOutputIsDeterministic(t *testing.T) {
	path := filepath.Join("..", "conformance", "testdata", "invalid-projection")
	runOnce := func() string {
		_, text, _ := runIn([]string{"validate", path})
		return text
	}
	if a, b := runOnce(), runOnce(); a != b {
		t.Error("CLI output differs between runs")
	}
}

// TestCompletionCommand verifies the standard Cobra completion command:
// every shell subcommand generates a script on stdout and exits 0.
func TestCompletionCommand(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		var out, errb bytes.Buffer
		code := Execute([]string{"completion", shell}, strings.NewReader(""), &out, &errb)
		if code != 0 {
			t.Errorf("completion %s: exit = %d, want 0\nstderr: %s", shell, code, errb.String())
		}
		if out.Len() == 0 {
			t.Errorf("completion %s: output must not be empty", shell)
		}
	}
}

// --- eka init CLI-level tests -------------------------------------------

// chdirInto changes the working directory for the duration of the test.
func chdirInto(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(old) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

func TestInitHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"init", "-h"}, {"init", "--help"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		if !strings.Contains(text, "eka init") {
			t.Errorf("args %v: init help text missing usage", args)
		}
	}
}

func TestInitUnknownFlagExitsTwo(t *testing.T) {
	for _, flag := range []string{"-n", "--force", "-x"} {
		code, _, errText := runIn([]string{"init", flag})
		if code != 2 {
			t.Errorf("flag %q: exit = %d, want 2", flag, code)
		}
		if !strings.Contains(errText, flag) {
			t.Errorf("flag %q: stderr must name the unknown flag, got %q", flag, errText)
		}
	}
}

func TestInitTooManyArgumentsExitsTwo(t *testing.T) {
	code, _, _ := runIn([]string{"init", "a", "b"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestInitRejectsPathSeparators(t *testing.T) {
	code, _, errText := runIn([]string{"init", "a/b"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "path separators") {
		t.Errorf("stderr must explain the separator rejection, got %q", errText)
	}
}

func TestInitRejectsFileTarget(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "occupied")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	code, _, errText := runIn([]string{"init", "occupied"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "not a directory") {
		t.Errorf("stderr must explain the file target rejection, got %q", errText)
	}
}

// TestInitDryRunWritesNothing verifies `eka init --dry-run` prints the plan
// and creates nothing.
func TestInitDryRunWritesNothing(t *testing.T) {
	dir := t.TempDir()
	chdirInto(t, dir)
	code, text, _ := runIn([]string{"init", "--dry-run"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(text, "EKA Bootstrap Plan (dry-run)") {
		t.Errorf("output must show the plan header:\n%s", text)
	}
	if !strings.Contains(text, "create dir: docs/") {
		t.Errorf("output must contain plan lines:\n%s", text)
	}
	if !strings.Contains(text, "Dry-run: no changes were written.") {
		t.Errorf("output must state that nothing was written:\n%s", text)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("dry-run must not create files, found: %v", entries)
	}
}

// TestInitNamedDryRun verifies the flag works after the project name and
// that the named target is not created.
func TestInitNamedDryRun(t *testing.T) {
	dir := t.TempDir()
	chdirInto(t, dir)
	code, text, _ := runIn([]string{"init", "myproj", "--dry-run"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(text, "create dir: myproj") {
		t.Errorf("plan must target myproj:\n%s", text)
	}
	if _, err := os.Stat(filepath.Join(dir, "myproj")); !os.IsNotExist(err) {
		t.Error("dry-run must not create the project directory")
	}
}

// TestInitCurrentDirectory verifies `eka init` bootstraps "." and the
// result validates.
func TestInitCurrentDirectory(t *testing.T) {
	dir := t.TempDir()
	chdirInto(t, dir)
	code, text, errText := runIn([]string{"init"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s\nstdout: %s", code, errText, text)
	}
	if !strings.Contains(text, "Validation Result: PASS") {
		t.Errorf("output must report PASS:\n%s", text)
	}
	for _, want := range []string{"docs/operating/protocol.md", "docs/exchange/validation.md", "README.md"} {
		if _, err := os.Stat(filepath.Join(dir, want)); err != nil {
			t.Errorf("expected generated file %s: %v", want, err)
		}
	}
	// The generated repository must pass `eka validate`.
	code, text, _ = runIn([]string{"validate", dir})
	if code != 0 {
		t.Errorf("generated repository must validate, exit = %d:\n%s", code, text)
	}
}

// TestInitNamedDirectory verifies `eka init <name>` creates the directory
// relative to the CWD.
func TestInitNamedDirectory(t *testing.T) {
	dir := t.TempDir()
	chdirInto(t, dir)
	code, text, _ := runIn([]string{"init", "myproj"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0:\n%s", code, text)
	}
	if _, err := os.Stat(filepath.Join(dir, "myproj", "docs")); err != nil {
		t.Errorf("expected myproj/docs to exist: %v", err)
	}
	if !fieldIs(text, "Repository Type:", "new") {
		t.Errorf("output must classify the repo as new:\n%s", text)
	}
}

// fieldIs reports whether a line in text starts with label and its value
// (after the colon) equals want, ignoring alignment padding.
func fieldIs(text, label, want string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, label) {
			value := strings.TrimSpace(strings.TrimPrefix(line, label))
			return value == want
		}
	}
	return false
}

// TestInitNonConformantExitsOne pre-places a malformed artifact and
// verifies the init exits 1 with a clear message.
func TestInitNonConformantExitsOne(t *testing.T) {
	dir := t.TempDir()
	workItems := filepath.Join(dir, "docs", "operating", "work-items")
	if err := os.MkdirAll(workItems, 0o755); err != nil {
		t.Fatal(err)
	}
	bad := "---\nnamespace: eka-cli\n"
	bad += "type: sto\n" // type without id violates the artifact rule (R0)
	bad += "---\n# Bad\n"
	if err := os.WriteFile(filepath.Join(workItems, "sto-bad.md"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	chdirInto(t, dir)
	code, text, errText := runIn([]string{"init"})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout: %s\nstderr: %s", code, text, errText)
	}
	if !strings.Contains(errText, "non-conformant") {
		t.Errorf("stderr must explain the non-conformant result, got %q", errText)
	}
	if !strings.Contains(text, "Validation Result: FAIL") {
		t.Errorf("output must report FAIL:\n%s", text)
	}
}

// TestInitTwiceIsIdempotent verifies the second run is a no-op and the
// repository still validates.
func TestInitTwiceIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	chdirInto(t, dir)
	code, _, errText := runIn([]string{"init"})
	if code != 0 {
		t.Fatalf("first init: exit = %d, want 0\nstderr: %s", code, errText)
	}
	before, err := snapshot(dir)
	if err != nil {
		t.Fatal(err)
	}
	code, text, _ := runIn([]string{"init"})
	if code != 0 {
		t.Fatalf("second init: exit = %d, want 0:\n%s", code, text)
	}
	if !strings.Contains(text, "Already initialized") {
		t.Errorf("second run must report already initialized:\n%s", text)
	}
	after, err := snapshot(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Errorf("second init changed the tree: %d files before, %d after", len(before), len(after))
	}
	for path, data := range before {
		if got, ok := after[path]; !ok || !bytes.Equal(got, data) {
			t.Errorf("second init modified %s", path)
		}
	}
	code, _, _ = runIn([]string{"validate", dir})
	if code != 0 {
		t.Errorf("repository must still validate after second init")
	}
}

// snapshot walks dir and returns relative path -> content.
func snapshot(dir string) (map[string][]byte, error) {
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
		out[rel] = data
		return nil
	})
	return out, err
}
