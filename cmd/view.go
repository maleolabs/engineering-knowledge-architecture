package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/conformance"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
	"github.com/spf13/cobra"
)

// newViewCommand builds the `eka view` command: project the Engineering
// Knowledge Model of the repository rooted at the current directory.
// All projection logic lives in the view package (the Knowledge
// Projection Engine); this command only validates arguments, runs the
// conformance gate, renders the projection and maps the result to the
// exit code contract.
//
// The projections are domain-first: discovery, architecture, planning,
// execution and operations render one Engineering Domain each; ticket
// renders a single ticket. The former sprint and wave projections
// remain registered as aliases of execution (identical output).
//
// Exit codes:
//
//	0  projection produced (including empty projections: no active
//	   container, no domain artifacts, no tickets)
//	1  repository validation failed: no projection is produced (the
//	   full report is printed)
//	2  usage or internal error (unknown projection, missing or unknown
//	   ticket target, unreadable root)
func newViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "view [projection] [target]",
		Short: "Project the Engineering Knowledge Model",
		Long: `Project the Engineering Knowledge Model of the repository rooted at the
current directory: read-only views over the repository's artifacts and
their relationships, derived from the model — never from file text.

The repository is validated against the conformance rules first
(R0-R12); a repository with blocking violations is refused and no
projection is produced. Warnings never block a projection.

Projections (domain-first):

  discovery    the Discovery domain: vis-, str-, req-, fnd- artifacts
               grouped by type with their content states
  architecture the Architecture domain: adr-, dec-, arc-, spec-, std-,
               gls- artifacts grouped by type with their content states
               (Decisions merge adr- and dec-)
  planning     the Planning domain: scp-, epc-, plan-, trc- artifacts
               grouped by type with content state, planning state and
               phase context
  execution    the active execution container: its tickets with the
               status projected from their work items, and its work
               items grouped by execution state
               (planned/todo/in-progress/in-review/done)
  operations   the Operations domain: run-, rel- artifacts grouped by
               type with their content states
  ticket       one ticket's projected status, derived from the
               referenced work item's execution state (ticket body
               content is never read)

Aliases:

  sprint, wave resolve to the execution projection (identical output)

The target argument is required by the ticket projection only
(a bare ticket id, tkt-<id> or tkt:<id>); the domain and execution
projections ignore it.

With no arguments the available projections are listed.

Exit codes:
  0  projection produced
  1  repository validation failed (no projection produced)
  2  usage or internal error (unknown projection, missing or unknown
     ticket target)`,
		Example: `  eka view
  eka view execution
  eka view planning
  eka view ticket tkt-sto-alpha
  eka view ticket sto-alpha`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				printViewLanding(styleFor(cmd))
				return nil
			}
			name := args[0]
			if !view.IsProjection(name) {
				return fmt.Errorf("unknown projection %q — available projections: %s",
					name, view.HelpList())
			}
			target := ""
			if len(args) == 2 {
				target = args[1]
			}
			if name == "ticket" && target == "" {
				return fmt.Errorf("the ticket projection requires a target: eka view ticket <tkt-id>")
			}
			s := styleFor(cmd)
			// Validation gate FIRST: only conformant repositories may be
			// projected. Blocking violations print the full report and
			// exit 1 — no projection is produced.
			report, err := conformance.Validate(".")
			if err != nil {
				return fmt.Errorf("view failed: %w", err)
			}
			if !report.Pass() {
				printReport(s, report)
				fmt.Fprintf(cmd.ErrOrStderr(), "eka: view refused: repository is not conformant\n")
				return &exitError{code: exitFail}
			}
			// One scan, one graph: the projection engine is synchronous
			// and stateless, so a future loading state can wrap the
			// whole call without restructuring.
			artifacts, err := conformance.Scan(".")
			if err != nil {
				return fmt.Errorf("view failed: %w", err)
			}
			g := view.NewGraph(".", artifacts)
			proj, err := view.Build(name, g, target)
			if err != nil {
				if errors.Is(err, view.ErrUnknownProjection) {
					return fmt.Errorf("unknown projection %q — available projections: %s",
						name, view.HelpList())
				}
				return err // TargetNotFoundError etc. map to exit 2.
			}
			renderView(s, g, proj)
			return nil
		},
	}
}

