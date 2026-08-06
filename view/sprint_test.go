package view

import (
	"reflect"
	"testing"
)

// sprintColumnIDs returns the identities of one column.
func sprintColumnIDs(col StateColumn) []string {
	out := make([]string, 0, len(col.WorkItems))
	for _, wi := range col.WorkItems {
		out = append(out, wi.Identity)
	}
	return out
}

func TestSprintProjection(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("sprint", g, "")
	if err != nil {
		t.Fatalf("Build(sprint): %v", err)
	}
	sprint, ok := p.(*SprintProjection)
	if !ok {
		t.Fatalf("Build(sprint) = %T, want *SprintProjection", p)
	}
	if sprint.Container == nil || sprint.Container.Identity != validForm+"ctr:wave-1" {
		t.Fatalf("Container = %+v, want ctr:wave-1", sprint.Container)
	}
	if sprint.MultipleActive {
		t.Error("MultipleActive must be false on the valid fixture")
	}
	if sprint.Tickets != 8 {
		t.Errorf("Tickets = %d, want 8 (seven resolved + one unresolved)", sprint.Tickets)
	}
	if sprint.Total != 5 {
		t.Errorf("Total = %d, want 5 (dedup keeps each work item once)", sprint.Total)
	}

	// Fixed column order.
	wantOrder := []string{"planned", "todo", "in-progress", "in-review", "done"}
	gotOrder := make([]string, len(sprint.Columns))
	for i, col := range sprint.Columns {
		gotOrder[i] = col.State
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Errorf("column order = %v, want %v", gotOrder, wantOrder)
	}

	// Per-column membership.
	cases := map[string][]string{
		"planned":     {validForm + "sto:alpha"},
		"todo":        {validForm + "sto:beta"},
		"in-progress": {validForm + "ts:gamma"},
		"in-review":   {validForm + "bug:delta"},
		"done":        {validForm + "ch:epsilon"},
	}
	for _, col := range sprint.Columns {
		want := cases[col.State]
		if got := sprintColumnIDs(col); !reflect.DeepEqual(got, want) {
			t.Errorf("column %q = %v, want %v", col.State, got, want)
		}
	}

	// Counts derive from the columns.
	for state, n := range map[string]int{"planned": 1, "todo": 1, "in-progress": 1, "in-review": 1, "done": 1} {
		if got := sprint.Columns.Count(state); got != n {
			t.Errorf("Count(%q) = %d, want %d", state, got, n)
		}
	}
	if got := sprint.Columns.Count("nope"); got != 0 {
		t.Errorf("Count(unknown) = %d, want 0", got)
	}

	// Work items carry the owner state and informational dimension.
	for _, col := range sprint.Columns {
		for _, wi := range col.WorkItems {
			if wi.State != col.State {
				t.Errorf("work item %s in column %q but owns state %q", wi.Identity, col.State, wi.State)
			}
		}
	}
	var alpha *WorkItem
	for _, col := range sprint.Columns {
		for i := range col.WorkItems {
			if col.WorkItems[i].ID == "alpha" {
				alpha = &col.WorkItems[i]
			}
		}
	}
	if alpha == nil || !alpha.HasDimension || alpha.Dimension != "operations" {
		t.Errorf("sto:alpha must carry its informational dimension, got %+v", alpha)
	}
}

func TestSprintProjectionNoActiveContainer(t *testing.T) {
	g := loadFixture(t, "no-active")
	p, err := Build("sprint", g, "")
	if err != nil {
		t.Fatalf("Build(sprint): %v", err)
	}
	sprint := p.(*SprintProjection)
	if sprint.Container != nil {
		t.Errorf("Container = %+v, want nil", sprint.Container)
	}
	if sprint.MultipleActive {
		t.Error("MultipleActive must be false")
	}
	if len(sprint.Columns) != 5 {
		t.Fatalf("columns = %d, want the fixed five", len(sprint.Columns))
	}
	for _, col := range sprint.Columns {
		if len(col.WorkItems) != 0 {
			t.Errorf("column %q must be empty, got %v", col.State, sprintColumnIDs(col))
		}
	}
	if sprint.Total != 0 || sprint.Tickets != 0 {
		t.Errorf("Total = %d, Tickets = %d; want 0/0", sprint.Total, sprint.Tickets)
	}
}

func TestSprintProjectionMultipleActive(t *testing.T) {
	g := loadFixture(t, "multi-active")
	p, err := Build("sprint", g, "")
	if err != nil {
		t.Fatalf("Build(sprint): %v", err)
	}
	sprint := p.(*SprintProjection)
	if !sprint.MultipleActive {
		t.Error("MultipleActive must be true on the multi-active fixture")
	}
	if sprint.Container == nil || sprint.Container.ID != "wave-1" {
		t.Errorf("Container = %+v, want the lexicographically smallest (wave-1)", sprint.Container)
	}
}
