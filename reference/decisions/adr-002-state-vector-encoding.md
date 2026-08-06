---
namespace: eka-ref-impl
type: adr
id: 002-state-vector-encoding
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

# ADR-002 — State Vector Encoding: five frontmatter fields + single-writer + change-log

## Context

The initial implementation ran **seven status machines** (living documents, ADR, decisions, roadmap, sprints, sessions, work items — EKA 7.3) with three-way status synchronization: metadata in the document table, status in the sprint table, and status in the ticket document — without a single writer. This violates P6 (Single Writer) and creates status duplication that never stays consistent. EKA 7.1 rejects a unified State domain; EKA 7.2 establishes five independent State domains; EKA 7.4 establishes the State Vector = the tuple of domains **owned** per type.

## Decision

1. **Five frontmatter fields, one per owned state domain**: `content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state` (EKA 7.2). A field is present only for domains owned by the Artifact type (EKA 7.4, 10); **absence = not-applicable**.
2. **Single-writer per field** (P6): every state field has exactly one owner (Content State → Knowledge Layer; Execution/Planning/Container/Existence → Operating Layer). No other state-writing path exists.
3. **`change-log`**: `{date, domain, from, to, by}` array in the frontmatter — the mandatory chronological record of every state transition (EKA 5.2).
4. **Canonical lowercase values** (EKA 7.2 value sets):
   - `content-state` (per-type variants): living `draft | review | approved | amended`; ADR `proposed | accepted | superseded`; decision record `draft | accepted | superseded`
   - `execution-state`: `planned | todo | in-progress | in-review | done`
   - `planning-state`: `draft | approved | immutable`
   - `container-state`: `active | completed` (completed = derived transition)
   - `existence-state`: `active | archived | retired`
5. **Legacy value mapping**:

| Legacy machine | Legacy value | New value | Domain |
|---|---|---|---|
| Living documents | Draft | `draft` | content-state |
| Living documents | Review | `review` | content-state |
| Living documents | Approved | `approved` | content-state |
| Living documents | Amended | `amended` | content-state |
| ADR | Proposed | `proposed` | content-state |
| ADR | Accepted | `accepted` | content-state |
| ADR | Superseded | `superseded` | content-state |
| Decision Record | Draft | `draft` | content-state |
| Decision Record | Accepted | `accepted` | content-state |
| Decision Record | Superseded | `superseded` | content-state |
| Roadmap | Draft | `draft` | planning-state |
| Roadmap | Approved | `approved` | planning-state |
| Roadmap | Immutable | `immutable` | planning-state |
| Sprints | Active | `active` | container-state |
| Sprints | Completed | `completed` | container-state (derived) |
| Sessions | Active | `active` | existence-state |
| Sessions | Archived | `archived` | existence-state |
| Work items | Planned | `planned` | execution-state |
| Work items | Todo | `todo` | execution-state |
| Work items | In Progress | `in-progress` | execution-state |
| Work items | In Review | `in-review` | execution-state |
| Work items | Done | `done` | execution-state |

6. **Derived conditions** (e.g., container/session "Completed") are not domain values — they are computed from the aggregate of owner state (EKA 7.2) and never written as facts (see ADR-003).

## Consequences

- **Positive**: status duplication eliminated — one field, one owner, one source of truth (P6).
- **Positive**: mechanical validation becomes possible: domain values validated against the value set; `change-log` consistent with the current state; forward-only transitions checkable (P7).
- **Negative (intentional)**: legacy status consumers (`Status:` column, three-way synchronization) break — see `breaking-changes.md` #10–11.
- **Negative**: every state transition must now be recorded in `change-log` — a new discipline for the Operating Layer.

## Alternatives Considered

- **Unified State domain (Artifact State)** — rejected: a State monolith is the root of status duplication (EKA 7.1); seven semantics cannot be expressed by one machine (EKA 7.3).
- **Keeping the three-way synchronization** — rejected: no writer exists; violates P6.
- **Fully derived state (no owned fields)** — rejected: the OS needs owned state to run the Protocol; a projection without an owner has no source of truth.

## References

- EKA 7.1 (domain candidate evaluation), 7.2 (formal domains + value sets), 7.3 (legacy machine mapping), 7.4 (State Vector), 5.2 (change-log)
- Principles P2 (Explicit State), P6 (Single Writer), P7 (Forward-Only Transitions)
- Related: [ADR-001](adr-001-identity-serialization.md), [ADR-003](adr-003-projection-model.md)
