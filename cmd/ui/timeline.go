package ui

import (
	"fmt"
	"strings"
)

// Timeline is the vertical activity/milestone primitive: rows prefixed
// with a │ connector, each row carrying a marker glyph (▸ / • / ●)
// and text colored as one span, plus optional milestone separator
// lines (─────). Rendering is a pure function of the added rows, so
// output is deterministic on TTY and non-TTY alike.
//
// The primitive knows nothing about domain data: it renders marker and
// text exactly as given. The connector is plain in every environment;
// colors apply to the "marker text" span only (dim separators on a
// TTY), and non-TTY output is plain UTF-8.
type Timeline struct {
	s    *Style
	rows []timelineRow
}

// timelineRow is one row of the timeline.
type timelineRow struct {
	marker string
	text   string
	color  func(string) string
	sep    bool
}

// NewTimeline starts a vertical timeline for the given style.
func NewTimeline(s *Style) *Timeline { return &Timeline{s: s} }

// Add appends one row: marker glyph (▸ / • / ●), text, and the color
// function applied to "marker text" as one span (nil = plain).
func (t *Timeline) Add(marker, text string, color func(string) string) *Timeline {
	t.rows = append(t.rows, timelineRow{marker: marker, text: text, color: color})
	return t
}

// Separator appends a milestone separator line ("│ ────...").
func (t *Timeline) Separator() *Timeline {
	t.rows = append(t.rows, timelineRow{sep: true})
	return t
}

// Render prints the timeline rows. A timeline without rows renders
// nothing.
func (t *Timeline) Render() {
	s := t.s
	for _, r := range t.rows {
		if r.sep {
			fmt.Fprintln(s.W, s.Dim("│ "+strings.Repeat("─", 24)))
			continue
		}
		line := r.marker + " " + r.text
		if r.color != nil {
			line = r.color(line)
		}
		fmt.Fprintln(s.W, "│ "+line)
	}
}
