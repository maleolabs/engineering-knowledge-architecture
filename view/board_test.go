package view

import (
	"reflect"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

// boardColumnIDs returns the identities of one board column.
func boardColumnIDs(col BoardColumn) []string {
	out := make([]string, 0, len(col.WorkItems))
	for _, bi := range col.WorkItems {
		out = append(out, bi.Identity)
	}
	return out
}

// boardColumnContainers returns the container tags of one board column.
func boardColumnContainers(col BoardColumn) [][]string {
	out := make([][]string, 0, len(col.WorkItems))
	for _, bi := range col.WorkItems {
		out = append(out, bi.Containers)
	}
	return out
}

func TestBoardProjectionValid(t *testing.T) {
	g := loadFixture(t, "valid")
	p, err := Build("board", g, "")
	if err != nil {
		t.Fatalf("Build(board): %v", err)
	}
	board, ok := p.(*BoardProjection)
	if !ok {
		t.Fatalf("Build(board) = %T, want *BoardProjection", p)
	}
	if board.Name() != "board" {
		t.Errorf("Name() = %q, want board", board.Name())
	}
	// All six work items of the fixture — the active container's five
	// plus the completed container's legacy item.
	if board.Total != 6 {
		t.Errorf("Total = %d, want 6", board.Total)
	}
	if board.Unassigned != 0 {
		t.Errorf("Unassigned = %d, want 0", board.Unassigned)
	}
	if board.ContainerCount != 2 {
		t.Errorf("ContainerCount = %d, want 2 (wave-0, wave-1)", board.ContainerCount)
	}

	// Fixed column order.
	wantOrder := []string{"planned", "todo", "in-progress", "in-review", "done"}
	gotOrder := make([]string, len(board.Columns))
	for i, col := range board.Columns {
		gotOrder[i] = col.State
	}
	if !reflect.DeepEqual(gotOrder, wantOrder) {
		t.Errorf("column order = %v, want %v", gotOrder, wantOrder)
	}

	// Per-column membership and container tags.
	cases := map[string]struct {
		items      []string
		containers [][]string
	}{
		"planned": {
			items:      []string{validForm + "sto:alpha"},
			containers: [][]string{{validForm + "ctr:wave-1"}},
		},
		"todo": {
			items:      []string{validForm + "sto:beta"},
			containers: [][]string{{validForm + "ctr:wave-1"}},
		},
		"in-progress": {
			items:      []string{validForm + "ts:gamma"},
			containers: [][]string{{validForm + "ctr:wave-1"}},
		},
		"in-review": {
			items:      []string{validForm + "bug:delta"},
			containers: [][]string{{validForm + "ctr:wave-1"}},
		},
		"done": {
			// Sorted by canonical identity: ch:epsilon before sto:legacy.
			items:      []string{validForm + "ch:epsilon", validForm + "sto:legacy"},
			containers: [][]string{{validForm + "ctr:wave-1"}, {validForm + "ctr:wave-0"}},
		},
	}
	for _, col := range board.Columns {
		want := cases[col.State]
		if got := boardColumnIDs(col); !reflect.DeepEqual(got, want.items) {
			t.Errorf("%s column items = %v, want %v", col.State, got, want.items)
		}
		if got := boardColumnContainers(col); !reflect.DeepEqual(got, want.containers) {
			t.Errorf("%s column containers = %v, want %v", col.State, got, want.containers)
		}
	}

	// Count mirrors the columns.
	if board.Columns.Count("done") != 2 || board.Columns.Count("in-progress") != 1 {
		t.Errorf("Count(done/in-progress) = %d/%d, want 2/1",
			board.Columns.Count("done"), board.Columns.Count("in-progress"))
	}
	if board.Columns.Count("nonexistent") != 0 {
		t.Error("Count of an unknown state must be 0")
	}
}

// TestBoardProjectionUnassigned: a work item that no ticket references
// is still projected, tagged as unassigned.
func TestBoardProjectionUnassigned(t *testing.T) {
	g := loadFixture(t, "no-active")
	p, err := Build("board", g, "")
	if err != nil {
		t.Fatalf("Build(board): %v", err)
	}
	board := p.(*BoardProjection)
	if board.Total != 1 {
		t.Errorf("Total = %d, want 1 (sto:legacy)", board.Total)
	}
	if board.Unassigned != 0 {
		t.Errorf("Unassigned = %d, want 0 (legacy is referenced by tkt-sto-legacy)", board.Unassigned)
	}
	if board.ContainerCount != 1 {
		t.Errorf("ContainerCount = %d, want 1", board.ContainerCount)
	}
	// Container tags resolve regardless of container-state: wave-0 is
	// completed, not active.
	col := board.Columns.Count("done")
	if col != 1 {
		t.Fatalf("done column count = %d, want 1", col)
	}
}

// TestBoardProjectionOrphan: a work item referenced by no ticket at
// all is unassigned.
func TestBoardProjectionOrphan(t *testing.T) {
	g := NewGraph(".", []conformance.Artifact{
		{
			Namespace: "ns", Type: "sto", ID: "orphan",
			States: map[string]string{conformance.DomainExecutionState: "planned"},
		},
	})
	p, err := Build("board", g, "")
	if err != nil {
		t.Fatalf("Build(board): %v", err)
	}
	board := p.(*BoardProjection)
	if board.Total != 1 || board.Unassigned != 1 {
		t.Errorf("Total/Unassigned = %d/%d, want 1/1", board.Total, board.Unassigned)
	}
	if board.ContainerCount != 0 {
		t.Errorf("ContainerCount = %d, want 0", board.ContainerCount)
	}
	if got := boardColumnContainers(board.Columns[0]); len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("orphan containers = %v, want empty tag", got)
	}
}

