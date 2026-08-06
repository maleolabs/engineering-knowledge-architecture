# docs/records/ — Records Dimension

> Anchor EKA: Knowledge Layer — **records** dimension (EKA 8).
> Engineering Domain: **Operations** (stratum 5).

## Purpose

The records dimension houses immutable chronological records: release records, event notes, and the project's audit trail. Unlike working documents, which keep evolving, records here are factual — they record what happened.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `rel-` | Release Record | `rel-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `rel-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A. After a release, a `rel-` is practically `approved` and is no longer edited (fact changes = new instance with `amends`).

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which release/event is recorded.
- `## Content` — summary of execution and the release.

## Release Record = Execution Aggregate + Release Gates

A `rel-` is an **aggregate** of execution outcomes (completed work items, sessions that occurred) and **release gates** (the approval and readiness gates passed before release). Both are referenced via relationships (`derives-from`/`validates`), not re-quoted — keeping the record concise and its trail traceable.

## Naming Conventions

`rel-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `rel-2026-08-mvp.md`.

## Ownership

| Role | Responsibility |
|---|---|
| DevOps | owner of `rel-` (release execution) |
| Tech Lead | signer of the technical release gate |
| Product Owner | signer of the product release gate |

## Related

- [operations/](../operations/) — release procedures executed (`run-`).
- [quality/](../quality/) — release gates verified by `rvw-`.
- [decisions/](../decisions/) — release-triggered decisions are recorded as `dec-`.
- [operating/work-items/](../operating/work-items/) — execution outcomes aggregated by `rel-`.
