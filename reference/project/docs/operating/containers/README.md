# containers/ — Execution Container (`ctr-`)

> Anchor EKA: Operating Layer — State Domain **Container State**.

## Purpose

A container is an execution unit that shelters a set of work items under one locked plan. The container holds the aggregation relationship (which work items are in this wave) and projects them into a table.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `ctr-` | `operating/containers/` | `container-state`, `existence-state` |

`container-state` values: `active → completed`. `existence-state` values: `active → archived → retired`. Work items do not store container state; containers do not store work item Execution State (that is owned by the work item).

## Exactly One Active Container

Mutual exclusion: only **one** `ctr-` may have `container-state: active` at a time. A new container is born only after the previous one is `completed`. Container creation is **atomic** with locking of its supporting plan (`plan-` → `immutable`; see [protocol.md](../protocol.md) §4).

## `completed` = Derived Transition

`container-state: completed` is not an arbitrary value — it is a **derived transition** triggered by the aggregate Execution State of all work items inside: every work item is `done`. When the aggregate is satisfied, the container owner writes the `active → completed` transition to `change-log`.

## Work Item Table = Projection

The `## Work Items` section is a **projection** — a snapshot of the owning work items' state:

```
> Generated — State Projection. Do NOT edit state here; refresh on read.
```

- Projections are refreshed on read; do not edit state in this table.
- Conflict between the table and owner state: **owner state wins** (see validation, Rule 8).

## Snapshot Semantics

The container records a context snapshot at creation: the locked plan (with its `instance-version`), the initial work item list, and the scope. The snapshot does not change for the container's lifetime; forward changes happen on subsequent artifacts/instances.

## Required Content Structure

- `## Objective` — the execution objective of this wave.
- `## Work Items` — projection table of work items (token, id, summary, Execution State).
- `## Change Log` — record of container state transitions.

## Naming Conventions

`ctr-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `ctr-wave-1.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | single owner of container state; writes the `completed` transition |
| Engineers | execute work items within the container |

## Related

- [planning/](../../planning/) — container locks `plan-` (lock-atomic-with-generation).
- [work-items/](../work-items/) — aggregated units of work.
- [projections/](../projections/) — `tkt-` projects work items to per-wave status.
- [sessions/](../sessions/) — execution within a container is recorded per session.
