# Reviews

Review documents — code reviews, architecture reviews, and audit findings. Reviews validate completed work and capture feedback.

## What Goes Here

- Code review summaries
- Architecture review documents
- Audit findings
- Review outcomes and action items

## What Does NOT Go Here

- Architecture decisions → `../adr/`
- Operational decisions → `../decisions/`
- Work item reviews (status updates) → `../work-items/`

## Naming Convention

| Type | Pattern | Example |
|---|---|---|
| Review | `nnn-<review-topic>.md` | `001-auth-module-review.md` |

Sequential numbering with descriptive slugs.

## Ownership

| Role | Responsibility |
|---|---|
| Reviewers | Create review documents |
| Tech Lead | Reviews findings and assigns actions |
| Engineers | Address review feedback |

## Related Folders

- `../decisions/` — Reviews may produce decisions
- `../work-items/` — Reviews validate completed work items
- `../sessions/` — Reviews validate session outputs
