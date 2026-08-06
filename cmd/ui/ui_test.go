package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// bufStyle returns a non-TTY, non-color Style writing into buf.
func bufStyle(buf *bytes.Buffer) *Style {
	return &Style{Color: false, TTY: false, Verbose: false, W: buf}
}

func TestIsTTYFalseForBuffers(t *testing.T) {
	if IsTTY(&bytes.Buffer{}) {
		t.Error("bytes.Buffer must not be a TTY")
	}
	// A regular file is not a terminal.
	f, err := os.CreateTemp(t.TempDir(), "ui")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if IsTTY(f) {
		t.Error("a regular file must not be a TTY")
	}
}

func TestColorOffIsPlain(t *testing.T) {
	s := bufStyle(&bytes.Buffer{})
	for _, got := range []string{
		s.Info("x"), s.Success("x"), s.Warning("x"), s.Error("x"),
		s.Progress("x"), s.Dim("x"), s.Accent("x"),
	} {
		if got != "x" {
			t.Errorf("color-disabled render = %q, want plain %q", got, "x")
		}
	}
}

func TestColorOnWrapsSGR(t *testing.T) {
	s := &Style{Color: true, TTY: true, W: &bytes.Buffer{}}
	want := func(code, text string) string { return "\x1b[" + code + "m" + text + "\x1b[0m" }
	cases := []struct {
		got  string
		code string
		text string
	}{
		{s.Info("i"), ColorInfo, "i"},
		{s.Success("s"), ColorSuccess, "s"},
		{s.Warning("w"), ColorWarning, "w"},
		{s.Error("e"), ColorError, "e"},
		{s.Progress("p"), ColorProgress, "p"},
		{s.Dim("d"), ColorDim, "d"},
		{s.Accent("a"), ColorAccent, "a"},
	}
	for _, c := range cases {
		if c.got != want(c.code, c.text) {
			t.Errorf("render = %q, want %q", c.got, want(c.code, c.text))
		}
	}
	// Accent must be the Info hue (headings use Info).
	if ColorAccent != ColorInfo {
		t.Error("ColorAccent must equal ColorInfo")
	}
}

func TestIconsAreSmallAndUTF8(t *testing.T) {
	icons := []string{IconDone, IconBullet, IconArrow, IconDown, TreeBranch, TreeLast, TreeVert, TreeSpace}
	for _, ic := range icons {
		if strings.ToValidUTF8(ic, "\uFFFD") != ic {
			t.Errorf("icon %q must be valid UTF-8", ic)
		}
	}
	if len(SpinnerFrames) == 0 {
		t.Error("spinner frames must not be empty")
	}
}

func TestStepPrefix(t *testing.T) {
	if got := Step(1, 5); got != "[1/5] " {
		t.Errorf("Step(1,5) = %q, want %q", got, "[1/5] ")
	}
	if got := Step(7, 7); got != "[7/7] " {
		t.Errorf("Step(7,7) = %q, want %q", got, "[7/7] ")
	}
}

// TestTreeNonTTYSequence verifies the deterministic sequential emission
// contract: pending/running nodes emit nothing, completed nodes emit
// "├── label" plus "│   ✓ detail", and re-rendering emits nothing new.
func TestTreeNonTTYSequence(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	tree := NewTree(s, "Repository export")

	n1 := tree.Add("[1/2] Discover repository")
	n2 := tree.Add("[2/2] Load Engineering Knowledge")

	tree.Render() // nothing completed: no output
	if buf.String() != "Repository export\n" {
		t.Fatalf("early render must emit only the root, got %q", buf.String())
	}

	n1.Running()
	tree.Render() // running emits nothing on non-TTY
	if buf.String() != "Repository export\n" {
		t.Fatalf("running node must not emit on non-TTY, got %q", buf.String())
	}

	n1.Done("scanned 6 artifacts")
	tree.Render()
	want := "Repository export\n" +
		"├── [1/2] Discover repository\n" +
		"│   ✓ scanned 6 artifacts\n"
	if buf.String() != want {
		t.Fatalf("after first done:\ngot  %q\nwant %q", buf.String(), want)
	}

	n2.Done("loaded 6 units")
	tree.Render()
	want += "├── [2/2] Load Engineering Knowledge\n" +
		"│   ✓ loaded 6 units\n"
	if buf.String() != want {
		t.Fatalf("after second done:\ngot  %q\nwant %q", buf.String(), want)
	}

	tree.Render() // idempotent: nothing new
	if buf.String() != want {
		t.Fatalf("re-render must be a no-op, got %q", buf.String())
	}
}

