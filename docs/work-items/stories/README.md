# Stories

User stories — work items that deliver user-facing value. Stories describe functionality from the user's perspective.

## What Goes Here

- User stories
- Acceptance criteria
- Story implementation notes

## What Does NOT Go Here

- Technical implementation details → `../technical-stories/`
- Bug reports → `../bugs/`
- Technical debt → `../tech-debt/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Story | `st-nnn-nnn-<slug>.md` | `st-012-001-user-login-flow.md` |

The double-number pattern (`st-nnn-nnn`) allows grouping stories under epics or themes.

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
| Product Owner | Defines story requirements |
| Engineers | Implement stories |
| QA | Verifies story acceptance criteria |

## Related Folders

- `../technical-stories/` — Technical implementation stories
- `../bugs/` — Bugs that may be found in stories
- `../../epics/` — Stories are grouped under epics
