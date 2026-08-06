package view

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
)

// loadFixture scans a fixture repository and builds its Knowledge Graph.
func loadFixture(t *testing.T, name string) *Graph {
	t.Helper()
	artifacts, err := conformance.Scan(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("scan fixture %s: %v", name, err)
	}
	return NewGraph(".", artifacts)
}

// validForm is the canonical identity form used across the valid
// fixture assertions.
const validForm = "eka-view-fixture/"

func TestGraphBuild(t *testing.T) {
	g := loadFixture(t, "valid")

	if g.Root() != "." {
		t.Errorf("Root() = %q, want %q", g.Root(), ".")
	}

	// The identity index resolves canonical line forms.
	if a := g.ByLineForm(validForm + "ctr:wave-1"); a == nil {
		t.Error("ByLineForm(ctr:wave-1) must resolve")
	} else if a.Type != "ctr" || a.ID != "wave-1" {
		t.Errorf("ByLineForm(ctr:wave-1) = %s/%s:%s", a.Namespace, a.Type, a.ID)
	}
	if a := g.ByLineForm(validForm + "tkt:ts-gamma"); a == nil || a.Type != "tkt" {
		t.Error("ByLineForm(tkt:ts-gamma) must resolve to the ticket line")
	}
	if a := g.ByLineForm(validForm + "sto:ghost"); a != nil {
		t.Error("ByLineForm of an unknown line must be nil")
	}

	// Relationship resolution: a reference parses and resolves within
	// the referrer's namespace.
	ref, err := conformance.ParseReference("ts:gamma", "eka-view-fixture", "tkt")
	if err != nil {
		t.Fatalf("ParseReference: %v", err)
	}
	if a := g.Resolve(ref); a == nil || a.Type != "ts" || a.ID != "gamma" {
		t.Errorf("Resolve(ts:gamma) = %+v, want ts:gamma", a)
	}
	versioned, err := conformance.ParseReference("ctr:wave-1:1", "eka-view-fixture", "tkt")
	if err != nil {
		t.Fatalf("ParseReference versioned: %v", err)
	}
	if a := g.Resolve(versioned); a == nil || a.InstanceVersion != 1 {
		t.Errorf("Resolve(ctr:wave-1:1) must resolve the exact instance")
	}

	// Active container: wave-1, exactly one.
	container, multiple := g.ActiveContainer()
	if container == nil || container.ID != "wave-1" || container.State != "active" {
		t.Fatalf("ActiveContainer() = %+v, want wave-1/active", container)
	}
	if multiple {
		t.Error("ActiveContainer() must not report multiple active containers")
	}

	// Containers, sorted by canonical identity.
	wantContainers := []string{validForm + "ctr:wave-0", validForm + "ctr:wave-1"}
	gotContainers := make([]string, 0, len(g.Containers()))
	for _, c := range g.Containers() {
		gotContainers = append(gotContainers, c.Identity)
	}
	if !reflect.DeepEqual(gotContainers, wantContainers) {
		t.Errorf("Containers() = %v, want %v", gotContainers, wantContainers)
	}

	// Tickets deriving from the active container, sorted. Eight tickets:
	// the five regular ones plus tkt:sto-alpha-dup and tkt:sto-beta-multi
	// (dedup / first-resolvable-wins fixtures) and tkt:unresolved.
	wantTickets := []string{
		validForm + "tkt:bug-delta",
		validForm + "tkt:ch-epsilon",
		validForm + "tkt:sto-alpha",
		validForm + "tkt:sto-alpha-dup",
		validForm + "tkt:sto-beta",
		validForm + "tkt:sto-beta-multi",
		validForm + "tkt:ts-gamma",
		validForm + "tkt:unresolved",
	}
	gotTickets := make([]string, 0, 8)
	for _, tkt := range g.TicketsForContainer(validForm + "ctr:wave-1") {
		gotTickets = append(gotTickets, tkt.Identity)
	}
	if !reflect.DeepEqual(gotTickets, wantTickets) {
		t.Errorf("TicketsForContainer(wave-1) = %v, want %v", gotTickets, wantTickets)
	}

	// Work items of the active container, deduplicated and sorted: the
	// duplicate ticket for sto:alpha must not double the member.
	wantItems := []string{
		validForm + "bug:delta",
		validForm + "ch:epsilon",
		validForm + "sto:alpha",
		validForm + "sto:beta",
		validForm + "ts:gamma",
	}
	gotItems := make([]string, 0, 5)
	for _, wi := range g.WorkItemsForContainer(validForm + "ctr:wave-1") {
		gotItems = append(gotItems, wi.Identity)
	}
	if !reflect.DeepEqual(gotItems, wantItems) {
		t.Errorf("WorkItemsForContainer(wave-1) = %v, want %v", gotItems, wantItems)
	}

	// Ticket -> work item and ticket -> container resolution.
	ticket := g.ByLineForm(validForm + "tkt:ts-gamma")
	wi := g.WorkItemForTicket(ticket)
	if wi == nil || wi.Identity != validForm+"ts:gamma" || wi.State != "in-progress" {
		t.Errorf("WorkItemForTicket(ts-gamma) = %+v, want ts:gamma/in-progress", wi)
	}
	c := g.ContainerForTicket(ticket)
	if c == nil || c.Identity != validForm+"ctr:wave-1" {
		t.Errorf("ContainerForTicket(ts-gamma) = %+v, want ctr:wave-1", c)
	}

	// Target resolution forms.
	for target, wantID := range map[string]string{
		"tkt-ts-gamma":  "ts-gamma",
		"ts-gamma":      "ts-gamma",
		"tkt:bug-delta": "bug-delta",
	} {
		if a := g.TicketByTarget(target); a == nil || a.ID != wantID {
			t.Errorf("TicketByTarget(%q) = %+v, want id %q", target, a, wantID)
		}
	}
	if a := g.TicketByTarget("tkt-ghost"); a != nil {
		t.Errorf("TicketByTarget(tkt-ghost) = %+v, want nil", a)
	}

	// Ticket ids are the available targets, sorted by identity.
	wantIDs := []string{"bug-delta", "ch-epsilon", "sto-alpha", "sto-alpha-dup", "sto-beta", "sto-beta-multi", "sto-legacy", "ts-gamma", "unresolved"}
	if got := g.TicketIDs(); !reflect.DeepEqual(got, wantIDs) {
		t.Errorf("TicketIDs() = %v, want %v", got, wantIDs)
	}
}

