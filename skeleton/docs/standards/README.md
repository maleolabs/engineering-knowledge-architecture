# docs/standards/ — Standards Dimension

> Anchor EKA: Knowledge Layer — **standards** dimension (EKA 8).
> Engineering Domain: **Architecture** (stratum 2).

## Purpose

The standards dimension houses standards and guidelines that bind behavior, style, and quality of work: code conventions, document formats, processes, and quality criteria. Standards establish **rules to be followed**, unlike runbooks, which describe **execution procedures**.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `std-` | Standard/Guideline | `std-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `std-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which area is governed.
- `## Content` — rules/conventions; may take the form of a conformance checklist.

## Naming Conventions

`std-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `std-frontmatter-identity.md`.

## Conventions ≠ Procedures

Standards establish **conventions** (what is allowed/not allowed). Step-by-step procedures are executed by runbooks (`run-`) in the operations dimension; if a document explains "how to do it", it is not a standard.

## Ownership

| Role | Responsibility |
|---|---|
| Tech Lead | single owner of technical standards |
| DevOps | owner of operational/deployment standards |
| All roles | must conform to approved `std-` |

## Related

- [operations/](../operations/) — procedures (`run-`) are separated from standards here.
- [quality/](../quality/) — standard conformance is verified by `rvw-`.
- [specifications/](../specifications/) — specifications follow standards.
- [vocabulary/](../vocabulary/) — terms used in standards must be defined.