func TestTreeNonTTYFail(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	tree := NewTree(s, "Repository initialization")
	n := tree.Add("[1/1] Validate")
	n.Fail("FAIL (2 errors, 1 warning)")
	tree.Finish()
	want := "Repository initialization\n" +
		"├── [1/1] Validate\n" +
		"│   failed: FAIL (2 errors, 1 warning)\n"
	if buf.String() != want {
		t.Fatalf("fail render:\ngot  %q\nwant %q", buf.String(), want)
	}
}

func TestTreeTTYRedraw(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: true, TTY: true, W: &buf}
	tree := NewTree(s, "Root")
	n := tree.Add("Step one")
	n.Running()
	tree.Render()
	if !strings.Contains(buf.String(), "\x1b[38;5;80m⠋ Step one\x1b[0m") {
		t.Fatalf("running node must show the spinner frame, got %q", buf.String())
	}
	n.Done("done")
	tree.Render()
	if !strings.Contains(buf.String(), "\x1b[1A") {
		t.Fatalf("redraw must move the cursor up, got %q", buf.String())
	}
	if !strings.Contains(buf.String(), "\x1b[38;5;114m✓ Step one\x1b[0m") {
		t.Fatalf("done node must show the check mark, got %q", buf.String())
	}
}

func TestTreeTTYNoANSIWhenColorOff(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: false, TTY: true, W: &buf}
	tree := NewTree(s, "Root")
	tree.Add("Step").Done("ok")
	tree.Finish()
	if strings.Contains(buf.String(), "\x1b") {
		t.Errorf("TTY without color must not emit ANSI escapes, got %q", buf.String())
	}
}

func TestSpinnerNonTTY(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	sp := NewSpinner(s, "Loading Engineering Knowledge...")
	if buf.String() != "Loading Engineering Knowledge...\n" {
		t.Fatalf("non-TTY spinner must print the message once, got %q", buf.String())
	}
	sp.Stop()
	if buf.String() != "Loading Engineering Knowledge...\n" {
		t.Fatalf("Stop must print nothing extra on non-TTY, got %q", buf.String())
	}
}

func TestSpinnerTTYFinalState(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: true, TTY: true, W: &buf}
	sp := NewSpinner(s, "Working...")
	sp.Stop()
	out := buf.String()
	if !strings.HasPrefix(out, "\r\x1b[K") {
		t.Errorf("spinner frames must start with carriage return, got %q", out)
	}
	if !strings.HasSuffix(out, "\x1b[38;5;114m✓\x1b[0m \x1b[38;5;80mWorking...\x1b[0m\n") {
		t.Errorf("final state must be a deterministic check line, got %q", out)
	}
	if strings.Contains(out, "⠋") && !strings.Contains(out, "\x1b[38;5;80m⠋") {
		t.Errorf("frames must be color-wrapped, got %q", out)
	}
}

func TestSummaryNonTTY(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	NewSummary(s).
		Add("Artifacts", "6").
		Add("Status", "Repository conforms to EKA v1").
		Render()
	want := "\nSummary:\n" +
		"└── Artifacts: 6\n" +
		"└── Status: Repository conforms to EKA v1\n"
	if buf.String() != want {
		t.Fatalf("summary render:\ngot  %q\nwant %q", buf.String(), want)
	}
}

