package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/maleolabs/engineering-knowledge-architecture/cmd/ui"
	"github.com/maleolabs/engineering-knowledge-architecture/view"
)

// This file implements the Board projection renderer: every work item
// in the repository — across all execution containers (active and
// completed) and outside any container — as a kanban board. Where the
// execution board answers "what is currently being worked on?", the
// board answers "what is the total work in the repository?". Each item
// keeps its container context as a tag; items without a referencing
// ticket container render as unassigned.

// renderBoardProjection renders the Board projection: the context
// header, the scope line, the five-column kanban board with per-item
// container tags, the unassigned warning when present, and the insight
// summary.
func renderBoardProjection(s *ui.Style, g *view.Graph, p *view.BoardProjection) {
	ui.NewHeader(s, "Board").
		Add("Container", "all").
		Add("Repository", g.Root()).
		Add("Knowledge", "EKA v"+standardVersion).
		Add("Domain", "Execution").
		Pipeline("View").
		Render()
	fmt.Fprintln(s.W)
	if p.Total == 0 {
		// Empty projection: a calm line, still exit 0 with the summary.
		fmt.Fprintf(s.W, "%s\n", s.Dim("No work items."))
	} else {
		fmt.Fprintf(s.W, "%s\n", s.Dim(plural(p.Total, "work item", "work items")+
			" across "+plural(p.ContainerCount, "container", "containers")))
	}
	fmt.Fprintln(s.W)
	renderBoardColumns(s, p.Columns)
	if p.Unassigned > 0 {
		fmt.Fprintln(s.W)
		fmt.Fprintf(s.W, "%s\n", s.Warning(plural(p.Unassigned, "work item", "work items")+
			" not referenced by any ticket container"))
	}
	renderBoardInsights(s, p)
}

// renderBoardColumns renders the work board: the fixed five
// execution-state columns with the short ids of their work items as
// cell labels, each tagged with its container context.
func renderBoardColumns(s *ui.Style, cols view.BoardColumns) {
	board := ui.NewBoard(s)
	var all []view.BoardItem
	for _, col := range cols {
		all = append(all, col.WorkItems...)
	}
	short := shortWorkItemID(boardItemWorkItems(all))
	shortCtr := shortContainerID(all)
	for _, col := range cols {
		labels := make([]string, 0, len(col.WorkItems))
		for _, bi := range col.WorkItems {
			labels = append(labels, boardCellLabel(short(bi.WorkItem), shortCtr(bi)))
		}
		board.AddColumn(boardTitle(col.State), stateColor(s, col.State), labels)
	}
	board.Render()
}

// boardCellLabel composes an item label with its container tag and
// truncates it to the board cell budget with the tag kept intact:
// "markdown-s… (wave-7)" — the container context is the point of the
// board and must never be truncated away. The ellipsis ends the id.
func boardCellLabel(id, tag string) string {
	label := id + " (" + tag + ")"
	if utf8.RuneCountInString(label) <= ui.BoardMaxItemWidth {
		return label
	}
	tagPart := " (" + tag + ")"
	remain := ui.BoardMaxItemWidth - utf8.RuneCountInString(tagPart) - 1 // ellipsis
	if remain < 1 {
		return tagPart
	}
	runes := []rune(id)
	if remain >= len(runes) {
		return label
	}
	return string(runes[:remain]) + "…" + tagPart
}

// boardItemWorkItems extracts the embedded work items of a board item
// set, preserving order (for the shared short-id ambiguity logic).
func boardItemWorkItems(items []view.BoardItem) []view.WorkItem {
	out := make([]view.WorkItem, len(items))
	for i, bi := range items {
		out[i] = bi.WorkItem
	}
	return out
}

// shortContainerID renders the container tag of a board item: the bare
// container id ("wave-7"), or the full canonical identity when the id
// is ambiguous across distinct containers; "(unassigned)" when no
// container references the item. Multiple containers join
// comma-separated.
func shortContainerID(items []view.BoardItem) func(view.BoardItem) string {
	// id → the distinct full identities sharing that id. A bare id is
	// ambiguous only when TWO DIFFERENT containers share it — not when
	// one container references many items (the tag frequency).
	idToForms := make(map[string][]string)
	for _, bi := range items {
		for _, c := range bi.Containers {
			id := shortID(c)
			if !stringSliceContains(idToForms[id], c) {
				idToForms[id] = append(idToForms[id], c)
			}
		}
	}
	return func(bi view.BoardItem) string {
		if len(bi.Containers) == 0 {
			return "unassigned"
		}
		parts := make([]string, 0, len(bi.Containers))
		for _, c := range bi.Containers {
			id := shortID(c)
			if len(idToForms[id]) > 1 {
				parts = append(parts, c)
			} else {
				parts = append(parts, id)
			}
		}
		return strings.Join(parts, ", ")
	}
}

// stringSliceContains reports whether s contains v.
func stringSliceContains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

// shortID renders the bare id of a canonical identity form
// ("<namespace>/<type>:<id>" → "<id>"). Ids never contain colons
// (identity contract), so the last colon is the separator.
func shortID(form string) string {
	if i := strings.LastIndex(form, ":"); i >= 0 {
		return form[i+1:]
	}
	return form
}

// renderBoardInsights renders the board summary: total work, active
// work (in progress + in review), completed work, the review queue,
// unassigned items, and overall progress.
func renderBoardInsights(s *ui.Style, p *view.BoardProjection) {
	inProgress := p.Columns.Count("in-progress")
	inReview := p.Columns.Count("in-review")
	done := p.Columns.Count("done")
	percent := "0%"
	if p.Total > 0 {
		percent = strconv.Itoa(done*100/p.Total) + "%"
	}
	progress := ui.ProgressBar(s, done, p.Total) + " " + fmt.Sprintf("%d/%d (%s)", done, p.Total, percent)
	ui.NewSummary(s).
		Add("Total Work Items", strconv.Itoa(p.Total)).
		Add("Active Work", strconv.Itoa(inProgress+inReview)).
		Add("Completed Work", strconv.Itoa(done)).
		Add("Review Queue", strconv.Itoa(inReview)).
		Add("Unassigned", strconv.Itoa(p.Unassigned)).
		Add("Overall Progress", progress).
		Render()
}
