# Engineering Knowledge Architecture (EKA)

**An open engineering knowledge standard — a canonical architecture for how engineering knowledge is identified, structured, exchanged, and operated.**

EKA is not a documentation template, not a Markdown repository scheme, and not a project management tool. It is a **standard**: a set of contracts — Identity, State, Knowledge Taxonomy, Layer Contracts, and Exchange — that any implementation can follow. This repository is one implementation of that standard, built to prove it works.

---

## Why EKA?

Engineering teams accumulate knowledge in an unstructured way: documents, decisions, plans, tickets, and records scattered across tools, each with its own conventions. Three problems recur:

1. **Identity is unstable** — documents are referenced by location ("the file in the decisions folder") instead of by a stable identity, so renames and reorganizations silently break references.
2. **State is implicit or duplicated** — the same work item has "status" in a ticket, a spreadsheet, and a document, with no single writer and no transition history.
3. **Knowledge cannot move between systems** — exporting from one tool and importing into another loses identity, state, and relationships.

EKA solves these by defining the conceptual model *before* any tool: a small, stable set of concepts (Artifact, Identity, State, Relationship) with explicit contracts, so that knowledge survives changes in storage, tooling, and organization.

## What is EKA for?

- **Engineers and engineering leads** who want durable, queryable, governed engineering knowledge.
- **Tool builders** who want an interoperable knowledge model instead of bespoke formats.
- **Platforms and AI agents** that need deterministic, machine-readable engineering knowledge.
- **Organizations** that want knowledge to outlive any single tool.

## Core Principles

| # | Principle | Meaning |
|---|---|---|
| P3 | **Stable Identity** | Identity is immutable and independent of location, storage, state, and classification. References always use Identity, never location. |
| P6 | **Single Writer** | Every State field has exactly one owner. Everything else is a projection, generated or validated — never independently edited. |
| P7 | **Forward-Only Transitions** | State domains move forward without regression. Corrections happen through new instances + relationships, not mutation. |
| P13 | **Lossless Exchange** | Exchange between systems never loses or duplicates Identity, State, Content, or Relationships. |
| P14 | **Minimum Canonical Core** | The standard defines concepts and contracts; implementations choose mechanisms. |
| P16 | **Enforcement Capability Varies, Invariants Don't** | Different implementations enforce differently; the invariants they enforce are identical. |

## Architecture Overview

EKA is organized as three layers, bound by Artifact Identity:

| Layer | Role | Owns | Owns nothing |
|---|---|---|---|
| **Knowledge Layer (KB)** | Stores knowledge | Content, classification, relationships, history, Identity Registry | Process state, execution protocol |
| **Operating Layer (OS)** | Executes engineering process | Execution/Planning/Container/Existence State, protocol, gates, commands | Content (never edits content) |
| **Exchange Layer (EX)** | Boundary transformation | Exchange contracts, round-trip rules, conformance validation | Content and State (never an owner) |

Every Artifact carries an **Identity** `(Namespace, Type, ID, InstanceVersion)`, a **State Vector** of owned state domains (Content, Execution, Planning, Container, Existence), **Content**, and **Relationships** (supersedes, amends, derives-from, depends-on, validates) — all by Identity.

## Specifications

| Specification | Status | Contents |
|---|---|---|
| **EKA Core Specification** v1.1 | Ratified | The canonical conceptual model: principles (P1–P16), Identity Model, State Taxonomy, Knowledge/Execution/Artifact Taxonomies, Layer Contracts, Import/Export Model, Extension Model; v1.1 adds §8.1 Engineering Domains and Knowledge Stratification (five Engineering Domains, Stratum Authority Invariant, Representation Aliases). |
| **EKA Exchange Specification** v1.1 | Ratified | The canonical exchange protocol: Exchange Units, Exchange Package Object Model, identity/relationship/state representation, versioning, import/export/synchronization semantics, Conformance Rules (R1–R9; v1.1 adds the Engineering Domain to unit classification and aligns the rule set to R0–R12), round-trip guarantees. Serialization-independent. |
| **EKA Naming and Terminology Specification** v1.1 | Ratified | The official naming system of the ecosystem: product identity, specification families, reference components, tooling, repository naming, canonical terminology (incl. Engineering Domain, Knowledge Stratum, Representation Alias), Conformance Rules R0–R12, deprecated terms. |
| **EKA Reference Serialization Format (RSF)** v1.1 | Reference (not normative) | One canonical serialization projection of the Exchange Package Object Model (incl. the Engineering Domain in unit classification) — the format used by `eka export` and `eka import`. |

The standard is deliberately serialization-independent: Git+Markdown is one implementation, and relational databases, graph stores, or future platforms are equally valid as long as they honor the contracts.

### Engineering Domains

Every artifact belongs to exactly one of five **Engineering Domains** (Core v1.1 §8.1–8.2) — the stratum-aligned category of engineering knowledge it holds: **Discovery** (intent, requirements, research) → **Architecture** (architecture, decisions, specifications, standards, vocabulary) → **Planning** (planning) → **Execution** (quality + operating) → **Operations** (operations, records). The domains form a strict authority chain (stratum 1 highest → 5): lower-stratum knowledge must not contradict higher-stratum knowledge in force, and never supersedes or amends upward (Stratum Authority Invariant). PRD, ADR, Epic, Sprint, Ticket, Release, and similar terms are **Representation Aliases** onto canonical tokens — see the [Representation Alias Registry](standard/representation-alias-registry-v1.1.md); methodologies (Scrum, Kanban, …) are convention layers, independent of the Core standard. How knowledge flows through these domains: [Engineering Knowledge Lifecycle](skeleton/docs/lifecycle.md). New users should start with the [Engineering Operating Guide](skeleton/docs/workflow-guide.md).

