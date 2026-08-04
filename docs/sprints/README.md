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

## Work Item Status

Every work item in the Sprint's Work Items table uses one of the following statuses. No other status is allowed.

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when the sprint is generated from the roadmap | Sprint generation (initial state) |
| `Todo` | Will be worked on — the item has entered an execution wave | Item's wave becomes active in the ticket document |
| `In Progress` | Implementation is underway | Assignee starts working on the item |
| `In Review` | Pull Request has been created | PR opened for the item's implementation |
| `Done` | Implementation is merged | PR merged to the target branch |

### Status Rules

- Status transitions are strictly sequential: `Planned` → `Todo` → `In Progress` → `In Review` → `Done`.
- Never skip a status. Never revert to a previous status.
- `Planned` is the only valid initial status — every work item starts here.
- `Done` is the only valid terminal status — there is no status beyond `Done`.
- The sprint's Work Items table is the single source of truth for work item status.
- Status changes are recorded in the Sprint's Change Log.

## Status Lifecycle

```
Active → Completed
```

- **Active**: Sprint is in progress; tickets are being worked on.
- **Completed**: All tickets are Done; sprint is closed.

## Active Sprint Rule

**ONLY ONE active sprint at a time.**

- A sprint becomes active when generated from the roadmap.
- The next sprint CANNOT be generated until the active sprint is completed (all work items Done).
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
