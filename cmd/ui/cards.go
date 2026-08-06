package ui

import (
	"fmt"
	"strings"
)

// Cards renders a block of compact boxed info cards (┌─┐ / │ / └─┘).
// By default one box per item, stacked vertically. With Grid(), the
// cards render as a horizontal grid: cells of uniform width arranged
// row-major (left to right, top to bottom) inside shared row boxes
// (┌─┬─┐ / │a│b│ / └─┴─┘), so many cards stay compact instead of
// growing the view vertically.
//
// The primitive knows nothing about domain data: it renders header and
// body lines exactly as given (truncated to the cell width when the
// grid is active). Colors are applied to text spans only (borders and
// body lines are dim on a TTY); non-TTY output is plain UTF-8.
type Cards struct {
	s      *Style
	cards  []card
	grid   bool
	budget int
}

// card is one boxed info card.
type card struct {
	header string
	color  func(string) string
	body   []string
}

// Card grid bounds and the width budget. The budget approximates a
// typical terminal width; it is a fixed constant so grid layout is
// deterministic on every environment (TTY and non-TTY alike) — column
// count adapts to CONTENT, never to the terminal.
const (
	cardMinWidth = 24
	cardMaxWidth = 48
	cardBudget   = 100
	cardMaxCols  = 4
)

// NewCards starts an info-card block for the given style.
func NewCards(s *Style) *Cards { return &Cards{s: s} }

// Add appends one card: header text, the color function applied to the
// header (nil = plain), and the body lines (rendered dim, plain when
// colors are off).
func (c *Cards) Add(header string, color func(string) string, body []string) *Cards {
	c.cards = append(c.cards, card{header: header, color: color, body: body})
	return c
}

// Grid switches the block to horizontal grid rendering. The column
// count is derived from the content: cell width = the widest content
// line (clamped to [cardMinWidth, cardMaxWidth], truncated beyond),
// columns = clamp(cardBudget / cellWidth, 1, cardMaxCols). A block
// without cards renders nothing either way.
func (c *Cards) Grid() *Cards {
	c.grid = true
	return c
}

// Render prints every card. Without Grid(), one box per item, stacked
// vertically, all boxes sharing the widest content line. With Grid(),
// cards render row-major inside shared row boxes. A block without
// cards renders nothing.
func (c *Cards) Render() {
	if len(c.cards) == 0 {
		return
	}
	if c.grid {
		c.renderGrid()
		return
	}
	s := c.s
	width := 0
	for _, cd := range c.cards {
		if w := displayWidth(cd.header); w > width {
			width = w
		}
		for _, line := range cd.body {
			if w := displayWidth(line); w > width {
				width = w
			}
		}
	}
	for _, cd := range c.cards {
		fmt.Fprintln(s.W, s.Dim("┌"+strings.Repeat("─", width+2)+"┐"))
		header := cd.header
		if cd.color != nil {
			header = cd.color(header)
		}
		fmt.Fprintln(s.W, s.Dim("│"), header+strings.Repeat(" ", width-displayWidth(cd.header)), s.Dim("│"))
		for _, line := range cd.body {
			fmt.Fprintln(s.W, s.Dim("│"), s.Dim(line)+strings.Repeat(" ", width-displayWidth(line)), s.Dim("│"))
		}
		fmt.Fprintln(s.W, s.Dim("└"+strings.Repeat("─", width+2)+"┘"))
	}
}

// cellWidth returns the uniform grid cell width: the widest content
// line (header or body) clamped to [cardMinWidth, cardMaxWidth].
func (c *Cards) cellWidth() int {
	width := 0
	for _, cd := range c.cards {
		if w := displayWidth(cd.header); w > width {
			width = w
		}
		for _, line := range cd.body {
			if w := displayWidth(line); w > width {
				width = w
			}
		}
	}
	if width < cardMinWidth {
		width = cardMinWidth
	}
	if width > cardMaxWidth {
		width = cardMaxWidth
	}
	return width
}

// renderGrid prints the cards as horizontal grid rows. Each row is one
// box containing up to columns cells; cells are the card stacks
// (header line + body lines), padded and truncated to the uniform cell
// width. The last row pads with blank cells so the box stays
// rectangular.
func (c *Cards) renderGrid() {
	s := c.s
	width := c.cellWidth()
	cols := cardBudget / width
	if cols < 1 {
		cols = 1
	}
	if cols > cardMaxCols {
		cols = cardMaxCols
	}
	// Never pad a row with empty cells beyond the card count: a
	// two-card group renders two cells, not three.
	if cols > len(c.cards) {
		cols = len(c.cards)
	}
	height := 0
	for _, cd := range c.cards {
		if h := 1 + len(cd.body); h > height {
			height = h
		}
	}
	for start := 0; start < len(c.cards); start += cols {
		row := c.cards[start:min(start+cols, len(c.cards))]
		// Top border: ┌───┬───┐
		var top strings.Builder
		for i := 0; i < cols; i++ {
			if i == 0 {
				top.WriteString("┌")
			} else {
				top.WriteString("┬")
			}
			top.WriteString(strings.Repeat("─", width+2))
		}
		top.WriteString("┐")
		fmt.Fprintln(s.W, s.Dim(top.String()))

		// Cell lines: header (colored) then body lines (dim); cells
		// shorter than the row height pad with blank lines.
		for line := 0; line < height; line++ {
			var sb strings.Builder
			for i := 0; i < cols; i++ {
				sb.WriteString(s.paint(ColorDim, "│ "))
				var text string
				if i < len(row) {
					cd := row[i]
					switch {
					case line == 0:
						text = cd.header
					case line-1 < len(cd.body):
						text = cd.body[line-1]
					}
				}
				text = truncate(text, width)
				if i < len(row) && line == 0 && row[i].color != nil {
					text = row[i].color(text)
				} else {
					text = s.paint(ColorDim, text)
				}
				sb.WriteString(text + strings.Repeat(" ", width-displayWidth(text)))
				sb.WriteString(s.paint(ColorDim, " "))
			}
			sb.WriteString(s.paint(ColorDim, "│"))
			fmt.Fprintln(s.W, sb.String())
		}

		// Bottom border: └───┴───┘
		var bottom strings.Builder
		for i := 0; i < cols; i++ {
			if i == 0 {
				bottom.WriteString("└")
			} else {
				bottom.WriteString("┴")
			}
			bottom.WriteString(strings.Repeat("─", width+2))
		}
		bottom.WriteString("┘")
		fmt.Fprintln(s.W, s.Dim(bottom.String()))
	}
}
