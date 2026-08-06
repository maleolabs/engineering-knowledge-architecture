package ui

// The EKA CLI icon set: a small, deliberately minimal set of Unicode
// glyphs (no emojis). All bytes are valid UTF-8. Icons are decoration;
// text always carries the meaning.
const (
	// IconDone marks a completed step or a finished spinner.
	IconDone = "✓"
	// IconBullet prefixes items in detail lists.
	IconBullet = "•"
	// IconArrow marks a relationship direction; used sparingly.
	IconArrow = "→"
	// IconDown marks the pipeline separator in the context header.
	IconDown = "↓"

	// Tree connectors for deterministic non-TTY tree lines.
	TreeBranch = "├──"
	TreeLast   = "└──"
	TreeVert   = "│"
	TreeSpace  = "   "
)

// SpinnerFrames is the Braille spinner animation cycle (TTY only).
// Non-TTY output never contains a frame.
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
