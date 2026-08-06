# docs/research/ — Research Dimension

> Anchor EKA: Knowledge Layer — **research** dimension (EKA 8), token extension `fnd-` (EKA 14.1).
> Engineering Domain: **Discovery** (stratum 1, highest authority).

## Purpose

The research dimension houses investigation and research findings: experiments, benchmarks, literature searches, and technical findings that are not (yet) decisions. Research is a **knowledge** dimension — not a substitute for decisions.

## What Lives Here

| Token | Type | Name format |
|---|---|---|
| `fnd-` | Research Finding | `fnd-<id>.md` |

Extension per EKA 14.1: the `fnd-` token is added to the standard token table to house research findings, with `dimension: research` and home folder `research/`.

## State Vector

| Type | Owned state domains |
|---|---|
| `fnd-` | `content-state`, `existence-state` |

`content-state` values: `draft → review → approved → amended`. `existence-state` values: `active → archived → retired`. Fields not listed = N/A.

## Good Content Structure

Required structure (knowledge document family, with research extensions):

- `## Purpose` — the research question answered.
- `## Content` — account of the findings.
- `## Investigation Summary` — summary of methods and evidence.
- `## Conclusion` — conclusions and recommendations.

## Mandatory Distillation (EKA 11.4)

Research findings **must not stop** as `fnd-`. When a finding affects project direction, it must be distilled into the decisions dimension:

1. Adopted conclusions → new `dec-`/`adr-` (or an amendment of a relevant decision), with `derives-from: [fnd:<id>]`.
2. The original `fnd-` remains intact as an evidence trail; it may only be archived after distillation completes.

Findings that are not adopted are simply archived with a note of the rejection reason.

## Naming Conventions

`fnd-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `fnd-plan-lock-approach.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Engineers | authors of `fnd-` (technical investigations) |
| Tech Lead | reviewer and distributor of distillation into decisions |
| All roles | may open new research |

## Related

- [decisions/](../decisions/) — **mandatory destination** of research distillation (EKA 11.4).
- [specifications/](../specifications/) — adopted findings become `spec-`.
- [operating/work-items/spikes/](../operating/work-items/spikes/) — spikes produce research material; `fnd-` and `spk-` reference each other.
- [vocabulary/](../vocabulary/) — new terms from research are registered in `gls-`.
