// Package view implements the EKA Knowledge Projection Engine: the
// application-layer views over the Engineering Knowledge Model that the
// `eka view` command renders.
//
// The engine is pure data in, pure data out. It contains no terminal
// knowledge, no command framework, and no output: it reads one
// conformance.Scan result and produces projection models. Every
// projection is derived exclusively from the Engineering Knowledge
// Model (artifact identity, state fields, relationships) — never from
// markdown content. The engine is synchronous and stateless: one graph
// build per scan, one projection build per view, so a loading state can
// wrap the whole call without restructuring.
//
// Determinism contract: all ordering is canonical. Artifacts are
// ordered by their canonical line identity form (<namespace>/<type>:<id>)
// and instance-version; execution-state columns always follow the fixed
// value order planned, todo, in-progress, in-review, done; tickets keep
// relationship file order in their reference lists. There are no maps in
// output ordering and no time-dependent values.
//
// Membership derivation (the single source of membership — relationships
// only, never file text):
//
//	A work item is a member of an execution container C iff a ticket
//	(tkt-) whose derives-from resolves to C's identity line also
//	resolves to the work item's identity line. A ticket belongs to C iff
//	one of its derives-from references resolves to C. Work items are
//	identified by their type owning the Execution State domain
//	(conformance.OwnedDomains); the ticket itself is never parsed beyond
//	its frontmatter relationships.
package view

