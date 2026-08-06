package view

import (
	"reflect"
	"testing"
)

func TestWaveProjection(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("wave", g, "")
	if err != nil {
		t.Fatalf("Build(wave): %v", err)
	}
	wave, ok := p.(*WaveProjection)
	if !ok {
		t.Fatalf("Build(wave) = %T, want *WaveProjection", p)
	}
	if wave.Container == nil || wave.Container.Identity != validForm+"ctr:wave-1" {
		t.Fatalf("Container = %+v, want ctr:wave-1", wave.Container)
	}
	if wave.MultipleActive {
		t.Error("MultipleActive must be false on the valid fixture")
	}
	if wave.Total != 5 {
		t.Errorf("Total = %d, want 5", wave.Total)
	}

	// Tickets deriving from the active container, sorted by identity,
	// each carrying the status projected from its work item. The
	// duplicate ticket (sto-alpha-dup) and the multi-reference ticket
	// (sto-beta-multi, first resolvable reference wins) are members too.
	wantTickets := []Ticket{
		{Identity: validForm + "tkt:bug-delta", Type: "tkt", ID: "bug-delta", Projected: "in-review"},
		{Identity: validForm + "tkt:ch-epsilon", Type: "tkt", ID: "ch-epsilon", Projected: "done"},
		{Identity: validForm + "tkt:sto-alpha", Type: "tkt", ID: "sto-alpha", Projected: "planned"},
		{Identity: validForm + "tkt:sto-alpha-dup", Type: "tkt", ID: "sto-alpha-dup", Projected: "planned"},
		{Identity: validForm + "tkt:sto-beta", Type: "tkt", ID: "sto-beta", Projected: "todo"},
		{Identity: validForm + "tkt:sto-beta-multi", Type: "tkt", ID: "sto-beta-multi", Projected: "todo"},
		{Identity: validForm + "tkt:ts-gamma", Type: "tkt", ID: "ts-gamma", Projected: "in-progress"},
		{Identity: validForm + "tkt:unresolved", Type: "tkt", ID: "unresolved", Projected: "unresolved"},
	}
	if !reflect.DeepEqual(wave.Tickets, wantTickets) {
		t.Errorf("Tickets = %+v\nwant %+v", wave.Tickets, wantTickets)
	}

	// Progress counts per state, fixed column order.
	wantOrder := []string{"planned", "todo", "in-progress", "in-review", "done"}
	for i, col := range wave.Columns {
		if col.State != wantOrder[i] {
			t.Errorf("column %d = %q, want %q", i, col.State, wantOrder[i])
		}
	}
	for state, n := range map[string]int{"planned": 1, "todo": 1, "in-progress": 1, "in-review": 1, "done": 1} {
		if got := wave.Columns.Count(state); got != n {
			t.Errorf("Count(%q) = %d, want %d", state, got, n)
		}
	}
}

func TestWaveProjectionNoActiveContainer(t *testing.T) {
	g := loadFixture(t, "no-active")
	p, err := Build("wave", g, "")
	if err != nil {
		t.Fatalf("Build(wave): %v", err)
	}
	wave := p.(*WaveProjection)
	if wave.Container != nil {
		t.Errorf("Container = %+v, want nil", wave.Container)
	}
	if len(wave.Tickets) != 0 {
		t.Errorf("Tickets = %v, want none", wave.Tickets)
	}
	for _, col := range wave.Columns {
		if len(col.WorkItems) != 0 {
			t.Errorf("column %q must be empty", col.State)
		}
	}
	if wave.Total != 0 {
		t.Errorf("Total = %d, want 0", wave.Total)
	}
}
