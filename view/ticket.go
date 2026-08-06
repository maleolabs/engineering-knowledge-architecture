package view

import "fmt"

// This file implements the ticket projection: one ticket (tkt-) with
// its projected status.
//
// Derivation:
//   - Target: a ticket identity, required. The target accepts a bare
//     ticket id, "tkt-<id>" or "tkt:<id>".
//   - Projected status: DERIVED from the referenced work item's
//     execution-state — the projection semantics read the owner state,
//     never the ticket's own `## Projected Status` content. A ticket
//     without a resolvable work item is a valid projection with the
//     explicit status "unresolved".
//   - Container: the ticket's derives-from ctr- reference, when it
//     resolves.
//   - References: the ticket's derives-from reference strings in file
//     order (the relationship list).

// TicketProjection is the view over one ticket.
type TicketProjection struct {
	// Ticket is the projected ticket.
	Ticket Ticket
	// Container is the container the ticket derives from, or nil when
	// the ctr- reference does not resolve.
	Container *Container
	// WorkItem is the work item the ticket derives from, or nil when
	// the ticket has no resolvable work item reference.
	WorkItem *WorkItem
	// Projected is the ticket's projected status: the referenced work
	// item's execution-state, or "unresolved".
	Projected string
	// References are the derives-from reference strings in file order.
	References []string
}

// Name returns the registry name of the projection.
func (p *TicketProjection) Name() string { return "ticket" }

func buildTicket(g *Graph, target string) (Projection, error) {
	if target == "" {
		return nil, fmt.Errorf("the ticket projection requires a target: eka view ticket <tkt-id>")
	}
	ticket := g.TicketByTarget(target)
	if ticket == nil {
		return nil, &TargetNotFoundError{
			Projection: "ticket",
			Target:     target,
			Available:  g.TicketIDs(),
		}
	}
	container, workItem := g.ticketTargets(ticket)
	p := &TicketProjection{
		Ticket: Ticket{
			Identity: LineForm(ticket.Namespace, ticket.Type, ticket.ID),
			Type:     ticket.Type,
			ID:       ticket.ID,
		},
		Container:  container,
		References: ticket.Relations["derives-from"],
	}
	if workItem != nil {
		p.WorkItem = workItem
		p.Projected = workItem.State
	} else {
		p.Projected = "unresolved"
	}
	p.Ticket.Projected = p.Projected
	return p, nil
}
