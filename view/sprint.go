package view

// This file implements the sprint projection: the active execution
// container's work items, grouped by execution state.
//
// Derivation (see package doc for the membership rule):
//   - Target: the ACTIVE Execution Container (container-state "active";
//     exactly-one-active per protocol). No active container yields a
//     projection with zero work items — an empty projection is a valid
//     state, not an error. Several active containers (invalid state)
//     resolve to the lexicographically smallest canonical identity,
//     noted on the model.
//   - Members: the work items referenced by the tickets whose
//     derives-from includes the active container's identity line
//     (ticket -> work item). Deduplicated by identity line.
//   - Grouping: fixed execution-state column order planned, todo,
//     in-progress, in-review, done.
//
// The sprint projection ignores the optional target argument.

// SprintProjection is the sprint view over the active container.
type SprintProjection struct {
	// Container is the active container, or nil when no container is
	// active.
	Container *Container
	// MultipleActive reports the invalid state where several containers
	// are active; Container then holds the lexicographically smallest
	// identity.
	MultipleActive bool
	// Columns are the fixed execution-state columns (always the full
	// five-column set, with zero-item columns when empty).
	Columns StateColumns
	// Total is the number of work items placed in the columns.
	Total int
	// Tickets is the number of tickets deriving from the active
	// container.
	Tickets int
}

// Name returns the registry name of the projection.
func (p *SprintProjection) Name() string { return "sprint" }

func buildSprint(g *Graph, target string) (Projection, error) {
	container, multiple := g.ActiveContainer()
	p := &SprintProjection{MultipleActive: multiple}
	if container == nil {
		// Empty projection: no active container. The columns keep the
		// fixed order so the projection stays shape-stable.
		p.Columns = groupByState(nil)
		return p, nil
	}
	p.Container = container
	items := g.WorkItemsForContainer(container.Identity)
	p.Columns = groupByState(items)
	for _, col := range p.Columns {
		p.Total += len(col.WorkItems)
	}
	p.Tickets = len(g.TicketsForContainer(container.Identity))
	return p, nil
}
