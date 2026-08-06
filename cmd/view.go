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
// Exit codes:
//
//	0  projection produced (including empty projections: no active
//	   container, no tickets)
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
(R0-R9); a repository with blocking violations is refused and no
projection is produced. Warnings never block a projection.

Projections:

  sprint   the active execution container's work items, grouped by
           execution state (planned/todo/in-progress/in-review/done)
  wave     the active container's tickets and work item progress
  ticket   one ticket's projected status, derived from the referenced
           work item's execution state (ticket body content is never
           read)

The target argument is required by the ticket projection only
(a bare ticket id, tkt-<id> or tkt:<id>); sprint and wave ignore it.

With no arguments the available projections are listed.

Exit codes:
  0  projection produced
  1  repository validation failed (no projection produced)
  2  usage or internal error (unknown projection, missing or unknown
     ticket target)`,
		Example: `  eka view
  eka view sprint
  eka view wave
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
					name, strings.Join(view.Projections(), ", "))
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
						name, strings.Join(view.Projections(), ", "))
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
	"sprint": "active container sprint board (work items by execution state)",
	"wave":   "active container wave: tickets and work item progress",
	"ticket": "one ticket's projected status from its work item",
}

// printViewLanding renders the calm no-argument orientation: the
// available projections and usage pointers. Informational output —
// exits 0, deterministic.
func printViewLanding(s *ui.Style) {
	fmt.Fprintln(s.W, s.Accent("Knowledge Projections"))
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "  The EKA Knowledge Projection Engine: read-only views over the")
	fmt.Fprintln(s.W, "  Engineering Knowledge Model — repository artifacts and their")
	fmt.Fprintln(s.W, "  relationships, projected by state.")
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Projections")
	for _, name := range view.Projections() {
		fmt.Fprintf(s.W, "  %-10s %s\n", name, viewDescriptions[name])
	}
	fmt.Fprintln(s.W)
	fmt.Fprintln(s.W, "Usage")
	fmt.Fprintln(s.W, "  Run 'eka view <projection>' for a projection,")
	fmt.Fprintln(s.W, "  or 'eka view ticket <tkt-id>' for one ticket.")
	fmt.Fprintln(s.W, "  Run 'eka view <projection> --help' for details.")
}

// renderView dispatches to the concrete projection renderer. The
// registry is closed over the three projections; an unknown concrete
// type is a programming error, not user input.
func renderView(s *ui.Style, g *view.Graph, p view.Projection) {
	switch p := p.(type) {
	case *view.SprintProjection:
		renderSprint(s, g, p)
	case *view.WaveProjection:
		renderWave(s, g, p)
	case *view.TicketProjection:
		renderTicket(s, g, p)
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

func renderSprint(s *ui.Style, g *view.Graph, p *view.SprintProjection) {
	container := "none"
	if p.Container != nil {
		container = p.Container.Identity
	}
	ui.NewHeader(s, "Sprint").
		Add("Container", container).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v1").
		Pipeline("View").
		Render()
	if p.MultipleActive {
		fmt.Fprintf(s.W, "%s\n", s.Warning("Multiple active containers — showing "+p.Container.Identity))
	}
	if p.Container == nil {
		// Empty projection: a calm line, still exit 0 with the summary.
		fmt.Fprintf(s.W, "%s\n", s.Dim("No active container."))
	} else {
		renderStateColumns(s, p.Columns)
	}
	renderSprintSummary(s, p)
}

func renderSprintSummary(s *ui.Style, p *view.SprintProjection) {
	container, status := "none", "no active container"
	if p.Container != nil {
		container = p.Container.Identity
		status = p.Container.State
	}
	ui.NewSummary(s).
		Add("Container", container).
		Add("Work items", strconv.Itoa(p.Total)).
		Add("In progress", strconv.Itoa(p.Columns.Count("in-progress"))).
		Add("Done", strconv.Itoa(p.Columns.Count("done"))).
		Add("Tickets", strconv.Itoa(p.Tickets)).
		Add("Status", status).
		Render()
}

func renderWave(s *ui.Style, g *view.Graph, p *view.WaveProjection) {
	container := "none"
	if p.Container != nil {
		container = p.Container.Identity
	}
	ui.NewHeader(s, "Wave").
		Add("Container", container).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v1").
		Pipeline("View").
		Render()
	if p.MultipleActive {
		fmt.Fprintf(s.W, "%s\n", s.Warning("Multiple active containers — showing "+p.Container.Identity))
	}
	if p.Container == nil {
		fmt.Fprintf(s.W, "%s\n", s.Dim("No active container."))
	} else {
		fmt.Fprintf(s.W, "%s\n", s.Info(fmt.Sprintf("Tickets (%d)", len(p.Tickets))))
		for _, t := range p.Tickets {
			fmt.Fprintf(s.W, "  %s %s %s\n", stateMark(s, t.Projected), t.Identity,
				stateColor(s, t.Projected)("("+t.Projected+")"))
		}
		fmt.Fprintf(s.W, "\n%s\n", s.Info("Progress"))
		for _, col := range p.Columns {
			fmt.Fprintf(s.W, "  %-11s %d\n", stateColor(s, col.State)(col.State), len(col.WorkItems))
		}
	}
	renderWaveSummary(s, p)
}

func renderWaveSummary(s *ui.Style, p *view.WaveProjection) {
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

func renderTicket(s *ui.Style, g *view.Graph, p *view.TicketProjection) {
	ui.NewHeader(s, "Ticket").
		Add("Ticket", p.Ticket.Identity).
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v1").
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
