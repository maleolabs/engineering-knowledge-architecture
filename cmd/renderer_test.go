package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
)

// rendererTestContext builds the plain style + an empty graph used by
// the renderer unit tests.
func rendererTestContext(t *testing.T) (*ui.Style, *bytes.Buffer, *view.Graph) {
	t.Helper()
	var buf bytes.Buffer
	s := ui.NewStyle(&buf, false)
	return s, &buf, view.NewGraph(".", nil)
}

// graphWith builds a graph over the given artifact lines (relations
// come from the artifact's own frontmatter-equivalent fields).
func graphWith(arts ...conformance.Artifact) *view.Graph {
	return view.NewGraph(".", arts)
}

// graphWithContainer builds a graph where each of the given work item
// forms is a member of container feather/ctr:wave-7 (one ticket per
// item), so renderer tests resolve container tags to "wave-7".
func graphWithContainer(forms ...string) *view.Graph {
	arts := []conformance.Artifact{
		{Namespace: "feather", Type: "ctr", ID: "wave-7",
			States: map[string]string{conformance.DomainContainerState: "active"}},
	}
	for i, form := range forms {
		parts := strings.SplitN(form, "/", 2)
		ns, rest := parts[0], parts[1]
		typeID := strings.SplitN(rest, ":", 2)
		token, id := typeID[0], typeID[1]
		arts = append(arts,
			conformance.Artifact{Namespace: ns, Type: token, ID: id,
				States: map[string]string{conformance.DomainExecutionState: "todo"}},
			conformance.Artifact{Namespace: ns, Type: "tkt", ID: fmt.Sprintf("tkt-%d", i),
				Relations: map[string][]string{"derives-from": {"ctr:wave-7", token + ":" + id}}},
		)
	}
	return view.NewGraph(".", arts)
}

// TestRenderBoardProjection: the board renders every work item with its
// container tag; unassigned items carry the unassigned tag.
func TestRenderBoardProjection(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.BoardProjection{
		Columns: view.BoardColumns{
			{State: "planned", WorkItems: []view.BoardItem{
				{WorkItem: view.WorkItem{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "planned"}, Containers: []string{"feather/ctr:wave-7"}},
			}},
			{State: "todo", WorkItems: []view.BoardItem{
				{WorkItem: view.WorkItem{Identity: "feather/sto:orphan", Type: "sto", ID: "orphan", State: "todo"}},
			}},
			{State: "in-progress", WorkItems: []view.BoardItem{
				{WorkItem: view.WorkItem{Identity: "feather/sto:beta", Type: "sto", ID: "beta", State: "in-progress"}, Containers: []string{"feather/ctr:wave-7"}},
			}},
			{State: "in-review", WorkItems: nil},
			{State: "done", WorkItems: []view.BoardItem{
				{WorkItem: view.WorkItem{Identity: "feather/ch:gamma", Type: "ch", ID: "gamma", State: "done"}, Containers: []string{"feather/ctr:wave-7"}},
			}},
		},
		Total:          4,
		Unassigned:     1,
		ContainerCount: 1,
	}
	renderBoardProjection(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"Board",
		"Container    all",
		"Domain       Execution",
		"4 work items across 1 container",
		"│ Planned (1)",
		"│ Todo (1)",
		"│ In Progress (1)",
		"│ In Review (0)",
		"│ Done (1)",
		"│ ▸ alpha (wave-7)",
		"│ ▸ orphan (unassigned)",
		"│ ▸ beta (wave-7)",
		"│ ▸ gamma (wave-7)",
		"—", // empty column
		"1 work item not referenced by any ticket container",
		"Total Work Items: 4",
		"Active Work: 1",
		"Completed Work: 1",
		"Review Queue: 0",
		"Unassigned: 1",
		"Overall Progress: ██░░░░░░░░ 1/4 (25%)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("board output must contain %q:\n%s", want, out)
		}
	}
}

// TestRenderBoardProjectionEmpty: no work items renders a calm empty
// projection with the full five-column shape and a zero summary.
func TestRenderBoardProjectionEmpty(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.BoardProjection{Columns: emptyBoardColumns()}
	renderBoardProjection(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"No work items.",
		"│ Planned (0)",
		"│ Todo (0)",
		"│ In Progress (0)",
		"│ In Review (0)",
		"│ Done (0)",
		"Total Work Items: 0",
		"Overall Progress: ░░░░░░░░░░ 0/0 (0%)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("empty board output must contain %q:\n%s", want, out)
		}
	}
}

