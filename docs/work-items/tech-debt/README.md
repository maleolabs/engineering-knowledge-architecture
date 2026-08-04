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

## Status Lifecycle

```
Draft → Ready → In Progress → In Review → Verified → Done
```

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
