package view

import (
	"reflect"
	"testing"
)

// groupArtifactStates returns the (identity, content-state,
// planning-state, phase) tuples of a group, in order — the projection's
// observable per-entry state.
func groupArtifactStates(group Group) [][4]string {
	out := make([][4]string, 0, len(group.Artifacts))
	for _, a := range group.Artifacts {
		ps, ph := "", ""
		if a.HasPlanningState {
			ps = a.PlanningState
		}
		if a.HasPhase {
			ph = a.Phase
		}
		out = append(out, [4]string{a.Identity, a.ContentState, ps, ph})
	}
	return out
}

func TestPlanningProjection(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("planning", g, "")
	if err != nil {
		t.Fatalf("Build(planning): %v", err)
	}
	planning, ok := p.(*PlanningProjection)
	if !ok {
		t.Fatalf("Build(planning) = %T, want *PlanningProjection", p)
	}

	// Fixed group order.
	wantOrder := []string{"Scope Definitions", "Epics", "Plans", "Traceability"}
	gotOrder := make([]string, len(planning.Groups))
	for i, gr := range planning.Groups {
		gotOrder[i] = gr.Name
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Errorf("group order = %v, want %v", gotOrder, wantOrder)
	}

	// Per-group membership and state values: scp- carries its phase
	// context, plan- its planning-state and phase.
	wantGroups := map[string][][4]string{
		"Scope Definitions": {
			{validForm + "scp:wave-2", "approved", "", "mvp"},
		},
		"Epics": {
			{validForm + "epc:auth", "review", "", ""},
		},
		"Plans": {
			{validForm + "plan:roadmap-2026", "approved", "approved", "release"},
		},
		"Traceability": {
			{validForm + "trc:spec-trace", "draft", "", ""},
		},
	}
	for _, gr := range planning.Groups {
		want := wantGroups[gr.Name]
		if got := groupArtifactStates(gr); !reflect.DeepEqual(got, want) {
			t.Errorf("group %q = %+v, want %+v", gr.Name, got, want)
		}
	}

	// Plans by planning-state: fixed value order, plan:roadmap-2026 is
	// approved.
	wantByState := []StateCount{{"draft", 0}, {"approved", 1}, {"immutable", 0}}
	if !reflect.DeepEqual(planning.PlansByState, wantByState) {
		t.Errorf("PlansByState = %+v, want %+v", planning.PlansByState, wantByState)
	}
}

// TestPlanningProjectionEmptyDomain: a repository without Planning
// artifacts yields an empty projection — all groups empty, zero plan
// counts, still exit-0-shaped.
func TestPlanningProjectionEmptyDomain(t *testing.T) {
	g := loadFixture(t, "no-active")
	p, err := Build("planning", g, "")
	if err != nil {
		t.Fatalf("Build(planning): %v", err)
	}
	planning := p.(*PlanningProjection)
	if len(planning.Groups) != 4 {
		t.Fatalf("groups = %d, want the fixed four", len(planning.Groups))
	}
	for _, gr := range planning.Groups {
		if len(gr.Artifacts) != 0 {
			t.Errorf("group %q must be empty, got %v", gr.Name, gr.Artifacts)
		}
	}
	for _, sc := range planning.PlansByState {
		if sc.Count != 0 {
			t.Errorf("PlansByState %q = %d, want 0", sc.State, sc.Count)
		}
	}
}
