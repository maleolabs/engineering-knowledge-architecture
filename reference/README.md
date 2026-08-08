# Reference Zone — Meta-Documentation of the EKA v1.1 Implementation

This zone contains the **meta-documentation** of the EKA v1.1 Reference Implementation: how this repository serializes the standard, why implementation decisions were made, what changed from the legacy structure, and the decision trail (ADR).

## Documentation status

| Document | Status | Notes |
|---|---|---|
| [`reference-architecture.md`](reference-architecture.md) | Active | Repository serialization architecture (zone → layers, serialization conventions, artifact rule). |
| [`runtime-architecture.md`](runtime-architecture.md) | Active | Knowledge Runtime Architecture (v0.2.0): EKA Workspace + canonical store + Knowledge Snapshots + synchronization protocol; Authoring · Compiler · CKO · Presentation pipeline; Git/EKA/Atrium separation; CLI behavior review; future roadmap; known limitations. Convention document (ADR-009/ADR-010/ADR-011/ADR-012 summarized). |
| [`cko-specification.md`](cko-specification.md) | Active | The Canonical Knowledge Object (CKO) specification (v0.2.0): the one canonical internal representation of one Engineering Knowledge Object (`exchange.Unit` / `unit.json` + representation-tagged content) — field reference, derived values (domain/stratum), integrity metadata (object_hash = RSF per-unit digest), representation independence, the Runtime Contract, relationship to the RSF. Convention document (ADR-012/ADR-011 summarized). |
| [`migration-guide.md`](migration-guide.md) | Active | Complete migration map legacy structure → EKA structure + step-by-step strategy. |
| [`migration-report-engineering-domains-v1.1.md`](migration-report-engineering-domains-v1.1.md) | Active | Migration classification of the Engineering Domain ontology v1.0 → v1.1 (zero breaking changes). |
| [`migration-report-runtime-v0.2.0.md`](migration-report-runtime-v0.2.0.md) | Active | Migration report of the Knowledge Runtime (v0.2.0): before/after storage model, compatibility guarantees, docs → sync migration path, verification evidence, open risks. |
| [`philosophy.md`](philosophy.md) | Active | Narrative of why EKA exists and why this repository is structured this way. |
| [`terminology-glossary.md`](terminology-glossary.md) | Active | Glossary of implementation-level terms (canonical terms live in `standard/glossary.md`). |
| [`breaking-changes.md`](breaking-changes.md) | Active | Summary of the 14 breaking changes against the legacy structure. |
| [`adr-summary.md`](adr-summary.md) | Active | Index of the 12 Implementation ADRs (zone `decisions/`). |
| [`traceability-matrix.md`](traceability-matrix.md) | Active | Traceability matrix: every repository element → EKA anchor. |
| [`ratification-notes.md`](ratification-notes.md) | Active | EKA v1.0 ratification notes (verbatim from the stabilization pass). |
| [`cli.md`](cli.md) | Active | Official `eka` CLI documentation: philosophy, `eka init` (5-stage bootstrapper), `eka export` (Exchange Package, RSF), `eka import`, `eka view`, `eka watch`, `eka sync` (Knowledge Runtime synchronization), `eka project` (registry), `eka status`, `eka validate`, exit codes, shell completion, CLI architecture (Cobra adapter), contribution guide for new commands, roadmap. |
| [`conformance-notes.md`](conformance-notes.md) | Active | Conformance implementation notes: 29 interpretation decisions + known gaps (traceability tables consolidated into the Conformance Traceability Matrix). |
| [`conformance-traceability-matrix.md`](conformance-traceability-matrix.md) | Active | Single source of truth for conformance coverage: REQ→Spec→Rule→Impl→Test→Coverage→Notes (R0–R12, 67 cited tests, 19 requirements (13 enforced + 6 governance-only)). |
| [`reference-project.md`](reference-project.md) | Active | The EKA Reference Project guide: purpose, Feather project overview (namespace `feather`, ~50% complete), repository tree, Engineering Domain map, knowledge relationship walkthrough, CLI demonstration scenarios with expected output, screenshot checklist, integration notes. |
| [`project/`](project/) | Active | The EKA Reference Project itself — a complete, conformant example repository ("Feather", a markdown blogging platform): the full `eka init` skeleton tree plus 37 populated artifacts (6 Discovery, 8 Architecture, 5 Planning, 15 Execution, 3 Operations), validating with 0 errors and 0 warnings, plus a committed Knowledge Snapshot at `exchange/snapshots/` (label `rsf-repo-feather-1.1`, 37 units). Demonstrates Engineering Knowledge in practice: every token family in use, R10 satisfied for every non-Discovery artifact, tickets/sessions as projections, export/import round-trip ready, sync-ready. |
| [`eka-reference-serialization-format-v1.1.md`](eka-reference-serialization-format-v1.1.md) | Reference (not normative) | **RSF v1.1** — one canonical serialization projection of the Exchange Package Object Model (Exchange Spec §4.4): Package Model, Unit Entry (Classification incl. Engineering Domain), Content Representation (EKA Structured Text), Attachment Model, Manifest, deterministic serialization principles, round-trip mapping, compatibility, conceptual examples, implementation recommendations for `eka export`/`eka import`. |
| [`decisions/`](decisions/) | Active | 12 Implementation ADRs (serialization architecture decisions ADR-001…008 + Knowledge Runtime ADR-009…012). |

> Contribution governance + the definition of "incomplete implementation" live in [`../CONTRIBUTING.md`](../CONTRIBUTING.md) (repo root).

## Cross-zone references

| Zone | Role | Entry point |
|---|---|---|
| **A — `../standard/`** | EKA v1.1 canonical standard (pre-layer) | [`../standard/README.md`](../standard/README.md) |
| **B — `../skeleton/`** | Copyable project structure (Git+Markdown serialization) | [`../skeleton/README.md`](../skeleton/README.md) |
| **C — `../reference/`** | Meta-documentation of this implementation | this document |
| **D — `project/`** | The EKA Reference Project: a fully populated example repository (Feather) | [`reference-project.md`](reference-project.md) |