// emptyBoardColumns returns the fixed five empty board columns.
func emptyBoardColumns() view.BoardColumns {
	cols := make(view.BoardColumns, 0, 5)
	for _, state := range []string{"planned", "todo", "in-progress", "in-review", "done"} {
		cols = append(cols, view.BoardColumn{State: state})
	}
	return cols
}

// TestRenderBoardCellLabel: the composed label truncates the id, never
// the container tag — the tag is the board's context and must always
// stay visible. The budget mirrors the default non-TTY column width.
func TestRenderBoardCellLabel(t *testing.T) {
	budget := ui.BoardItemBudget(0, 5)
	cases := []struct {
		id, tag string
		want    string
	}{
		// Fits: full id + tag.
		{"alpha", "wave-7", "alpha (wave-7)"},
		// Does not fit: id truncated, tag intact.
		{"markdown-syntax-highlighting", "wave-7", "markdown-syntax-high… (wave-7)"},
		// Unassigned tag kept too.
		{"markdown-syntax-highlighting", "unassigned", "markdown-syntax-… (unassigned)"},
	}
	for _, c := range cases {
		if got := boardCellLabel(c.id, c.tag, budget); got != c.want {
			t.Errorf("boardCellLabel(%q, %q, %d) = %q, want %q", c.id, c.tag, budget, got, c.want)
		}
	}
}

// TestRenderBoardAdaptiveBudget: on a narrower terminal the label
// budget shrinks with the column width; the tag still survives.
func TestRenderBoardAdaptiveBudget(t *testing.T) {
	// 80-cell terminal: (80-16)/5 = 12 per column, budget 10.
	budget := ui.BoardItemBudget(80, 5)
	if budget != 10 {
		t.Fatalf("BoardItemBudget(80, 5) = %d, want 10", budget)
	}
	got := boardCellLabel("markdown-syntax-highlighting", "wave-7", budget)
	want := "… (wave-7)"
	if got != want {
		t.Errorf("boardCellLabel on 80-col = %q, want %q (tag intact)", got, want)
	}
}

func TestRenderExecutionBoard(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.ExecutionProjection{
		Container: &view.Container{Identity: "feather/ctr:wave-7", State: "active"},
		Tickets: []view.Ticket{
			{Identity: "feather/tkt:a", Projected: "done"},
			{Identity: "feather/tkt:b", Projected: "in-progress"},
		},
		Columns: view.StateColumns{
			{State: "planned", WorkItems: []view.WorkItem{{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "planned"}}},
			{State: "todo", WorkItems: nil},
			{State: "in-progress", WorkItems: []view.WorkItem{{Identity: "feather/sto:beta", Type: "sto", ID: "beta", State: "in-progress"}}},
			{State: "in-review", WorkItems: nil},
			{State: "done", WorkItems: []view.WorkItem{{Identity: "feather/ch:gamma", Type: "ch", ID: "gamma", State: "done"}}},
		},
		Total: 3,
	}
	// The items resolve to container wave-7 through the graph, so their
	// labels carry the container tag — same rule as the board.
	g = graphWithContainer("feather/sto:alpha", "feather/sto:beta", "feather/ch:gamma")
	renderExecution(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"Execution",
		"• feather/ctr:wave-7  (active)",
		"┌",
		"│ Planned (1)",
		"│ Todo (0)",
		"│ In Progress (1)",
		"│ In Review (0)",
		"│ Done (1)",
		"│ ▸ alpha (wave-7)",
		"│ ▸ beta (wave-7)",
		"│ ▸ gamma (wave-7)",
		"—", // empty columns
		"2 tickets project these work items",
		"Active Work: 1",
		"Completed Work: 1",
		"Review Queue: 0",
		"Overall Progress: ███░░░░░░░ 1/3 (33%)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("execution board output must contain %q:\n%s", want, out)
		}
	}
}

