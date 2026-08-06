# docs/planning/ — Planning Dimension

> Anchor EKA: Knowledge Layer — **planning** dimension (EKA 8) + State Domain **Planning State**.

## Purpose

The planning dimension houses planning artifacts: scope definitions, capabilities, execution plans, and relationship artifacts. This is the only dimension where `phase` (discovery | mvp | milestone | release | growth | maturity | sunset) is used — as a **context attribute** on `scp-`/`plan-`, not a folder.

## What Lives Here

| Token | Type | Name format | Versioned |
|---|---|---|---|
| `scp-` | Scope Definition | `scp-<id>-v<instance-version>.md` | yes (phase) |
| `epc-` | Epic | `epc-<id>.md` | no |
| `plan-` | Plan | `plan-<id>-v<instance-version>.md` | yes (phase, Planning State) |
| `trc-` | Traceability/Relationship artifact | `trc-<id>.md` | no |

## State Vector

| Type | Owned state domains |
|---|---|
| `scp-`, `epc-`, `trc-` | `content-state`, `existence-state` |
| `plan-` | `content-state`, `planning-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `planning-state` values: `draft → approved → immutable`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (planning artifact family):

- `## Objective` — the objective of this artifact.
- `## Scope` — what is included.
- `## Out of Scope` — what is deliberately excluded.

## Naming Conventions

Versioned (always, including v1): `scp-<id>-v<instance-version>.md` and `plan-<id>-v<instance-version>.md`. Non-versioned: `epc-<id>.md`, `trc-<id>.md`. `instance-version` is required in frontmatter for `scp-`/`plan-`.

## Per-Type Notes

### `scp-` — Scope Definition
Carries `phase` as a context attribute. Phase changes are recorded in `change-log` with `domain: phase`. An approved scope becomes the basis for an execution container.

### `plan-` — Plan
Carries `phase` and **Planning State** (`draft → approved → immutable`). **Lock-atomic-with-generation:** once a `plan-` becomes `immutable`, any change — including fixes — must not edit that instance; create a new `instance-version` (`plan-<id>-v<nn+1>.md`). The transition to `immutable` happens atomically with container creation (see [operating/containers/](../operating/containers/) and [operating/protocol.md](../operating/protocol.md)).

### `epc-` — Epic
A capability that realizes scope. Does not carry `phase`; reference its parent `scp-` with `derives-from` when needed.

### `trc-` — Traceability/Relationship artifact
An artifact that **carries relationships only** (e.g. requirement→specification→work item matrices, dimension maps). Its main content is a list of references that must resolve; it does not replace relationships written on the referring artifacts.

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | owner of `scp-`, `epc-`, `plan-` (scope & priority) |
| Tech Lead | feasibility reviewer and owner of technical `plan-` |
| Engineers | contributors of estimates and plan details |

## Related

- [requirements/](../requirements/) — scope selects `req-`.
- [intent/](../intent/) — `scp-` elaborates strategy into phased context.
- [operating/containers/](../operating/containers/) — `plan-` locks atomically with the birth of `ctr-`.
- [operating/projections/](../operating/projections/) — tickets represent work items within a container.
- [records/](../records/) — `rel-` records plan execution results.
