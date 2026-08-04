# Tickets

Wave/ticket execution documents — deterministic wave decomposition of active sprints. Tickets are consumed by execution commands and track the granular execution of sprint work.

## What Goes Here

- Wave ticket documents
- Ticket execution tracking
- Deterministic decomposition of sprint work

## What Does NOT Go Here

- Sprint definitions → `../sprints/`
- Individual work items → `../work-items/`
- Roadmap planning → `../roadmap/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Wave tickets | `<sprint-id>-wave-tickets.md` | `mvp-001-s01-wave-tickets.md` |

## Ticket Status

Every ticket uses one of the following statuses. No other status is allowed. Ticket status mirrors the sprint work item status — when a work item changes status, its corresponding ticket changes too.

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when the ticket document is generated | Ticket document generation (initial state) |
| `Todo` | Will be worked on — the ticket has entered an execution wave | Wave becomes active |
| `In Progress` | Implementation is underway | Assignee starts working on the ticket |
| `In Review` | Pull Request has been created | PR opened for the ticket's implementation |
| `Done` | Implementation is merged | PR merged to the target branch |

### Status Rules

- Status transitions are strictly sequential: `Planned` → `Todo` → `In Progress` → `In Review` → `Done`.
- Never skip a status. Never revert to a previous status.
- `Planned` is the only valid initial status.
- `Done` is the only valid terminal status.
- Ticket status and its corresponding sprint work item status must always agree.
- Status changes are recorded in the Ticket document's Change Log.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Generates tickets from sprints |
| Engineers | Execute tickets |

## Related Folders

- `../sprints/` — Tickets are generated from sprints
- `../work-items/` — Work items are referenced by tickets
