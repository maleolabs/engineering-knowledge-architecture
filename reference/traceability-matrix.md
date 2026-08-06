# Traceability Matrix — Repository Elements → EKA v1.0 Anchor

Complete matrix: every repository element (root README, all `standard/` files, all `reference/` files, the 7 ADRs, and — referencing sibling zones — every `skeleton/docs/` element) is mapped to its EKA anchor (section/principle/taxonomy) with rationale.

Type legend: **artifact** = file with `type` + `id` frontmatter (artifact rule); **convention** = convention document (no `type`/`id`).

## Standard Zone elements (A)

| Element | Zone | Type | EKA Anchor | Rationale |
|---|---|---|---|---|
| `standard/README.md` | Standard Zone | convention | 1.2, 4 | Explains the pre-layer zone; the standard defines the layers, not project artifacts. |
| `standard/eka-specification-v1.0.md` | Standard Zone | convention (canonical text) | entire document (1–16) | Verbatim copy of the canonical standard; conformance reference. |
| `standard/eka-exchange-specification-v1.0.md` | Standard Zone | convention (canonical text) | 13, 5.4 (invariant 7), P13, 16.1 (milestone 2) | Exchange Contract v1 (Ratified): operationalizes Section 13 + invariant 5.4.7 into a round-trip, idempotency, referential integrity, schema versioning, conformance suite contract (R1–R9); the refinement pass adds the Exchange Package Object Model (§4.4), Capability Declaration (§4.5), single declaration location, complete rule coverage; subordinate terms defined in §4.1 without amending the canonical text. |
| `reference/ratification-notes-exchange-v1.0.md` | Reference Zone | convention | 16.1 (milestone 2), P13 | Exchange Specification v1.0 ratification report: refinement summary, terminology, new conceptual model, architectural readiness confirmation. |
| `reference/eka-reference-serialization-format-v1.0.md` | Reference Zone | convention | Exchange Spec 4.4, 10, 15, 17.1; Naming 6.1, 7 | RSF v1.0 — canonical serialization projection of the Exchange Package Object Model (reference, not normative): Package Model, Unit Entry, Content Representation Model, Attachment Model, Manifest, serialization principles, round-trip mapping, compatibility, conceptual examples, import/export implementation recommendations. |
| `standard/eka-naming-and-terminology-specification-v1.0.md` | Standard Zone | convention (canonical text) | 3, 14.2, P14 | Meta-specification (Ratified): official ecosystem naming — EKA product identity, "EKA \<Family\> Specification" pattern, reference components, tooling (`eka` + verb subcommands), repository naming (`eka-<component>`), canonical + deprecated term table, migration; does not amend the canonical text. |
| `standard/glossary.md` | Standard Zone | convention | 3 | Exact definitions of Core Concepts terms; not paraphrased. |

## Reference Zone elements (C)

| Element | Zone | Type | EKA Anchor | Rationale |
|---|---|---|---|---|
| `README.md` (root) | root | convention | 1.3, 16.2 | Repository identity: one serialization, not the architecture. |
| `reference/README.md` | Reference Zone | convention | 1.3 | Index of implementation meta-documentation. |
| `reference/reference-architecture.md` | Reference Zone | convention | 1.3, 12.4, 6.4 | Explains the serialization: zone→layers, conventions, artifact rule. |
| `reference/migration-guide.md` | Reference Zone | convention | 16.2, 6.4 | Legacy→new map + strategy; serialization conformance. |
| `reference/philosophy.md` | Reference Zone | convention | 1.1, 16.3 | Dual-layer narrative, first-class protocol, phase-as-context, single-writer. |
| `reference/terminology-glossary.md` | Reference Zone | convention | 3 | Implementation-level terms; canonical terms linked, not redefined. |
| `reference/breaking-changes.md` | Reference Zone | convention | 6.4, P9 | 14 deliberate breaking changes so legacy consumers cannot read Identity from location. |
| `reference/adr-summary.md` | Reference Zone | convention | 14.2.5 | Governance index of implementation decisions (accepted). |
| `reference/traceability-matrix.md` | Reference Zone | convention | 12.2, 12.4 | Conformance evidence: every element → standard anchor. |
| `reference/ratification-notes.md` | Reference Zone | convention | 16.1 (milestone 1) | EKA v1.0 ratification notes (stabilization pass, verbatim). |

### Implementation ADRs (`reference/decisions/`) — all artifacts (`type: adr`, `dimension: decisions`, status accepted)