import (
	"sort"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

// LineForm renders the canonical line identity form of an artifact:
// "<namespace>/<type>:<id>". It identifies an artifact line (all
// instance-versions) and is the ordering key used throughout the
// package.
func LineForm(ns, typeToken, id string) string {
	return ns + "/" + typeToken + ":" + id
}

// lineKey is the internal index key for an identity line. The \x00
// separators cannot collide with the identity components (namespaces,
// type tokens and ids may contain hyphens but not colons, slashes or
// NULs).
func lineKey(ns, typeToken, id string) string {
	return ns + "\x00" + typeToken + "\x00" + id
}

// Graph is the Knowledge Graph built from one conformance.Scan result:
// an identity index (canonical line form -> artifact), instance lines
// (grouped by identity, ordered by instance-version), relationship
// resolution (reference -> artifact), and the membership helpers used by
// the projections. It is immutable after construction.
type Graph struct {
	root      string
	artifacts []conformance.Artifact
	// byLine indexes artifacts by identity line; each bucket holds every
	// instance sorted by instance-version (line resolution returns the
	// lowest instance, matching the validator).
	byLine map[string][]*conformance.Artifact
	// byForm maps the canonical line identity form to the line's lowest
	// instance (the line-level resolution target).
	byForm map[string]*conformance.Artifact
	// byType maps type tokens to their artifact lines, sorted by
	// canonical identity form (then instance-version, for robustness on
	// non-conformant input).
	byType map[string][]*conformance.Artifact
}

// NewGraph builds the Knowledge Graph from a scan result. The artifact
// slice is copied; the graph never aliases caller-owned memory.
func NewGraph(root string, artifacts []conformance.Artifact) *Graph {
	g := &Graph{
		root:      root,
		artifacts: artifacts,
		byLine:    make(map[string][]*conformance.Artifact, len(artifacts)),
		byForm:    make(map[string]*conformance.Artifact, len(artifacts)),
		byType:    make(map[string][]*conformance.Artifact),
	}
	for i := range g.artifacts {
		a := &g.artifacts[i]
		key := lineKey(a.Namespace, a.Type, a.ID)
		g.byLine[key] = append(g.byLine[key], a)
		g.byType[a.Type] = append(g.byType[a.Type], a)
	}
	for _, bucket := range g.byLine {
		sort.Slice(bucket, func(i, j int) bool {
			return bucket[i].InstanceVersion < bucket[j].InstanceVersion
		})
		g.byForm[LineForm(bucket[0].Namespace, bucket[0].Type, bucket[0].ID)] = bucket[0]
	}
	for _, bucket := range g.byType {
		sort.Slice(bucket, func(i, j int) bool {
			a, b := LineForm(bucket[i].Namespace, bucket[i].Type, bucket[i].ID),
				LineForm(bucket[j].Namespace, bucket[j].Type, bucket[j].ID)
			if a != b {
				return a < b
			}
			return bucket[i].InstanceVersion < bucket[j].InstanceVersion
		})
	}
	return g
}

// Root returns the scan root the graph was built from.
func (g *Graph) Root() string { return g.root }

// Artifacts returns every artifact of the graph, sorted by canonical
// identity form (then instance-version).
func (g *Graph) Artifacts() []conformance.Artifact {
	out := make([]conformance.Artifact, len(g.artifacts))
	copy(out, g.artifacts)
	sort.Slice(out, func(i, j int) bool {
		a, b := LineForm(out[i].Namespace, out[i].Type, out[i].ID),
			LineForm(out[j].Namespace, out[j].Type, out[j].ID)
		if a != b {
			return a < b
		}
		return out[i].InstanceVersion < out[j].InstanceVersion
	})
	return out
}

// ByLineForm resolves a canonical line identity form to the line's
// lowest instance, or nil when the line is not in the graph.
func (g *Graph) ByLineForm(form string) *conformance.Artifact {
	return g.byForm[form]
}

// Resolve resolves a parsed reference to its target artifact: a
// versioned reference resolves to the exact instance, a line reference
// to the lowest instance. It returns nil when the reference does not
// resolve. Semantics match the validator's Rule 5 resolution.
func (g *Graph) Resolve(ref conformance.Reference) *conformance.Artifact {
	bucket := g.byLine[lineKey(ref.Namespace, ref.Type, ref.ID)]
	if len(bucket) == 0 {
		return nil
	}
	if ref.HasVersion {
		for _, a := range bucket {
			if a.InstanceVersion == ref.Version {
				return a
			}
		}
		return nil
	}
	return bucket[0]
}

// Container is one execution container (ctr-) line in the graph.
type Container struct {
	// Identity is the canonical line identity form.
	Identity string
	Type     string
	ID       string
	// State is the container-state value ("active" or "completed").
	State string
}

func containerFor(a *conformance.Artifact) *Container {
	return &Container{
		Identity: LineForm(a.Namespace, a.Type, a.ID),
		Type:     a.Type,
		ID:       a.ID,
		State:    a.States[conformance.DomainContainerState],
	}
}

// ActiveContainer returns the single active execution container per the
// exactly-one-active protocol. It returns nil when no container is
// active (a valid, empty projection state). When the repository is in an
// invalid state — several containers with container-state "active" — it
// returns the lexicographically smallest canonical identity and reports
// the anomaly through multiple.
func (g *Graph) ActiveContainer() (*Container, bool) {
	var active []*conformance.Artifact
	for _, a := range g.byType["ctr"] {
		if a.States[conformance.DomainContainerState] == "active" {
			active = append(active, a)
		}
	}
	if len(active) == 0 {
		return nil, false
	}
	// The byType bucket is sorted by canonical identity: the first
	// element is the lexicographically smallest.
	return containerFor(active[0]), len(active) > 1
}

// Containers returns every execution container line, sorted by
// canonical identity form.
func (g *Graph) Containers() []Container {
	out := make([]Container, 0, len(g.byType["ctr"]))
	for _, a := range g.byType["ctr"] {
		out = append(out, *containerFor(a))
	}
	return out
}

// Ticket is one ticket (tkt-) line in the graph, with the status
// projected from its referenced work item.
type Ticket struct {
	Identity string
	Type     string
	ID       string
	// Projected is the execution-state of the ticket's referenced work
	// item, or "unresolved" when the ticket has no resolvable work item
	// reference. It is derived from the owner artifact's state — the
	// ticket's own body content is never read.
	Projected string
}

// WorkItem is one work item (sto-/ts-/bug-/td-/ch-/spk-) line in the
// graph.
type WorkItem struct {
	Identity string
	Type     string
	ID       string
	// State is the execution-state value.
	State string
	// Dimension is the informational dimension field, when present.
	Dimension    string
	HasDimension bool
}

// isWorkItemType reports whether a type token owns the Execution State
// domain — the six work item types, characterized without duplicating
// the type table.
func isWorkItemType(token string) bool {
	for _, d := range conformance.OwnedDomains(token) {
		if d == conformance.DomainExecutionState {
			return true
		}
	}
	return false
}

func workItemFor(a *conformance.Artifact) *WorkItem {
	return &WorkItem{
		Identity:     LineForm(a.Namespace, a.Type, a.ID),
		Type:         a.Type,
		ID:           a.ID,
		State:        a.States[conformance.DomainExecutionState],
		Dimension:    a.Dimension,
		HasDimension: a.HasDimension,
	}
}

// ticketTargets resolves the container and the work item a ticket
// derives from, scanning the ticket's derives-from references in file
// order: the first resolvable ctr- reference is the container, the
// first resolvable execution-state-owning reference is the work item.
func (g *Graph) ticketTargets(t *conformance.Artifact) (*Container, *WorkItem) {
	if t == nil {
		return nil, nil
	}
	var container *Container
	var workItem *WorkItem
	for _, raw := range t.Relations["derives-from"] {
		ref, err := conformance.ParseReference(raw, t.Namespace, t.Type)
		if err != nil {
			continue // malformed references are reported by Rule 5.
		}
		target := g.Resolve(ref)
		if target == nil {
			continue
		}
		switch {
		case container == nil && target.Type == "ctr":
			container = containerFor(target)
		case workItem == nil && isWorkItemType(target.Type):
			workItem = workItemFor(target)
		}
	}
	return container, workItem
}

// ticketsFor returns the ticket artifacts deriving from the container
// identity form, sorted by canonical identity.
func (g *Graph) ticketsFor(form string) []*conformance.Artifact {
	var out []*conformance.Artifact
	for _, a := range g.byType["tkt"] {
		container, _ := g.ticketTargets(a)
		if container != nil && container.Identity == form {
			out = append(out, a)
		}
	}
	return out
}

// TicketsForContainer returns the tickets deriving from the container
// identity form (canonical line identity "<ns>/<type>:<id>"), sorted by
// canonical identity. Membership is relationship-only: a ticket belongs
// to the container iff one of its derives-from references resolves to
// the container's identity line.
func (g *Graph) TicketsForContainer(form string) []Ticket {
	out := make([]Ticket, 0, len(g.byType["tkt"]))
	for _, a := range g.ticketsFor(form) {
		_, workItem := g.ticketTargets(a)
		projected := "unresolved"
		if workItem != nil {
			projected = workItem.State
		}
		out = append(out, Ticket{
			Identity:  LineForm(a.Namespace, a.Type, a.ID),
			Type:      a.Type,
			ID:        a.ID,
			Projected: projected,
		})
	}
	return out
}

// WorkItemsForContainer returns the work items of the container's
// tickets — the members of the container — deduplicated by identity
// line and sorted by canonical identity. A work item is a member iff a
// ticket deriving from the container also derives from the work item;
// tickets without a resolvable work item contribute nothing.
func (g *Graph) WorkItemsForContainer(form string) []WorkItem {
	seen := make(map[string]bool)
	var out []WorkItem
	for _, a := range g.ticketsFor(form) {
		_, workItem := g.ticketTargets(a)
		if workItem == nil || seen[workItem.Identity] {
			continue
		}
		seen[workItem.Identity] = true
		out = append(out, *workItem)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Identity < out[j].Identity })
	return out
}

