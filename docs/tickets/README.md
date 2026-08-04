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

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Generates tickets from sprints |
| Engineers | Execute tickets |

## Related Folders

- `../sprints/` — Tickets are generated from sprints
- `../work-items/` — Work items are referenced by tickets
