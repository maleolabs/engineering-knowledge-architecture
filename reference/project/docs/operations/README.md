# docs/operations/ — Operations Dimension

> Anchor EKA: Knowledge Layer — **operations** dimension (EKA 8).
> Engineering Domain: **Operations** (stratum 5).

## Purpose

The operations dimension houses operational knowledge for running and maintaining the system: runbooks, deployment procedures, recovery, and routine tasks. Documents here are **procedures** ("how to do"), separated from **standards** ("rules to follow"), which live in the standards dimension.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `run-` | Runbook | `run-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `run-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which situation/procedure is described.
- `## Content` — procedure steps, prerequisites, and expected results.

## Naming Conventions

`run-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `run-deploy-staging.md`.

## Procedures vs Standards

The separation is made at EKA 8: **standards** (`std-`) establish rules and conventions; **operations** (`run-`) explain execution procedures. Runbooks follow standards, they do not replace them. If a runbook establishes a new rule, that rule must be promoted to `std-`.

## Ownership

| Role | Responsibility |
|---|---|
| DevOps | single owner of runbooks |
| Tech Lead | reviewer of runbooks touching architecture |
| Engineers | authors of runbooks for the components they build |

## Related

- [standards/](../standards/) — procedures here are subject to `std-`.
- [records/](../records/) — `rel-` records releases executed with runbooks.
- [quality/](../quality/) — procedure effectiveness is verified via `rvw-`.
- [sessions/](../../operating/sessions/) — operational findings can be distilled into new `run-`.
