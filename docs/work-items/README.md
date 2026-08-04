# Work Items

Work item definitions — the implementation units within sprints. Work items are organized by type in subfolders. Each type has its own naming convention and lifecycle.

## Subfolders

| Subfolder | Purpose | File Naming | Example |
|---|---|---|---|
| `bugs/` | Bug reports and fixes | `bug-nnn-<slug>.md` | `bug-001-description-field-null.md` |
| `stories/` | User stories | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |
| `technical-stories/` | Technical implementation stories | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |
| `tech-debt/` | Technical debt items | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |
| `chores/` | Maintenance and operational tasks | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |
| `spikes/` | Research and investigation tasks | `sp-nnn-<slug>.md` | `sp-001-evaluate-cache-strategies.md` |
| `planning/` | Planning-related work items | Descriptive kebab-case | `planning-epic-breakdown.md` |

## General Naming Rules

1. Work item files use their ID prefix in **lowercase**.
2. All filenames use **kebab-case**.
3. The ID prefix identifies the work item type.
4. Sequential numbering within each type.

## Work Item Status

Every work item uses one of the following statuses. No other status is allowed.

```
Planned → Todo → In Progress → In Review → Done
```

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when assigned to a sprint | Sprint generation (initial state) |
| `Todo` | Will be worked on — the item has entered an execution wave | Item's wave becomes active |
| `In Progress` | Implementation is underway | Assignee starts working |
| `In Review` | Pull Request has been created | PR opened |
| `Done` | Implementation is merged | PR merged to the target branch |

### Status Rules

- Status transitions are strictly sequential: `Planned` → `Todo` → `In Progress` → `In Review` → `Done`.
- Never skip a status. Never revert to a previous status.
- `Planned` is the only valid initial status.
- `Done` is the only valid terminal status.
- Work item status, sprint table status, and ticket status must always agree.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | Creates and prioritizes work items |
| Engineers | Implement and update work items |
| QA | Verifies completed work items |

## Related Folders

- `../sprints/` — Work items are assigned to sprints
- `../tickets/` — Tickets reference work items
- `../sessions/` — Implementation sessions work on work items
- `../reviews/` — Reviews validate completed work items