| Element | Zone | Type | EKA Anchor | Rationale |
|---|---|---|---|---|
| `adr-001-identity-serialization.md` | Reference Zone (decisions/) | artifact | 6.2, 6.4, 6.3, P3, P9 | Identity in frontmatter; filename as projection; 26 ambiguity-free tokens. |
| `adr-002-state-vector-encoding.md` | Reference Zone (decisions/) | artifact | 7.1, 7.2, 7.3, 7.4, 5.2, P2, P6, P7 | 5 frontmatter state fields per owned domain; absence = not-applicable; change-log; legacy mapping. |
| `adr-003-projection-model.md` | Reference Zone (decisions/) | artifact | 7.4, 7.5, 9, 10, 15.5, P6, P9 | Ticket/container tables = State Projections; empty State Vector; on-read refresh. |
| `adr-004-phase-as-metadata.md` | Reference Zone (decisions/) | artifact | 3, 7.1, 7.5, 11.2, 11.3, P3 | Phase = frontmatter field on scp-/plan-; phase change = context update via gate. |
| `adr-005-dimension-layout.md` | Reference Zone (decisions/) | artifact | 8, 4.1, 14.2, P1, P9, P15 | 12 folders = 12 dimensions 1:1 + operating/ + exchange/; `dimension == folder`. |
| `adr-006-exchange-conventions.md` | Reference Zone (decisions/) | artifact | 13.1, 13.2, 13.3, P13, P16 | Exchange seam: validation.md (9 rules) + transfer.md (round-trip, Identity conflict, idempotency, schema versioning). |
| `adr-007-extension-research-finding.md` | Reference Zone (decisions/) | artifact | 8, 10, 11.4, 14.1, 14.2, P12 | Extension type `fnd-`: complete owned State Vector (Content, Existence); exchangeable. |

### Tooling elements — Go implementation (EKA validator)

| Element | Zone | Type | EKA Anchor | Rationale |
|---|---|---|---|---|
| `go.mod` | root (tooling) | convention (build) | P16 | Module `github.com/maleolabs/engineering-knowledge-architecture` (Go 1.24+); dependencies: `gopkg.in/yaml.v3`, `spf13/cobra` (CLI adapter), `golang.org/x/term` (wizard TTY detection); `go install`/`go build` CLI entry points. |
| `cmd/` (command layer) | root (tooling) | convention (CLI) | 13.3, P16, Naming §7 | Pure Cobra command definitions: `root.go` (root command + `Execute(args, stdin, stdout, stderr) int` + exit 0/1/2 mapping), `validate.go`, `init.go`; no domain logic; standard Cobra help + completion. |
| `cmd/eka/` | root (tooling) | convention (CLI) | 13.3, P16 | Thin entry point: `os.Exit(cmd.Execute(...))`; executable name `eka`. |
| `bootstrap/` | root (tooling) | convention (engine) | Naming §7, Skeleton | `eka init` engine (application layer, public package): Workspace Discovery, Bootstrap Planning, Interactive Wizard (adaptive, deterministically non-interactive via `x/term`), Repository Generation from the Reference Skeleton, post-generation validation. |
| `skeletonembed.go` | root (tooling) | convention (embed) | Naming §6.1 | `//go:embed skeleton` — canonical Reference Skeleton embedded for `eka init` (standalone binary, no hardcoded directory). |
| `conformance/` | root (tooling) | convention (engine) | 13, P16 | Canonical Conformance Rules implementation (validation.md): reusable public engine, independent of the CLI; entry `Validate(root) (*Report, error)`, `Scan(root) ([]Artifact, error)`, `ParseReference` (additive for exchange consumers). |
| `exchange/` | root (tooling) | convention (engine) | Exchange Spec 4.4, 10, 11, 12, 15; RSF | Import/export engine (application layer, public package): discovery, loading (via `conformance.Scan`), scope resolution (repo/line/instance/collection), external reference declaration, deterministic RSF projection, deserialization + integrity verification, identity/relationship resolution, conflict analyzer, integration engine (staged commit + rollback), identity charset guard (RSF §5.2.3). |
| `reference/cli.md` | Reference Zone | convention | 13, P16, Naming §7 | Official CLI documentation: philosophy, installation, `eka init` (5 stages, adaptive wizard, idempotency, dry-run, post-generation validation), `eka validate`, exit codes, shell completion, CLI architecture (Cobra adapter + application layer), contribution guide for new commands, roadmap. |
| `reference/conformance-notes.md` | Reference Zone | convention | 13, P16 | Traceability record of rules R0–R9 → EKA anchor → implementation location + 29 interpretation decisions (policy: documented before implementation). |
| `.gitignore` | root (tooling) | convention (hygiene) | — | Implementation hygiene: the built `eka` binary never enters VCS. |

## Skeleton Zone elements (B) — sibling zone, referenced

