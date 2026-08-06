---
namespace: eka-ref-impl
type: adr
id: 003-projection-model
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

# ADR-003 — Projection Model: container and ticket tables are State Projections

## Context

The sprint table and ticket documents (wave docs) in the initial implementation stored status **duplicated** from the work item: status in the sprint table, status in the ticket document, and status in the work item file — three copies, kept manually in sync with no single writer. EKA 7.4 resolves this formally: derived representations are **State Projections** (derived views, owning no State of their own, never writers — P6, P9), and Artifacts whose entire state is projected have an **empty State Vector** (example: Ticket = `(∅)`). The projection validation-timing policy is still an open question (EKA 15.5: event-driven vs on-read).

## Decision

1. **The work item table on containers and the ticket artifact are State Projections** (EKA 7.4, 9 — Projection Semantics): their status is derived from the owner (the referenced work item), validated through Projection Refresh, and **never edited as an independent fact**.
2. **Ticket (`tkt-`) has an empty State Vector**: no state fields in the frontmatter; all ticket state is a projection of the referenced work item (EKA 10 — Ticket = `∅`).
3. **Relationships are written by Identity**: `derives-from: [ctr:<id>]` in the ticket frontmatter; the container points to work items; the projection chain always terminates at an owner (EKA 6.2.7).
4. **Generated artifacts carry an explicit header**:

   > Generated — State Projection. Do NOT edit state here; refresh on read.

5. **Default refresh policy: on-read** — projections are validated against the owner when read (EKA 15.5); the invariant "projections never write" is absolute (EKA 5.5).

## Consequences

- **Positive**: single-writer preserved — no second writer on status (P6); formal status duplication removed.
- **Positive**: container/ticket tables can be regenerated at any time from the owner without information loss.
- **Negative**: projection readers may observe stale status until an on-read refresh is performed — a consequence of the chosen policy, compensated by the warning header.
- **Negative**: legacy tooling that wrote status into tables/tickets breaks intentionally (`breaking-changes.md` #4–5).

## Alternatives Considered

- **Tables/tickets as the authoritative status source** — rejected: two writers per state field; violates P6; repeats the legacy duplication.
- **Automatic owner-state → projection sync script** — rejected: enforcement is an implementation capability (EKA 12.3, P16), but truth remains in the owner; a projection stays a projection, not a second fact.
- **Ticket with its own owned state** — rejected: EKA 10 establishes Ticket = `(∅)`; ticket status is a projection over the referenced work item (ratification resolution of Issue #1).

## References

- EKA 7.4 (State Vector; State Projection), 7.5 (interactions), 9 (Execution Taxonomy — Projection Semantics), 10 (Artifact Taxonomy — Ticket)
- Principles P6 (Single Writer), P9 (Structure as Projection of State)
- EKA 15.5 (open question: Projection Refresh policy)
- Related: [ADR-002](adr-002-state-vector-encoding.md)
