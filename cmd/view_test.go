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
// calm mini-landing listing the canonical projections and their
// aliases — exit 0.
func TestViewNoArgsListsProjections(t *testing.T) {
	code, out, errText := runIn([]string{"view"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstderr: %s", code, errText)
	}
	if errText != "" {
		t.Errorf("stderr must be empty, got %q", errText)
	}
	for _, want := range []string{
		"Knowledge Projections",
		"discovery", "architecture", "planning", "execution", "operations", "ticket",
		"Aliases", "sprint", "wave",
		"eka view ticket <tkt-id>",
	} {
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
		for _, want := range []string{"eka view", "discovery", "architecture", "planning", "execution", "operations", "ticket", "sprint", "wave"} {
			if !strings.Contains(text, want) {
				t.Errorf("args %v: help missing %q:\n%s", args, want, text)
			}
		}
	}
}

// TestViewUnknownProjectionExitsTwo: an unregistered projection is a
// usage error with the available list (canonical + aliases) — exit 2,
// no repository access.
func TestViewUnknownProjectionExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"view", "bogus"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "unknown projection \"bogus\"") {
		t.Errorf("stderr must name the projection, got %q", errText)
	}
	if !strings.Contains(errText, "available projections: architecture, discovery, execution, operations, planning, ticket (aliases: sprint, wave)") {
		t.Errorf("stderr must list canonical projections and aliases, got %q", errText)
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
	code, _, _ := runIn([]string{"view", "execution", "a", "b"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

// TestViewExecutionHappyPath: the execution projection of a conformant
// fixture — header, tickets with projected status, state columns and
// the merged summary — exit 0.
func TestViewExecutionHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "execution"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Execution",
		"Container    eka-view-fixture/ctr:wave-1",
		"Repository   .",
		"Knowledge    EKA v1",
		"Domain       Execution",
		"↓ View",
		"Tickets (8)",
		"eka-view-fixture/tkt:ts-gamma (in-progress)",
		"eka-view-fixture/tkt:unresolved (unresolved)",
		"planned (1)",
		"todo (1)",
		"in-progress (1)",
		"in-review (1)",
		"done (1)",
		"  • eka-view-fixture/sto:alpha",
		"  → eka-view-fixture/ts:gamma",
		"  ✓ eka-view-fixture/ch:epsilon",
		"Summary:",
		"Container: eka-view-fixture/ctr:wave-1",
		"Tickets: 8",
		"Work items: 5",
		"In progress: 1",
		"Done: 1",
		"Status: active",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewExecutionAliasesIdentical: the sprint and wave aliases render
// byte-identical output to the canonical execution projection.
func TestViewExecutionAliasesIdentical(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	runOnce := func(args ...string) string {
		_, out, _ := runIn(args)
		return out
	}
	execution := runOnce([]string{"view", "execution"}...)
	for _, alias := range []string{"sprint", "wave"} {
		if got := runOnce([]string{"view", alias}...); got != execution {
			t.Errorf("view %s output must be byte-identical to view execution", alias)
		}
	}
}

// TestViewMultipleActiveWarning: the multi-active container anomaly is
// surfaced at the CLI — the warning line names the deterministically
// chosen container (lexicographically smallest canonical identity) and
// the command still exits 0.
func TestViewMultipleActiveWarning(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "multi-active"))
	code, out, errText := runIn([]string{"view", "execution"})
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
	_, out2, _ := runIn([]string{"view", "execution"})
	if out != out2 {
		t.Error("multi-active execution output is not deterministic")
	}
}

// TestViewPlanningHappyPath: the planning projection — artifact groups
// in fixed order with content state, planning state and phase context.
func TestViewPlanningHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "planning"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Planning",
		"Domain       Planning",
		"Scope Definitions",
		"eka-view-fixture/scp:wave-2  (approved, phase mvp)",
		"Epics",
		"eka-view-fixture/epc:auth  (review)",
		"Plans",
		"eka-view-fixture/plan:roadmap-2026  (approved, planning-state approved, phase release)",
		"Traceability",
		"eka-view-fixture/trc:spec-trace  (draft)",
		"Summary:",
		"Scope Definitions: 1",
		"Epics: 1",
		"Plans: 1",
		"Traceability: 1",
		"Plans by state: draft 0, approved 1, immutable 0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewArchitectureHappyPath: the architecture projection — the
// Decisions group merges adr-/dec- (including the superseded ADR), the
// remaining groups are single-token.
func TestViewArchitectureHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "architecture"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Architecture",
		"Domain       Architecture",
		"Decisions",
		"eka-view-fixture/adr:001-login-serialization  (accepted)",
		"eka-view-fixture/adr:002-session-encoding  (superseded)",
		"eka-view-fixture/adr:003-token-format  (accepted)",
		"eka-view-fixture/dec:001-api-shape  (accepted)",
		"Architecture Descriptions",
		"eka-view-fixture/arc:system-architecture  (approved)",
		"Specifications",
		"eka-view-fixture/spec:auth-flow  (draft)",
		"Standards & Guidelines",
		"eka-view-fixture/std:gofmt  (review)",
		"Vocabulary",
		"eka-view-fixture/gls:domain-terms  (amended)",
		"Summary:",
		"Decisions: 4",
		"Architecture Descriptions: 1",
		"Specifications: 1",
		"Standards & Guidelines: 1",
		"Vocabulary: 1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewDiscoveryHappyPath: the discovery projection — vision,
// strategy, requirements and research findings with content states.
func TestViewDiscoveryHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "discovery"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Discovery",
		"Domain       Discovery",
		"Vision",
		"eka-view-fixture/vis:product-vision  (draft)",
		"Strategy",
		"eka-view-fixture/str:go-to-market  (review)",
		"Requirements",
		"eka-view-fixture/req:onboarding  (approved)",
		"Research Findings",
		"eka-view-fixture/fnd:market-research  (approved)",
		"Summary:",
		"Vision: 1",
		"Strategy: 1",
		"Requirements: 1",
		"Research Findings: 1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewOperationsHappyPath: the operations projection — runbooks and
// release records with content states.
func TestViewOperationsHappyPath(t *testing.T) {
	chdirInto(t, viewFixtureAbs(t, "valid"))
	code, out, errText := runIn([]string{"view", "operations"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{
		"Operations",
		"Domain       Operations",
		"Runbooks",
		"eka-view-fixture/run:deploy  (approved)",
		"Release Records",
		"eka-view-fixture/rel:release-1  (review)",
		"Summary:",
		"Runbooks: 1",
		"Release Records: 1",
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
	code, out, errText := runIn([]string{"view", "execution"})
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
// conformant; the execution projection renders a calm "No active
// container" line and still exits 0.
func TestViewEmptyProjectionExitsZero(t *testing.T) {
	chdirInto(t, t.TempDir())
	code, out, errText := runIn([]string{"view", "execution"})
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout: %s\nstderr: %s", code, out, errText)
	}
	for _, want := range []string{"No active container.", "Work items: 0", "Tickets: 0", "Status: no active container"} {
		if !strings.Contains(out, want) {
			t.Errorf("output must contain %q:\n%s", want, out)
		}
	}
}

// TestViewEmptyDomainExitsZero: a repository without artifacts of a
// domain renders a calm "No <Domain> artifacts." line per domain and
// still exits 0.
func TestViewEmptyDomainExitsZero(t *testing.T) {
	chdirInto(t, t.TempDir())
	for domain, want := range map[string]string{
		"planning":     "No Planning artifacts.",
		"architecture": "No Architecture artifacts.",
		"discovery":    "No Discovery artifacts.",
		"operations":   "No Operations artifacts.",
	} {
		code, out, errText := runIn([]string{"view", domain})
		if code != 0 {
			t.Fatalf("view %s: exit = %d, want 0\nstdout: %s\nstderr: %s", domain, code, out, errText)
		}
		if !strings.Contains(out, want) {
			t.Errorf("view %s must render %q:\n%s", domain, want, out)
		}
		// The summary block follows the calm line; assert its shape.
		if !strings.Contains(out, "Summary:") {
			t.Errorf("view %s must still render the summary:\n%s", domain, out)
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
		{"view", "discovery"},
		{"view", "architecture"},
		{"view", "planning"},
		{"view", "execution"},
		{"view", "operations"},
		{"view", "ticket", "tkt-ts-gamma"},
		{"view", "ticket", "tkt-unresolved"},
		{"view", "sprint"},
		{"view", "wave"},
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
		{"view", "discovery"},
		{"view", "architecture"},
		{"view", "planning"},
		{"view", "execution"},
		{"view", "operations"},
		{"view", "ticket", "tkt-ts-gamma"},
		{"view", "sprint"},
		{"view", "wave"},
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
