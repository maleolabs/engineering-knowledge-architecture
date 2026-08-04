# PRD (Product Requirements Documents)

Product Requirements Documents define what the product must do and why. PRDs are the starting point of the documentation workflow.

## What Goes Here

- Product requirement documents
- PRD amendments (changes to approved PRDs)

## What Does NOT Go Here

- MVP definitions → `../mvp/`
- Architecture decisions → `../adr/`
- Implementation details → `../work-items/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| PRD | `prd-nnn-<product-name>.md` | `prd-001-anvil.md` |
| Amendment | `prd-nnn-amendment-nn-<slug>.md` | `prd-001-amendment-01-scope-change.md` |

## Status Lifecycle

```
Draft → Review → Approved → Amended
```

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | Creates and maintains PRDs |
| Tech Lead | Reviews technical feasibility |
| Stakeholders | Approve PRD content |

## Related Folders

- `../mvp/` — MVP definitions derive from PRDs
- `../epics/` — Epics break down PRD requirements
- `../manifesto/` — Product vision and principles
