# Reference Architecture — EKA v1.0 Serialization in Git + Markdown

This document explains **how this repository serializes EKA v1.0**: zone → layer mapping, implemented serialization conventions, artifact rules, and key implementation decisions.

This repository is **one serialization (Git+Markdown) of the standard — not its architecture** (EKA 1.3).

## 1. Zone → layer mapping

| Zone | Contents | EKA Layer | Notes |
|---|---|---|---|
| `standard/` | Canonical EKA v1.0 texts + canonical glossary | **Pre-layer** | The standard defines the layers; it is not an artifact of any project. |
| `skeleton/docs/` — 12 knowledge dimension folders | intent, requirements, architecture, decisions, specifications, standards, operations, quality, planning, records, research, vocabulary | **Knowledge Layer (KB)** | Content, classification, Relationship, Records, Identity. |
| `skeleton/docs/operating/` | containers (`ctr-`), work-items (`sto-`/`ts-`/`bug-`/`td-`/`ch-`/`spk-`), projections (`tkt-`), sessions (`ses-`), protocol | **Operating Layer (OS)** | State Domains (Execution, Planning, Container, Existence), Protocol, Gate, Command. |
| `skeleton/docs/exchange/` | `validation.md`, `transfer.md` | **Exchange Layer (EX)** | Round-trip contract, conformance validation, schema versioning. |
| `reference/` | Meta-documentation of this implementation | — | Convention documentation, not part of the project serialization. |

The three layers are bound by Identity `(Namespace, Type, ID, InstanceVersion)` and the 7 global invariants (EKA 5.4).

## 2. Serialization conventions

### 2.1 Identity encoding

- **Location of truth**: artifact frontmatter — `namespace`, `type`, `id`, `instance-version`, `revision` (EKA 6.4, P3, P9).
- **Filename is a projection**: pattern `<type-token>-<id>[-v<nn>].md` — explicit type token, ID collision-free per `(Namespace, Type)`. The `-v<nn>` suffix is **mandatory** for versioned types (`scp-`, `plan-`) — always, including v1 — and **forbidden** for other types. The filename is only an Identity projection for human navigation + consistency validation; the true Identity lives in frontmatter.
- **26-type-token table** (ambiguity-free: no token is a prefix of another; anti-prefix pairs: `sto-`/`str-`, `spk-`/`spec-`):

| Token | Artifact Type | Dimension | Token | Artifact Type | Dimension |
|---|---|---|---|---|---|
| `vis-` | Vision / Manifesto | Product Intent | `ses-` | Session | — (OS) |
| `str-` | Strategy | Product Intent | `rvw-` | Review | Governance & Quality |
| `req-` | Requirement (PRD) | Requirements | `adr-` | ADR | Decisions |
| `scp-` | Scope Definition | Planning + Requirements | `dec-` | Decision Record | Decisions |
| `epc-` | Epic | Planning Knowledge | `arc-` | Architecture Description | Architecture |
| `plan-` | Plan (roadmap) | Planning Knowledge | `spec-` | Specification | Specifications |
| `ctr-` | Execution Container | — (OS) | `std-` | Standard / Guideline | Standards & Guidelines |
| `tkt-` | Ticket | — (OS, projection) | `run-` | Runbook / Operational Guide | Operational Knowledge |
| `sto-` | Work Item: Story | Requirements / Records / Research | `rel-` | Release Record | Records |
| `ts-` | Work Item: Technical Story | Requirements / Records / Research | `gls-` | Glossary / Term | Vocabulary |
| `bug-` | Work Item: Bug | Requirements / Records / Research | `trc-` | Traceability / Relationship | Planning Knowledge |
| `td-` | Work Item: Tech Debt | Requirements / Records / Research | `fnd-` | Research Finding (extension, ADR-007) | Research |
| `ch-` | Work Item: Chore | Requirements / Records / Research | | | |
| `spk-` | Work Item: Spike | Requirements / Records / Research | | | |

### 2.2 State Vector encoding

