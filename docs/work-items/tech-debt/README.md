# Tech Debt

Technical debt items — work items that address accumulated technical compromises. Tech debt items track shortcuts, deprecated patterns, and refactoring needs.

## What Goes Here

- Technical debt documentation
- Refactoring plans
- Deprecated code tracking
- Migration work items

## What Does NOT Go Here

- Bug reports → `../bugs/`
- Feature stories → `../stories/`
- Maintenance chores → `../chores/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Tech Debt | `td-nnn-<slug>.md` | `td-002-remove-legacy-auth.md` |

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
| Engineers | Identify and document tech debt |
| Tech Lead | Prioritizes tech debt items |
| Engineers | Resolve tech debt items |

## Related Folders

- `../chores/` — Maintenance tasks
- `../stories/` — Feature work that may introduce tech debt
- `../../adr/` — Architecture decisions that may create tech debt