// TestGraphDedupByIdentity: two tickets deriving from the same container
// may reference the same work item; membership deduplicates by identity
// line so the work item appears exactly once. Regression: the original
// fixture carried one ticket per work item, so dedup was never exercised.
func TestGraphDedupByIdentity(t *testing.T) {
	g := loadFixture(t, "valid")
	alpha := validForm + "sto:alpha"

	a := g.ByLineForm(validForm + "tkt:sto-alpha")
	dup := g.ByLineForm(validForm + "tkt:sto-alpha-dup")
	if a == nil || dup == nil {
		t.Fatal("fixture must carry tkt:sto-alpha and tkt:sto-alpha-dup")
	}
	for name, tkt := range map[string]*conformance.Artifact{"tkt:sto-alpha": a, "tkt:sto-alpha-dup": dup} {
		if _, wi := g.ticketTargets(tkt); wi == nil || wi.Identity != alpha {
			t.Errorf("%s must resolve to sto:alpha, got %+v", name, wi)
		}
	}

	items := g.WorkItemsForContainer(validForm + "ctr:wave-1")
	n := 0
	for _, wi := range items {
		if wi.Identity == alpha {
			n++
		}
	}
	if n != 1 {
		t.Errorf("sto:alpha appears %d times, want exactly 1 (dedup by identity)", n)
	}
	if !reflect.DeepEqual(items, g.WorkItemsForContainer(validForm+"ctr:wave-1")) {
		t.Error("WorkItemsForContainer is not deterministic across repeated calls")
	}
}