// TestRenderExecutionSharedContainerTag: an item referenced by tickets
// of two containers shows both tags on the active container's board —
// the same tag rule as the board projection.
func TestRenderExecutionSharedContainerTag(t *testing.T) {
	arts := []conformance.Artifact{
		{Namespace: "feather", Type: "ctr", ID: "wave-7", States: map[string]string{conformance.DomainContainerState: "active"}},
		{Namespace: "feather", Type: "ctr", ID: "wave-0", States: map[string]string{conformance.DomainContainerState: "completed"}},
		{Namespace: "feather", Type: "sto", ID: "shared", States: map[string]string{conformance.DomainExecutionState: "in-progress"}},
		{Namespace: "feather", Type: "tkt", ID: "one", Relations: map[string][]string{"derives-from": {"ctr:wave-7", "sto:shared"}}},
		{Namespace: "feather", Type: "tkt", ID: "two", Relations: map[string][]string{"derives-from": {"ctr:wave-0", "sto:shared"}}},
	}
	s, buf, _ := rendererTestContext(t)
	g := graphWith(arts...)
	p := &view.ExecutionProjection{
		Container: &view.Container{Identity: "feather/ctr:wave-7", State: "active"},
		Columns: view.StateColumns{
			{State: "in-progress", WorkItems: []view.WorkItem{{Identity: "feather/sto:shared", Type: "sto", ID: "shared", State: "in-progress"}}},
		},
		Total: 1,
	}
	renderExecution(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"│ ▸ shared (wave-0, wave-7)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("execution board output must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderExecutionNoContainer(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.ExecutionProjection{
		Columns: view.StateColumns{
			{State: "planned"}, {State: "todo"}, {State: "in-progress"},
			{State: "in-review"}, {State: "done"},
		},
	}
	renderExecution(s, g, p)
	out := buf.String()
	for _, want := range []string{"No active container.", "—", "0/0 (0%)"} {
		if !strings.Contains(out, want) {
			t.Errorf("empty execution output must contain %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "tickets project") {
		t.Errorf("no-container output must not claim tickets:\n%s", out)
	}
}

func TestRenderExecutionShortIDAmbiguity(t *testing.T) {
	// Same id across two work item types keeps the type prefix.
	items := []view.WorkItem{
		{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "planned"},
		{Identity: "feather/ts:alpha", Type: "ts", ID: "alpha", State: "done"},
		{Identity: "feather/sto:beta", Type: "sto", ID: "beta", State: "done"},
	}
	short := shortWorkItemID(items)
	if got := short(items[0]); got != "sto:alpha" {
		t.Errorf("ambiguous id must keep the type prefix, got %q", got)
	}
	if got := short(items[1]); got != "ts:alpha" {
		t.Errorf("ambiguous id must keep the type prefix, got %q", got)
	}
	if got := short(items[2]); got != "beta" {
		t.Errorf("unique id must render bare, got %q", got)
	}
}

func TestRenderPlanningRoadmap(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.PlanningProjection{
		Groups: []view.Group{
			{Name: "Scope Definitions", Artifacts: []view.DomainArtifact{
				{Identity: "feather/scp:mvp-v1", Type: "scp", ID: "mvp-v1", ContentState: "approved", HasContentState: true, Phase: "mvp", HasPhase: true},
			}},
			{Name: "Epics", Artifacts: []view.DomainArtifact{
				{Identity: "feather/epc:authoring", Type: "epc", ID: "authoring", ContentState: "review", HasContentState: true},
				{Identity: "feather/epc:distribution", Type: "epc", ID: "distribution", ContentState: "draft", HasContentState: true},
			}},
			{Name: "Plans", Artifacts: []view.DomainArtifact{
				{Identity: "feather/plan:roadmap-v1", Type: "plan", ID: "roadmap-v1", ContentState: "approved", HasContentState: true, PlanningState: "approved", HasPlanningState: true, Phase: "mvp", HasPhase: true},
			}},
			{Name: "Traceability", Artifacts: []view.DomainArtifact{
				{Identity: "feather/trc:feather-trace", Type: "trc", ID: "feather-trace", ContentState: "approved", HasContentState: true},
			}},
		},
		PlansByState: []view.StateCount{
			{State: "draft", Count: 0}, {State: "approved", Count: 1}, {State: "immutable", Count: 0},
		},
	}
	renderPlanning(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"Planning",
		"✓ feather/plan:roadmap-v1  (approved, planning-state approved, phase mvp)",
		"────",
		"│ ▸ feather/scp:mvp-v1  (approved, phase mvp)",
		"│ ▸ feather/epc:authoring  (review)",
		"│ ▸ feather/epc:distribution  (draft)",
		"│ ▸ traceability: feather/trc:feather-trace (approved)",
		"Committed: 1",
		"Exploring: 1",
		"Next milestone: mvp",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("planning roadmap output must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderPlanningNoPlan(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.PlanningProjection{
		Groups: []view.Group{
			{Name: "Epics", Artifacts: []view.DomainArtifact{
				{Identity: "feather/epc:distribution", Type: "epc", ID: "distribution", ContentState: "draft", HasContentState: true},
			}},
		},
		PlansByState: []view.StateCount{{State: "draft"}, {State: "approved"}, {State: "immutable"}},
	}
	renderPlanning(s, g, p)
	out := buf.String()
	for _, want := range []string{"no plan yet — roadmap undefined", "Next milestone: —"} {
		if !strings.Contains(out, want) {
			t.Errorf("plan-less roadmap must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderArchitectureTree(t *testing.T) {
	s, buf, _ := rendererTestContext(t)
	g := graphWith(
		conformance.Artifact{Namespace: "feather", Type: "adr", ID: "content-storage", Revision: 1,
			Relations: map[string][]string{"depends-on": {"fnd:markdown-editor-options"}}},
		conformance.Artifact{Namespace: "feather", Type: "fnd", ID: "markdown-editor-options", Revision: 1},
		conformance.Artifact{Namespace: "feather", Type: "arc", ID: "feather-system", Revision: 1},
	)
	p := &view.ArchitectureProjection{
		Groups: []view.Group{
			{Name: "Decisions", Artifacts: []view.DomainArtifact{
				{Identity: "feather/adr:content-storage", Type: "adr", ID: "content-storage", ContentState: "accepted", HasContentState: true},
				{Identity: "feather/dec:reverse-proxy", Type: "dec", ID: "reverse-proxy", ContentState: "review", HasContentState: true},
			}},
			{Name: "Architecture Descriptions", Artifacts: []view.DomainArtifact{
				{Identity: "feather/arc:feather-system", Type: "arc", ID: "feather-system", ContentState: "approved", HasContentState: true},
			}},
			{Name: "Specifications", Artifacts: nil},
			{Name: "Standards & Guidelines", Artifacts: []view.DomainArtifact{
				{Identity: "feather/std:definition-of-done", Type: "std", ID: "definition-of-done", ContentState: "approved", HasContentState: true},
			}},
			{Name: "Vocabulary", Artifacts: []view.DomainArtifact{
				{Identity: "feather/gls:feather-terms", Type: "gls", ID: "feather-terms", ContentState: "amended", HasContentState: true},
			}},
		},
	}
	renderArchitecture(s, g, p)
	out := buf.String()
	for _, want := range []string{
		"feather/arc:feather-system  (approved)",
		"├── Decisions",
		"│  ├── ✓ feather/adr:content-storage  (accepted) (depends-on fnd:markdown-editor-options)",
		"│  └── • feather/dec:reverse-proxy  (review)",
		"├── Standards & Guidelines",
		"│  └── ✓ feather/std:definition-of-done  (approved)",
		"└── Vocabulary",
		"   └── • feather/gls:feather-terms  (amended)",
		"Accepted decisions: 1",
		"Open items: 1",
		"Superseded: 0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("architecture tree output must contain %q:\n%s", want, out)
		}
	}
	// The empty Specifications group is skipped.
	if strings.Contains(out, "Specifications") {
		t.Errorf("empty group must be skipped:\n%s", out)
	}
}

func TestRenderArchitectureNoDescription(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.ArchitectureProjection{
		Groups: []view.Group{
			{Name: "Decisions", Artifacts: []view.DomainArtifact{
				{Identity: "feather/adr:one", Type: "adr", ID: "one", ContentState: "accepted", HasContentState: true},
			}},
		},
	}
	renderArchitecture(s, g, p)
	out := buf.String()
	if !strings.Contains(out, "Architecture\n") {
		t.Errorf("no arc- artifact must root the tree at the Architecture node:\n%s", out)
	}
	// A single child subtree renders as the last branch.
	if !strings.Contains(out, "└── Decisions") {
		t.Errorf("decisions subtree must render under the root:\n%s", out)
	}
}

func TestRenderDiscoveryCards(t *testing.T) {
	s, buf, _ := rendererTestContext(t)
	g2 := graphWith(
		conformance.Artifact{Namespace: "feather", Type: "vis", ID: "feather-vision", Revision: 3},
		conformance.Artifact{Namespace: "feather", Type: "req", ID: "comments-phase2", Revision: 1},
	)
	p := &view.DiscoveryProjection{
		Groups: []view.Group{
			{Name: "Vision", Artifacts: []view.DomainArtifact{
				{Identity: "feather/vis:feather-vision", Type: "vis", ID: "feather-vision", ContentState: "approved", HasContentState: true},
			}},
			{Name: "Strategy", Artifacts: nil},
			{Name: "Requirements", Artifacts: []view.DomainArtifact{
				{Identity: "feather/req:comments-phase2", Type: "req", ID: "comments-phase2", ContentState: "draft", HasContentState: true},
			}},
			{Name: "Research Findings", Artifacts: nil},
		},
	}
	renderDiscovery(s, g2, p)
	out := buf.String()
	for _, want := range []string{
		"Vision",
		"┌",
		"│ ✓ feather/vis:feather-vision │",
		"│ approved · revision 3",
		"└",
		"Requirements",
		"│ ○ feather/req:comments-phase2 │",
		"│ draft · revision 1",
		"Committed direction: 1",
		"Exploring: 1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("discovery cards output must contain %q:\n%s", want, out)
		}
	}
	// All content rows share one box width per group; the box width is
	// the widest content line + 4 (bars and pads).
	wantWidths := map[int]bool{}
	for _, s := range []string{"✓ feather/vis:feather-vision", "○ feather/req:comments-phase2"} {
		wantWidths[len([]rune(s))+4] = true
	}
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if !strings.HasPrefix(line, "│") {
			continue
		}
		if !wantWidths[len([]rune(line))] {
			t.Errorf("discovery card row %q spans %d cells, want one of %v", line, len([]rune(line)), wantWidths)
		}
	}
	if strings.Contains(out, "Strategy") || strings.Contains(out, "Research Findings") {
		t.Errorf("empty groups must be skipped:\n%s", out)
	}
}

func TestRenderOperationsRelease(t *testing.T) {
	s, buf, _ := rendererTestContext(t)
	g2 := graphWith(
		conformance.Artifact{Namespace: "feather", Type: "rel", ID: "v090", Revision: 1,
			Relations: map[string][]string{"derives-from": {"plan:roadmap-v1:1"}}},
		conformance.Artifact{Namespace: "feather", Type: "plan", ID: "roadmap-v1", InstanceVersion: 1, Revision: 1},
	)
	p := &view.OperationsProjection{
		Groups: []view.Group{
			{Name: "Runbooks", Artifacts: []view.DomainArtifact{
				{Identity: "feather/run:deploy-feather", Type: "run", ID: "deploy-feather", ContentState: "approved", HasContentState: true},
				{Identity: "feather/run:backup-feather", Type: "run", ID: "backup-feather", ContentState: "draft", HasContentState: true},
			}},
			{Name: "Release Records", Artifacts: []view.DomainArtifact{
				{Identity: "feather/rel:v090", Type: "rel", ID: "v090", ContentState: "approved", HasContentState: true},
			}},
		},
	}
	renderOperations(s, g2, p)
	out := buf.String()
	for _, want := range []string{
		"Release Records",
		"┌",
		"│ ✓ feather/rel:v090",
		"│ approved",
		"│ derives-from plan:roadmap-v1:1",
		"Runbooks",
		"▸ feather/run:deploy-feather  (approved)",
		"│ ▸ feather/run:backup-feather  (draft)",
		"Releases delivered: 1",
		"Runbooks maintained: 2",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("operations output must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderTicketDetail(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.TicketProjection{
		Ticket:     view.Ticket{Identity: "feather/tkt:sto-draft-autosave", Type: "tkt", ID: "sto-draft-autosave", Projected: "in-progress"},
		Container:  &view.Container{Identity: "feather/ctr:wave-7", State: "active"},
		WorkItem:   &view.WorkItem{Identity: "feather/sto:draft-autosave", Type: "sto", ID: "draft-autosave", State: "in-progress"},
		Projected:  "in-progress",
		References: []string{"ctr:wave-7", "sto:draft-autosave"},
	}
	renderTicket(s, g, p)
	out := buf.String()
	// The supporting rows align labels to the widest label ("Derives
	// From", 12 cells) plus 3 spaces, and the card pads every row to
	// its width — computed, not hand-counted.
	rows := [][2]string{
		{"Work Item", "feather/sto:draft-autosave (in-progress)"},
		{"Container", "feather/ctr:wave-7"},
		{"Derives From", "ctr:wave-7, sto:draft-autosave"},
	}
	width := len([]rune("feather/tkt:sto-draft-autosave"))
	for _, r := range rows {
		if w := 12 + 3 + len([]rune(r[1])); w > width {
			width = w
		}
	}
	ticketRow := func(label, value string) string {
		content := fmt.Sprintf("%-12s   %s", label, value)
		return "│ " + content + strings.Repeat(" ", width-len([]rune(content))) + " │"
	}
	for _, want := range []string{
		"Ticket",
		"Projected Status  → in-progress",
		"┌",
		"│ feather/tkt:sto-draft-autosave",
		ticketRow(rows[0][0], rows[0][1]),
		ticketRow(rows[1][0], rows[1][1]),
		ticketRow(rows[2][0], rows[2][1]),
		"└",
		"Projected status: in-progress",
		"Work item: feather/sto:draft-autosave (in-progress)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("ticket detail output must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderTicketUnresolved(t *testing.T) {
	s, buf, g := rendererTestContext(t)
	p := &view.TicketProjection{
		Ticket:    view.Ticket{Identity: "feather/tkt:ghost", Projected: "unresolved"},
		Projected: "unresolved",
	}
	renderTicket(s, g, p)
	out := buf.String()
	rows := [][2]string{
		{"Work Item", "unresolved"},
		{"Container", "unresolved"},
		{"Derives From", "—"},
	}
	width := len([]rune("feather/tkt:ghost"))
	for _, r := range rows {
		if w := 12 + 3 + len([]rune(r[1])); w > width {
			width = w
		}
	}
	ticketRow := func(label, value string) string {
		content := fmt.Sprintf("%-12s   %s", label, value)
		return "│ " + content + strings.Repeat(" ", width-len([]rune(content))) + " │"
	}
	for _, want := range []string{
		"Projected Status  • unresolved",
		ticketRow(rows[0][0], rows[0][1]),
		ticketRow(rows[1][0], rows[1][1]),
		ticketRow(rows[2][0], rows[2][1]),
		"Projected status: unresolved",
		"Work item: unresolved",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("unresolved ticket output must contain %q:\n%s", want, out)
		}
	}
}

func TestRenderEmptyDomains(t *testing.T) {
	cases := []struct {
		name    string
		render  func(*ui.Style, *view.Graph)
		line    string
		summary string
	}{
		{"planning", func(s *ui.Style, g *view.Graph) {
			renderPlanning(s, g, &view.PlanningProjection{Groups: nil, PlansByState: []view.StateCount{{State: "draft"}, {State: "approved"}, {State: "immutable"}}})
		}, "No Planning artifacts.", "Committed: 0"},
		{"architecture", func(s *ui.Style, g *view.Graph) {
			renderArchitecture(s, g, &view.ArchitectureProjection{Groups: nil})
		}, "No Architecture artifacts.", "Accepted decisions: 0"},
		{"discovery", func(s *ui.Style, g *view.Graph) {
			renderDiscovery(s, g, &view.DiscoveryProjection{Groups: nil})
		}, "No Discovery artifacts.", "Committed direction: 0"},
		{"operations", func(s *ui.Style, g *view.Graph) {
			renderOperations(s, g, &view.OperationsProjection{Groups: nil})
		}, "No Operations artifacts.", "Releases delivered: 0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			s := ui.NewStyle(&buf, false)
			g := view.NewGraph(".", nil)
			tc.render(s, g)
			out := buf.String()
			if !strings.Contains(out, tc.line) {
				t.Errorf("must render %q:\n%s", tc.line, out)
			}
			if !strings.Contains(out, tc.summary) {
				t.Errorf("must render insight %q:\n%s", tc.summary, out)
			}
			if !strings.Contains(out, "Summary:") {
				t.Errorf("must render the summary block:\n%s", out)
			}
		})
	}
}

func TestRenderersDeterministic(t *testing.T) {
	build := func() string {
		var buf bytes.Buffer
		s := ui.NewStyle(&buf, false)
		g := graphWithContainer("feather/sto:alpha")
		renderExecution(s, g, &view.ExecutionProjection{
			Container: &view.Container{Identity: "feather/ctr:wave-7", State: "active"},
			Columns: view.StateColumns{
				{State: "done", WorkItems: []view.WorkItem{{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "done"}}},
			},
			Total: 1,
		})
		renderPlanning(s, g, &view.PlanningProjection{Groups: nil, PlansByState: nil})
		renderArchitecture(s, g, &view.ArchitectureProjection{Groups: nil})
		renderDiscovery(s, g, &view.DiscoveryProjection{Groups: nil})
		renderOperations(s, g, &view.OperationsProjection{Groups: nil})
		renderTicket(s, g, &view.TicketProjection{Ticket: view.Ticket{Identity: "x"}, Projected: "unresolved"})
		renderBoardProjection(s, g, &view.BoardProjection{
			Columns: view.BoardColumns{
				{State: "done", WorkItems: []view.BoardItem{
					{WorkItem: view.WorkItem{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "done"}, Containers: []string{"feather/ctr:wave-7"}},
				}},
			},
			Total: 1, ContainerCount: 1,
		})
		return buf.String()
	}
	if build() != build() {
		t.Error("renderer output must be deterministic")
	}
}

func TestRenderersNoANSI(t *testing.T) {
	var buf bytes.Buffer
	s := ui.NewStyle(&buf, false)
	g := graphWithContainer("feather/sto:alpha")
	renderExecution(s, g, &view.ExecutionProjection{
		Container: &view.Container{Identity: "feather/ctr:wave-7", State: "active"},
		Columns: view.StateColumns{
			{State: "in-progress", WorkItems: []view.WorkItem{{Identity: "feather/sto:alpha", Type: "sto", ID: "alpha", State: "in-progress"}}},
			{State: "done", WorkItems: []view.WorkItem{{Identity: "feather/ch:beta", Type: "ch", ID: "beta", State: "done"}}},
		},
		Total: 2,
	})
	renderPlanning(s, g, &view.PlanningProjection{Groups: []view.Group{
		{Name: "Plans", Artifacts: []view.DomainArtifact{
			{Identity: "feather/plan:roadmap", Type: "plan", ID: "roadmap", ContentState: "approved", HasContentState: true, PlanningState: "approved", HasPlanningState: true, Phase: "mvp", HasPhase: true},
		}},
	}, PlansByState: []view.StateCount{{State: "draft"}, {State: "approved", Count: 1}, {State: "immutable"}}})
	renderArchitecture(s, g, &view.ArchitectureProjection{Groups: []view.Group{
		{Name: "Decisions", Artifacts: []view.DomainArtifact{
			{Identity: "feather/adr:one", Type: "adr", ID: "one", ContentState: "accepted", HasContentState: true},
		}},
	}})
	renderDiscovery(s, g, &view.DiscoveryProjection{Groups: []view.Group{
		{Name: "Vision", Artifacts: []view.DomainArtifact{
			{Identity: "feather/vis:v", Type: "vis", ID: "v", ContentState: "approved", HasContentState: true},
		}},
	}})
	renderOperations(s, g, &view.OperationsProjection{Groups: []view.Group{
		{Name: "Runbooks", Artifacts: []view.DomainArtifact{
			{Identity: "feather/run:r", Type: "run", ID: "r", ContentState: "approved", HasContentState: true},
		}},
	}})
	renderTicket(s, g, &view.TicketProjection{Ticket: view.Ticket{Identity: "feather/tkt:t", Projected: "done"},
		WorkItem: &view.WorkItem{Identity: "feather/sto:w", Type: "sto", ID: "w", State: "done"}, Projected: "done"})
	renderBoardProjection(s, g, &view.BoardProjection{
		Columns: view.BoardColumns{
			{State: "done", WorkItems: []view.BoardItem{
				{WorkItem: view.WorkItem{Identity: "feather/sto:b", Type: "sto", ID: "b", State: "done"}, Containers: []string{"feather/ctr:wave-7"}},
			}},
		},
		Total: 1, ContainerCount: 1,
	})
	if strings.Contains(buf.String(), "\x1b") {
		t.Error("non-TTY renderer output must not contain ANSI escapes")
	}
}
