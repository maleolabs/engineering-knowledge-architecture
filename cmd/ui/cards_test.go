package ui

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCardsLayout(t *testing.T) {
	s, buf := testStyle()
	NewCards(s).
		Add("✓ feather/vis:feather-vision", nil, []string{"approved · revision 1"}).
		Add("○ feather/req:comments", nil, []string{"draft", "note"}).
		Render()
	out := buf.String()

	// Header 1 is the widest line (29 cells incl. ✓ + space); the box
	// width is the max display width + 2 pads + 2 bars.
	width := 28 // "feather/vis:feather-vision" is the widest content (28 cells)
	want := fmt.Sprintf(
		"┌%s┐\n│ %s%s │\n│ %s%s │\n└%s┘\n"+
			"┌%s┐\n│ %s%s │\n│ %s%s │\n│ %s%s │\n└%s┘\n",
		strings.Repeat("─", width+2),
		"✓ feather/vis:feather-vision", "",
		"approved · revision 1", strings.Repeat(" ", width-len([]rune("approved · revision 1"))),
		strings.Repeat("─", width+2),
		strings.Repeat("─", width+2),
		"○ feather/req:comments", strings.Repeat(" ", width-len([]rune("○ feather/req:comments"))),
		"draft", strings.Repeat(" ", width-len([]rune("draft"))),
		"note", strings.Repeat(" ", width-len([]rune("note"))),
		strings.Repeat("─", width+2),
	)
	if got := buf.String(); got != want {
		t.Errorf("cards layout mismatch:\n got: %q\nwant: %q", got, want)
	}

	// Every content row spans the same display width (28 + 2 + 2).
	rows := 0
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if strings.HasPrefix(line, "│") {
			rows++
			if w := len([]rune(line)); w != 32 {
				t.Errorf("content row %q spans %d cells, want 32", line, w)
			}
		}
	}
	if rows != 5 {
		t.Errorf("content rows = %d, want 5 (2 headers + 3 bodies)", rows)
	}
}

func TestCardsEmptyBody(t *testing.T) {
	s, buf := testStyle()
	NewCards(s).Add("header", nil, nil).Render()
	out := buf.String()
	if !strings.Contains(out, "│ header │") {
		t.Errorf("header-only card must render:\n%s", out)
	}
}

func TestCardsNoCardsRendersNothing(t *testing.T) {
	s, buf := testStyle()
	NewCards(s).Render()
	if buf.Len() != 0 {
		t.Errorf("a cards block without cards must render nothing, got %q", buf.String())
	}
}

func TestCardsHeaderColored(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: true, W: &buf}
	NewCards(s).Add("header", s.Success, nil).Render()
	if !strings.Contains(buf.String(), "\x1b[38;5;114mheader\x1b[0m") {
		t.Errorf("header must be colored by its color function:\n%q", buf.String())
	}
}