// viewDescriptions are the one-line projection descriptions used by the
// no-argument landing.
var viewDescriptions = map[string]string{
	"discovery":    "Discovery domain artifacts (vis-, str-, req-, fnd-)",
	"architecture": "Architecture domain artifacts (adr-, dec-, arc-, spec-, std-, gls-)",
	"planning":     "Planning domain artifacts (scp-, epc-, plan-, trc-)",
	"execution":    "active container: tickets and work items by execution state",
	"operations":   "Operations domain artifacts (run-, rel-)",
	"ticket":       "one ticket's projected status from its work item",
}

// printViewLanding renders the calm no-argument orientation: the
// available projections (canonical + aliases) and usage pointers.
// Informational output — exits 0, deterministic.
func printViewLanding(s *ui.Style) {
	fmt.Fprintln(s.W, s.Accent("Knowledge Projections"))
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "  The EKA Knowledge Projection Engine: read-only views over the")
	fmt.Fprintln(s.W, "  Engineering Knowledge Model — repository artifacts and their")
	fmt.Fprintln(s.W, "  relationships, projected by domain and state.")
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Projections")
	for _, name := range view.Projections() {
		fmt.Fprintf(s.W, "  %-12s %s\n", name, viewDescriptions[name])
	}
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Aliases")
	for _, alias := range view.Aliases() {
		fmt.Fprintf(s.W, "  %-12s → %s\n", alias, view.AliasTarget(alias))
	}
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Usage")
	fmt.Fprintln(s.W, "  Run 'eka view <projection>' for a projection,")
	fmt.Fprintln(s.W, "  or 'eka view ticket <tkt-id>' for one ticket.")
	fmt.Fprintln(s.W, "  Run 'eka view <projection> --help' for details.")
}

// renderView dispatches to the concrete projection renderer. The
// registry is closed over the six canonical projections; an unknown
// concrete type is a programming error, not user input.
func renderView(s *ui.Style, g *view.Graph, p view.Projection) {
	switch p := p.(type) {
	case *view.ExecutionProjection:
		renderExecution(s, g, p)
	case *view.TicketProjection:
		renderTicket(s, g, p)
	case *view.PlanningProjection:
		renderPlanning(s, g, p)
	case *view.ArchitectureProjection:
		renderArchitecture(s, g, p)
	case *view.DiscoveryProjection:
		renderDiscovery(s, g, p)
	case *view.OperationsProjection:
		renderOperations(s, g, p)
	default:
		fmt.Fprintln(s.W, s.Error("cannot render projection"))
	}
}

// stateColor returns the presentation color of an execution state
// value: planned dim, todo info, in-progress progress, in-review
// warning, done success. "unresolved" reads as a warning (amber).
func stateColor(s *ui.Style, state string) func(string) string {
	switch state {
	case "planned":
		return s.Dim
	case "todo":
		return s.Info
	case "in-progress":
		return s.Progress
	case "in-review":
		return s.Warning
	case "done":
		return s.Success
	case "unresolved":
		return s.Warning
	default:
		return s.Dim
	}
}

// contentStateColor returns the presentation color of a content-state
// value: draft dim, review info, approved success, amended warning,
// proposed info, accepted success, superseded warning.
func contentStateColor(s *ui.Style, state string) func(string) string {
	switch state {
	case "draft":
		return s.Dim
	case "review":
		return s.Info
	case "approved":
		return s.Success
	case "amended":
		return s.Warning
	case "proposed":
		return s.Info
	case "accepted":
		return s.Success
	case "superseded":
		return s.Warning
	default:
		return s.Dim
	}
}

// planningStateColor returns the presentation color of a planning-state
// value: draft dim, approved success, immutable warning.
func planningStateColor(s *ui.Style, state string) func(string) string {
	switch state {
	case "draft":
		return s.Dim
	case "approved":
		return s.Success
	case "immutable":
		return s.Warning
	default:
		return s.Dim
	}
}

