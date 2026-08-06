# chores/ — Chore (`ch-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `ch-`.

## Purpose

A chore is an administrative or maintenance unit of work that does not change product behavior: dependency updates, configuration, cleanup, and routine project tasks.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `ch-` | `work-items/chores/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — the task being done.
- `## Acceptance Criteria` — conditions proving the task is complete.

## Naming Conventions

`ch-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `ch-update-ci-dependencies.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Assigner | whoever assigns (PO/Tech Lead/DevOps) |
| Engineer / DevOps (implementer) | single owner of state; execution |

## Related

- [operations/](../../../operations/) — routine tasks can reference `run-`.
- [containers/](../../containers/) — chores are owned by the active container.
