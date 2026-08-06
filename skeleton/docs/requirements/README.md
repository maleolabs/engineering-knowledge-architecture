# docs/requirements/ — Requirements Dimension

> Anchor EKA: Knowledge Layer — **requirements** dimension (EKA 8).

## Purpose

The requirements dimension houses the needs the product must fulfill, derived from intent and agreed with stakeholders. Each `req-` is a single requirement whose fulfillment can be verified by `rvw-` or a work item.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `req-` | Requirement | `req-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `req-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which requirement is described.
- `## Content` — the requirement statement, acceptance criteria, and its context.

## Naming Conventions

`req-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `req-login-email.md`.

## Amendment = New Instance

Changes to an approved requirement do **not** blindly edit the old document: create a new instance `req-<new-id>.md` with `content-state: amended` on the old document (or archive the old one) and the `amends: [req:<old-id>]` field on the new one. The amendment chain is traceable through `amends`.

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | single owner of requirements content |
| Engineers | reviewers of technical feasibility at `review` |
| All roles | propose new requirements |

## Related

- [intent/](../intent/) — source of requirement derivation.
- [specifications/](../specifications/) — `req-` is detailed into `spec-`.
- [quality/](../quality/) — `rvw-` can validate requirement fulfillment.
- [planning/](../planning/) — `scp-` selects requirements into scope.