// stateIcon returns the icon of an execution state value: ✓ done,
// → in progress, • everything else. Icons decorate; the state word
// carries the meaning.
func stateIcon(state string) string {
	switch state {
	case "done":
		return ui.IconDone
	case "in-progress":
		return ui.IconArrow
	default:
		return ui.IconBullet
	}
}

// stateMark renders the colored state icon.
func stateMark(s *ui.Style, state string) string {
	return stateColor(s, state)(stateIcon(state))
}

// renderStateColumns renders the fixed execution-state columns: one
// heading line per state (name + count, colored by state) followed by
// the plain identity of each work item with its colored icon. Empty
// columns render their heading only, keeping the shape stable.
func renderStateColumns(s *ui.Style, cols view.StateColumns) {
	for _, col := range cols {
		fmt.Fprintf(s.W, "%s\n", stateColor(s, col.State)(fmt.Sprintf("%s (%d)", col.State, len(col.WorkItems))))
		for _, wi := range col.WorkItems {
			fmt.Fprintf(s.W, "  %s %s\n", stateMark(s, wi.State), wi.Identity)
		}
	}
}

func renderExecution(s *ui.Style, g *view.Graph, p *view.ExecutionProjection) {
	container := "none"
	if p.Container != nil {
		container = p.Container.Identity
	}
	ui.NewHeader(s, "Execution").
		Add("Container", container).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v"+standardVersion).
		Add("Domain", "Execution").
		Pipeline("View").
		Render()
	if p.MultipleActive {
		fmt.Fprintf(s.W, "%s\n", s.Warning("Multiple active containers — showing "+p.Container.Identity))
	}
	if p.Container == nil {
		// Empty projection: a calm line, still exit 0 with the summary.
		fmt.Fprintf(s.W, "%s\n", s.Dim("No active container."))
	} else {
		fmt.Fprintf(s.W, "%s\n", s.Info(fmt.Sprintf("Tickets (%d)", len(p.Tickets))))
		for _, t := range p.Tickets {
			fmt.Fprintf(s.W, "  %s %s %s\n", stateMark(s, t.Projected), t.Identity,
				stateColor(s, t.Projected)("("+t.Projected+")"))
		}
		fmt.Fprintln(s.W)
		renderStateColumns(s, p.Columns)
	}
	renderExecutionSummary(s, p)
}

func renderExecutionSummary(s *ui.Style, p *view.ExecutionProjection) {
	container, status := "none", "no active container"
	if p.Container != nil {
		container = p.Container.Identity
		status = p.Container.State
	}
	ui.NewSummary(s).
		Add("Container", container).
		Add("Tickets", strconv.Itoa(len(p.Tickets))).
		Add("Work items", strconv.Itoa(p.Total)).
		Add("In progress", strconv.Itoa(p.Columns.Count("in-progress"))).
		Add("Done", strconv.Itoa(p.Columns.Count("done"))).
		Add("Status", status).
		Render()
}

// renderDomainArtifact renders one artifact line of a domain group:
// identity plus the state values relevant to the group, colored by
// state value. The part order is fixed: content-state, planning-state,
// phase — never map iteration.
func renderDomainArtifact(s *ui.Style, a view.DomainArtifact) {
	var parts []string
	if a.HasContentState {
		parts = append(parts, contentStateColor(s, a.ContentState)(a.ContentState))
	}
	if a.HasPlanningState {
		parts = append(parts, "planning-state "+planningStateColor(s, a.PlanningState)(a.PlanningState))
	}
	if a.HasPhase {
		parts = append(parts, "phase "+a.Phase)
	}
	line := "  " + ui.IconBullet + " " + a.Identity
	if len(parts) > 0 {
		line += "  (" + strings.Join(parts, ", ") + ")"
	}
	fmt.Fprintln(s.W, line)
}

