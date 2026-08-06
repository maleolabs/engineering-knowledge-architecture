# docs/quality/ — Quality Dimension

> Anchor EKA: Knowledge Layer — **quality** dimension (EKA 8).

## Purpose

The quality dimension houses review and quality-verification results: code reviews, architecture reviews, and product reviews. Each `rvw-` validates one or more other artifacts through the `validates` relationship, so quality is traceable back to what was verified.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `rvw-` | Review | `rvw-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `rvw-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family, with review extensions):

- `## Purpose` — the object and scope of the review.
- `## Content` — general account of the review results.
- `## Findings` — findings: problems, risks, standard violations.
- `## Action Items` — required follow-ups.

The `validates: [<type>:<id>]` relationship must point to the reviewed artifacts (e.g. `spec:login`, `std:frontmatter`).

## Naming Conventions

`rvw-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `rvw-frontmatter-serialization.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | peer reviews (code/specifications) |
| Tech Lead | architecture reviews; owner of the quality dimension |
| DevOps | operational/deployment reviews |
| Product Owner | product reviews (requirement fulfillment) |

## Related

- [standards/](../standards/) — conformance is verified by `rvw-`.
- [specifications/](../specifications/) — `rvw-` validates `spec-`.
- [decisions/](../decisions/) — review findings can give rise to `dec-`.
- [operating/work-items/](../operating/work-items/) — action items can become new work items.