| Element | Zone | Type | EKA Anchor | Rationale |
|---|---|---|---|---|
| `skeleton/docs/README.md` | Skeleton Zone | convention | 1.3 | Project serialization entry point; source of truth for the structure. |
| `skeleton/docs/intent/` (`vis-`, `str-`) | Skeleton Zone | KB (folder) | 8 (Product Intent), 10 | Home of vision/manifesto + strategy (new types). |
| `skeleton/docs/requirements/` (`req-`) | Skeleton Zone | KB (folder) | 8 (Requirements), 10 | Requirement + amendments as instances with `amends`. |
| `skeleton/docs/architecture/` (`arc-`) | Skeleton Zone | KB (folder) | 8 (Architecture), 10 | Architecture Description. |
| `skeleton/docs/decisions/` (`adr-`, `dec-`) | Skeleton Zone | KB (folder) | 8 (Decisions), 10 | Single Decisions dimension: ADR + decision record. |
| `skeleton/docs/specifications/` (`spec-`) | Skeleton Zone | KB (folder) | 8 (Specifications), 10 | New dimension; separate from Vocabulary. |
| `skeleton/docs/standards/` (`std-`) | Skeleton Zone | KB (folder) | 8 (Standards & Guidelines), 10 | Conventions, separate from operational. |
| `skeleton/docs/operations/` (`run-`) | Skeleton Zone | KB (folder) | 8 (Operational Knowledge), 10 | Procedures only. |
| `skeleton/docs/quality/` (`rvw-`) | Skeleton Zone | KB (folder) | 8 (Governance & Quality), 10 | Review + findings; `validates`. |
| `skeleton/docs/planning/` (`scp-`, `epc-`, `plan-`, `trc-`) | Skeleton Zone | KB (folder) | 8 (Planning Knowledge), 10, 11.2 | Scope (with `phase`), epic, plan (Planning State), traceability. |
| `skeleton/docs/records/` (`rel-`) | Skeleton Zone | KB (folder) | 8 (Records), 10 | Release record; immutable. |
| `skeleton/docs/research/` (`fnd-`) | Skeleton Zone | KB (folder) | 8 (Research), 11.4, 14.1 | Extension type; mandatory Distillation path. |
| `skeleton/docs/vocabulary/` (`gls-`) | Skeleton Zone | KB (folder) | 8 (Vocabulary), 10 | Glossary; not Specifications. |
| `skeleton/docs/operating/` | Skeleton Zone | OS (folder) | 4.1, 5.2 | OS layer: state machine & protocol. |
| `skeleton/docs/operating/containers/` (`ctr-`) | Skeleton Zone | OS (folder) | 10, 7.2, 7.5, 9 | Execution Container; Container State; exactly-one-active. |
| `skeleton/docs/operating/work-items/` (`sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-`) | Skeleton Zone | OS (folder) | 10, 7.2, 9 | Single-writer Execution State. |
| `skeleton/docs/operating/projections/` (`tkt-`) | Skeleton Zone | OS (folder) | 7.4, 10, P6 | Ticket: empty State Vector, projection. |
| `skeleton/docs/operating/sessions/` (`ses-`) | Skeleton Zone | OS (folder) | 10, 11.4 | Existence State; Distillation mandatory before Archived. |
| `skeleton/docs/operating/protocol.md` | Skeleton Zone | convention | 9, 11, 5.2 | Protocol definition, ordering, gates, commands. |
| `skeleton/docs/exchange/validation.md` | Skeleton Zone | convention | 13, P16 | 9 Conformance Rules. |
| `skeleton/docs/exchange/transfer.md` | Skeleton Zone | convention | 13.2, P13 | Round-trip, Identity conflict policy, idempotency, schema versioning. |

## Serialization conventions (implemented rules)

| Convention | EKA Anchor | Rationale |
|---|---|---|
| Identity encoding (frontmatter + filename `<type-token>-<id>[-v<nn>]`) | 6.2, 6.4, P3 | Identity canonical, unambiguous, machine-parseable; decoupled from location. |
| State encoding (5 fields per owned domain; absence = not-applicable) | 7.2, P6 | Explicit state, single-writer; non-applicable domains not serialized. |
| Phase as metadata | 11.2 | Context attribute, not category/state; authorized by the readiness gate. |
| Relationships by Identity (`supersedes`, `amends`, `derives-from`, `depends-on`, `validates`) | 6.2.7, 13.2.3 | References never by location; cross-system referential integrity. |
| Classification `dimension == folder` | 8, P15 | Classification is an artifact property; reclassification does not break references. |
| Projections (tickets, container tables) | 7.4, P6, P9 | Projections are not writers; refreshed against the owner. |
| Artifact rule (`type` + `id` ⇒ Artifact) | 3, 5.1 | Distinguishes artifacts from convention documents; layer ownership identity. |
| Change-log (`{date, domain, from, to, by}`) | 5.2 | Mandatory record of all state transitions. |
| Exactly-one-active (container) | 9 | Execution Container mutual-exclusion concurrency. |
| Lock-atomic-with-generation | 5.2 | Container creation atomic with plan lock (Planning State → Immutable). |
| On-read refresh | 15.5 | Default Projection Refresh policy; the projections-never-write invariant stays absolute. |
| Exchange validation (9 rules) | 13 | Validation before commit; conformance to the standard. |
| Extension `fnd-` (Research Finding) | 14.1 | New type via the extension mechanism; complete owned State Vector (14.2.6). |
