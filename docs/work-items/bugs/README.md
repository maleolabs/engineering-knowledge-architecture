# Bugs

Bug reports and fixes. Each bug document describes a defect, its impact, and the resolution.

## What Goes Here

- Bug reports
- Bug fix documentation
- Root cause analysis

## What Does NOT Go Here

- Technical debt → `../tech-debt/`
- Feature requests → `../stories/`
- Chores → `../chores/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Bug | `bug-nnn-<slug>.md` | `bug-001-description-field-null.md` |

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
| Engineers | Report and fix bugs |
| QA | Verifies bug fixes |
| Tech Lead | Prioritizes bugs |

## Related Folders

- `../stories/` — User stories that may introduce bugs
- `../tech-debt/` — Technical debt that may cause bugs
