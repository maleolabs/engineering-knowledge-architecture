---
namespace: eka-ref-impl
type: adr
id: 005-dimension-layout
instance-version: 1
revision: 1
content-state: accepted
existence-state: active
dimension: decisions
author: Engineering Architecture
created: 2026-08-05
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on:
  - 001-identity-serialization
change-log:
  - date: 2026-08-05
    domain: existence-state
    from: "-"
    to: active
    by: Engineering Architecture
  - date: 2026-08-05
    domain: content-state
    from: proposed
    to: accepted
    by: Engineering Architecture
---

# ADR-005 — Dimension Layout: 12 knowledge folders = 12 Knowledge Dimensions, 1:1

## Context

The legacy taxonomy was mixed: 16 folders with dimension-mixed contents (operations mixed with standards, planning as a catch-all, work-items/planning, specification-corpus as a vocabulary misnomer). EKA 8 establishes 12 Knowledge Dimensions with strict separation (Operational ≠ Standards, Vocabulary ≠ Specifications, etc.); P15 establishes classification as an Artifact property; P1 establishes Separation of Concerns between knowledge and execution. The legacy folders mixed all three: classification, pipeline, and layers.

## Decision

1. **12 knowledge folders = 12 Knowledge Dimensions, 1:1** (EKA 8): `intent`, `requirements`, `architecture`, `decisions`, `specifications`, `standards`, `operations`, `quality`, `planning`, `records`, `research`, `vocabulary`.
2. **`operating/`** — OS-layer folder (containers, work-items, projections, sessions, protocol): the state machine and protocol live here, not in dimension folders (P1, EKA 4.1).
3. **`exchange/`** — EX-layer folder (`validation.md`, `transfer.md`): exchange contracts (EKA 13).
4. **Location rule**: knowledge artifacts live in their dimension folder; validation enforces **`dimension == folder`** (the `dimension` field in the frontmatter must equal the folder where the file resides) — P15, P9.
5. **Operating artifacts are exempt** from the `dimension == folder` rule (a work item has dimension Requirements/Records/Research yet lives in `operating/work-items/` — the OS dimension determines its home, not its Knowledge Dimension).
6. **Catch-alls dissolved**: no folder hosts a mix of dimensions (EKA 14.2, P15); content is distributed to the correct dimension.

## Consequences

- **Positive**: stable, mechanically validated classification — `dimension == folder` (P15); reclassification never breaks references because references are by Identity (P3).
- **Positive**: strict layer separation (P1): knowledge in dimension folders, execution in `operating/`, exchange in `exchange/`.
- **Positive**: legacy catch-alls gone; every artifact has a clear taxonomic home.
- **Negative (intentional)**: all legacy paths change (`breaking-changes.md` #1, #8, #9); migration follows `migration-guide.md`.

## Alternatives Considered

- **Keeping the 16 legacy folders** — rejected: mixes dimensions and layers; violates P1, P15, EKA 8.
- **Subfolder per dimension under one knowledge folder** — rejected: 1:1 folder↔dimension was chosen so that `dimension == folder` validation stays simple and deterministic (P16: enforcement varies, invariants identical).
- **Classification only in frontmatter, free-form folders** — rejected: the folder is a classification projection (P9); without the projection, navigation and validation weaken.

## References

- EKA 8 (Knowledge Taxonomy — 12 dimensions), 14.2 (extensions; core closed, taxonomy open), 4.1 (Layer Model)
- Principles P1 (Separation of Concerns), P9 (Structure as Projection of State), P15 (Classification is Property, Not Identity)
- Related: [ADR-001](adr-001-identity-serialization.md), [ADR-006](adr-006-exchange-conventions.md), [ADR-007](adr-007-extension-research-finding.md)