func TestBulletsEmptyRendersNothing(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	s.Bullets("Units:", nil)
	if buf.Len() != 0 {
		t.Errorf("empty bullets must render nothing, got %q", buf.String())
	}
	s.Bullets("Units:", []string{"a", "b"})
	want := "\nUnits:\n  • a\n  • b\n"
	if buf.String() != want {
		t.Fatalf("bullets render:\ngot  %q\nwant %q", buf.String(), want)
	}
}

// TestNoTimeDependence verifies that every deterministic renderer
// produces byte-identical output across repeated calls.
func TestNoTimeDependence(t *testing.T) {
	render := func() string {
		var buf bytes.Buffer
		s := bufStyle(&buf)
		tree := NewTree(s, "Repository initialization")
		tree.Add("[1/2] Plan").Done("5 actions planned")
		tree.Add("[2/2] Validate").Done("PASS (0 errors, 0 warnings)")
		tree.Finish()
		NewSummary(s).Add("Status", "ok").Render()
		return buf.String()
	}
	if a, b := render(), render(); a != b {
		t.Error("deterministic renderers must be byte-identical across runs")
	}
}

// TestHeaderNonTTYLayout verifies the deterministic non-TTY header
// bytes: object kind line, aligned label column, pipeline separator.
func TestHeaderNonTTYLayout(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	NewHeader(s, "Repository").
		Add("Name", "myproj").
		Add("Namespace", "eka-cli").
		Add("Knowledge", "EKA v1").
		Pipeline("Bootstrap").
		Render()
	want := "Repository\n" +
		"Name        myproj\n" +
		"Namespace   eka-cli\n" +
		"Knowledge   EKA v1\n" +
		"↓ Bootstrap\n"
	if buf.String() != want {
		t.Fatalf("header render:\ngot  %q\nwant %q", buf.String(), want)
	}
}

func TestHeaderNoPipeline(t *testing.T) {
	var buf bytes.Buffer
	s := bufStyle(&buf)
	NewHeader(s, "Object").
		Add("Only", "row").
		Render()
	// The label column is always aligned: max label width + 3.
	want := "Object\nOnly   row\n"
	if buf.String() != want {
		t.Fatalf("header without pipeline:\ngot  %q\nwant %q", buf.String(), want)
	}
}

func TestHeaderTTYColorAccentOnly(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: true, TTY: true, W: &buf}
	NewHeader(s, "Repository").
		Add("Name", "myproj").
		Pipeline("Bootstrap").
		Render()
	// Only the object kind line is colored; rows and pipeline stay
	// plain so the identity values remain readable.
	want := "\x1b[38;5;75mRepository\x1b[0m\n" +
		"Name   myproj\n" +
		"↓ Bootstrap\n"
	if buf.String() != want {
		t.Fatalf("header TTY render:\ngot  %q\nwant %q", buf.String(), want)
	}
}

func TestHeaderNoANSIWhenColorOff(t *testing.T) {
	for _, tty := range []bool{false, true} {
		var buf bytes.Buffer
		s := &Style{Color: false, TTY: tty, W: &buf}
		NewHeader(s, "Repository").
			Add("Name", "myproj").
			Pipeline("Bootstrap").
			Render()
		if strings.Contains(buf.String(), "\x1b") {
			t.Errorf("TTY=%v color off: header must not emit ANSI, got %q", tty, buf.String())
		}
	}
}

func TestSpinnerTTYNoColorIsStatic(t *testing.T) {
	var buf bytes.Buffer
	s := &Style{Color: false, TTY: true, W: &buf}
	sp := NewSpinner(s, "Working...")
	// TTY without color behaves exactly like non-TTY: one deterministic
	// line, no animation, no "\r", no erase codes.
	if buf.String() != "Working...\n" {
		t.Fatalf("spinner start must print the message once, got %q", buf.String())
	}
	sp.Stop()
	if buf.String() != "Working...\n" {
		t.Fatalf("Stop must be a no-op without color, got %q", buf.String())
	}
	if strings.Contains(buf.String(), "\x1b") || strings.Contains(buf.String(), "\r") {
		t.Errorf("no-color spinner must not emit control sequences, got %q", buf.String())
	}
}
