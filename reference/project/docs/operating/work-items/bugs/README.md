# bugs/ — Bug (`bug-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `bug-`.

## Purpose

A bug is a unit of work for behavior deviating from expectations: functional failure, non-conformance to a specification, or regression. Every bug must be reproducible and traceable to the artifact it violates.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `bug-` | `work-items/bugs/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — symptom, reproduction steps, and expected behavior.
- `## Impact` — impact on users/system and severity.
- `## Root Cause` (optional) — root cause once identified.

## Naming Conventions

`bug-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `bug-ticket-stale-after-refresh.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Reporter | describes symptoms and impact |
| Engineer (implementer) | single owner of state; diagnosis and fix |
| Tech Lead | fix reviewer at `in-review` |

## Related

- [specifications/](../../../specifications/) — correct behavior is referenced from `spec-`.
- [quality/](../../../quality/) — fixes are verified by `rvw-`.
- [containers/](../../containers/) — bugs are owned by the active container.
