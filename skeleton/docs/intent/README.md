# docs/intent/ — Intent Dimension

> Anchor EKA: Knowledge Layer — **intent** dimension (EKA 8).

## Purpose

The intent dimension houses the product/project's direction and reason for existence: vision, manifesto, and strategy. Documents here answer "why does this project exist" and "where is it heading", and serve as the justification reference for other dimensions.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `vis-` | Vision/Manifesto | `vis-<id>.md` |
| `str-` | Strategy | `str-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `vis-` | `content-state`, `existence-state` |
| `str-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = not applicable (N/A). Every transition is recorded in `change-log` by the single owner.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — the purpose of this document.
- `## Content` — the vision/manifesto or strategy content.

## Naming Conventions

`vis-<id>.md` and `str-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix (non-versioned types). Examples: `vis-product-core.md`, `str-market-entry-2026.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | single owner of intent content |
| Tech Lead | contributor for technical `str-` |
| All roles | may propose changes; approved content changes become a new instance with `amends` |

## Related

- [requirements/](../requirements/) — intent is derived into requirements (`req-`).
- [vocabulary/](../vocabulary/) — key intent terms must be defined in `gls-`.
- [decisions/](../decisions/) — strategic decisions (`dec-`) refer back to intent.
- [planning/](../planning/) — `scp-` elaborates the phased context of the strategy.
