package view

import "github.com/maleolabs/engineering-knowledge-architecture/conformance"

// This file implements the planning projection: the Planning domain
// (scp-, epc-, plan-, trc-) grouped by type in fixed order. Each entry
// carries its content-state; scope definitions and plans also carry
// their phase context, and plans their planning-state. The summary
// counts artifacts per group and plans per planning-state.
//
// The planning projection ignores the optional target argument.

// PlanningProjection is the Planning domain view.
type PlanningProjection struct {
	// Groups are the fixed artifact groups in order: Scope Definitions,
	// Epics, Plans, Traceability.
	Groups []Group
	// PlansByState counts plans per planning-state in the fixed value
	// order draft, approved, immutable.
	PlansByState []StateCount
}

// Name returns the registry name of the projection.
func (p *PlanningProjection) Name() string { return "planning" }

func buildPlanning(g *Graph, target string) (Projection, error) {
	groups := domainGroups(g, conformance.Planning, []groupDef{
		{[]string{"scp"}, "Scope Definitions"},
		{[]string{"epc"}, "Epics"},
		{[]string{"plan"}, "Plans"},
		{[]string{"trc"}, "Traceability"},
	})
	p := &PlanningProjection{Groups: groups}
	for _, state := range conformance.DomainValues(conformance.DomainPlanningState, "plan") {
		p.PlansByState = append(p.PlansByState, StateCount{State: state})
	}
	for _, group := range groups {
		for _, a := range group.Artifacts {
			if a.Type != "plan" {
				continue
			}
			for i := range p.PlansByState {
				if p.PlansByState[i].State == a.PlanningState {
					p.PlansByState[i].Count++
				}
			}
		}
	}
	return p, nil
}