// renderDomainProjection renders one domain projection: the header
// (Domain row = the projection's domain), the artifact groups in fixed
// order, the calm empty-domain line, and the summary. Empty groups are
// skipped; a domain with no artifacts at all renders a single calm line
// and still exits 0.
func renderDomainProjection(s *ui.Style, g *view.Graph, domain string, groups []view.Group, renderSummary func()) {
	ui.NewHeader(s, domain).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v"+standardVersion).
		Add("Domain", domain).
		Pipeline("View").
		Render()
	if view.GroupTotal(groups) == 0 {
		fmt.Fprintf(s.W, "%s\n", s.Dim("No "+domain+" artifacts."))
	} else {
		for _, gr := range groups {
			if len(gr.Artifacts) == 0 {
				continue
			}
			fmt.Fprintf(s.W, "%s\n", s.Info(gr.Name))
			for _, a := range gr.Artifacts {
				renderDomainArtifact(s, a)
			}
			fmt.Fprintln(s.W)
		}
	}
	renderSummary()
}

// renderGroupSummary renders the per-group artifact counts.
func renderGroupSummary(s *ui.Style, groups []view.Group) {
	sm := ui.NewSummary(s)
	for _, gr := range groups {
		sm.Add(gr.Name, strconv.Itoa(len(gr.Artifacts)))
	}
	sm.Render()
}

func renderPlanning(s *ui.Style, g *view.Graph, p *view.PlanningProjection) {
	renderDomainProjection(s, g, "Planning", p.Groups, func() {
		byState := make([]string, 0, len(p.PlansByState))
		for _, sc := range p.PlansByState {
			byState = append(byState, fmt.Sprintf("%s %d", sc.State, sc.Count))
		}
		sm := ui.NewSummary(s)
		for _, gr := range p.Groups {
			sm.Add(gr.Name, strconv.Itoa(len(gr.Artifacts)))
		}
		sm.Add("Plans by state", strings.Join(byState, ", "))
		sm.Render()
	})
}

func renderArchitecture(s *ui.Style, g *view.Graph, p *view.ArchitectureProjection) {
	renderDomainProjection(s, g, "Architecture", p.Groups, func() {
		renderGroupSummary(s, p.Groups)
	})
}

func renderDiscovery(s *ui.Style, g *view.Graph, p *view.DiscoveryProjection) {
	renderDomainProjection(s, g, "Discovery", p.Groups, func() {
		renderGroupSummary(s, p.Groups)
	})
}

func renderOperations(s *ui.Style, g *view.Graph, p *view.OperationsProjection) {
	renderDomainProjection(s, g, "Operations", p.Groups, func() {
		renderGroupSummary(s, p.Groups)
	})
}

func renderTicket(s *ui.Style, g *view.Graph, p *view.TicketProjection) {
	ui.NewHeader(s, "Ticket").
		Add("Ticket", p.Ticket.Identity).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v"+standardVersion).
		Add("Domain", "Execution").
		Pipeline("View").
		Render()

	workItem := "unresolved"
	if p.WorkItem != nil {
		workItem = p.WorkItem.Identity + " " + stateColor(s, p.WorkItem.State)("("+p.WorkItem.State+")")
	}
	container := "unresolved"
	if p.Container != nil {
		container = p.Container.Identity
	}
	derives := "—"
	if len(p.References) > 0 {
		derives = strings.Join(p.References, ", ")
	}
	rows := [][2]string{
		{"Projected Status", stateColor(s, p.Projected)(p.Projected)},
		{"Work Item", workItem},
		{"Container", container},
		{"Derives From", derives},
	}
	width := 0
	for _, r := range rows {
		if len(r[0]) > width {
			width = len(r[0])
		}
	}
	for _, r := range rows {
		fmt.Fprintf(s.W, "%-*s   %s\n", width, r[0], r[1])
	}

	workItemValue := "unresolved"
	if p.WorkItem != nil {
		workItemValue = p.WorkItem.Identity + " (" + p.WorkItem.State + ")"
	}
	ui.NewSummary(s).
		Add("Ticket", p.Ticket.Identity).
		Add("Work item", workItemValue).
		Add("Container", container).
		Add("Status", p.Projected).
		Render()
}
