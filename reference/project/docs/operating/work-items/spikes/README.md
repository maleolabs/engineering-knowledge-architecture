# spikes/ — Spike (`spk-`)

> Anchor EKA: Operating Layer — State Domain Execution State; work item subtype `spk-`.

## Purpose

A spike is a time-boxed investigation unit of work to reduce uncertainty: proving feasibility, exploring approaches, or gathering data before a decision. A spike's result is **not** a decision — its findings must be distilled into the knowledge store.

## Token & State Vector

| Token | Folder | Owned state |
|---|---|---|
| `spk-` | `work-items/spikes/` | `execution-state`, `existence-state` |

`execution-state` values: `planned → todo → in-progress → in-review → done` (forward-only, never skip, never revert). `existence-state` values: `active → archived → retired`.

## Required Content Structure

- `## Description` — the question under investigation and its time-box/scope.
- `## Investigation Notes` — the investigation trail: experiments, data, sources.
- `## Conclusion` — must contain links to the knowledge distillation destination (e.g. `fnd-` in research/ or `dec-`/`adr-` in decisions/) — distill before archiving (EKA 11.4).

## Naming Conventions

`spk-<id>.md`, kebab-case, unique within (namespace, type). No `-v<nn>` suffix. Example: `spk-ticket-projection-feasibility.md`.

## Ownership

| Role | Responsibility |
|---|---|
| Engineer (implementer) | single owner of state; investigation |
| Tech Lead | reviewer of conclusions and distillation path |

## Related

- [research/](../../../research/) — research results are distilled into `fnd-` (EKA 14.1).
- [decisions/](../../../decisions/) — adopted conclusions become `dec-`/`adr-`.
- [specifications/](../../../specifications/) — proven findings become `spec-`.
