# MVP (Minimum Viable Product Definitions)

Minimum Viable Product definitions establish scope boundaries per release milestone. MVPs translate PRD requirements into achievable release increments.

## What Goes Here

- MVP scope definitions
- MVP amendments (changes to approved MVPs)
- Release milestone boundaries

## What Does NOT Go Here

- Product requirements → `../prd/`
- Epic breakdowns → `../epics/`
- Sprint planning → `../roadmap/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| MVP | `mvp-nnn-<product-name>-<version>.md` | `mvp-001-anvil-v1.md` |
| Amendment | `mvp-nnn-amendment-nn-<slug>.md` | `mvp-001-amendment-01-scope-adjustment.md` |

## Status Lifecycle

```
Draft → Review → Approved → Amended
```

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Defines MVP scope |
| Tech Lead | Validates technical feasibility |
| Stakeholders | Approve MVP boundaries |

## Related Folders

- `../prd/` — PRDs define what the product must do
- `../epics/` — Epics break MVPs into capability areas
- `../roadmap/` — Roadmaps plan work within an MVP
