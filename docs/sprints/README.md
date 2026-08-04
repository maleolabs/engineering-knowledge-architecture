# Sprints

Sprint documents — execution snapshots from the latest roadmap. Only one sprint can be active at a time. Sprints are generated from roadmaps and drive ticket execution.

## What Goes Here

- Sprint execution documents
- Sprint change logs
- Sprint completion records

## What Does NOT Go Here

- Roadmap planning → `../roadmap/`
- Wave/ticket decomposition → `../tickets/`
- Individual work items → `../work-items/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Sprint | `mvp-nnn-s<nn>.md` | `mvp-001-s01.md` |

## Status Lifecycle

```
Active → Completed
```

- **Active**: Sprint is in progress; tickets are being worked on.
- **Completed**: All tickets are Verified; sprint is closed.

## Active Sprint Rule

**ONLY ONE active sprint at a time.**

- A sprint becomes active when generated from the roadmap.
- The next sprint CANNOT be generated until the active sprint is completed (all tickets Verified).
- Sprint completion is recorded in the sprint's Change Log.
- If a sprint is partially completed, deferred items roll to the next sprint.

## Ownership

| Role | Responsibility |
|---|---|
| Scrum Master | Manages sprint execution |
| Tech Lead | Generates sprints from roadmap |
| Engineers | Execute sprint work items |

## Related Folders

- `../roadmap/` — Sprints are generated from roadmaps
- `../tickets/` — Tickets are generated from sprints
- `../work-items/` — Work items are the implementation units
