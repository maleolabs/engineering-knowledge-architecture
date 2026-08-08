package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/sync"
	"github.com/maleolabs/engineering-knowledge-architecture/workspace"
)

// watch tests: the live projection command. Two layers:
//
//   - usage/terminal gate tests run through Execute with buffer writers
//     (never a TTY): every watch invocation with bad arguments or a
//     non-terminal stdout exits 2 with a deterministic message.
//   - frame tests call the pure helpers (renderWatchFrame, frameChanged,
//     writeClearScreen) directly — the infinite loop is not exercised,
//     only the per-cycle functions, which is where all the state lives.
//     The projection source is the workspace canonical store, so every
//     frame test seeds a fresh workspace (EKA_HOME + sync) first.
//
// watch requires a TTY, so runIn-based success-path tests (like view's)
// are impossible by construction; the frame helpers take the workspace
// explicitly and are tested without a terminal.

// watchFrameStyle returns a non-TTY Style for frame rendering. The
// writer is irrelevant: renderWatchFrame replaces it with its own
// buffer, keeping the Color/TTY/Verbose settings (a non-TTY style
// yields ANSI-free frames, matching the deterministic contract).
func watchFrameStyle() *ui.Style {
	return ui.NewStyle(io.Discard, false)
}

// --- usage and terminal gate -------------------------------------------

// TestWatchNonTTYExitsTwo: stdout is not a terminal — deterministic
// error on stderr, empty stdout, exit 2. This is the CI/pipe contract.
func TestWatchNonTTYExitsTwo(t *testing.T) {
	code, out, errText := runIn([]string{"watch", "execution"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout must be empty on a non-TTY watch, got %q", out)
	}
	if !strings.Contains(errText, "watch requires a terminal (stdout is not a TTY); use 'eka view' for one-shot output") {
		t.Errorf("stderr must explain the terminal requirement deterministically, got %q", errText)
	}
}

// TestWatchUnknownProjectionExitsTwo: same helpful message as view.
func TestWatchUnknownProjectionExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"watch", "bogus"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "unknown projection \"bogus\"") {
		t.Errorf("stderr must name the projection, got %q", errText)
	}
	if !strings.Contains(errText, "available projections: architecture, board, discovery, execution, operations, planning, ticket (aliases: sprint, wave)") {
		t.Errorf("stderr must list canonical projections and aliases, got %q", errText)
	}
}

// TestWatchTicketMissingTargetExitsTwo: the ticket projection requires
// its target, like view.
func TestWatchTicketMissingTargetExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"watch", "ticket"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "requires a target") {
		t.Errorf("stderr must explain the requirement, got %q", errText)
	}
}

// TestWatchNoProjectionExitsTwo: unlike view, watch has no landing —
// a projection is required.
func TestWatchNoProjectionExitsTwo(t *testing.T) {
	code, _, errText := runIn([]string{"watch"})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "watch requires a projection") {
		t.Errorf("stderr must ask for a projection, got %q", errText)
	}
	if !strings.Contains(errText, "available projections:") {
		t.Errorf("stderr must list the projections, got %q", errText)
	}
}

// TestWatchInvalidIntervalExitsTwo: interval below 1 is a usage error
// (exit 2); a non-numeric interval is a flag parse error (exit 2).
func TestWatchInvalidIntervalExitsTwo(t *testing.T) {
	for _, args := range [][]string{
		{"watch", "execution", "--interval", "0"},
		{"watch", "execution", "--interval", "-1"},
		{"watch", "execution", "--interval=-1"},
	} {
		code, _, errText := runIn(args)
		if code != 2 {
			t.Errorf("%v: exit = %d, want 2", args, code)
		}
		if !strings.Contains(errText, "invalid interval") {
			t.Errorf("%v: stderr must explain the invalid interval, got %q", args, errText)
		}
	}
	code, _, errText := runIn([]string{"watch", "execution", "--interval", "abc"})
	if code != 2 {
		t.Errorf("non-numeric interval: exit = %d, want 2", code)
	}
	if !strings.Contains(errText, "--interval") {
		t.Errorf("non-numeric interval: stderr must name the flag, got %q", errText)
	}
}

