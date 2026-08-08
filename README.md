# Engineering Knowledge Architecture (EKA)

**An open engineering knowledge standard — a canonical architecture for how engineering knowledge is identified, structured, exchanged, and operated.**

EKA is not a documentation template, not a Markdown repository scheme, and not a project management tool. It is a **standard**: a set of contracts — Identity, State, Knowledge Taxonomy, Layer Contracts, and Exchange — that any implementation can follow. This repository is one implementation of that standard, built to prove it works.

### Knowledge Runtime (v0.2)

Milestone **EKA v0.2.0** adds the **Knowledge Runtime Architecture**: a local **EKA Workspace** (`~/.eka/` or `$EKA_HOME`; `workspace.json` + `eka.db`, an embedded SQLite canonical store) where canonical Engineering Knowledge lives, with repositories as **synchronization endpoints** carrying deterministic **Knowledge Snapshots** (`exchange/snapshots/`, RSF directory packages). `eka sync` (pull/push), `eka project register`/`list`, `eka status`, and `eka integrity check` operate the runtime; all pre-existing commands behave exactly as before. The runtime consumes **Canonical Knowledge Objects** — compiled from Markdown via the Knowledge Compiler (`compile/`); Markdown is the authoring format, not the runtime representation (ADR-012). Projections (`eka view` / `eka watch`) read those CKOs from the workspace canonical store — run `eka sync` first (the projection covers the whole project, the union of its registered repositories; unregistered → exit 1) (ADR-013). The store implements the **Immutable Engineering Knowledge Model** (ADR-011): knowledge objects are content-addressed (`object_hash` = SHA-256(unit.json ‖ content), insert-only, never updated) with mutable references only, and `eka integrity check` verifies the store by recomputing every content-derived hash (manual database modification is detected, not prevented). Git stays the VCS — synchronization is explicit, no hooks. Experimental: workspace/snapshot/sync terminology is not finalized. See [Knowledge Runtime Architecture](reference/runtime-architecture.md), the [Canonical Knowledge Object specification](reference/cko-specification.md), and [ADR-009](reference/decisions/adr-009-knowledge-runtime-architecture.md) / [ADR-010](reference/decisions/adr-010-synchronization-model.md) / [ADR-011](reference/decisions/adr-011-immutable-engineering-knowledge-model.md) / [ADR-012](reference/decisions/adr-012-canonical-knowledge-object-runtime.md) / [ADR-013](reference/decisions/adr-013-store-backed-projections.md).

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

### Git · EKA Runtime · Atrium

Since v0.2.0, three runtimes share the one knowledge model above:

- **Git** — source code version control, never replaced: it versions code and the knowledge transport (snapshots commit like any other content).
- **EKA Knowledge Runtime** — the local, indexed runtime of canonical Engineering Knowledge: the EKA Workspace (`~/.eka/` or `$EKA_HOME`) holds the canonical store (`eka.db`); repositories are synchronization endpoints carrying deterministic Knowledge Snapshots (`exchange/snapshots/`), moved by `eka sync`. Canonical storage lives in the workspace; the repository is the transport.
- **Atrium** — the future unified project runtime: a consumer of the complete Engineering Knowledge of a multi-repository project (e.g. `api`/`web`/`mobile` under one project) from the runtime. Not implemented in v0.2; the architecture is shaped for it.

The full runtime architecture, sync protocol, and known limitations: [Knowledge Runtime Architecture](reference/runtime-architecture.md) (ADR-009, ADR-010).

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
- **The Reference Project** ([`reference/project/`](reference/project/)) — a complete, conformant example repository ("Feather", a markdown blogging platform): 37 artifacts across all five Engineering Domains, zero validation warnings, screenshots and CLI scenarios in [reference/reference-project.md](reference/reference-project.md). It demonstrates Engineering Knowledge in practice — read it as the worked example of the standard, not as software architecture complexity.
- **The Reference Validator** — the `eka` CLI, the canonical executable form of the Conformance Rules.

## CLI Overview

The `eka` CLI is the official interface of the architecture (Cobra-based command layer over reusable application packages):