- **Five frontmatter fields**, one per **owned** state domain (EKA 7.4): `content-state`, `execution-state`, `planning-state`, `container-state`, `existence-state`.
- **Absence = not-applicable**: fields are present only for domains the artifact type owns (example: ADR = `content-state` + `existence-state`; work item = `execution-state` + `existence-state`; ticket = no state fields at all).
- **Value sets** (EKA 7.2, lowercase values):
  - `content-state` (variant per type, EKA 7.2): living `draft | review | approved | amended`; ADR `proposed | accepted | superseded`; decision record `draft | accepted | superseded`
  - `execution-state`: `planned | todo | in-progress | in-review | done`
  - `planning-state`: `draft | approved | immutable`
  - `container-state`: `active | completed` (completed = derived transition)
  - `existence-state`: `active | archived | retired`
- **Single-writer per field** (P6): every state field has exactly one owner; other views are projections.
- **`change-log`**: array `{date, domain, from, to, by}` — mandatory record of all state transitions (EKA 5.2).

### 2.3 Phase as context

- Field `phase` on `scp-`/`plan-` artifacts only, values `discovery|mvp|milestone|release|growth|maturity|sunset` (EKA 11.2, ADR-004).
- Phase change = context update authorized by the readiness gate; recorded in `change-log` with `domain: phase`. No phase folders.

### 2.4 Relationship

- Relations encoded by Identity in frontmatter: `supersedes`, `amends`, `derives-from`, `depends-on`, `validates` (EKA 6.2.7, 13.2.3). References are never by location (P3).

### 2.5 Classification

- Field `dimension` in frontmatter; knowledge artifacts live in their dimension folder — the `dimension == folder` rule is enforced by validation (EKA 8, P15, ADR-005). Operating artifacts (`operating/`) are exempt.

### 2.6 Projection

- Container and ticket tables (`tkt-`) are State Projections (EKA 7.4): tickets carry an empty State Vector + `derives-from`, generated artifacts carry the header "Generated — State Projection. Do NOT edit state here; refresh on read."; default refresh policy on-read (EKA 15.5, ADR-003).

### 2.7 Well-formed content

- Content follows the per-artifact-type structure (skeleton per folder) so it is machine-parseable and deterministic (EKA 3, 5.3).

## 3. Artifact rule vs convention documents

> **A file is an Artifact iff its frontmatter carries `type` AND `id`.**

- **Artifact**: carries full Identity + State Vector per its type; managed by the Operating Layer; exchangeable.
- **Convention document** (examples: `README.md`, `operating/protocol.md`, `exchange/validation.md`, `exchange/transfer.md`): a document that **explains** conventions — no `type`/`id`, carries no Identity, not part of the state machine, not exchanged as an Artifact.

Convention documents are recognizable by the absence of the `type`+`id` pair in frontmatter.

## 4. Summary of key implementation decisions

| Decision | EKA Anchor | ADR |
|---|---|---|
| Identity in frontmatter; filename = projection; explicit tokens | 6.4, P3, P9 | [ADR-001](decisions/adr-001-identity-serialization.md) |
| State = 5 frontmatter fields per owned domain; change-log | 7.2, 5.2, P6 | [ADR-002](decisions/adr-002-state-vector-encoding.md) |
| Ticket/container tables = projections; on-read refresh | 7.4, 15.5, P6 | [ADR-003](decisions/adr-003-projection-model.md) |
| Phase = frontmatter metadata, not folders | 11.2, P3 | [ADR-004](decisions/adr-004-phase-as-metadata.md) |
| 12 folders = 12 dimensions 1:1 + operating/ + exchange/ | 8, P15 | [ADR-005](decisions/adr-005-dimension-layout.md) |
| Exchange seam = validation.md + transfer.md | 13, P13 | [ADR-006](decisions/adr-006-exchange-conventions.md) |
| Extension type `fnd-` (Research Finding) | 14.1, 14.2, 11.4 | [ADR-007](decisions/adr-007-extension-research-finding.md) |

## References

- Canonical standard: [`../standard/eka-specification-v1.0.md`](../standard/eka-specification-v1.0.md)
- Copyable structure: [`../skeleton/README.md`](../skeleton/README.md)
- Migration map: [`migration-guide.md`](migration-guide.md)
- Breaking changes: [`breaking-changes.md`](breaking-changes.md)