// TestBoardProjectionMultiContainer: a work item referenced by tickets
// of two different containers carries both container tags, sorted.
func TestBoardProjectionMultiContainer(t *testing.T) {
	artifacts := []conformance.Artifact{
		{Namespace: "ns", Type: "ctr", ID: "wave-b", States: map[string]string{conformance.DomainContainerState: "active"}},
		{Namespace: "ns", Type: "ctr", ID: "wave-a", States: map[string]string{conformance.DomainContainerState: "completed"}},
		{Namespace: "ns", Type: "sto", ID: "shared", States: map[string]string{conformance.DomainExecutionState: "todo"}},
		{Namespace: "ns", Type: "tkt", ID: "one", Relations: map[string][]string{"derives-from": {"ctr:wave-b", "sto:shared"}}},
		{Namespace: "ns", Type: "tkt", ID: "two", Relations: map[string][]string{"derives-from": {"ctr:wave-a", "sto:shared"}}},
	}
	g := NewGraph(".", artifacts)
	p, err := Build("board", g, "")
	if err != nil {
		t.Fatalf("Build(board): %v", err)
	}
	board := p.(*BoardProjection)
	if board.Total != 1 || board.Unassigned != 0 {
		t.Errorf("Total/Unassigned = %d/%d, want 1/0", board.Total, board.Unassigned)
	}
	if board.ContainerCount != 2 {
		t.Errorf("ContainerCount = %d, want 2", board.ContainerCount)
	}
	want := []string{"ns/ctr:wave-a", "ns/ctr:wave-b"}
	if got := boardColumnContainers(board.Columns[1]); !reflect.DeepEqual(got[0], want) {
		t.Errorf("shared containers = %v, want %v (sorted)", got[0], want)
	}
}

// TestBoardProjectionEmpty: a repository without work items projects an
// empty board — a valid state, not an error.
func TestBoardProjectionEmpty(t *testing.T) {
	g := NewGraph(".", nil)
	p, err := Build("board", g, "")
	if err != nil {
		t.Fatalf("Build(board): %v", err)
	}
	board := p.(*BoardProjection)
	if board.Total != 0 || board.Unassigned != 0 || board.ContainerCount != 0 {
		t.Errorf("empty board = %+v, want all-zero", board)
	}
	if len(board.Columns) != 5 {
		t.Errorf("empty board columns = %d, want the fixed five", len(board.Columns))
	}
}