| Command | Purpose |
|---|---|
| `eka init` | Repository Bootstrapper: analyzes the workspace, adaptively configures, generates an EKA repository from the Reference Skeleton, validates the result. Idempotent; `--dry-run` supported. |
| `eka validate` | Conformance validator: runs Conformance Rules R0–R12 mechanically, with deterministic output and exit codes (0/1/2). |
| `eka export` | Exports engineering knowledge into a canonical Exchange Package following the RSF — deterministic, validated before export, scopes: repository / line / instance / collection. The same object model as the Knowledge Snapshot used by `eka sync` (directory layout at `exchange/snapshots/`); `.ekapkg` remains the on-demand package form. |
| `eka import` | Integrates an Exchange Package into an existing repository — atomic, conservative merge, conflict detection, rollback, post-import validation. Unchanged under the runtime; synchronization uses the same verification path (`exchange.LoadPackage`) rather than import. |
| `eka view` | Knowledge projections: `execution`, `planning`, `architecture`, `discovery`, `operations`, `ticket`, `board` — read-only, relationship-derived projections over the Engineering Knowledge Model (never markdown text), read from the workspace canonical store — the repository must be registered and synced first (`eka sync`; unregistered → exit 1); `sprint` / `wave` are CLI aliases of `execution`; deterministic. `board` shows every work item across all containers with per-item container tags. |
| `eka watch` | Realtime projection viewer: the same projections as `eka view`, refreshed in place by polling the canonical store (`--interval`, default 2s, min 1s); TTY-only; live refusal frames (repository not registered) with auto-recovery; Ctrl-C to stop. Requires a registered + synced repository (`eka sync` first). |
| `eka sync [path]` | Knowledge Runtime synchronization: pull (verify snapshot → seed canonical store; or conformance-gated seed from the docs tree when no snapshot exists) then push (store → deterministic snapshot at `exchange/snapshots/`). Idempotent; deletions never applied; auto-registers the repository. |
| `eka project register [path] [--name NAME]` | Registers a repository in the EKA workspace under a project; same `--name` = same project (multi-repository projects). |
| `eka project list` | Lists the workspace's registered projects and repositories (deterministic). |
| `eka status` | EKA workspace status: path, schema version, store totals (objects/references, immutable payloads, attachments), per-repository last sync. Read-only probe; never creates the workspace. |
| `eka integrity check` | Verifies the workspace canonical store: recomputes every payload hash, strict-decodes every payload, verifies every reference (target + derived index columns), recomputes attachment digests, checks the repository registry. Read-only; unreferenced payloads count as history, never violations; manual database modification is detected, not prevented; exit codes 0 (clean) / 1 (violations) / 2 (internal). |
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

# 6. Knowledge Runtime: register the repository under a project
eka project register . --name my-project

# 7. Sync: seed the canonical store (~/.eka/eka.db) and write the
#    Knowledge Snapshot (exchange/snapshots/) — idempotent, re-run safe
eka sync

# 8. Inspect the workspace
eka status
```

The snapshot at `exchange/snapshots/` is ordinary repository content — commit it with normal Git workflows (no hooks; synchronization is explicit). More: [Knowledge Runtime Architecture](reference/runtime-architecture.md).

New to EKA? Read the [Engineering Operating Guide](skeleton/docs/workflow-guide.md) first — the primary onboarding document: the twelve-part journey from mental model and knowledge lifecycle through engineering domains, daily and AI workflows, projections, and the CLI.

## Installation

**Linux / macOS** — install the prebuilt binary from the GitHub Release asset registry (checksum-verified):

```sh
curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh
```

Install a specific version or custom directory:

```sh
curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh -s -- --version v0.1.0
curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh -s -- --to ~/.local/bin
```

The installer downloads `eka-<os>-<arch>` (Linux/macOS × amd64/arm64) from the release assets, verifies it against `SHA256SUMS.txt` (fail-closed: no unverified installs), and installs it to `/usr/local/bin` (or `--to`).

**Windows** (PowerShell 5.1+ or 7+) — install the prebuilt binary from the GitHub Release asset registry (checksum-verified):

```powershell
irm https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.ps1 | iex
```

Specific version or custom directory:

```powershell
$s = irm https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.ps1
iex "$s -Version v0.2.0"
iex "$s -To 'C:\tools\bin'"
```

The installer downloads `eka-windows-<arch>.exe` (amd64/arm64) from the release assets, verifies it against `SHA256SUMS.txt` (fail-closed), and installs it to `%LOCALAPPDATA%\Programs\eka` (or `-To`).

**Requires Go 1.24+** to build from source:

```sh
go install github.com/maleolabs/engineering-knowledge-architecture/cmd/eka@latest
# or:
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
cmd/               CLI command layer (Cobra): root, init, validate, export, import, view, watch, sync, project, status
bootstrap/         Application layer: eka init engine (public package)
exchange/          Application layer: export/import engine (public package)
conformance/       Application layer: validation engine (public package)
compile/           Application layer: the Knowledge Compiler — authoring → Canonical Knowledge Objects (public package)
view/              Application layer: knowledge projection engine (public package)
workspace/         Application layer: EKA workspace + project/repository registry (public package)
store/             Application layer: canonical store (SQLite, schema v2 — immutable content-addressed payloads + mutable references) (public package)
sync/              Application layer: synchronization engine (pull/push) (public package)
skeletonembed.go   Embedded Reference Skeleton (go:embed)
```

The application packages (`bootstrap/`, `exchange/`, `conformance/`, `compile/`, `view/`, `workspace/`, `store/`, `sync/`) are public and reusable independently of the CLI — by SDKs, MCP integrations, or other tools.

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
| `eka init`, `eka validate`, `eka export`, `eka import`, `eka view`, `eka watch` | Implemented |
| Knowledge Runtime (v0.2): `eka sync`, `eka project`, `eka status`, workspace + canonical store | Implemented (experimental) |
| `eka diagnose`, `eka graph`, sync strategies (replace, forward-only reconciliation), cloud sync, deletion protocol | Future |

## Contributing

EKA is an open standard. Contributions must follow the terminology governance (Naming and Terminology Specification) and the conformance governance (`CONTRIBUTING.md`): every change to specification, rules, implementation, or tests must keep the Conformance Traceability Matrix in sync, and every Pull Request must pass the validator and test suite.

## License

Apache-2.0. See [LICENSE](LICENSE).
