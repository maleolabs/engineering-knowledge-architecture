package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// viewFixtureAbs resolves the absolute path of a view test fixture.
func viewFixtureAbs(t *testing.T, name string) string {
	t.Helper()
	abs, err := filepath.Abs(filepath.Join("..", "view", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// TestViewNoArgsListsProjections: `eka view` without arguments is a
// calm mini-landing listing the available projections — exit 0.
func TestViewNoArgsListsProjections(t *testing.T) {
	code, out, errText := runIn([]string{"view"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errText)
	}
	if errText != "" {
		t.Errorf("stderr must be empty, got %q", errText)
	}
	for _, want := range []string{"Knowledge Projections", "sprint", "wave", "ticket", "eka view ticket <tkt-id>"} {
		if !strings.Contains(out, want) {
			t.Errorf("landing missing %q:\n%s", want, out)
		}
	}
	_, out2, _ := runIn([]string{"view"})
	if out != out2 {
		t.Error("view landing is not deterministic")
	}
}

// TestViewHelpExitsZero covers the help entry points.
func TestViewHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"view", "-h"}, {"view", "--help"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		for _, want := range []string{"eka view", "sprint", "wave", "ticket"} {
			if !strings.Contains(text, want) {
				t.Errorf("args %v: help missing %q:\n%s", args, want, text)
			}
		}
	}
}

// TestViewUnknownProjectionExitsTwo: an unregistered projection is a
// usage error with the available list — exit 2, no repository access.
func TestViewUnknownProjectionExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"view", "bogus"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "unknown projection \"bogus\"") {
		t.Errorf("stderr must name the projection, got %q", errText)
	}
	if !strings.Contains(errText, "available projections: sprint, ticket, wave") {
		t.Errorf("stderr must list the available projections, got %q", errText)
	}
}

// TestViewTicketMissingTargetExitsTwo: the ticket projection requires
// its target.
func TestViewTicketMissingTargetExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"view", "ticket"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "requires a target") {
		t.Errorf("stderr must explain the requirement, got %q", errText)
	}
}

