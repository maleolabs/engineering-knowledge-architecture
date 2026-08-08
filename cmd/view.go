package cmd

import (
	"errors"
	"fmt"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/compile"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
	"github.com/spf13/cobra"
)

// newViewCommand builds the `eka view` command: project the Engineering
// Knowledge Model of the repository rooted at the current directory.
// All projection logic lives in the view package (the Knowledge
// Projection Engine); this command only validates arguments, compiles
// the repository through the Knowledge Compiler (conformance gate +
// CKO assembly), renders the projection and maps the result to the
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

The repository is compiled from the authoring tree via the Knowledge
Compiler first (conformance-gated: the authoring rules R0-R12 run
before any projection); a repository with blocking violations is
refused and no projection is produced. Warnings never block a
projection.

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
  board        every work item in the repository across all execution
               containers, grouped by execution state, each item tagged
               with its container (unassigned when none)
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
			name, target, err := parseProjectionArgs(args)
			if err != nil {
				return err
			}
			s := styleFor(cmd)
			// Compile FIRST: the Knowledge Compiler runs the authoring
			// conformance gate and assembles the Canonical Knowledge
			// Objects. Blocking violations print the full report and
			// exit 1 — no projection is produced.
			res, err := compile.Compile(".")
			if err != nil {
				var ve *compile.ValidationError
				if errors.As(err, &ve) {
					printReport(s, ve.Report)
					fmt.Fprintf(cmd.ErrOrStderr(), "eka: view refused: repository is not conformant\n")
					return &exitError{code: exitFail}
				}
				return fmt.Errorf("view failed: %w", err)
			}
			// One compile, one graph: the projection engine is
			// synchronous and stateless, so a future loading state can
			// wrap the whole call without restructuring.
			g := view.NewGraph(".", res.CKOs)
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

// parseProjectionArgs validates the projection+target argument pair
// shared by view and watch: the projection must be registered
// (canonical or alias) and the ticket projection requires its target.
// Errors are usage-class (exit 2) with the same helpful messages in
// both commands. It assumes args is non-empty (both commands guard
// the no-argument case before calling).
func parseProjectionArgs(args []string) (name, target string, err error) {
	name = args[0]
	if !view.IsProjection(name) {
		return "", "", fmt.Errorf("unknown projection %q — available projections: %s",
			name, view.HelpList())
	}
	if len(args) == 2 {
		target = args[1]
	}
	if name == "ticket" && target == "" {
		return "", "", fmt.Errorf("the ticket projection requires a target: eka view ticket <tkt-id>")
	}
	return name, target, nil
}

// viewDescriptions are the one-line projection descriptions used by the
// no-argument landing.
var viewDescriptions = map[string]string{
	"discovery":    "Discovery domain artifacts (vis-, str-, req-, fnd-)",
	"architecture": "Architecture domain artifacts (adr-, dec-, arc-, spec-, std-, gls-)",
	"planning":     "Planning domain artifacts (scp-, epc-, plan-, trc-)",
	"execution":    "active container: tickets and work items by execution state",
	"board":        "all work items across every container, by execution state",
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
	case *view.BoardProjection:
		renderBoardProjection(s, g, p)
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
