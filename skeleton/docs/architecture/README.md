# docs/architecture/ — Architecture Dimension

> Anchor EKA: Knowledge Layer — **architecture** dimension (EKA 8).
> Engineering Domain: **Architecture** (stratum 2).

## Purpose

The architecture dimension houses system architecture descriptions: component structure, interactions, technical constraints, and the rationale of architectural decisions. `arc-` documents describe "how the system is structured", while "why it is structured that way" lives in `adr-` in the decisions dimension.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `arc-` | Architecture Description | `arc-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `arc-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — the scope of the architecture described.
- `## Content` — component, interaction, and constraint descriptions (diagrams may be referenced by path).

## Naming Conventions

`arc-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `arc-identity-namespace.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | single owner of architecture content |
| Engineers | contributors of architecture segments |
| DevOps | contributors for deployment-infrastructure aspects |

## Related

- [decisions/](../decisions/) — `adr-` explains the decisions behind `arc-`.
- [specifications/](../specifications/) — `arc-` is realized as detailed specifications.
- [standards/](../standards/) — `std-` binds technical style and conventions.
- [requirements/](../requirements/) — architecture fulfills requirements.
