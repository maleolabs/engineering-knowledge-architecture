package view

// This file implements the wave projection: the active execution
// container, its tickets (tkt- deriving from it) and the work item
// progress counts per execution state.
//
// Derivation: identical to the sprint projection's membership rule —
// the active container, its deriving tickets, and the work items
// referenced by those tickets. Each ticket carries the status projected
// from its referenced work item ("unresolved" when the work item does
// not resolve).
//
// The wave projection ignores the optional target argument.

// WaveProjection is the wave view over the active container.
type WaveProjection struct {
	// Container is the active container, or nil when no container is
	// active.
	Container *Container
	// MultipleActive reports the invalid state where several containers
	// are active; Container then holds the lexicographically smallest
	// identity.
	MultipleActive bool
	// Tickets are the tickets deriving from the active container,
	// sorted by canonical identity, each carrying its projected status.
	Tickets []Ticket
	// Columns are the work item progress counts per execution state
	// (fixed column order).
	Columns StateColumns
	// Total is the number of work items placed in the columns.
	Total int
}

// Name returns the registry name of the projection.
func (p *WaveProjection) Name() string { return "wave" }

func buildWave(g *Graph, target string) (Projection, error) {
	container, multiple := g.ActiveContainer()
	p := &WaveProjection{MultipleActive: multiple}
	if container == nil {
		p.Columns = groupByState(nil)
		return p, nil
	}
	p.Container = container
	p.Tickets = g.TicketsForContainer(container.Identity)
	items := g.WorkItemsForContainer(container.Identity)
	p.Columns = groupByState(items)
	for _, col := range p.Columns {
		p.Total += len(col.WorkItems)
	}
	return p, nil
}
