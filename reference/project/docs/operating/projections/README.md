# projections/ — Ticket (`tkt-`)

> Anchor EKA: Operating Layer — pure projection; **owns** no State Domain.

## Purpose

A ticket is a projection of one work item into a compact form: executable commands and projected status. A ticket is **not a state owner** — it is a shadow of the work item's owner state behind it.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `tkt-` | `operating/projections/` | **none (empty state vector)** |

Tickets carry no state fields at all (`content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state` are all N/A). The status visible on a ticket is only a read copy.

## Relationship

A ticket represents one work item: `derives-from: [ctr:<id>, <type-work-item>:<id>]` — pointing to the active container and its source work item. References must resolve; if the source work item changes state, the ticket must be refreshed (see below).

## Content

- `## Commands` — deterministic commands per work item type (e.g. "run the tests for X", "update the changelog") — identical on every refresh.
- `## Projected Status` — status projected from the work item's owner state (Execution State, etc.) — repopulated on every refresh.

Required header at the top:

```
> Generated — State Projection. Do NOT edit state here; refresh on read.
```

## Refresh on Read

Projections are refreshed on every read (default). Automated event-driven mechanisms (projection updated when state changes) are left to future tooling; until such tooling exists, **a projection is the result of the read at that moment** — not a source of truth.

## Do Not Edit State in Projections

Changing status in a ticket **does not change** the work item and violates single-writer (P6). Status changes are made only by the work item's state owner; the ticket is refreshed afterwards.

## Naming Conventions

`tkt-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `tkt-sto-login-email.md`.

## Ownership

| Role | Responsibility |
|---|---|
| OS (projection) | "writer" of ticket content — regenerated, never manually edited |
| All roles | read-only; report anomalies (ticket ≠ owner state) to the state owner |

## Related

- [work-items/](../work-items/) — the projected owner-state source.
- [containers/](../containers/) — container work item tables are a similar projection.
- [validation.md](../../exchange/validation.md) — Rule 8: projections must not carry state owned by other artifacts.
