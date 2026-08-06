---
namespace: eka-ref-impl
type: adr
id: 004-phase-as-metadata
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
  - 002-state-vector-encoding
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

# ADR-004 — Phase as Metadata (context attribute), not folder

## Context

The initial implementation encoded phase as a **folder** (`docs/mvp/`, etc.): the product phase became part of the location, and because the legacy Identity was encoded via location, the phase also became part of the Identity representation. This violates EKA 11.2 (Phase is a **context attribute** on planning/scope Artifacts — not a category, not a State Domain) and P3 (Identity independent of location). EKA 7.1 also rejects "Lifecycle State" as a single domain: what survives is Existence State (domain) + Phase (context).

## Decision

1. **A `phase` field in the frontmatter**, only on `scp-` (Scope Definition) and `plan-` (Plan) artifacts — per EKA 11.2 (phase attaches to Scope Definition / Plan as a context attribute).
2. **Value set**: `discovery | mvp | milestone | release | growth | maturity | sunset` (EKA 11.2–11.3 values in lowercase).
3. **A phase change is a context update** authorized by the **readiness Gate** (EKA 11.2), evaluated over the owner's State aggregate (release-ready = all work items Done ∧ all Containers Completed ∧ plan Immutable ∧ review Gate passed ∧ Content approval Gate passed).
4. **Recording**: every phase change is recorded in `change-log` with `domain: phase` (e.g., `from: discovery, to: mvp`).
5. **No phase folders** — artifacts do not change location when the phase changes; Identity is untouched (P3).

## Consequences

- **Positive**: Identity decoupled from phase (P3) — a Scope Definition keeps the same identity when the product moves across phases (EKA 11.2).
- **Positive**: a phase change becomes an audited metadata operation (change-log), not a file-moving operation.
- **Positive**: per-phase globbing is no longer needed; phase queries become field queries.
- **Negative (intentional)**: the `mvp/` folder disappears; tooling that reads phase from the path breaks (`breaking-changes.md` #12).

## Alternatives Considered

- **Folder per phase** (legacy status quo) — rejected: phase becomes part of the location → part of the Identity representation; violates EKA 6.4, P3.
- **Phase as a State Domain** — rejected: EKA 7.1 rejects Lifecycle State; phase has no protocol transitions; it is a context that changes via a Gate (EKA 11.2).
- **Phase as classification (Knowledge Dimension)** — rejected: classification is a retrieval property (P15), not a temporal context; phase is not stable enough as a classification axis.

## References

- EKA 3 (Phase definition), 7.1 (Lifecycle State rejected), 7.5 (readiness Gate → phase change), 11.2 (Phase as context), 11.3 (product lifecycle)
- Principles P3 (Stable Identity), P9 (Structure as Projection of State)
- Related: [ADR-001](adr-001-identity-serialization.md), [ADR-002](adr-002-state-vector-encoding.md), [ADR-005](adr-005-dimension-layout.md)
