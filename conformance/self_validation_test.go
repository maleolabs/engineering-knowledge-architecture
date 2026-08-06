package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestReferenceImplementationConforms codifies the milestone:
// "the EKA repository is the first repository to pass its own conformance
// suite" (standard/eka-specification-v1.0.md §16.1).
//
// It locates the repository root by walking up from this test file until
// go.mod is found, runs Validate on the whole repository, and asserts zero
// blocking errors (warnings are allowed and do not affect compliance).
func TestReferenceImplementationConforms(t *testing.T) {
	root := findRepoRoot(t)
	report, err := Validate(root)
	if err != nil {
		t.Fatalf("Validate(%s): %v", root, err)
	}
	if report.Artifacts < 8 {
		t.Errorf("expected at least the 8 reference ADRs, got %d artifacts", report.Artifacts)
	}
	if report.ErrorCount() != 0 {
		t.Errorf("the EKA repository must pass its own conformance suite: %d errors:\n%s",
			report.ErrorCount(), dumpResults(report))
	}
}

// findRepoRoot walks up from the test source file until it finds go.mod.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found above test file")
		}
		dir = parent
	}
}

func TestFindRepoRoot(t *testing.T) {
	root := findRepoRoot(t)
	if root == "" || root == "/" {
		t.Fatalf("unexpected repo root %q", root)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("repo root %q has no go.mod", root)
	}
	fmt.Println("repository root:", root)
}
