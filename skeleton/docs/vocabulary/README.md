# docs/vocabulary/ — Vocabulary Dimension

> Anchor EKA: Knowledge Layer — **vocabulary** dimension (EKA 8).
> Engineering Domain: **Architecture** (stratum 2).

## Purpose

The vocabulary dimension houses terms and their definitions: product, technical, and domain terms. One term = one `gls-`; definitions are canonical and referenced by all other dimensions so cross-role conversation stays unambiguous.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `gls-` | Glossary/Term | `gls-<id>.md` |

## State Vector

| Type | Owned state domains |
|---|---|
| `gls-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family):

- `## Purpose` — which term is defined.
- `## Content` — definition, synonyms, non-terms (what the term does not mean), and usage examples.

## Naming Conventions

`gls-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Examples: `gls-namespace.md`, `gls-artifact.md`.

## Vocabulary ≠ Specifications

The vocabulary dimension **defines the meaning of terms**; the specifications dimension **establishes behavior and technical formats**. Behavioral rules are not vocabulary — and vocabulary definitions are not specifications. `gls-` answers "what does X mean"; `spec-` answers "how does X work". Do not write specifications into `gls-` (and vice versa).

## Ownership

| Role | Responsibility |
|---|---|
| Product Owner | owner of product/domain terms |
| Tech Lead | owner of technical terms |
| All roles | propose new terms |

## Related

- [specifications/](../specifications/) — specifications use terms defined here.
- [intent/](../intent/) — key vision/strategy terms are defined here.
- [standards/](../standards/) — standards use standard vocabulary.
- [research/](../research/) — new terms from research are registered here.
