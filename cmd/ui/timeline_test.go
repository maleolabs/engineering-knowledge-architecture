package ui

import (
	"strings"
	"testing"
)

func TestTimelineRows(t *testing.T) {
	s, buf := testStyle()
	NewTimeline(s).
		Add("●", "milestone", s.Success).
		Separator().
		Add("▸", "activity", s.Info).
		Render()
	out := buf.String()
	for _, want := range []string{
		"│ ● milestone",
		"│ ▸ activity",
		"│ ──",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("timeline output must contain %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "───") {
		t.Errorf("separator must be a dashed milestone line:\n%s", out)
	}
}

func TestTimelineNoRowsRendersNothing(t *testing.T) {
	s, buf := testStyle()
	NewTimeline(s).Render()
	if buf.Len() != 0 {
		t.Errorf("a timeline without rows must render nothing, got %q", buf.String())
	}
}

func TestTimelineDeterministic(t *testing.T) {
	build := func() string {
		s, buf := testStyle()
		NewTimeline(s).Add("▸", "a", nil).Separator().Render()
		return buf.String()
	}
	if build() != build() {
		t.Error("timeline output must be deterministic")
	}
}