// TestTicketTargetsFirstResolvableWins: a ticket with several work item
// references resolves deterministically to the FIRST resolvable one in
// file order. Regression: the original fixture's tickets each referenced
// a single work item.
func TestTicketTargetsFirstResolvableWins(t *testing.T) {
	g := loadFixture(t, "valid")
	tkt := g.ByLineForm(validForm + "tkt:sto-beta-multi")
	if tkt == nil {
		t.Fatal("fixture must carry tkt:sto-beta-multi")
	}
	container, workItem := g.ticketTargets(tkt)
	if container == nil || container.Identity != validForm+"ctr:wave-1" {
		t.Errorf("container = %+v, want ctr:wave-1", container)
	}
	if workItem == nil || workItem.Identity != validForm+"sto:beta" || workItem.State != "todo" {
		t.Errorf("work item = %+v, want sto:beta/todo (first resolvable reference wins)", workItem)
	}
	// Deterministic: repeated resolution picks the same targets.
	c2, w2 := g.ticketTargets(tkt)
	if !reflect.DeepEqual(container, c2) || !reflect.DeepEqual(workItem, w2) {
		t.Error("ticketTargets is not deterministic")
	}
	// The wave projection surfaces the same deterministic pick.
	p, err := Build("wave", g, "")
	if err != nil {
		t.Fatalf("Build(wave): %v", err)
	}
	wave := p.(*WaveProjection)
	for _, tkt := range wave.Tickets {
		if tkt.ID == "sto-beta-multi" && tkt.Projected != "todo" {
			t.Errorf("tkt:sto-beta-multi projects %q, want todo (sto:beta wins over ts:gamma)", tkt.Projected)
		}
	}
}

// TestFixtureConforms verifies every view fixture passes the conformance
// rules: the CLI gate refuses non-conformant repositories, so the
// fixtures must stay compliant.
func TestFixtureConforms(t *testing.T) {
	for _, name := range []string{"valid", "no-active", "multi-active"} {
		report, err := conformance.Validate(filepath.Join("testdata", name))
		if err != nil {
			t.Fatalf("validate fixture %s: %v", name, err)
		}
		if !report.Pass() {
			t.Errorf("fixture %s must be conformant, got %d errors:", name, report.ErrorCount())
			for _, res := range report.SortedResults() {
				t.Errorf("  [%s] %s %s: %s", res.Severity, res.Rule, res.File, res.Message)
			}
		}
	}
}

// TestGraphBuildDeterministic verifies that identical graphs produce
// identical model slices.
func TestGraphBuildDeterministic(t *testing.T) {
	a, b := loadFixture(t, "valid"), loadFixture(t, "valid")
	if !reflect.DeepEqual(a.Artifacts(), b.Artifacts()) {
		t.Error("Artifacts() differs between identical graphs")
	}
	if !reflect.DeepEqual(a.TicketsForContainer(validForm+"ctr:wave-1"), b.TicketsForContainer(validForm+"ctr:wave-1")) {
		t.Error("TicketsForContainer differs between identical graphs")
	}
	if !reflect.DeepEqual(a.WorkItemsForContainer(validForm+"ctr:wave-1"), b.WorkItemsForContainer(validForm+"ctr:wave-1")) {
		t.Error("WorkItemsForContainer differs between identical graphs")
	}
}

func TestGraphNoActiveContainer(t *testing.T) {
	g := loadFixture(t, "no-active")
	container, multiple := g.ActiveContainer()
	if container != nil {
		t.Errorf("ActiveContainer() = %+v, want nil", container)
	}
	if multiple {
		t.Error("no active container must not report multiple")
	}
	if got := g.TicketsForContainer(validForm + "ctr:wave-0"); len(got) != 1 {
		t.Errorf("completed container keeps its ticket: got %d, want 1", len(got))
	}
}

func TestGraphMultipleActiveContainers(t *testing.T) {
	g := loadFixture(t, "multi-active")
	container, multiple := g.ActiveContainer()
	if container == nil {
		t.Fatal("ActiveContainer() must pick a container")
	}
	if !multiple {
		t.Error("two active containers must report multiple")
	}
	// Lexicographically smallest canonical identity wins: wave-1 < wave-2.
	if container.ID != "wave-1" {
		t.Errorf("ActiveContainer() = %q, want the lexicographically smallest (wave-1)", container.ID)
	}
}