// TestWatchTooManyArgsExitsTwo: at most one projection and one target.
func TestWatchTooManyArgsExitsTwo(t *testing.T) {
	code, _, _ := runIn([]string{"watch", "execution", "a", "b"})
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

// TestWatchHelpExitsZero: help works without a terminal.
func TestWatchHelpExitsZero(t *testing.T) {
	for _, args := range [][]string{{"watch", "-h"}, {"watch", "--help"}} {
		code, text, _ := runIn(args)
		if code != 0 {
			t.Errorf("args %v: exit = %d, want 0", args, code)
		}
		for _, want := range []string{
			"eka watch", "execution", "--interval", "ticket",
			"must be a terminal", "sprint", "wave",
		} {
			if !strings.Contains(text, want) {
				t.Errorf("args %v: help missing %q:\n%s", args, want, text)
			}
		}
		if !strings.Contains(text, "EKA workspace") {
			t.Errorf("args %v: help must document the workspace canonical source:\n%s", args, text)
		}
	}
}

// --- frame rendering (unit, no loop) -----------------------------------

// TestWatchFrameProjection: one watch cycle over a synced fixture
// renders the projection (byte-identical to the one-shot view output,
// plus the watching footer) — not the refusal frame.
func TestWatchFrameProjection(t *testing.T) {
	seedViewRepo(t, "valid")
	ws := openWorkspace(t)
	frame, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatalf("renderWatchFrame: %v", err)
	}
	text := string(frame)
	for _, want := range []string{
		"Execution",
		"Container    eka-view-fixture/ctr:wave-1",
		"│ Planned (1)",
		"8 tickets project these work items",
		"Summary:",
		"watching — Ctrl-C to stop (interval 2s)",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("frame must contain %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "not registered") {
		t.Errorf("a synced repository must not render the refusal frame:\n%s", text)
	}
	if strings.Contains(text, "\x1b") {
		t.Errorf("non-TTY frame must not contain ANSI escapes:\n%s", text)
	}
}

// TestWatchFrameByteAgreementWithView: the projection part of a watch
// frame is the one-shot `eka view` output byte for byte (same
// renderers, same style settings); the frame only appends the footer.
func TestWatchFrameByteAgreementWithView(t *testing.T) {
	seedViewRepo(t, "valid")
	ws := openWorkspace(t)
	code, viewOut, errText := runIn([]string{"view", "execution"})
	if code != 0 {
		t.Fatalf("view execution: exit = %d\nstderr: %s", code, errText)
	}
	frame, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatalf("renderWatchFrame: %v", err)
	}
	footer := "watching — Ctrl-C to stop (interval 2s)\n"
	if !strings.HasSuffix(string(frame), footer) {
		t.Errorf("frame must end with the watching footer, got:\n%s", frame)
	}
	if got := strings.TrimSuffix(string(frame), footer); got != viewOut {
		t.Errorf("watch frame projection part must equal view output byte-for-byte\n--- frame ---\n%s\n--- view ---\n%s", got, viewOut)
	}
}

// TestWatchFrameDeterministic: identical store state produces identical
// frames (no clock, no timestamps in the frame content).
func TestWatchFrameDeterministic(t *testing.T) {
	seedViewRepo(t, "valid")
	ws := openWorkspace(t)
	a, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	b, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Error("identical store state must produce identical frames")
	}
	// Frames are pure content: the clear-screen sequence belongs to the
	// loop, never to the frame bytes.
	if bytes.Contains(a, []byte("\x1b[2J")) {
		t.Error("frame bytes must not carry the clear-screen sequence")
	}
}

// TestWatchFrameRefusal: a repository not registered in the EKA
// workspace renders the calm refusal frame — the refusal message, the
// sync hint and the watching footer — with no projection content and
// no exit.
func TestWatchFrameRefusal(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	chdirInto(t, t.TempDir())
	ws := openWorkspace(t)
	frame, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatalf("renderWatchFrame: %v", err)
	}
	text := string(frame)
	for _, want := range []string{
		"Repository",
		"↓ View",
		"Repository not registered in the EKA workspace.",
		"eka sync",
		"watching — run 'eka sync' to register the repository",
		"watching — Ctrl-C to stop (interval 2s)",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("refusal frame must contain %q:\n%s", want, text)
		}
	}
	// No projection content: the frame is the refusal state, not a
	// degraded projection.
	for _, forbidden := range []string{"Planned", "No active container.", "Summary:"} {
		if strings.Contains(text, forbidden) {
			t.Errorf("refusal frame must not contain projection content %q:\n%s", forbidden, text)
		}
	}
}

// TestWatchFrameRefusalRecovery: the same helper that renders the
// refusal frame renders the projection again once the repository is
// registered and synced — the frame flips back without any exit.
func TestWatchFrameRefusalRecovery(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	dir := t.TempDir()
	writeStory(t, filepath.Join(dir, "docs", "operating", "work-items", "stories", "sto-one.md"), "eka-watch", "one", "todo")
	chdirInto(t, dir)
	ws := openWorkspace(t)
	refusal, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(refusal), "not registered") {
		t.Fatalf("unregistered repository must render the refusal frame:\n%s", refusal)
	}
	// Register + seed the repository: the next cycle must flip back to
	// the projection.
	if _, err := sync.Run(ws, dir, sync.Options{Pull: true, Push: true}); err != nil {
		t.Fatal(err)
	}
	recovered, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(recovered), "not registered") {
		t.Errorf("registered repository must render the projection frame:\n%s", recovered)
	}
	if bytes.Equal(refusal, recovered) {
		t.Error("registration must change the frame bytes")
	}
}