// ContainerForTicket resolves the container a ticket derives from, or
// nil when the ticket has no resolvable ctr- reference.
func (g *Graph) ContainerForTicket(t *conformance.Artifact) *Container {
	container, _ := g.ticketTargets(t)
	return container
}

// WorkItemForTicket resolves the work item a ticket derives from, or
// nil when the ticket has no resolvable work item reference.
func (g *Graph) WorkItemForTicket(t *conformance.Artifact) *WorkItem {
	_, workItem := g.ticketTargets(t)
	return workItem
}

// TicketByTarget resolves a ticket from a user-supplied target string:
// a bare ticket id, "tkt-<id>" or "tkt:<id>" (the prefix is stripped).
// It returns nil when no ticket matches. When several namespaces hold a
// ticket with the same id, the lexicographically smallest canonical
// identity wins.
func (g *Graph) TicketByTarget(target string) *conformance.Artifact {
	id := strings.TrimPrefix(strings.TrimPrefix(target, "tkt-"), "tkt:")
	for _, a := range g.byType["tkt"] {
		if a.ID == id {
			return a
		}
	}
	return nil
}

// TicketIDs returns the ids of every ticket line, sorted by canonical
// identity — the available targets of the ticket projection.
func (g *Graph) TicketIDs() []string {
	out := make([]string, 0, len(g.byType["tkt"]))
	for _, a := range g.byType["tkt"] {
		out = append(out, a.ID)
	}
	return out
}

// StateColumn is one execution-state column of a projection: the fixed
// value plus its work items, sorted by canonical identity.
type StateColumn struct {
	State     string
	WorkItems []WorkItem
}

// StateColumns is the ordered execution-state column set. The order is
// the fixed value order planned, todo, in-progress, in-review, done —
// never map iteration.
type StateColumns []StateColumn

// Count returns the number of work items in the column for state, or 0
// when the column is absent.
func (c StateColumns) Count(state string) int {
	for _, col := range c {
		if col.State == state {
			return len(col.WorkItems)
		}
	}
	return 0
}

// groupByState groups work items into the fixed execution-state column
// order. Work items whose state is not a valid execution-state value
// (impossible behind the validation gate) are not placed in any column.
func groupByState(items []WorkItem) StateColumns {
	order := conformance.DomainValues(conformance.DomainExecutionState, "sto")
	cols := make(StateColumns, 0, len(order))
	for _, state := range order {
		col := StateColumn{State: state}
		for _, wi := range items {
			if wi.State == state {
				col.WorkItems = append(col.WorkItems, wi)
			}
		}
		cols = append(cols, col)
	}
	return cols
}
