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

## Status Lifecycle

```
Draft → Ready → In Progress → In Review → Verified → Done
```

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