// TestViewTooManyArgsExitsTwo: at most one projection and one target.
func TestViewTooManyArgsExitsTwo(t *testing.T) {
	code, _, _ := runIn([]string{"view", "sprint", "a", "b"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

// TestViewSprintHappyPath: the sprint projection of a conformant
// fixture — exit 0 with the header, columns and summary.
func TestViewSprintHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "sprint"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Sprint",
		"Container    eka-view-fixture/ctr:wave-1",
		"Repository   .",
		"Knowledge    EKA v1",
		"Domain       Execution",
		"↓ View",
		"planned (1)",
		"todo (1)",
		"in-progress (1)",
		"in-review (1)",
		"done (1)",
		"  • eka-view-fixture/sto:alpha",
		"  → eka-view-fixture/ts:gamma",
		"  ✓ eka-view-fixture/ch:epsilon",
		"Summary:",
		"Work items: 5",
		"In progress: 1",
		"Done: 1",
		"Tickets: 8",
		"Status: active",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewMultipleActiveWarning: the multi-active container anomaly is
// surfaced at the CLI — the warning line names the deterministically
// chosen container (lexicographically smallest canonical identity) and
// the command still exits 0. Regression: the warning was only covered at
// graph/projection level.
func TestViewMultipleActiveWarning(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "multi-active"))
	code, out, errText := runIn([]string{"view", "sprint"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if errText != "" {
		t.Errorf("stderr must be empty, got %q", errText)
	}
	for _, want := range []string{
		"Multiple active containers — showing eka-view-fixture/ctr:wave-1",
		"Container    eka-view-fixture/ctr:wave-1",
		"Work items: 0",
		"Tickets: 0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
	// Deterministic across runs, like every other projection.
	_, out2, _ := runIn([]string{"view", "sprint"})
	if out != out2 {
		t.Error("multi-active sprint output is not deterministic")
	}
}

// TestViewWaveHappyPath: the wave projection — tickets with projected
// status and progress counts.
func TestViewWaveHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "wave"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Wave",
		"Tickets (8)",
		"eka-view-fixture/tkt:ts-gamma (in-progress)",
		"eka-view-fixture/tkt:unresolved (unresolved)",
		"Domain       Execution",
		"Progress",
		"planned     1",
		"in-progress 1",
		"done        1",
		"Tickets: 8",
		"Work items: 5",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewTicketHappyPath: the ticket projection derives the projected
// status from the work item's owner state.
func TestViewTicketHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "ticket", "tkt-ts-gamma"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Ticket",
		"Ticket       eka-view-fixture/tkt:ts-gamma",
		"Domain       Execution",
		"Projected Status   in-progress",
		"Work Item          eka-view-fixture/ts:gamma (in-progress)",
		"Container          eka-view-fixture/ctr:wave-1",
		"Derives From       ctr:wave-1, ts:gamma",
		"Work item: eka-view-fixture/ts:gamma (in-progress)",
		"Status: in-progress",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewTicketBareID: a bare ticket id resolves like tkt-<id>.
func TestViewTicketBareID(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, _ := runIn([]string{"view", "ticket", "ts-gamma"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	if !strings.Contains(out, "eka-view-fixture/tkt:ts-gamma") {
		t.Errorf("bare id must resolve to the same ticket:\n%s", out)
	}
}

// TestViewTicketUnresolved: a ticket without a resolvable work item
// renders an explicit unresolved status — exit 0.
func TestViewTicketUnresolved(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, _ := runIn([]string{"view", "ticket", "tkt-unresolved"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	for _, want := range []string{
		"Projected Status   unresolved",
		"Work Item          unresolved",
		"Container          eka-view-fixture/ctr:wave-1",
		"Status: unresolved",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewTicketNotFoundExitsTwo: an unknown ticket target is a usage
// error with the available tickets.
func TestViewTicketNotFoundExitsTwo(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, _, errText := runIn([]string{"view", "ticket", "tkt-ghost"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "target \"tkt-ghost\" not found") {
		t.Errorf("stderr must explain the missing target, got %q", errText)
	}
	if !strings.Contains(errText, "available tickets: bug-delta, ch-epsilon, sto-alpha, sto-alpha-dup, sto-beta, sto-beta-multi, sto-legacy, ts-gamma, unresolved") {
		t.Errorf("stderr must list the available tickets, got %q", errText)
	}
}

// TestViewInvalidRepoExitsOne: the conformance gate runs first — a
// repository with blocking violations is refused, the report is
// printed, exit 1, no projection.
func TestViewInvalidRepoExitsOne(t *testing.T) {
	dir := t.TempDir()
	workItems := filepath.Join(dir, "docs", "operating", "work-items", "stories")
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
	code, out, errText := runIn([]string{"view", "sprint"})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout: %s\nstderr: %s", code, out, errText)
	}
	if !strings.Contains(out, "Verdict: FAIL") {
		t.Errorf("stdout must contain the validation report:\n%s", out)
	}
	if !strings.Contains(errText, "view refused") {
		t.Errorf("stderr must explain the refusal, got %q", errText)
	}
}

// TestViewEmptyProjectionExitsZero: an empty directory is trivially
// conformant; the sprint projection renders a calm "No active
// container" line and still exits 0.
func TestViewEmptyProjectionExitsZero(t *testing.T) {
	chdirInto(t, t.TempDir())
	code, out, errText := runIn([]string{"view", "sprint"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{"No active container.", "Work items: 0", "Tickets: 0", "Status: no active container"} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewDeterministicCLI: two runs of each projection produce
// byte-identical output.
func TestViewDeterministicCLI(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	runOnce := func(args ...string) string {
		_, out, _ := runIn(args)
		return out
	}
	for _, args := range [][]string{
		{"view"},
		{"view", "sprint"},
		{"view", "wave"},
		{"view", "ticket", "tkt-ts-gamma"},
		{"view", "ticket", "tkt-unresolved"},
	} {
		if a, b := runOnce(args...), runOnce(args...); a != b {
			t.Errorf("output differs between runs for %v", args)
		}
	}
}

// TestViewNoANSIEscapesInNonTTYOutput verifies the determinism
// contract for the view command: non-TTY output carries no ANSI
// escapes.
func TestViewNoANSIEscapesInNonTTYOutput(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	for _, args := range [][]string{
		{"view"},
		{"view", "sprint"},
		{"view", "wave"},
		{"view", "ticket", "tkt-ts-gamma"},
	} {
		var out, errb bytes.Buffer
		code := Execute(args, strings.NewReader(""), &out, &errb)
		if strings.Contains(out.String(), "\x1b") || strings.Contains(errb.String(), "\x1b") {
			t.Errorf("%v: non-TTY output must not contain ANSI escapes:\nstdout: %q\nstderr: %q",
				args, out.String(), errb.String())
		}
		if code != 0 {
			t.Errorf("%v: exit = %d, want 0", args, code)
		}
	}
}
