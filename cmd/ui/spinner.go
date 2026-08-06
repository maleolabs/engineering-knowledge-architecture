package ui

import (
	"fmt"
	"sync"
	"time"
)

// spinnerInterval is the animation period for the TTY spinner.
const spinnerInterval = 100 * time.Millisecond

// Spinner renders a contextual loading state. Every spinner MUST carry
// a message ("Discovering repository...", "Loading Engineering
// Knowledge...") — a bare spinner is never acceptable.
//
// On a TTY with color enabled it prints the message with an animated
// Braille frame via "\r" on a private goroutine (stopped by Stop). On a
// non-TTY — and on a TTY with color disabled — it prints the message
// once, deterministically, and Stop prints nothing extra: no frames,
// no "\r", no erase codes in plain output.
type Spinner struct {
	s   *Style
	msg string

	stop chan struct{}
	wg   sync.WaitGroup
}

// animated reports whether the frame animation may run: the writer must
// be a terminal and colors must be enabled, so no control sequences
// ever reach plain output.
func (sp *Spinner) animated() bool {
	return sp.s.TTY && sp.s.Color
}

// NewSpinner starts the spinner. On a non-TTY (or TTY without color)
// the deterministic "message" line is printed immediately; on a TTY
// with color the animation starts and the first frame is drawn right
// away.
func NewSpinner(s *Style, message string) *Spinner {
	sp := &Spinner{s: s, msg: message, stop: make(chan struct{})}
	if sp.animated() {
		sp.wg.Add(1)
		go sp.animate()
	} else {
		fmt.Fprintln(s.W, message)
	}
	return sp
}

// animate cycles the spinner frames on the spinner's own line until
// stop is closed. The frame index starts at the first frame so the
// first render is immediate.
func (sp *Spinner) animate() {
	defer sp.wg.Done()
	frame := 0
	draw := func() {
		fmt.Fprintf(sp.s.W, "\r\x1b[K%s %s",
			sp.s.Progress(SpinnerFrames[frame%len(SpinnerFrames)]), sp.msg)
	}
	draw()
	ticker := time.NewTicker(spinnerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-sp.stop:
			return
		case <-ticker.C:
			frame++
			draw()
		}
	}
}

// Stop stops the animation and prints the deterministic final line
// ("✓ <message>") on a TTY with color. Otherwise it is a no-op: the
// start line was already deterministic. Stop must be called exactly
// once.
func (sp *Spinner) Stop() {
	if !sp.animated() {
		return
	}
	close(sp.stop)
	sp.wg.Wait()
	fmt.Fprintf(sp.s.W, "\r\x1b[K%s %s\n", sp.s.Success(IconDone), sp.s.Progress(sp.msg))
}
