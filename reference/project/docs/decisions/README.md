# docs/decisions/ — Decisions Dimension

> Anchor EKA: Knowledge Layer — **decisions** dimension (EKA 8).
> Engineering Domain: **Architecture** (stratum 2).

## Purpose

The decisions dimension houses decisions that have been made, along with their context, alternatives, and consequences — so decisions do not get "lost" when the people who made them leave. This dimension has two artifact variants with different state variants.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `adr-` | Architecture Decision Record | `adr-<id>.md` |
| `dec-` | Decision Record (general/operational decision) | `dec-<id>.md` |

## State Vector

| Type | Owned state domains | `content-state` values |
|---|---|---|
| `adr-` | `content-state`, `existence-state` | `proposed → accepted → superseded` |
| `dec-` | `content-state`, `existence-state` | `draft → accepted → superseded` |

`existence-state` values: `active → archived → retired`. Note: the `adr-` variant does not use `draft`/`review`/`approved`/`amended`; decision status is expressed directly as `proposed`/`accepted`/`superseded`.

## Good Content Structure

Required structure (decision record family):

- `## Context` — background and the problem that triggered the decision.
- `## Decision` — the decision made.
- `## Consequences` — positive and negative impacts.
- `## Alternatives Considered` — alternatives evaluated and reasons for rejection.
- `## References` (optional) — additional references.

## Supersession

- **`adr-`: supersession is mandatory.** When an ADR is no longer valid, create a new ADR and set `supersedes: [adr:<old-id>]` on the new one and `content-state: superseded` on the old one. A superseded ADR without a replacement reference is a validation violation.
- **`dec-`: supersession is optional.** Operational decisions may be superseded but are not required to be; a decision can simply move to `existence-state: archived` when no longer relevant.

## Naming Conventions

`adr-<id>.md` and `dec-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Examples: `adr-identity-serialization.md`, `dec-git-workflow.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | owner of `adr-` (architectural decisions) |
| Product Owner | owner of `dec-` (product/scope decisions) |
| Engineers / DevOps | contributors of technical-operational `dec-` |

## Related

- [architecture/](../architecture/) — `adr-` explains why `arc-` is the way it is.
- [records/](../records/) — chronological decision records at release time.
- [sessions/](../../operating/sessions/) — session distillation results must flow into `dec-`/`adr-`.
- [research/](../research/) — research findings require a distillation path to here (EKA 11.4).
