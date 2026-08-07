package ui

import (
	"bytes"
	"strings"
	"testing"
)

// testStyle builds a plain (non-TTY, no color) style over a buffer.
func testStyle() (*Style, *bytes.Buffer) {
	var buf bytes.Buffer
	return NewStyle(&buf, false), &buf
}

func TestBoardLayout(t *testing.T) {
	s, buf := testStyle()
	NewBoard(s).
		AddColumn("Todo", nil, []string{"alpha", "beta"}).
		AddColumn("Done", nil, []string{"gamma"}).
		Render()
	want := "" +
		"┌──────────┬──────────┐\n" +
		"│ Todo (2) │ Done (1) │\n" +
		"├──────────┼──────────┤\n" +
		"│ ▸ alpha  │ ▸ gamma  │\n" +
		"│ ▸ beta   │          │\n" +
		"└──────────┴──────────┘\n"
	if got := buf.String(); got != want {
		t.Errorf("board layout mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestBoardEmptyColumnShowsDash(t *testing.T) {
	s, buf := testStyle()
	NewBoard(s).
		AddColumn("Planned", nil, []string{"alpha"}).
		AddColumn("Todo", nil, nil).
		Render()
	out := buf.String()
	if !strings.Contains(out, "│ ▸ alpha") {
		t.Errorf("filled column must show its item:\n%s", out)
	}
	if !strings.Contains(out, "—") {
		t.Errorf("empty column must show its \"—\" marker:\n%s", out)
	}
	// The "—" appears exactly once (the empty column's first row).
	if got := strings.Count(out, "—"); got != 1 {
		t.Errorf("empty marker count = %d, want 1:\n%s", got, out)
	}
}

func TestBoardAllColumnsEmpty(t *testing.T) {
	s, buf := testStyle()
	NewBoard(s).
		AddColumn("Planned", nil, nil).
		AddColumn("Done", nil, nil).
		Render()
	out := buf.String()
	if got := strings.Count(out, "—"); got != 2 {
		t.Errorf("every empty column must show \"—\", got %d:\n%s", got, out)
	}
}

func TestBoardMinWidth(t *testing.T) {
	s, buf := testStyle()
	NewBoard(s).AddColumn("A", nil, nil).Render()
	// Min width 8: "│ " + 8 cells + " │".
	if !strings.Contains(buf.String(), "│ A (0)    │") {
		t.Errorf("column must not narrow below the min width:\n%s", buf.String())
	}
}

func TestBoardTruncationCappedAtMaxWidth(t *testing.T) {
	s, buf := testStyle()
	long := strings.Repeat("x", 40)
	NewBoard(s).AddColumn("C", nil, []string{long}).Render()
	out := buf.String()
	if strings.Contains(out, long) {
		t.Errorf("long item must be truncated:\n%s", out)
	}
	if !strings.Contains(out, "…") {
		t.Errorf("truncated item must carry the ellipsis:\n%s", out)
	}
	// Display width of the truncated cell: "▸ " + 21 x + ellipsis = 24.
	if !strings.Contains(out, "▸ xxxxxxxxxxxxxxxxxxxxx…") {
		t.Errorf("truncation must cap at 24 display cells:\n%s", out)
	}
}

func TestBoardNoColumnsRendersNothing(t *testing.T) {
	s, buf := testStyle()
	NewBoard(s).Render()
	if buf.Len() != 0 {
		t.Errorf("a board without columns must render nothing, got %q", buf.String())
	}
}

func TestBoardDeterministic(t *testing.T) {
	build := func() string {
		var buf bytes.Buffer
		s := NewStyle(&buf, false)
		NewBoard(s).AddColumn("T", nil, []string{"a"}).Render()
		return buf.String()
	}
	if build() != build() {
		t.Error("board output must be deterministic")
	}
}

func TestBoardColorsOnTTY(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: true, W: &buf}
	NewBoard(s).AddColumn("T", s.Success, []string{"a"}).Render()
	out := buf.String()
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("color-enabled board must emit ANSI:\n%q", out)
	}
	if !strings.Contains(out, "\x1b[38;5;245m") {
		t.Errorf("borders must be dim:\n%q", out)
	}
}
