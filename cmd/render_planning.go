package cmd

import (
	"strconv"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
)

// This file implements the Planning projection renderer: the roadmap
// answering "what are we planning next?" — the plan instance as a
// milestone line, the scope definition and the epics as timeline rows,
// traceability as a footer note, and the insight summary.

// renderPlanning renders the Planning projection as a roadmap timeline.
func renderPlanning(s *ui.Style, g *view.Graph, p *view.PlanningProjection) {
	renderDomainHeader(s, g, "Planning")
	if view.GroupTotal(p.Groups) == 0 {
		renderDomainEmpty(s, "Planning")
		renderPlanningInsights(s, p)
		return
	}
	plans := groupByName(p.Groups, "Plans")
	scopes := groupByName(p.Groups, "Scope Definitions")
	epics := groupByName(p.Groups, "Epics")
	traces := groupByName(p.Groups, "Traceability")

	tl := ui.NewTimeline(s)
	if len(plans) == 0 {
		// No committed plan: the roadmap has no milestone yet.
		tl.Add("○", "no plan yet — roadmap undefined", s.Dim)
	} else {
		for _, plan := range plans {
			// Milestone line: the plan instance with its
			// planning-state icon, planning state and phase.
			tl.Add(planningStateIcon(plan.PlanningState), plan.Identity+artifactStateText(plan),
				planningStateColor(s, plan.PlanningState))
		}
	}
	tl.Separator()
	for _, scp := range scopes {
		tl.Add("▸", scp.Identity+artifactStateText(scp), contentStateColor(s, scp.ContentState))
	}
	for _, epc := range epics {
		tl.Add("▸", epc.Identity+artifactStateText(epc), contentStateColor(s, epc.ContentState))
	}
	for _, tr := range traces {
		// Traceability is a Planning artifact like any other: it gets a
		// normal timeline row (connected to the structure), with the
		// "traceability:" label for context.
		tl.Add("▸", "traceability: "+tr.Identity+" ("+tr.ContentState+")",
			contentStateColor(s, tr.ContentState))
	}
	tl.Render()
	renderPlanningInsights(s, p)
}

// renderPlanningInsights renders the planning summary: committed
// (approved plans), exploring (draft epics), and the next milestone
// (the phase of the first approved plan, "—" when none).
func renderPlanningInsights(s *ui.Style, p *view.PlanningProjection) {
	approved := 0
	for _, sc := range p.PlansByState {
		if sc.State == "approved" {
			approved = sc.Count
		}
	}
	draftEpics := 0
	for _, a := range groupByName(p.Groups, "Epics") {
		if a.ContentState == "draft" {
			draftEpics++
		}
	}
	next := "—"
	for _, a := range groupByName(p.Groups, "Plans") {
		if a.PlanningState != "approved" {
			continue
		}
		if a.HasPhase {
			next = a.Phase
		}
		break
	}
	ui.NewSummary(s).
		Add("Committed", strconv.Itoa(approved)).
		Add("Exploring", strconv.Itoa(draftEpics)).
		Add("Next milestone", next).
		Render()
}
