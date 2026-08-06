package ui

import (
	"fmt"
	"strings"
)

// Cards renders a block of compact boxed info cards (┌─┐ / │ / └─┘),
// one box per item, stacked directly. Each card has a header line (id
// plus state icon, colored by the item's state color) and up to a few
// supporting metadata lines (dim). Rendering is a pure function of the
// added cards, so output is deterministic on TTY and non-TTY alike.
//
// The primitive knows nothing about domain data: it renders header and
// body lines exactly as given. Colors are applied to text spans only
// (borders and body lines are dim on a TTY); non-TTY output is plain
// UTF-8.
type Cards struct {
	s     *Style
	cards []card
}

// card is one boxed info card.
type card struct {
	header string
	color  func(string) string
	body   []string
}

// NewCards starts an info-card block for the given style.
func NewCards(s *Style) *Cards { return &Cards{s: s} }

// Add appends one card: header text, the color function applied to the
// header (nil = plain), and the body lines (rendered dim, plain when
// colors are off).
func (c *Cards) Add(header string, color func(string) string, body []string) *Cards {
	c.cards = append(c.cards, card{header: header, color: color, body: body})
	return c
}

// Render prints every card, one box per item. All cards of the block
// share the widest content line, so stacked boxes align. A block
// without cards renders nothing.
func (c *Cards) Render() {
	if len(c.cards) == 0 {
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
