package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// run() unit tests: argument handling, output contract and exit codes.
//
// Exit code contract (documented):
//
//	0  fully compliant (warnings allowed)
//	1  blocking violations present
//	2  usage or internal error (unknown command, invalid path,
//	   unreadable scan root)

func TestExitCodeUsage(t *testing.T) {
	var out, errb bytes.Buffer
	if code := run(nil, &out, &errb); code != 2 {
		t.Errorf("no args: exit = %d, want 2", code)
	}
	out.Reset()
	errb.Reset()
	if code := run([]string{"doctor"}, &out, &errb); code != 2 {
		t.Errorf("unknown command: exit = %d, want 2", code)
	}
	if !strings.Contains(errb.String(), "validate") {
		t.Errorf("unknown command message must mention the validate-only scope, got %q", errb.String())
	}
	out.Reset()
	errb.Reset()
	if code := run([]string{"validate", "a", "b"}, &out, &errb); code != 2 {
		t.Errorf("too many args: exit = %d, want 2", code)
	}
}

func TestExitCodeBadPath(t *testing.T) {
	var out, errb bytes.Buffer
	code := run([]string{"validate", filepath.Join(t.TempDir(), "nope")}, &out, &errb)
	if code != 2 {
		t.Errorf("invalid path: exit = %d, want 2", code)
	}
}

func TestHelpExitsZero(t *testing.T) {
	var out, errb bytes.Buffer
	for _, args := range [][]string{{"-h"}, {"--help"}, {"help"}, {"validate", "--help"}} {
		out.Reset()
		if code := run(args, &out, &errb); code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		if !strings.Contains(out.String(), "eka validate") {
			t.Errorf("args %v: help text missing usage", args)
		}
	}
}

func TestValidateValidRepoExitsZero(t *testing.T) {
	var out, errb bytes.Buffer
	code := run([]string{"validate", filepath.Join("..", "..", "conformance", "testdata", "valid")}, &out, &errb)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errb.String())
	}
	text := out.String()
	if !strings.Contains(text, "PASS") {
		t.Errorf("output must contain PASS:\n%s", text)
	}
	if !strings.Contains(text, "Artifacts: 6") {
		t.Errorf("output must contain the artifact count:\n%s", text)
	}
}

func TestValidateInvalidRepoExitsOne(t *testing.T) {
	var out, errb bytes.Buffer
	code := run([]string{"validate", filepath.Join("..", "..", "conformance", "testdata", "invalid-dimension")}, &out, &errb)
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstderr: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "FAIL") {
		t.Errorf("output must contain FAIL:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "R6") {
		t.Errorf("output must list the R6 violations:\n%s", out.String())
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
	var out, errb bytes.Buffer
	code := run([]string{"validate", root}, &out, &errb)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (warnings must not block)\noutput:\n%s", code, out.String())
	}
	if !strings.Contains(out.String(), "Warnings:  1") {
		t.Errorf("output must count the warning:\n%s", out.String())
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
	var out, errb bytes.Buffer
	if code := run([]string{"validate"}, &out, &errb); code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errb.String())
	}
	if !strings.Contains(out.String(), "Root:      .") {
		t.Errorf("default root must be '.':\n%s", out.String())
	}
}

// TestOutputIsDeterministic verifies stable ordering in printed results.
func TestOutputIsDeterministic(t *testing.T) {
	path := filepath.Join("..", "..", "conformance", "testdata", "invalid-projection")
	runOnce := func() string {
		var out bytes.Buffer
		run([]string{"validate", path}, &out, &bytes.Buffer{})
		return out.String()
	}
	if a, b := runOnce(), runOnce(); a != b {
		t.Error("CLI output differs between runs")
	}
}
