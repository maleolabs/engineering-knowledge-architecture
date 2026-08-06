package cmd

import (
	"fmt"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
)

// This file implements the Ticket projection renderer: the ticket
// detail card. The projected status is the first focal point (colored,
// with its state icon); the ticket card then carries the work item,
// the container and the derives-from references as supporting rows.

// renderTicket renders one ticket projection as a detail card.
func renderTicket(s *ui.Style, g *view.Graph, p *view.TicketProjection) {
	ui.NewHeader(s, "Ticket").
		Add("Ticket", p.Ticket.Identity).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v"+standardVersion).
		Add("Domain", "Execution").
		Pipeline("View").
		Render()

	// The projected status leads: icon + state word as one colored
	// span, the first thing the eye lands on.
	fmt.Fprintf(s.W, "\n%s  %s\n", s.Accent("Projected Status"),
		stateColor(s, p.Projected)(stateIcon(p.Projected)+" "+p.Projected))

	workItem := "unresolved"
	if p.WorkItem != nil {
		workItem = p.WorkItem.Identity + " (" + p.WorkItem.State + ")"
	}
	container := "unresolved"
	if p.Container != nil {
		container = p.Container.Identity
	}
	derives := "—"
	if len(p.References) > 0 {
		derives = strings.Join(p.References, ", ")
	}
	rows := [][2]string{
		{"Work Item", workItem},
		{"Container", container},
		{"Derives From", derives},
	}
	width := 0
	for _, r := range rows {
		if len(r[0]) > width {
			width = len(r[0])
		}
	}
	body := make([]string, 0, len(rows))
	for _, r := range rows {
		body = append(body, fmt.Sprintf("%-*s   %s", width, r[0], r[1]))
	}
	ui.NewCards(s).
		Add(p.Ticket.Identity, stateColor(s, p.Projected), body).
		Render()

	workItemValue := "unresolved"
	if p.WorkItem != nil {
		workItemValue = p.WorkItem.Identity + " (" + p.WorkItem.State + ")"
	}
	ui.NewSummary(s).
		Add("Projected status", p.Projected).
		Add("Work item", workItemValue).
		Render()
}
