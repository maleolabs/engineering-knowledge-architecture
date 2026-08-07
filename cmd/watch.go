package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
	"github.com/spf13/cobra"
)

// This file implements the `eka watch` command: the live version of
// `eka view`. The projection surface, the conformance gate and the
// frame renderers are exactly the view ones — a watch frame of a
// conformant repository is the one-shot view output plus the watching
// footer, byte for byte. Watch adds only the presentation loop: a
// polling interval, change detection (identical frames are not
// redrawn), a clear-screen on open and on exit, and SIGINT handling.
//
// Validation failures are rendered as a calm in-place frame, not an
// exit: watch keeps polling and flips back to the projection
// automatically once the repository recovers. There is no fsnotify and
// no new dependency: the repository is re-read on a timer.
//
// Watch is a live display: stdout must be a terminal. The terminal
// gate runs after argument validation, so usage errors stay
// deterministic (and testable) on any writer; the non-TTY error keeps
// piped/CI behavior explicit.

// flagWatchInterval is the polling interval flag of `eka watch`
// (seconds).
const flagWatchInterval = "interval"

// clearScreen is the ANSI clear-screen + cursor-home sequence: wipe the
// terminal and put the cursor at the top-left so every frame starts on
// a clean surface. Emitted before the first frame and after the loop
// (SIGINT), so the shell prompt restarts clean. TTY-only by
// construction: watch never runs without a terminal.
const clearScreen = "\x1b[2J\x1b[H"

// newWatchCommand builds the `eka watch` command.
//
// Exit codes:
//
//	0  clean stop (Ctrl-C)
//	2  usage or internal error (unknown projection, missing ticket
//	   target, invalid interval, non-terminal stdout, unreadable root)
func newWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch <projection> [target]",
		Short: "Watch a projection live, redrawn on change",
		Long: `Live projection of the Engineering Knowledge Model of the repository
rooted at the current directory: 'eka watch' is 'eka view' in a loop —
the projection is re-rendered in place on a polling interval and
redrawn only when it changed (no flicker on stable state).

On every cycle the repository is re-validated against the conformance
rules (R0-R12). While blocking violations exist the watch shows a calm
validation-failure frame instead of the projection; it keeps polling
and flips back to the projection automatically once the repository
recovers. Warnings never block a projection.

Projections (the same surface as 'eka view'):

  discovery    the Discovery domain: vis-, str-, req-, fnd- artifacts
               grouped by type with their content states
  architecture the Architecture domain: adr-, dec-, arc-, spec-, std-,
               gls- artifacts grouped by type with their content states
  planning     the Planning domain: scp-, epc-, plan-, trc- artifacts
               grouped by type with content state, planning state and
               phase context
  execution    the active execution container: its tickets with the
               status projected from their work items, and its work
               items grouped by execution state
  operations   the Operations domain: run-, rel- artifacts grouped by
               type with their content states
  ticket       one ticket's projected status, derived from the
               referenced work item's execution state

Aliases:

  sprint, wave resolve to the execution projection (identical output)

The target argument is required by the ticket projection only.

watch is a live display: stdout must be a terminal. For one-shot
output use 'eka view'.

Exit codes:
  0  clean stop (Ctrl-C)
  2  usage or internal error (unknown projection, missing ticket
     target, invalid interval, non-terminal stdout)`,
		Example: `  eka watch execution
  eka watch execution --interval 5
  eka watch ticket tkt-sto-alpha`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("watch requires a projection — available projections: %s",
					view.HelpList())
			}
			name, target, err := parseProjectionArgs(args)
			if err != nil {
				return err
			}
			interval, err := cmd.Flags().GetInt(flagWatchInterval)
			if err != nil {
				return err
			}
			if interval < 1 {
				return fmt.Errorf("invalid interval %d: must be at least 1 second", interval)
			}
			// Terminal gate: watch is a live display. The deterministic
			// non-TTY error keeps piped/CI behavior explicit.
			if !ui.IsTTY(cmd.OutOrStdout()) {
				return fmt.Errorf("watch requires a terminal (stdout is not a TTY); use 'eka view' for one-shot output")
			}
			return runWatch(styleFor(cmd), name, target, interval)
		},
	}
	cmd.Flags().Int(flagWatchInterval, 2, "polling interval in seconds (minimum 1)")
	return cmd
}

