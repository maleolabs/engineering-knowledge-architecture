# tech-debt/ — Tech Debt (`td-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `td-`.

## Purpose

Tech Debt is a unit of work for identified technical debt: shortcuts, outdated components, and inconsistencies that burden subsequent development. Every `td-` must record why the debt exists (rationale) so its repayment can be consciously prioritized.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `td-` | `work-items/tech-debt/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — the form of the technical debt and its location.
- `## Acceptance Criteria` — conditions proving the debt is repaid.
- `## Debt Rationale` — why this debt was taken on (decision/timing), so future decisions stay contextual.

## Naming Conventions

`td-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `td-single-writer-migration.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Engineer (implementer) | single owner of state; execution |
| Tech Lead | decision-maker for debt priority |
| Product Owner | reviewer of impact on the plan |

## Related

- [decisions/](../../../decisions/) — `Debt Rationale` can reference the `dec-` that gave rise to it.
- [standards/](../../../standards/) — repaying debt restores `std-` conformance.
- [containers/](../../containers/) — tech debt is owned by the active container.
