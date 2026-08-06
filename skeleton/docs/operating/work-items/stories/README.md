# stories/ — Story (`sto-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `sto-`.

## Purpose

A story is a unit of work delivering observable user value, described from the stakeholder's point of view.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `sto-` | `work-items/stories/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — the user need and the value delivered.
- `## Acceptance Criteria` — measurable conditions satisfying the definition of done.

## Naming Conventions

`sto-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `sto-login-email.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | definition of value and acceptance criteria |
| Engineer (implementer) | single owner of state; execution |
| Tech Lead | technical reviewer at `in-review` |

## Related

- [requirements/](../../../requirements/) — stories realize `req-`.
- [containers/](../../containers/) — stories are owned by the active container.
- [projections/](../../projections/) — story status is projected to `tkt-`.
