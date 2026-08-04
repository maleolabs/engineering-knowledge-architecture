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

## Status Lifecycle

```
Draft → Ready → In Progress → In Review → Verified → Done
```

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