// runWatch runs the live loop: clear the screen, render the first
// frame immediately (no initial delay), then re-render every interval
// seconds, redrawing the terminal only when the frame changed. SIGINT
// clears the screen and stops cleanly (exit 0); a single handler is
// enough — the loop is short and the cleanup is one write. The signal
// handler is installed only on the TTY path (this function is only
// reached past watch's terminal gate).
func runWatch(s *ui.Style, projection, target string, interval int) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	writeClearScreen(s)
	prev, err := renderWatchFrame(s, ".", projection, target, interval)
	if err != nil {
		return err
	}
	io.Copy(s.W, bytes.NewReader(prev))
	flush(s.W)

	render := func() error {
		frame, err := renderWatchFrame(s, ".", projection, target, interval)
		if err != nil {
			return err
		}
		if !frameChanged(prev, frame) {
			return nil
		}
		writeClearScreen(s)
		io.Copy(s.W, bytes.NewReader(frame))
		flush(s.W)
		prev = frame
		return nil
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-interrupt:
			writeClearScreen(s)
			return nil
		case <-ticker.C:
			if err := render(); err != nil {
				return err
			}
		}
	}
}

// renderWatchFrame runs one watch cycle into a fresh buffer: the
// conformance gate (the same gate as view), then either the projection
// frame — byte-identical to the one-shot view output plus the watching
// footer — or the validation-failure frame. The frame is a pure
// function of (repository state, projection, target, interval): no
// clock, no timestamps, so identical states produce identical frames
// and the loop skips redraws by byte comparison.
//
// The base style is copied and its writer replaced with the buffer, so
// the frame carries exactly the color/TTY settings of the live stdout
// (a verified TTY for watch).
func renderWatchFrame(s *ui.Style, root, projection, target string, interval int) ([]byte, error) {
	frame := *s
	var buf bytes.Buffer
	frame.W = &buf

	report, err := conformance.Validate(root)
	if err != nil {
		return nil, fmt.Errorf("watch failed: %w", err)
	}
	if !report.Pass() {
		renderWatchFailure(&frame, report, interval)
		return buf.Bytes(), nil
	}
	artifacts, err := conformance.Scan(root)
	if err != nil {
		return nil, fmt.Errorf("watch failed: %w", err)
	}
	g := view.NewGraph(root, artifacts)
	proj, err := view.Build(projection, g, target)
	if err != nil {
		return nil, fmt.Errorf("watch failed: %w", err)
	}
	renderView(&frame, g, proj)
	renderWatchFooter(&frame, interval)
	return buf.Bytes(), nil
}

// renderWatchFailure renders the validation-failure frame: the live
// error state of the watch. It reuses the existing validation pieces —
// the context header, the report summary and the sorted findings — in
// a calm red-tinted frame, and ends with a dim hint that the watch
// keeps polling (this is a rendered state, not an exit).
func renderWatchFailure(s *ui.Style, r *conformance.Report, interval int) {
	ui.NewHeader(s, "Repository").
		Add("Path", r.Root).
		Add("Knowledge", "EKA v"+standardVersion).
		Pipeline("Validate").
		Render()
	fmt.Fprintln(s.W)
	fmt.Fprintf(s.W, "%s\n", s.Error("Repository validation failed."))
	fmt.Fprintf(s.W, "%s\n", validationDetail(r))
	renderFindings(s, r)
	fmt.Fprintln(s.W)
	fmt.Fprintf(s.W, "%s\n", s.Dim("watching — fix the repository to resume the projection"))
	renderWatchFooter(s, interval)
}

// renderWatchFooter prints the one dim footer line shared by both
// frame kinds.
func renderWatchFooter(s *ui.Style, interval int) {
	fmt.Fprintf(s.W, "%s\n", s.Dim(fmt.Sprintf("watching — Ctrl-C to stop (interval %ds)", interval)))
}

// frameChanged reports whether the freshly rendered frame differs from
// the previous one. Identical frames are not redrawn: no clear-screen,
// no write, no flicker on stable state.
func frameChanged(prev, next []byte) bool {
	return !bytes.Equal(prev, next)
}

// writeClearScreen emits the clear-screen sequence and flushes. It is
// the single exit/entry point for terminal wiping (open and SIGINT
// paths both go through it).
func writeClearScreen(s *ui.Style) {
	fmt.Fprint(s.W, clearScreen)
	flush(s.W)
}

// flush pushes buffered output to the destination. os.File writes are
// unbuffered (no-op); buffered writers (*bufio.Writer) need the flush.
func flush(w io.Writer) {
	if bw, ok := w.(*bufio.Writer); ok {
		bw.Flush()
	}
}
