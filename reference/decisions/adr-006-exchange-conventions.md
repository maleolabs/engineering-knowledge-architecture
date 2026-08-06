---
namespace: eka-ref-impl
type: adr
id: 006-exchange-conventions
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
  - 005-dimension-layout
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

# ADR-006 — Exchange Conventions: validation.md + transfer.md as the EX-layer seam

## Context

The Knowledge OS vision (EKA 1.4, 16.1) requires an exchange seam defined at the standard level: the repository must be importable/exportable without losing Identity, State, Content, or Relationship (EKA 13.1, P13). Without exchange conventions, this repository is just a file store — not ready to be a consumer/producer of Artifacts for external systems.

## Decision

The exchange seam is realized as two convention documents in `skeleton/docs/exchange/`:

1. **`validation.md`** — 9 Conformance Rules (R1–R9) (extracted from the standard contracts):
   1. Complete and unique Identity: `(namespace, type, id)` unique; `instance-version` unique per Line (6.2.2).
   2. Frontmatter Identity canonical and machine-parseable (6.2.6).
   3. State values valid against each domain's value set (7.2).
   4. Forward-only transitions; `change-log` consistent with current state (P7, 5.2).
   5. Single-writer: no two owners for one state field (P6).
   6. References by Identity; no dangling references (6.2.3, 5.1).
   7. Classification: `dimension == folder` for knowledge artifacts (8, P15).
   8. Phase: valid value and only on `scp-`/`plan-` (11.2).
   9. Content Well-formed per artifact type (3, 5.3).
2. **`transfer.md`** — the transfer contract, following EKA 13.2:
   - **Lossless round-trip** (13.2.1) and **idempotent**: re-import = no-op (13.2.2);
   - **Referential integrity** across systems (13.2.3);
   - **Identity conflict policy**: importing an already-existing Identity = **reject or explicit re-namespace** — never silent merge (13.2.4);
   - **Validation before commit** (13.2.5);
   - **Schema versioning**: versioned exchange contract; import/export declares the contract version complied with (13.2.6).

Both documents are convention documents (without `type`/`id`) — not Artifacts; they describe the contracts and carry no state.

## Consequences

- **Positive**: the repository is import/export-ready without redesign — the exchange seam is defined from the start (EKA 13).
- **Positive**: a mechanical validator can be built from the 9 `validation.md` rules (P16).
- **Positive**: Identity conflicts are never resolved silently — the lossless round-trip invariant is preserved (P13).
- **Negative**: these conventions bind — every new structure must pass validation; export/import must declare the contract version.

## Alternatives Considered

- **No exchange seam** — rejected: EKA 1.3/13.1 require lossless exchange support on every implementation.
- **Bespoke exporter/importer written later** — rejected: the seam must be defined at the contract level from the start, not retrofitted; Knowledge OS integration requires an explicit boundary (EKA 4.2).
- **Exchange conventions in the standard, not in the repository** — rejected: the standard establishes the contracts; the repository establishes their concrete serialization (EKA 12.4, 13.3).

## References

- EKA 1.3, 1.4, 4.2, 13.1 (what must be preserved), 13.2 (round-trip requirements), 13.3 (serialization format contracts)
- Principles P13 (Lossless Exchange), P16 (Enforcement Capability Varies, Invariants Don't)
- Related: [ADR-005](adr-005-dimension-layout.md), [ADR-007](adr-007-extension-research-finding.md)
