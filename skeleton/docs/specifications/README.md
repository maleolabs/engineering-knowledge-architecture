# docs/specifications/ — Specifications Dimension

> Anchor EKA: Knowledge Layer — **specifications** dimension (EKA 8).

## Purpose

The specifications dimension houses implementable specifications: behavior, interface, format, and technical rule details precise enough to serve as the basis for implementation and testing. Specifications bridge requirements (`req-`) toward implementation.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `spec-` | Specification | `spec-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `spec-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which part of the system is specified.
- `## Content` — the detailed specification: behavior, input/output, constraints, format.

## Naming Conventions

`spec-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `spec-ticket-projection.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | single owner of technical specifications |
| Engineers | contributors and implementers |
| Product Owner | reviewer of fit with requirements |

## Related

- [requirements/](../requirements/) — specifications derive from `req-`.
- [architecture/](../architecture/) — specifications are subject to `arc-` constraints.
- [vocabulary/](../vocabulary/) — terms used in specifications must be defined (Vocabulary ≠ Specifications).
- [quality/](../quality/) — `rvw-` validates specifications via `validates`.
- [standards/](../standards/) — specifications follow `std-`.
