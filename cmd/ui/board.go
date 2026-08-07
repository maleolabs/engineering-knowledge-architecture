package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Board is the kanban presentation primitive: a fixed grid of box-drawn
// columns (┌─┬─┐ / ├─┼─┤ / └─┴─┘ / │) rendered side by side with
// aligned rows. Each column owns its width — max(header with count,
// longest item, 8), capped at 16, overflow truncated with "…" — and
// renders its item labels in order; an empty column shows "—" in its
// first item row. Rendering is a pure function of the added columns,
// so output is deterministic on TTY and non-TTY alike.
//
// The primitive knows nothing about domain data: it renders titles and
// item labels exactly as given. Colors are applied to text spans only
// (borders are dim on a TTY); non-TTY output is plain UTF-8.
type Board struct {
	s       *Style
	columns []boardColumn
}

// boardColumn is one column of the board.
type boardColumn struct {
	title string
	color func(string) string
	items []string
}

// Board width bounds: columns never narrow below boardMinWidth and
// never exceed boardMaxWidth (longer content is truncated with "…").
const (
	boardMinWidth = 8
	boardMaxWidth = 24
)

// BoardMaxItemWidth is the display budget a single item label may use
// inside a board cell: the column width cap minus the "▸ " prefix
// (two display cells). Renderers that compose item labels (e.g. a
// label plus a context tag that must stay visible) truncate against
// this budget.
const BoardMaxItemWidth = boardMaxWidth - 2

// itemPrefix marks each item label on the board so list membership is
// readable at a glance; it is part of the cell content and counts
// toward the column width.
const itemPrefix = "▸ " // two display cells

// NewBoard starts a kanban board for the given style.
func NewBoard(s *Style) *Board { return &Board{s: s} }

// AddColumn appends one column: title, the color function applied to
// the header text and item labels (nil = plain), and the item labels
// in display order.
func (b *Board) AddColumn(title string, color func(string) string, items []string) *Board {
	b.columns = append(b.columns, boardColumn{title: title, color: color, items: items})
	return b
}

// Render prints the board: top border, header row (title + count),
// separator, item rows (one row per maximum column item count; empty
// cells stay blank), bottom border. A board without columns renders
// nothing.
func (b *Board) Render() {
	if len(b.columns) == 0 {
		return
	}
	s := b.s
	widths := make([]int, len(b.columns))
	headers := make([]string, len(b.columns))
	for i, c := range b.columns {
		headers[i] = fmt.Sprintf("%s (%d)", c.title, len(c.items))
		w := len(headers[i])
		for _, it := range c.items {
			// Each item carries the "▸ " prefix inside its cell.
			if displayWidth(it)+displayWidth(itemPrefix) > w {
				w = displayWidth(it) + displayWidth(itemPrefix)
			}
		}
		if w < boardMinWidth {
			w = boardMinWidth
		}
		if w > boardMaxWidth {
			w = boardMaxWidth
		}
		widths[i] = w
	}

	// Top border: ┌───┬───┐
	var top strings.Builder
	for i, w := range widths {
		if i == 0 {
			top.WriteString("┌")
		} else {
			top.WriteString("┬")
		}
		top.WriteString(strings.Repeat("─", w+2))
	}
	top.WriteString("┐")
	fmt.Fprintln(s.W, s.Dim(top.String()))

	// Header row.
	fmt.Fprintln(s.W, boardRow(s, widths, func(i int) (string, func(string) string) {
		return headers[i], b.columns[i].color
	}))

	// Separator: ├───┼───┤
	var mid strings.Builder
	for i, w := range widths {
		if i == 0 {
			mid.WriteString("├")
		} else {
			mid.WriteString("┼")
		}
		mid.WriteString(strings.Repeat("─", w+2))
	}
	mid.WriteString("┤")
	fmt.Fprintln(s.W, s.Dim(mid.String()))

	// Item rows. At least one row renders, so an all-empty board still
	// shows its "—" markers.
	rows := 1
	for _, c := range b.columns {
		if len(c.items) > rows {
			rows = len(c.items)
		}
	}
	for r := 0; r < rows; r++ {
		fmt.Fprintln(s.W, boardRow(s, widths, func(i int) (string, func(string) string) {
			c := b.columns[i]
			if r >= len(c.items) {
				// An empty column signals its emptiness with "—" in the
				// first item row; later rows stay blank.
				if len(c.items) == 0 && r == 0 {
					return "—", c.color
				}
				return "", nil
			}
			return itemPrefix + c.items[r], c.color
		}))
	}

	// Bottom border: └───┴───┘
	var bottom strings.Builder
	for i, w := range widths {
		if i == 0 {
			bottom.WriteString("└")
		} else {
			bottom.WriteString("┴")
		}
		bottom.WriteString(strings.Repeat("─", w+2))
	}
	bottom.WriteString("┘")
	fmt.Fprintln(s.W, s.Dim(bottom.String()))
}

// truncate shortens text to the display width, appending "…" when it
// does not fit. Titles and ids are ASCII by the reference grammar; item
// cells may carry the "▸ " prefix and the "…" ellipsis — the only
// multi-byte glyphs the board emits — so truncation operates on runes
// (display cells), never on bytes.
func truncate(text string, width int) string {
	if displayWidth(text) <= width {
		return text
	}
	runes := []rune(text)
	return string(runes[:width-1]) + "…"
}

// displayWidth returns the terminal display width of text. Every glyph
// the UI emits (box drawing, ✓, ○, ▸, •, ·, —, …) occupies exactly one
// cell, so the display width is the rune count; byte length would
// over-count the multi-byte glyphs.
func displayWidth(text string) int {
	return utf8.RuneCountInString(text)
}

// boardRow renders one grid row: "│ <cell> │ ... │" with each cell
// padded to its column width. The cell content (already truncated) is
// colored by its column's color function; the bars and padding are dim
// on a TTY. Padding is computed on the plain display width, so colored
// and truncated cells stay aligned.
func boardRow(s *Style, widths []int, cell func(int) (string, func(string) string)) string {
	var sb strings.Builder
	for i, w := range widths {
		text, color := cell(i)
		sb.WriteString(s.paint(ColorDim, "│ "))
		text = truncate(text, w)
		padded := text + strings.Repeat(" ", w-displayWidth(text))
		if color != nil {
			padded = color(padded)
		}
		sb.WriteString(padded)
		sb.WriteString(s.paint(ColorDim, " "))
	}
	sb.WriteString(s.paint(ColorDim, "│"))
	return sb.String()
}
