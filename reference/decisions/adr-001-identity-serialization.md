---
namespace: eka-ref-impl
type: adr
id: 001-identity-serialization
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
depends-on: []
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

# ADR-001 — Identity Serialization: explicit type token + full Identity in frontmatter

## Context

The initial implementation encoded four Artifact types (scope definition, plan, Execution Container, ticket) in **one shared ID space** with the same prefix (`mvp-nnn`). As a result, the Type could not be determined deterministically from the Identity representation: the same representation could be read as either a scope definition or a plan. This is a case study of an EKA 6.4 violation: rules 6.2.1–6.2.2 are violated (Type not explicit → ID not unique per `(Namespace, Type)`), and rule 6.2.3 is violated conceptually (Identity encoded via location and representation conventions, not Artifact properties). Binding lesson: **Identity must not be encoded in location, process stage, or representation conventions** (EKA 6.4).

## Decision

Identity serialization on this repository:

1. **Full Identity lives in the frontmatter**: `namespace`, `type`, `id`, `instance-version`, `revision` (EKA 6.4, P3, P9). The frontmatter is the source of truth for Identity; references are always by Identity.
2. **The filename is a projection of Identity**, not Identity itself (P9): pattern `<type-token>-<id>[-v<nn>].md`, with `<type-token>` = explicit type token and `<id>` = ID unique within `(Namespace, Type)`. The `-v<nn>` suffix is **mandatory** for versioned types (`scp-`, `plan-`) — always, including v1 — and **forbidden** for other types.
3. **Table of 26 type tokens** (ambiguity-free tokens: no token is a prefix of another token; corrected anti-prefix pairs: `sto-`/`str-`, `spk-`/`spec-`):

| Token | Artifact Type | Dimension |
|---|---|---|
| `vis-` | Vision / Manifesto | Product Intent |
| `str-` | Strategy | Product Intent |
| `req-` | Requirement (PRD) | Requirements |
| `scp-` | Scope Definition | Planning + Requirements |
| `epc-` | Epic | Planning Knowledge |
| `plan-` | Plan (roadmap) | Planning Knowledge |
| `ctr-` | Execution Container | — (OS) |
| `tkt-` | Ticket | — (OS, projection) |
| `sto-` | Work Item: Story | Requirements / Records / Research |
| `ts-` | Work Item: Technical Story | Requirements / Records / Research |
| `bug-` | Work Item: Bug | Requirements / Records / Research |
| `td-` | Work Item: Tech Debt | Requirements / Records / Research |
| `ch-` | Work Item: Chore | Requirements / Records / Research |
| `spk-` | Work Item: Spike | Requirements / Records / Research |
| `ses-` | Session | — (OS) |
| `rvw-` | Review | Governance & Quality |
| `adr-` | ADR | Decisions |
| `dec-` | Decision Record | Decisions |
| `arc-` | Architecture Description | Architecture |
| `spec-` | Specification | Specifications |
| `std-` | Standard / Guideline | Standards & Guidelines |
| `run-` | Runbook / Operational Guide | Operational Knowledge |
| `rel-` | Release Record | Records |
| `gls-` | Glossary / Term | Vocabulary |
| `trc-` | Traceability / Relationship Artifact | Planning Knowledge |
| `fnd-` | Research Finding (extension — ADR-007) | Research |

4. **Deterministic parsing**: the Identity representation can be parsed unambiguously from the frontmatter; the filename is validated as consistent with the frontmatter, not the other way around.

## Consequences

- **Positive**: deterministic Identity parsing; the `mvp-nnn` collision is resolved (ID unique per `(Namespace, Type)` — EKA 6.2.2); Identity decoupled from location, process stage, and phase (P3, P9).
- **Positive**: human navigation stays easy (type token in the filename) without sacrificing Identity correctness.
- **Negative (intentional)**: all legacy naming patterns break (`mvp-*`, `sp-*`, etc.) — legacy consumers must migrate (see `reference/breaking-changes.md`).
- **Negative**: every file now carries an Identity frontmatter whose consistency must be maintained — closed by mechanical validation (`dimension == folder`, valid tokens, etc., ADR-005/006).

## Alternatives Considered

- **Shared prefix + folder discriminator** (legacy status quo) — rejected: Identity still encoded via location; violates P3/P9 and EKA 6.4.
- **Type as suffix** (e.g., `id-type.md`) — rejected: parsing ambiguous with multi-word kebab-case IDs; per-type globbing becomes difficult.
- **Type only in frontmatter, free-form filename** — rejected: human navigation and validation weaken; a consistent filename projection aids determinism without becoming a source of truth.

## References

- EKA 6.1 (Identity composition), 6.2 (Identity rules), 6.3 (version semantics), 6.4 (collision case study)
- Principles P3 (Stable Identity), P9 (Structure as Projection of State)
- Related: [ADR-002](adr-002-state-vector-encoding.md) (state in frontmatter), [ADR-005](adr-005-dimension-layout.md) (folder = dimension), [ADR-007](adr-007-extension-research-finding.md) (`fnd-` token)
