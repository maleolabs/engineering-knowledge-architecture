package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Board is the kanban presentation primitive: a fixed grid of box-drawn
// columns (┌─┬─┐ / ├─┼─┤ / └─┴─┘ / │) rendered side by side with
// aligned rows. Each column owns its width — max(header with count,
// longest card line, 8), capped at the adaptive column width, overflow
// truncated with "…" — and renders its item cards in order; an empty
// column shows "—" in its first item row.
//
// An item is a Card: several display lines, each with its own color
// function. A card renders as a vertical block (the "▸ " marker on its
// first line, the following lines indented to the text column), and a
// row's height is the tallest card in that row — shorter cards and
// empty columns pad with blank lines. Consecutive rows are separated by
// a blank gap row, so cards stay visually distinct.
//
// Rendering is a pure function of the added columns, so output is
// deterministic on TTY and non-TTY alike.
//
// The primitive knows nothing about domain data: it renders titles and
// item labels exactly as given. Colors are applied to text spans only
// (borders are dim on a TTY); non-TTY output is plain UTF-8.
type Board struct {
	s       *Style
	columns []boardColumn
}

// CardLine is one display line of a card: its text plus its own color
// function (nil = the column's default color).
type CardLine struct {
	Text string
	// Color colors this line only; nil falls back to the column color.
	Color func(string) string
}

// Card is a multi-line item label. The first line carries the "▸ "
// marker; following lines render indented to the text column.
type Card []CardLine

// boardColumn is one column of the board.
type boardColumn struct {
	title string
	color func(string) string
	items []Card
}

// Board width bounds: columns never narrow below boardMinWidth and
// never exceed boardMaxWidth (longer content is truncated with "…").
// On a TTY the columns adapt to the terminal width (BoardColumnWidth);
// with an unknown width (pipes, CI, tests) every column renders at the
// fixed maximum, keeping non-TTY output byte-identical.
const (
	boardMinWidth = 8
	boardMaxWidth = 32
)

// BoardColumnWidth returns the column width the board renders for a
// terminal of the given display width and column count. A termWidth of
// 0 (unknown — pipes, CI, tests) falls back to the fixed maximum. The
// result is clamped to [boardMinWidth, boardMaxWidth]: very narrow
// terminals keep the minimum width (the grid overflows — unavoidable
// without wrapping), wide terminals never exceed the cap.
func BoardColumnWidth(termWidth, columns int) int {
	if termWidth <= 0 || columns <= 0 {
		return boardMaxWidth
	}
	// Grid overhead: one border cell per column boundary (columns+1)
	// plus the two-cell padding of every column.
	overhead := (columns + 1) + columns*2
	w := (termWidth - overhead) / columns
	if w < boardMinWidth {
		w = boardMinWidth
	}
	if w > boardMaxWidth {
		w = boardMaxWidth
	}
	return w
}

// BoardItemBudget returns the display budget a single item label may
// use inside one board column: BoardColumnWidth minus the "▸ " prefix.
// Renderers that compose item labels (e.g. a label plus a context tag
// that must stay visible) truncate against this budget so the board
// primitive never re-truncates their tag away.
func BoardItemBudget(termWidth, columns int) int {
	return BoardColumnWidth(termWidth, columns) - displayWidth(itemPrefix)
}

// itemPrefix marks each item label on the board so list membership is
// readable at a glance; it is part of the cell content and counts
// toward the column width.
const itemPrefix = "▸ " // two display cells

// NewBoard starts a kanban board for the given style.
func NewBoard(s *Style) *Board { return &Board{s: s} }

// AddColumn appends one column of single-line items: title, the color
// function applied to the header text and item labels (nil = plain),
// and the item labels in display order. Multi-line labels ("a\nb")
// split into card lines; use AddCards for per-line colors.
func (b *Board) AddColumn(title string, color func(string) string, items []string) *Board {
	cards := make([]Card, len(items))
	for i, it := range items {
		lines := strings.Split(it, "\n")
		card := make(Card, len(lines))
		for j, line := range lines {
			card[j] = CardLine{Text: line}
		}
		cards[i] = card
	}
	return b.AddCards(title, color, cards)
}

// AddCards appends one column of item cards: title, the color function
// applied to the header text and to card lines without their own color
// (nil = plain), and the cards in display order.
func (b *Board) AddCards(title string, color func(string) string, cards []Card) *Board {
	b.columns = append(b.columns, boardColumn{title: title, color: color, items: cards})
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
	colWidth := BoardColumnWidth(s.Width, len(b.columns))
	widths := make([]int, len(b.columns))
	headers := make([]string, len(b.columns))
	for i, c := range b.columns {
		headers[i] = fmt.Sprintf("%s (%d)", c.title, len(c.items))
		w := len(headers[i])
		for _, card := range c.items {
			for _, line := range card {
				// The first line carries the "▸ " prefix inside its cell.
				if displayWidth(line.Text)+displayWidth(itemPrefix) > w {
					w = displayWidth(line.Text) + displayWidth(itemPrefix)
				}
			}
		}
		if w < boardMinWidth {
			w = boardMinWidth
		}
		if w > colWidth {
			w = colWidth
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
	// shows its "—" markers. Each row is as tall as its tallest card;
	// shorter cards pad with blank lines, and a blank gap row separates
	// consecutive rows so cards stay visually distinct.
	rows := 1
	for _, c := range b.columns {
		if len(c.items) > rows {
			rows = len(c.items)
		}
	}
	for r := 0; r < rows; r++ {
		height := 1
		for _, c := range b.columns {
			if r < len(c.items) {
				if n := len(c.items[r]); n > height {
					height = n
				}
			}
		}
		for h := 0; h < height; h++ {
			fmt.Fprintln(s.W, boardRow(s, widths, func(i int) (string, func(string) string) {
				c := b.columns[i]
				if r >= len(c.items) {
					// An empty column signals its emptiness with "—" in
					// the first item row; later rows stay blank.
					if len(c.items) == 0 && r == 0 {
						return "—", c.color
					}
					return "", nil
				}
				card := c.items[r]
				if h >= len(card) {
					// A shorter card pads with a blank line.
					return "", nil
				}
				line := card[h]
				lineColor := line.Color
				if lineColor == nil {
					lineColor = c.color
				}
				// The "▸ " marker belongs to the card's first line; the
				// following lines indent to the text column.
				if h == 0 {
					return itemPrefix + line.Text, lineColor
				}
				return strings.Repeat(" ", displayWidth(itemPrefix)) + line.Text, lineColor
			}))
		}
		if r < rows-1 {
			// Row gap: a blank separator row across all columns.
			fmt.Fprintln(s.W, boardRow(s, widths, func(int) (string, func(string) string) {
				return "", nil
			}))
		}
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
