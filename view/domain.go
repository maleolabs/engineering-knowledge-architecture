package view

import (
	"sort"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

// This file implements the shared model of the domain projections
// (discovery, architecture, planning, operations): unit lines
// grouped by type in a fixed, projection-specific order. The model is
// plain, deterministically ordered data; rendering conventions (state
// coloring) are the caller's concern.

// Group is one unit group of a domain projection: the type tokens
// it aggregates, its display name, and its unit lines sorted by
// canonical identity.
type Group struct {
	// Name is the group's display name ("Scope Definitions",
	// "Decisions", ...).
	Name string
	// Artifacts are the group's unit lines, sorted by canonical
	// identity (line-level: instance-versions collapse to the line).
	Artifacts []DomainArtifact
}

// DomainArtifact is one unit line of a domain projection: its
// canonical identity plus the state values relevant to its group —
// content-state on every knowledge unit; planning-state and phase
// on Planning units. Presence flags distinguish an absent field
// from an empty value (a CKO omits absent fields).
type DomainArtifact struct {
	Identity         string
	Type             string
	ID               string
	ContentState     string
	HasContentState  bool
	PlanningState    string
	HasPlanningState bool
	Phase            string
	HasPhase         bool
}

// StateCount is one (state, count) pair of a summary, in the fixed
// value order of its domain (e.g. plans by planning-state: draft,
// approved, immutable).
type StateCount struct {
	State string
	Count int
}

// groupDef defines one artifact group: the type tokens it aggregates
// and its display name.
type groupDef struct {
	tokens []string
	name   string
}

// domainGroups collects the unit lines of each group in the fixed
// definition order. Within a group the lines are sorted by canonical
// identity and collapsed to line level (the lowest instance, matching
// the validator's line resolution). Units whose home domain differs
// from the projection's domain are skipped (defensive: the validator
// already guarantees the mapping).
func domainGroups(g *Graph, domain conformance.Domain, defs []groupDef) []Group {
	groups := make([]Group, 0, len(defs))
	for _, def := range defs {
		group := Group{Name: def.name}
		seen := make(map[string]bool)
		for _, token := range def.tokens {
			for _, u := range g.byType[token] {
				home, ok := conformance.DomainForToken(u.Identity.Type)
				if !ok || home != domain {
					continue
				}
				identity := LineForm(u.Identity.Namespace, u.Identity.Type, u.Identity.ID)
				if seen[identity] {
					continue
				}
				seen[identity] = true
				group.Artifacts = append(group.Artifacts, DomainArtifact{
					Identity:         identity,
					Type:             u.Identity.Type,
					ID:               u.Identity.ID,
					ContentState:     u.StateVector.ContentState,
					HasContentState:  u.StateVector.ContentState != "",
					PlanningState:    u.StateVector.PlanningState,
					HasPlanningState: u.StateVector.PlanningState != "",
					Phase:            u.Phase,
					HasPhase:         u.Phase != "",
				})
			}
		}
		sort.Slice(group.Artifacts, func(i, j int) bool {
			return group.Artifacts[i].Identity < group.Artifacts[j].Identity
		})
		groups = append(groups, group)
	}
	return groups
}

// GroupTotal returns the number of artifact lines across all groups —
// the domain's artifact count.
func GroupTotal(groups []Group) int {
	total := 0
	for _, gr := range groups {
		total += len(gr.Artifacts)
	}
	return total
}
