---
namespace: eka-ref-impl
type: adr
id: 007-extension-research-finding
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
  - 005-dimension-layout
  - 006-exchange-conventions
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

# ADR-007 — Extension: `fnd-` Artifact type (Research Finding)

## Context

The **Research** dimension (EKA 8) — "investigation findings, technical research results; mandatory Distillation path to durable dimensions" — **has no Artifact type** in the EKA 10 Artifact Taxonomy. In the initial implementation, investigation results (spikes, technical research) had no durable home: findings were lost or buried in ephemeral sessions. EKA 14.1 provides a lightweight extension mechanism for new Artifact types; EKA 14.2.6 requires new types to declare a **complete** owned State Vector (no implicit default inheritance).

## Decision

Register the **Research Finding** (`fnd-`) extension type via the EKA 14.1 extension mechanism:

| Aspect | Declaration |
|---|---|
| Type token | `fnd-` (enters the 26-token table, ADR-001) |
| Knowledge Dimension | Research |
| Owned State Vector (complete) | `(Content State, Existence State)` — other domains not-applicable |
| Identity rules | Line + instance; ID unique within `(namespace, fnd)` |
| Folder | `research/` (`dimension == folder` rule, ADR-005) |
| Relationship | `derives-from` (e.g., from spike `spk-`); Distillation output toward durable dimensions (decisions, ADR, Record) |

Governance consequences: the extension is registered as part of the standard (EKA 14.2.5: proposal → review → acceptance) and **is exchangeable** (EKA 14.2.4) — covered by the exchange conventions' schema versioning (ADR-006).

## Consequences

- **Positive**: the spike → durable knowledge Distillation path becomes explicit (EKA 11.4): research findings (`fnd-`) are distilled into decisions/ADRs/Records instead of evaporating in sessions.
- **Positive**: the Research dimension (EKA 8) now has an artifact home; investigation results are preserved (P12).
- **Positive**: the extension is legitimate and exchangeable — it does not deviate from the invariants (EKA 14.2.1); backward compatible (14.2.2).
- **Negative**: one new type must be kept compliant — the complete owned State Vector declaration (Content, Existence) must not change implicitly (14.2.6).

## Alternatives Considered

- **Using existing types (`spec-`/`std-`/`rel-`)** — rejected: Research ≠ Specifications/Standards/Records (EKA 8); Research is cumulative, not immutable, and mandates a Distillation path.
- **No dedicated type (findings only in sessions)** — rejected: sessions are ephemeral by design (EKA 10); violates P12 (Preservation Over Deletion) and EKA 11.4.
- **Research as a dimension without an artifact type** — rejected: a dimension without a type cannot be produced/exported as an Artifact (EKA 14.1, 13.1).

## References

- EKA 8 (Research dimension), 10 (Artifact Taxonomy), 11.4 (Distillation lifecycle), 14.1 (extension points), 14.2 (extension rules)
- Principles P12 (Preservation Over Deletion), P15 (Classification is Property, Not Identity)
- Related: [ADR-001](adr-001-identity-serialization.md), [ADR-005](adr-005-dimension-layout.md), [ADR-006](adr-006-exchange-conventions.md)