// --- frame-change detection --------------------------------------------

// TestWatchFrameChanged: the change helper — identical bytes mean no
// redraw, different bytes mean redraw.
func TestWatchFrameChanged(t *testing.T) {
	if frameChanged([]byte("same"), []byte("same")) {
		t.Error("identical frames must not be marked changed (no redraw)")
	}
	if !frameChanged([]byte("a"), []byte("b")) {
		t.Error("different frames must be marked changed (redraw)")
	}
	if !frameChanged(nil, []byte("first")) {
		t.Error("the first frame must always be drawn")
	}
}

// TestWatchFrameChangeDetectionOnRepoState: flipping a work item's
// execution state (with its change-log entry) and re-syncing the
// repository changes the frame bytes — the loop would redraw on the
// next interval.
func TestWatchFrameChangeDetectionOnRepoState(t *testing.T) {
	t.Setenv("EKA_HOME", t.TempDir())
	repo := copyFixture(t, viewFixtureAbs(t, "valid"))
	ws, err := workspace.Ensure()
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()
	if _, err := sync.Run(ws, repo, sync.Options{Pull: true, Push: true}); err != nil {
		t.Fatal(err)
	}
	chdirInto(t, repo)
	before, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	// delta is in-review: move it to done, keeping the change-log
	// consistent so the repository stays conformant. The new entry is
	// inserted inside the front matter, before the closing marker.
	delta := filepath.Join(repo, "docs", "operating", "work-items", "bugs", "bug-delta.md")
	orig, err := os.ReadFile(delta)
	if err != nil {
		t.Fatal(err)
	}
	modified := strings.Replace(string(orig), "execution-state: in-review", "execution-state: done", 1)
	modified = strings.Replace(modified,
		"    to: in-review\n    by: Engineering Architecture\n---",
		"    to: in-review\n    by: Engineering Architecture\n  - date: 2026-08-05\n    domain: execution-state\n    from: in-review\n    to: done\n    by: Engineering Architecture\n---", 1)
	if err := os.WriteFile(delta, []byte(modified), 0o644); err != nil {
		t.Fatal(err)
	}
	// Re-seed from the docs tree: the snapshot is unchanged, so a
	// snapshot-mode pull would skip the work (idempotent digest).
	if _, err := sync.Run(ws, repo, sync.Options{Pull: true, FromDocs: true, Push: true}); err != nil {
		t.Fatal(err)
	}
	after, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(before, after) {
		t.Error("flipping a work item state must change the frame bytes")
	}
	if !frameChanged(before, after) {
		t.Error("frameChanged must report the state flip as a redraw")
	}
	// The re-synced repository must still render a projection, not the
	// refusal frame (the edited repository stays conformant: the
	// change-log rule holds).
	if strings.Contains(string(after), "not registered") {
		t.Errorf("edited repository must stay registered:\n%s", after)
	}
}

// --- clear-screen emission ---------------------------------------------

// TestWatchClearScreenEmission: the clear+home sequence is emitted
// exactly by writeClearScreen — the helper used on the open path
// (before the first frame) and on the SIGINT exit path.
func TestWatchClearScreenEmission(t *testing.T) {
	var buf bytes.Buffer
	s := ui.NewStyle(&buf, false)
	writeClearScreen(s)
	if got := buf.String(); got != "\x1b[2J\x1b[H" {
		t.Errorf("clear sequence mismatch: got %q, want %q", got, "\x1b[2J\x1b[H")
	}
	// One emission per call: the open path emits once, the exit path
	// emits once.
	buf.Reset()
	writeClearScreen(s)
	writeClearScreen(s)
	if got := strings.Count(buf.String(), "\x1b[2J\x1b[H"); got != 2 {
		t.Errorf("expected one sequence per emission, got %d", got)
	}
	// The open path: the clear sequence leads the first frame.
	seedViewRepo(t, "valid")
	ws := openWorkspace(t)
	frame, err := renderWatchFrame(watchFrameStyle(), ws, "execution", "", 2)
	if err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	writeClearScreen(s)
	buf.Write(frame)
	if !strings.HasPrefix(buf.String(), "\x1b[2J\x1b[H") {
		t.Error("the open path must clear the screen before the first frame")
	}
	if strings.Count(buf.String(), "\x1b[2J\x1b[H") != 1 {
		t.Error("the open path must emit the clear sequence exactly once")
	}
}