## Reference Implementation

This repository is the official **Reference Implementation**: a working serialization (Git + Markdown) of the standard, plus:

- **The Standard Zone** (`standard/`) — the canonical specification texts.
- **The Skeleton Zone** (`skeleton/`) — the copyable project structure (`eka init` generates repositories from it).
- **The Reference Zone** (`reference/`) — implementation meta-documentation, ADRs, and traceability matrices.
- **The Reference Validator** — the `eka` CLI, the canonical executable form of the Conformance Rules.

## CLI Overview

The `eka` CLI is the official interface of the architecture (Cobra-based command layer over reusable application packages):

| Command | Purpose |
|---|---|
| `eka init` | Repository Bootstrapper: analyzes the workspace, adaptively configures, generates an EKA repository from the Reference Skeleton, validates the result. Idempotent; `--dry-run` supported. |
| `eka validate` | Conformance validator: runs Conformance Rules R0–R12 mechanically, with deterministic output and exit codes (0/1/2). |
| `eka export` | Exports engineering knowledge into a canonical Exchange Package following the RSF — deterministic, validated before export, scopes: repository / line / instance / collection. |
| `eka import` | Integrates an Exchange Package into an existing repository — atomic, conservative merge, conflict detection, rollback, post-import validation. |
| `eka view` | Knowledge projections: `execution`, `planning`, `architecture`, `discovery`, `operations`, `ticket` — read-only, relationship-derived projections over the Engineering Knowledge Model (never markdown text); `sprint` / `wave` are CLI aliases of `execution`; conformance-gated, deterministic. |
| `eka completion` | Shell completion (bash/zsh/fish/powershell). |
| `eka version` | CLI build version and the EKA standard version implemented. |
| `eka` | Product landing: a calm orientation with a compact command overview (help and version pointers). |

All commands are deterministic: identical input produces identical output.

## Quick Start

```sh
# 1. Initialize a new EKA repository
eka init my-project
cd my-project

# 2. Create your first artifact (e.g., an ADR under docs/decisions/)
#    — follow the conventions in docs/README.md

# 3. Validate conformance
eka validate

# 4. Export a portable knowledge package
eka export

# 5. Import it elsewhere
eka import ./rsf-repo-my-project-1.ekapkg
```

New to EKA? Read the [Engineering Operating Guide](skeleton/docs/workflow-guide.md) first — the primary onboarding document: the twelve-part journey from mental model and knowledge lifecycle through engineering domains, daily and AI workflows, projections, and the CLI.

## Installation

Requires **Go 1.24+**.

```sh
go install github.com/maleolabs/engineering-knowledge-architecture/cmd/eka@latest
# or build from source:
git clone <this-repository>
cd <repo>
go build -o eka ./cmd/eka
```

## Repository Structure

```
standard/          Canonical specification texts (Core, Exchange, Naming, Glossary)
skeleton/          Copyable project serialization (docs structure, conventions)
reference/         Reference Implementation meta-documentation: architecture,
                   ADRs, CLI docs, ratification notes, traceability matrices
cmd/               CLI command layer (Cobra): root, init, validate, export, import, view
bootstrap/         Application layer: eka init engine (public package)
exchange/          Application layer: export/import engine (public package)
conformance/       Application layer: validation engine (public package)
view/              Application layer: knowledge projection engine (public package)
skeletonembed.go   Embedded Reference Skeleton (go:embed)
```

The application packages (`bootstrap/`, `exchange/`, `conformance/`, `view/`) are public and reusable independently of the CLI — by SDKs, MCP integrations, or other tools.

## Example Workflow

1. **Adopt EKA** in an existing project: `eka init` detects the workspace, reuses what exists, and only generates what is missing.
2. **Record decisions** as ADRs with stable identities; **track work** as work items whose Execution State is single-writer; **preserve history** via Records (superseded and archived artifacts are never deleted).
3. **Exchange knowledge** between repositories, organizations, or future platforms via Exchange Packages — lossless by contract, verified by digest.
4. **Automate enforcement** in CI: `eka validate` blocks non-conformant commits with deterministic verdicts.

## Roadmap

| Milestone | Status |
|---|---|
| EKA Core Specification v1.1 (Engineering Domains + Knowledge Stratification, §8.1) | Ratified |
| EKA Exchange Specification v1.1 | Ratified |
| EKA Naming and Terminology Specification v1.1 | Ratified |
| Reference Serialization Format v1.1 | Reference |
| Reference Implementation + Validator (rules R0–R12) | Active |
| `eka init`, `eka validate`, `eka export`, `eka import`, `eka view` | Implemented |
| `eka diagnose`, `eka graph`, sync strategies (replace, forward-only reconciliation) | Future |

## Contributing

EKA is an open standard. Contributions must follow the terminology governance (Naming and Terminology Specification) and the conformance governance (`CONTRIBUTING.md`): every change to specification, rules, implementation, or tests must keep the Conformance Traceability Matrix in sync, and every Pull Request must pass the validator and test suite.

## License

Apache-2.0. See [LICENSE](LICENSE).
