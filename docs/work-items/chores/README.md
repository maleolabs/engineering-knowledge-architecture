# Chores

Maintenance and operational tasks. Chores are work items that don't add user-facing features but keep the project healthy.

## What Goes Here

- Dependency updates
- Build process maintenance
- Configuration changes
- Operational tasks

## What Does NOT Go Here

- Feature work → `../stories/`
- Bug fixes → `../bugs/`
- Technical debt → `../tech-debt/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Chore | `ch-nnn-nnn-<slug>.md` | `ch-001-001-update-dependencies.md` |

## Status

Every work item uses one of the following statuses. No other status is allowed.

```
Planned → Todo → In Progress → In Review → Done
```

| Status | Meaning | Condition |
|---|---|---|
| `Planned` | Default status when assigned to a sprint | Sprint generation (initial state) |
| `Todo` | Will be worked on — the item has entered an execution wave | Wave becomes active |
| `In Progress` | Implementation is underway | Assignee starts working |
| `In Review` | Pull Request has been created | PR opened |
| `Done` | Implementation is merged | PR merged to the target branch |

Status transitions are strictly sequential. Never skip a status. Never revert to a previous status.

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | Execute chores |
| Tech Lead | Prioritizes chores |

## Related Folders

- `../tech-debt/` — Technical debt items
- `../operations/` — Operational references
