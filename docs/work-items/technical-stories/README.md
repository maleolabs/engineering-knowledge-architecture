# Technical Stories

Technical implementation stories — work items that deliver technical value without direct user-facing impact. Technical stories describe implementation details, infrastructure changes, and system improvements.

## What Goes Here

- Technical implementation stories
- Infrastructure work items
- System improvement stories
- Non-functional requirement implementations

## What Does NOT Go Here

- User-facing stories → `../stories/`
- Bug reports → `../bugs/`
- Technical debt → `../tech-debt/`
- Research spikes → `../spikes/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Technical Story | `ts-nnn-nnn-<slug>.md` | `ts-012-001-database-migration.md` |

The double-number pattern (`ts-nnn-nnn`) allows grouping technical stories under epics or themes.

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
| Tech Lead | Defines technical stories |
| Engineers | Implement technical stories |
| QA | Verifies technical acceptance criteria |

## Related Folders

- `../stories/` — User stories that may depend on technical stories
- `../spikes/` — Spikes that may produce technical stories
- `../../architecture/` — Architecture docs that inform technical stories
