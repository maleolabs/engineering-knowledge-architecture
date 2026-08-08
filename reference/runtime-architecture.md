# Knowledge Runtime Architecture — EKA v0.2.0

> Convention document, not an artifact (no `type`/`id`). Implementation-focused companion to [`reference-architecture.md`](reference-architecture.md) (repository serialization) — it documents the **runtime** storage, synchronization, and consumption (compiler + Canonical Knowledge Object) layer added in milestone EKA v0.2.0.
> Authoritative decisions: [ADR-009](decisions/adr-009-knowledge-runtime-architecture.md) (workspace + embedded canonical store), [ADR-010](decisions/adr-010-synchronization-model.md) (synchronization protocol), [ADR-011](decisions/adr-011-immutable-engineering-knowledge-model.md) (the Immutable Engineering Knowledge Model — schema v2), [ADR-012](decisions/adr-012-canonical-knowledge-object-runtime.md) (the Canonical Knowledge Object Runtime — the consumption side, §2.1), and [ADR-013](decisions/adr-013-store-backed-projections.md) (store-backed projections — projections read the canonical store, §2.1/§8/§11). This document summarizes and explains; the ADRs are normative for this implementation.
> Terminology status: the names *workspace*, *Knowledge Snapshot*, *runtime*, and *sync* are **not finalized** — a future milestone may rename them. Everything here is experimental.

## 1. The mental model: three runtimes, one knowledge model

Since ADR-001, this repository's serialization (Git + Markdown) *was* the storage of Engineering Knowledge. The v0.2.0 runtime inverts that: the repository stops being the primary store and becomes a **transport medium**. Three runtimes now share one knowledge model:

```
┌─────────────────────────┐      ┌──────────────────────────────┐      ┌──────────────────────────┐
│         Git             │      │      EKA Knowledge           │      │        Atrium            │
│  source code version    │      │      Runtime (v0.2)          │      │   unified project        │
│  control                │      │                              │      │   runtime (future)       │
│                         │      │  canonical store (eka.db)    │      │                          │
│  code, review, release  │      │  workspace.json              │      │  consumes the complete   │
│  workflow               │      │  sync protocol               │      │  engineering knowledge   │
└───────────┬─────────────┘      └──────────────┬───────────────┘      │  of a project across     │
            │                                  │                       │  repositories            │
            │  commits the snapshot            │                       └──────────┬───────────────┘
            │  (ordinary repository content)   │  eka sync — pull/push           │
            ▼                                  ▼                                  │
┌───────────────────────────────────────────────────────────────────────────────────────┐
│                        One Engineering Knowledge Model                               │
│  Identity (Namespace, Type, ID, InstanceVersion) · State Vector · Content ·           │
│  Relationships · Classification (Dimension, Engineering Domain)                      │
│  Serialized as: Git+Markdown docs tree · RSF packages (.ekapkg) · Knowledge          │
│  Snapshots (RSF directory packages) · canonical store rows                            │
└───────────────────────────────────────────────────────────────────────────────────────┘
```

- **Git** remains the version control system — never replaced, never forked. It versions source code *and* the knowledge transport (the snapshot directory commits and reviews like any other content).
- **EKA Knowledge Runtime** (this milestone) is the local, indexed, queryable storage of canonical Engineering Knowledge: a workspace database on the machine, synchronized with repositories.
- **Atrium** is the future unified project runtime — a consumer that reads the complete Engineering Knowledge of a multi-repository project from the runtime (see [Future roadmap](#12-future-roadmap-not-implemented)). It is **not implemented** in v0.2; it is the architectural reason the runtime exists.

The division of responsibility is what changed, not the knowledge model: Identity, the State Vector, the RSF, the Knowledge Dimensions, and the Engineering Domains (ADR-008) are untouched by this milestone.

## 2. Core concepts

| Concept | Role | Location | Canonical? |
|---|---|---|---|
| **EKA Workspace** | machine-local runtime root: metadata + canonical store | `~/.eka/` (Linux/macOS) or `%USERPROFILE%/.eka/` (Windows), overridable via `EKA_HOME` | yes — the runtime's home |
| **Canonical Store** (`eka.db`) | SQLite database: immutable content-addressed payloads (`object_payloads`), mutable references + derived indexes (`object_refs`), attachments, sync log | inside the workspace | yes — canonical Engineering Knowledge |
| **Repository** | synchronization endpoint: `docs/` tree + Knowledge Snapshot | anywhere; registered in the workspace | no — transport |
| **Knowledge Snapshot** | RSF package in directory layout at `<repo>/exchange/snapshots/` | inside the repository, committed to Git | no — deterministic transport form |
| **Synchronization** | explicit pull/push between repository and store (`eka sync`) | protocol, not a location | — |
| **Docs tree** (`docs/`) | the human-readable Git+Markdown **authoring representation** — one authoring adapter for the Knowledge Compiler (Markdown is never a runtime representation) | inside the repository | no — authoring representation |

The pairing to remember: **canonical storage lives in the workspace; the repository is the transport**. Knowledge moves into the repository only as a snapshot, and moves out of it only through verified pull.

### 2.1 Authoring · Compiler · CKO · Presentation

The runtime's consumption path is representation-independent: **authoring is an input, the Canonical Knowledge Object (CKO) is the model, presentation is the output** ([ADR-012](decisions/adr-012-canonical-knowledge-object-runtime.md); the CKO schema is [`cko-specification.md`](cko-specification.md)). The pivot completes the inversion started by the immutable model: the store already persisted canonical objects (ADR-011), and the readers are now canonical too — Markdown is **one authoring adapter**, never a runtime representation. The pipeline (ADR-012 §Decision 7):

```
┌────────────────────┐
│      Authoring     │   Markdown docs/ — one authoring adapter (conformance scan/analyze)
└─────────┬──────────┘
          ▼
┌────────────────────┐
│ Knowledge Compiler │   compile/: parse → validate R0–R12 → normalize → generate CKO → verify
│                    │   (the authoring gateway — used by eka sync docs-mode / --from-docs)
└─────────┬──────────┘
          ▼
┌────────────────────┐
│        CKO         │   exchange.Unit — unit.json + representation-tagged content payload
└─────────┬──────────┘
          ▼
┌────────────────────┐
│     Runtime DB     │   object_payloads + object_refs (SQLite; CKO, never a Markdown cache)
└─────────┬──────────┘
          ▼
┌────────────────────┐
│      Resolver      │   store.UnitsByProject — object_refs → object_payloads → DecodeUnit
└─────────┬──────────┘
          ▼
┌────────────────────┐
│     Projection     │   view/ renderers — eka view / eka watch (store reads, zero Markdown)
└─────────┬──────────┘
          ▼
┌────────────────────┐
│ Human / Machine /  │   terminal projections today; MCP JSON · Atrium (future)
│        AI          │
└────────────────────┘
```

Role split — who parses what, and who validates what:

| Role | Location | Responsibility |
|---|---|---|
| **Authoring adapter** | `conformance/` (`Scan`/`analyzeFile`) | Markdown parsing + authoring conformance R0–R12 — the compiler's validation stage; the only Markdown-aware code in the runtime |
| **Knowledge Compiler** | `compile/` | the canonical gateway: read authoring → parse → validate syntax → normalize → generate CKO → integrity (package digest) → hand off (`compile.Compile`); representation-independent by construction; the **authoring gateway** — used by sync docs-mode/migration (ADR-010 §Decision 2), **no longer imported by `eka view` / `eka watch`** (ADR-013 §Decision 1) |
| **CKO consumers** | `view/` + `cmd` | projections over `exchange.Unit` read **from the workspace canonical store** — `eka view` / `eka watch` resolve the project via `store.UnitsByProject` (refs → immutable payloads → `exchange.DecodeUnit`, canonical-form order) and never read Markdown; the compile path is gone from `cmd/view.go` / `cmd/watch.go`; **zero Markdown parsing in the consumption path** (the `view` package imports no `conformance.Scan`/`Artifact`) |
| **Authoring validation** | `eka validate` | R0–R12 against the authoring representation (adapter layer; read-only, P6) |
| **Runtime validation** | `eka integrity check` | CKO integrity: payload hash, decode, references, attachments, workspace ([§6](#6-integrity)) |
| **CKO persistence** | `store/` | `object_payloads` (unit.json + content) + `object_refs` — never a Markdown cache ([§4](#4-canonical-store-schema-summary)) |

The ontology helpers remain shared and representation-independent: `conformance.ParseReference` (reference resolution), `DomainValues`/`OwnedDomains` (state-domain taxonomy), `DomainForToken`/`Stratum` (domain/stratum derivation). Presentation is a separate concern: terminal projections today; future presentations (MCP JSON, Atrium components) consume the same CKO contract. The docs-mode pull (sync) goes through the compiler — the authoring gateway; projections (`eka view` / `eka watch`) read the canonical store directly (`store.UnitsByProject`), so the runtime consumption path never touches Markdown and never re-runs the authoring gate per invocation (ADR-013).

## 3. Workspace layout

```
~/.eka/  (or $EKA_HOME — must be an absolute path)
├── workspace.json   metadata: schema_version, id, created
└── eka.db           SQLite canonical store (schema v2, the Immutable Engineering Knowledge Model)
```

- **`workspace.json`** — `schema_version: 1`, a **deterministic workspace id** (`"eka-" + first 12 hex chars of SHA-256(absolute workspace path)`), and the created date. The id derives from the path: stable across restarts, reproducible on any machine with the same path.
- **`eka.db`** — SQLite via `modernc.org/sqlite` (pure Go, no cgo; the build stays portable, including the CI-cross-compiled Windows binaries). Opened with `journal_mode(WAL)`, `busy_timeout(5000)`, `foreign_keys(1)`.
- Layout invariant: **one workspace, many projects, many repositories**. The workspace is machine-local: never committed, never exchanged, never part of any serialized output (determinism of all serialized bytes is unaffected).
- The workspace directory is created with mode `0700` on Unix and `workspace.json` with `0600` — the canonical store is private to the user.
- `EKA_HOME` must be absolute; a relative value is rejected (usage error, exit `2`).

## 4. Canonical store schema (summary)

Schema version `2` (the **Immutable Engineering Knowledge Model**, [ADR-011](decisions/adr-011-immutable-engineering-knowledge-model.md)), one project-aware database per workspace — explicitly **not** database-per-project (cross-project queries and the union model need one store). The full DDL is in [ADR-009 §Decision 3](decisions/adr-009-knowledge-runtime-architecture.md) as amended by ADR-011; the tables in one line each:

| Table | Role |
|---|---|
| `meta` | schema bookkeeping: `schema_version`, `workspace_id` |
| `projects` | registered projects (id = name) |
| `repos` | registered repositories, keyed `(project_id, name)`, with normalized absolute path |
| `object_payloads` | **the immutable unit of record**: one row per content-addressed Engineering Knowledge Object, keyed by `object_hash` = `SHA-256(unit.json ‖ content)` — byte-identical to the RSF per-unit digest; insert-only, never updated, never deleted; carries `prev_hash` lineage and first-write `created_at` |
| `object_refs` | **the mutable resolver**: one row per RSF Canonical Identity Form `<namespace>/<type>:<id>:<instance-version>`, pointing at the *current* immutable payload via `object_hash`; carries provenance (`project_id`, `source_repo`) and the derived index columns (`namespace`, `type`, `id`, `instance_version`, `revision`, `dimension`, `domain`, `phase`) |
| `attachments` | attachment payloads keyed by id, digest-verified |
| `sync_log` | every pull/push run: project, repo, direction, snapshot digest, unit count, timestamp |

Derived data is **not** stored separately: relationships and the change-log live inside the immutable `unit.json` payload (they are part of the RSF unit serialization) and are parsed at read time. The v1 `objects` / `relationships` / `change_log` tables are **removed**. No numeric canonical identity exists anywhere in the knowledge model — knowledge identity is the form (string) and the object hash (content); `sync_log`'s `AUTOINCREMENT` is operational run sequencing only.

Three design facts carry the philosophy forward:

- **Knowledge is immutable and content-addressed.** Every knowledge object is written once, keyed by its own content hash; the store has **no update path** for payloads — "modify an object" is not expressible (see [§5](#5-the-immutable-model)).
- **History is derived, not maintained.** Forward-only instance versions (new forms), `prev_hash` chains, and retained (unreferenced) payloads are the complete history — there is no dedicated history table to write, migrate, or keep consistent.
- **Indexes serve the future feature set.** Identity, dimension, domain, phase, and provenance index columns are point lookups; relationships are parsed from payloads at read time (no SQL-side graph table in v0.2); the content BLOB is ready for FTS5/`sqlite-vec` vector indexing. The schema leaves room — nothing beyond v0.2 is implemented (see [roadmap](#12-future-roadmap-not-implemented)).

## 5. The immutable model

ADR-011 made the immutability invariant construction-level: **Engineering Knowledge Objects are immutable and content-addressed**. This is not a storage optimization — it is what makes snapshots verifiable, history derivable, synchronization idempotent, and future replication content-addressed.

### 5.1 Content addressing: the object *is* its hash

- `object_hash` = `SHA-256(unit.json bytes ‖ content bytes)` — **byte-identical to the RSF per-unit digest** (`exchange/serialize.go`, `deserialize.go`; RSF §9.4). The same bytes hashed by the exchange layer, the store, or `eka integrity check` produce the same value, so object hashes agree with snapshot digests **by construction**.
- The hash is **content-derived, never DB-generated**: there is no database sequence, no trigger, no application-assigned id. The object *is* its address and its verification key.
- `content` is `NOT NULL` with a zero-length BLOB for empty payloads (never `NULL`): the digest covers `unit.json ‖ content` unconditionally, so the empty case is represented, not absent.

### 5.2 Insert-only store, mutable references only

- **`object_payloads` is insert-only**: `INSERT ... ON CONFLICT DO NOTHING` — a payload that already exists is never touched (**first-writer-wins**). No store API exposes an update or delete path for payload rows; the guarantee is in the model, not in SQLite constraints.
- **The only mutable rows in the knowledge store are `object_refs`**: a reference moves from one immutable payload to another (the current object of a form within its provenance). Everything else that mutates is operational, not knowledge: the `projects`/`repos` registry (registration bookkeeping), `sync_log` (run records), `attachments` (id-keyed payloads with digest verification).
- The refs row carries only identity, provenance, the derived classification columns (`dimension`, `domain`, `phase`), and the pointer. The full State Vector, author/created/updated, `content_representation`, and secondary dimensions are **not** duplicated: they live inside the immutable `unit_json` payload and are read from it.

### 5.3 Layering: SQLite is persistence only

```
Engineering Knowledge Model → Immutable Objects → Persistence (SQLite) → Indexes → Resolver → Projection / Context / AI
```

The persistence layer stores and indexes; it never defines Engineering Knowledge semantics. Concretely: the store has no update path for payloads, and hashes are content-derived — so a future storage engine swap changes nothing about the model or its verification. Every layer above persistence (indexes, resolver, projections) is a **derived view** over immutable objects — replaceable, rebuildable, and verifiable at any time. If immutability were enforced by database triggers, it would be coupled to SQLite and silently vanish on an engine swap; the model carries it intrinsically.

### 5.4 History is derived, not maintained

Knowledge history emerges from three sources, none of which requires maintenance:

1. **Forward-only instance versions = distinct forms.** A knowledge change produces a new instance version — a new form — never a mutation of a prior line (P3, P7). Form sequences are history by construction.
2. **Same-form evolution via `prev_hash` lineage.** Each payload embeds the parent payload hash (`prev_hash`; `""` for roots), forming a chain. Insert is first-writer-wins: there is only ever one payload per hash, and the lineage is the chain of distinct payloads actually written. **Documented gap:** when a superseded payload (already in the store as an unreferenced history payload) is *re-referenced* by a ref update, the ref points at a payload whose `prev_hash` chain does not continue from the immediately previous reference state. This is rare with EKA forward-only transitions (it requires an explicit identity regression) and is **recorded, not hidden**: integrity counts unreferenced payloads, and the refs row always names exactly which payload is current.
3. **Unreferenced payloads retained as the immutable history archive.** Payloads whose hash is no longer referenced are **never deleted** — they are the retained history. Integrity **counts** them (report) but does not **flag** them (not violations).

There is **no dedicated history table**: the v1 `change_log` table was removed because its data already lives inside the immutable `unit.json` payload (the RSF serialization includes the `change-log` array). A second, mutable copy of history would be redundant and would reintroduce the mutable-canonical-state problem. History is self-maintaining: forward-only forms + `prev_hash` chains + retained payloads.

## 6. Integrity

`eka integrity check` makes integrity a first-class runtime concern: a **read-only** scan (no writes; all SQL parameterized) that recomputes every content-derived value and compares it with the stored state, independent of the storage engine.

### 6.1 What it verifies

| Level | Check |
|---|---|
| **1. Payload integrity** | recompute `SHA-256(unit.json ‖ content)` for every `object_payloads` row and compare with `object_hash` |
| **2. Payload decode** | every `unit_json` payload strict-decodes (unknown fields rejected — the same decode path as the exchange layer, RSF §9.5) |
| **3. Reference integrity** | every `object_refs.object_hash` exists in `object_payloads` (target), **and** the refs' derived index columns plus the refs' `form` equal the referenced payload's identity fields — namespace, type, id, instance version, revision, dimension, domain, phase (index) |
| **4. Workspace integrity** | registry foreign keys (`repos` → `projects`), and attachment digests recomputed against `attachments.data` |

Violation kinds, in the report: `payload-hash`, `payload-decode`, `reference-target`, `reference-index`, `attachment-hash`, `registry` — sorted by (kind, subject) for deterministic output.

### 6.2 Detected, not prevented

Manual modification of the SQLite file is **DETECTED, not prevented** — the runtime verifies and reports inconsistencies (exit `1`); it does not pretend a hand-edited database is trusted. SQLite is a persistence layer, not a trust boundary: the check catches tampering after the fact; it cannot (and does not try to) stop it.

### 6.3 Orphans are history, not violations

Unreferenced payloads are the **immutable history archive** — the report counts them as `History payloads` and never reports them as violations. A store with superseded knowledge therefore reports clean with a non-zero history count.

### 6.4 Usage and exit codes

```
eka integrity check
```

No arguments; the workspace is the default `~/.eka` (or `$EKA_HOME`). Exit codes follow the CLI contract:

| Code | Meaning |
|---|---|
| `0` | clean — no integrity violations (history payloads may be present) |
| `1` | violations found — the report is printed; stderr carries `eka: integrity check found N violation(s)` |
| `2` | usage or internal error — workspace resolution or store read failure |

Example output (clean store):

```
Runtime
Workspace  /home/user/.eka
Schema     v2

Summary:
└── Payloads checked: 37
└── References checked: 37
└── Attachments checked: 0
└── History payloads: 0
└── Violations: 0
```

A violation renders as `• <kind> <subject>: <detail>` under the summary, e.g. `• payload-hash <hash>: recomputed SHA-256(unit.json || content) is <got>`; the command then exits `1`.

## 7. Projects, repositories, and registration

A **project** groups one or more **repositories**; a repository's *name* is the basename of its normalized absolute path (the unit key of the canonical store — objects are attributed via `source_repo = repo name`). The project name is the `--name` flag value or the same basename; the project row is created when missing (project id = name).

Registration rules (implemented in `workspace.RegisterRepo`):

- Registering an already-registered repository is a **no-op** (the stored path is refreshed); the report says `already registered`.
- A repository registered under one project **always resolves to that project** — the project that registered it first owns it.
- `eka sync` **auto-registers** the current repository (project name = directory basename) when not registered — registration is a convenience, not a prerequisite.

### Multi-repository projects: the Atrium example

One project = many repositories. Each repository carries only its relevant snapshot; the workspace database reconstructs the **complete** Engineering Knowledge as the union, partitioned by `source_repo` provenance:

```
workspace eka.db  ←  canonical union of the project's knowledge
   ├── source_repo = "api"     ←  exchange/snapshots/ in api/
   ├── source_repo = "web"     ←  exchange/snapshots/ in web/
   └── source_repo = "mobile"  ←  exchange/snapshots/ in mobile/
```

```sh
eka project register ./api    --name atrium     # project "atrium", repo "api"
eka project register ./web    --name atrium     # project "atrium", repo "web"
eka project register ./mobile --name atrium     # project "atrium", repo "mobile"
```

Repositories in one project may share namespaces — partition is **by provenance, never by namespace**; the identity namespace remains a pure identity concern (ADR-010 §Decision 3). Duplicate identity across repositories in one project resolves by deterministic **last-wins** overwrite, recorded in the sync report (a documented v0.2 limitation, §13).

## 8. The synchronization protocol

The full protocol is specified in [ADR-010](decisions/adr-010-synchronization-model.md). The behavioral contract, in short:

### 8.1 Knowledge Snapshot = RSF package in directory layout

A snapshot is an RSF package — `header.json`, `manifest.json`, `declarations.json`, `integrity.json`, `units/`, `attachments/` — written deterministically at `<repo>/exchange/snapshots/`. The directory layout (not the single-file `.ekapkg`) is the canonical transport form: it diffs cleanly in Git and is reviewable in pull requests. **Determinism contract:** identical workspace state → byte-identical snapshot. The snapshot is the same Exchange Package Object Model as `eka export`/`eka import` (ADR-006) — one format, two transports. Per-unit digests in the snapshot's `integrity.json` **equal object hashes by construction** (ADR-011 §Decision 6): a snapshot references the same immutable bytes the store keeps, so transport verification and store verification agree without translation. Mutable runtime state (`sync_log`, registry, refs) is never serialized into packages.

### 8.2 Pull (repository → store)

- **Snapshot mode** (default when a snapshot exists): the package is verified **byte-exact** — structure, strict JSON (unknown fields rejected), SHA-256 integrity (package, per-unit, per-attachment), manifest self-consistency — via the same `exchange.LoadPackage` verification path as import. Verified units are then stored as **immutable payloads** (`store.PutUnit`: the canonical `unit.json` entry is inserted verbatim, keyed by its content hash, `prev_hash` chaining onto the reference's current payload) and the reference is upserted in `object_refs`, attributed with `source_repo` = the repository name.
- **Idempotency:** when the last recorded pull digest equals the current snapshot digest, the upsert work is skipped and the run reports `unchanged`. The package is *still verified first* — a corrupt snapshot is always an error (`snapshot package refused`, exit `1`), never silently skipped.
- **Deletions are NOT applied in v0.2** — the snapshot is an additive transport; units missing from a new pull stay in the store. A future deletion protocol is reserved.

### 8.3 Docs mode (migration / re-seed)

When a repository has a `docs/` tree but no snapshot yet, pull **seeds** the store from the docs tree through the **Knowledge Compiler** (`compile.Compile` — read authoring → parse → validate R0–R12 → normalize → generate CKO → integrity verification → hand off; ADR-012 §Decision 2): the authoring conformance gate runs first — blocking violations **refuse the pull** (exit `1`, the full report is printed) — then the package is assembled exactly as `eka export` would assemble it, so migration-mode digests agree with normal exports, and its units and attachments are upserted. This is the **migration path for every existing repository** (see the migration report: [`migration-report-runtime-v0.2.0.md`](migration-report-runtime-v0.2.0.md)).

`eka sync pull --from-docs` forces this path regardless of snapshot presence — the explicit re-seed.

### 8.4 Push (store → repository)

The repository's references in the canonical store (`source_repo` = this repository's name) are resolved to their immutable payloads and assembled into a deterministic RSF package, written **atomically**: entries are staged in `<repo>/exchange/.snapshots-tmp`, then the old snapshot directory is removed and the staging directory renamed into place. A failed push leaves the previous snapshot untouched. A repository with no stored references is a **no-op** (nothing written, `no-op (no stored objects)` in the report). Every run is recorded in `sync_log`.

Namespace resolution for the package label (deterministic order): the existing snapshot's header namespace when a snapshot exists → else the most common namespace among the repository's objects (ties resolve to the lexicographically smallest) → else an error.

### 8.5 The full cycle and determinism

- `eka sync [path]` = **pull, then push** — the default workflow. Repositories are processed one at a time; the report is deterministic and machine-readable (workspace, project, repository, status, pull source and counts, push counts, snapshot label and 12-hex digest).
- Determinism: identical state → identical pull result, identical push bytes, identical report. The only time-dependent data in the entire runtime is operational bookkeeping — `sync_log` timestamps and digests plus the store's `created_at`/`updated_at` columns — which never enters serialized output.

### 8.6 Store-backed projections: the sync → projection link

`eka sync` is the **precondition of every projection**: it compiles the authoring (docs mode / `--from-docs`) or verifies the snapshot (snapshot mode), seeds the canonical store, and thereby makes the project's knowledge projectable. The projections consume the store directly — `store.UnitsByProject(projectID)` resolves every reference of every registered repository of the project to its immutable payload and strict-decodes each via `exchange.DecodeUnit`, in canonical-form order (ADR-013 §Decision 1–2). The projection scope is the **project union**: multi-repository projects (e.g. Atrium `api`/`web`/`mobile`) project as one knowledge set, partitioned by `source_repo` provenance ([§7](#7-projects-repositories-and-registration)). Projection input is deterministic: identical synced store state → identical unit sequence → identical projection.

The failure modes are deterministic and part of the CLI contract ([§11](#11-cli-behavior-review)):

- **Unregistered repository** — the current directory does not resolve to a registered repository: `eka view` is **refused** (exit `1`) with the deterministic message + hint; `eka watch` renders a refusal frame and keeps polling until the repository is registered and synced (ADR-013 §Decision 3, 5).
- **Registered project without synced knowledge** — `store.UnitsByProject` returns nothing: the projection renders empty with the informational note `no synced knowledge for project <id>; run 'eka sync' after editing docs`, exit `0` — consistent with the existing empty-projection behavior.

Authoring validation (R0–R12) runs inside sync (the compile gate); projections never re-validate authoring per invocation — the correctness of the projected units is the store's integrity contract, checked by `eka integrity check` ([§6](#6-integrity)). Authoring UX: **write Markdown → `eka sync` → `eka view`**.

## 9. Git integration: explicit synchronization

**Decision (v0.2):** users run `eka sync` explicitly; the snapshot is then committed, reviewed, pushed, and merged through **normal Git workflows — untouched**. The workspace database is never committed to any repository (the snapshot is the transport, the DB is the store).

Git hooks and wrapper commands were evaluated and **rejected for v0.2** (ADR-010 §Decision 4, §Alternatives):

| Option | Why rejected |
|---|---|
| Git hooks (post-commit / post-merge auto-sync) | surprise users, run outside their control, and break the explicit determinism contract; premature automation while the protocol is young |
| Wrapper commands (`git eka-*`, aliases) | fork the Git UX, add a maintenance surface, and obscure that synchronization is a knowledge-layer operation, not a VCS operation |
| Auto-sync on `eka validate` | `eka validate` is a read-only conformance gate (P6, projection philosophy); attaching write side effects violates the CLI contract |

The explicit model has a cost: the `docs/` tree and the snapshot can drift if sync discipline lapses. The reconciliation tooling is `eka sync pull --from-docs` (re-seed) and `eka status` (see drift, per-repository last sync).

## 10. Knowledge lifecycle extension

The existing lifecycle — Produce → Organize → Validate → Project → Exchange → Consume — evolves to include the runtime:

```
Draft → Validate → Publish → Synchronize → Project → Consume
```

The runtime inserts **Synchronize** between publishing and projecting: after knowledge is drafted, validated, and published (approved into the governed world), it is synchronized into the canonical store; projections then consume the **canonical model — the CKO** — read from the canonical store (`store.UnitsByProject`; ADR-013), instead of re-scanning Markdown or re-compiling per invocation (a projection reads `exchange.Unit` from the store, never files). **Terminology is not finalized** — stage names (Draft/Publish/Synchronize) and their mapping onto the current Produce/Organize/Exchange steps are implementation-informed and may change before the workflow is fixed. See [`../skeleton/docs/lifecycle.md`](../skeleton/docs/lifecycle.md) for the lifecycle document.

## 11. CLI behavior review

All commands are deterministic (identical input → identical output) and follow the exit-code contract `0` ok (warnings allowed) / `1` blocking violation or integrity failure / `2` usage or internal error.

| Command | Behavior under the runtime | Compatibility |
|---|---|---|
| `eka init` | unchanged — repository bootstrapper; runtime commands are separate | **unchanged** |
| `eka validate` | unchanged — read-only **authoring** conformance gate (R0–R12; the compiler's authoring-validation stage); never writes | **unchanged** |
| `eka export` | unchanged — `.ekapkg` or directory-layout RSF package on demand; the snapshot is a sibling transport with the same object model | **unchanged** |
| `eka import` | unchanged — package → repository integration | **unchanged** |
| `eka view` | **store-backed** — reads Canonical Knowledge Objects **from the workspace canonical store** (`store.UnitsByProject` → `exchange.DecodeUnit`), projecting the **complete knowledge of the project** (the union of every registered repository's units); requires a **registered and synced** repository (run `eka sync` first): unregistered → **refused**, exit `1` (deterministic message + hint); registered-but-unsynced → empty projection + informational note, exit `0`; `compile` is no longer imported (ADR-013) | **revised** — sync-first precondition (was: compile-on-demand); ADR-013 |
| `eka watch` | **store-backed** — polls the canonical store (`store.UnitsByProject` per tick, `--interval` unchanged); TTY contract and byte-comparison redraw unchanged; the **unregistered-repository refusal frame** replaces the compile-failure frame; keeps polling until registered + synced (ADR-013 §Decision 5) | **revised** — sync-first precondition (was: compile-on-demand); ADR-013 |
| `eka sync [path]` | **new** — pull then push; auto-registers the repository | new |
| `eka sync pull [path] [--from-docs]` | **new** — snapshot mode (verified, idempotent) or docs mode (conformance gate, migration/re-seed) | new |
| `eka sync push [path]` | **new** — store → snapshot, atomic temp-dir swap; no-op on empty store | new |
| `eka project register [path] [--name NAME]` | **new** — register a repository under a project; no-op re-registration | new |
| `eka project list` | **new** — deterministic listing of projects and repositories | new |
| `eka status` | **new** — workspace overview: path, schema version, id, created, store totals (**Objects** = references, **Payloads** = immutable objects, **Attachments**), per-repository last sync; never creates the workspace (read-only probe, exits `0` without one) | new |
| `eka integrity check` | **new** — read-only store integrity scan: recompute every payload hash, strict-decode every payload, verify every reference (target existence + derived index columns), recompute attachment digests, check the repository registry; unreferenced payloads counted as history, never violations; exits `0` clean / `1` violations / `2` internal | new |
| `eka version` | unchanged — CLI build version + EKA standard version (the standard remains **1.1**; `0.2.0` is the CLI/artifact version) | **unchanged** |
| `eka completion` | unchanged — Cobra-provided shell completion | **unchanged** |

The compatibility promise: every pre-runtime command keeps its pre-runtime behavior, **with one deliberate, documented revision**: `eka view` / `eka watch` are store-backed and require a registered + synced repository (ADR-013 — the resolution of ADR-012's own future-work note). The runtime adds commands; it changes nothing else.

## 12. Future roadmap (NOT implemented in v0.2)

The schema and protocol leave structural room for all of these; none of them exist yet — **except store-backed projections, implemented in v0.2** (ADR-013; see §2.1 and §8). One paragraph each on how the runtime supports the rest (ADR-011 §Decision 8 reserves the same room):

- **Store-Backed Projections — DONE (v0.2)** — `eka view` / `eka watch` read the canonical store (`store.UnitsByProject` → `exchange.DecodeUnit`) and project the project union; `eka sync` is the precondition; an unregistered repository is refused (exit `1`). This closes the "store-backed projections remain future work" note ADR-012 left open.
- **Knowledge History** — the Immutable Engineering Knowledge Model is the complete raw material: `prev_hash` chains (same-form evolution), forward-only instance versions (distinct forms), and retained unreferenced payloads (the immutable archive). A history query is a **read-side projection** over payloads; there is no dedicated history table, and nothing in v0.2 implements a history command.
- **Knowledge Timeline** — the change-log array inside each immutable `unit.json` payload carries the per-form transitions (date, domain, from, to, by) in occurrence order; a timeline is a projection over payloads. No timeline command exists.
- **Context Engine** — metadata queries over the `object_refs` index columns (`dimension`, `domain`, `phase`, provenance, identity tuple) are point lookups on existing indexes. v0.2 exposes only `eka status` counts.
- **Machine-readable API** — the store is a single SQLite file with a stable schema; a read-only API over payloads + refs is a future command or library, made trivially safe for concurrent readers by immutability. The schema is deliberately public (ADR-009) so the API can be added without redesign.
- **Knowledge Watch (event-driven)** — the store-backed polling watch is implemented (ADR-013: `eka watch` re-reads `store.UnitsByProject` per tick, refusal frame instead of compile-failure frame). What remains future is **event-driven** change detection — `sync_log` + WAL as the change trigger — replacing polling with change notification; the store-backed read makes that a presentation-loop change only (ADR-010 §Decision 6 reserves the room). Not implemented.
- **MCP Integration** — the application packages (`workspace/`, `store/`, `sync/`) are public Go packages independent of the CLI; an MCP server is a thin adapter over them. Not implemented.
- **Knowledge Graph** — relationships are parsed from payloads at read time (they live inside `unit.json`); there is **no SQL-side graph table** in v0.2. Recursive traversal over the five canonical relationship types is a read-time projection over the payload archive — not implemented.
- **Vector Search** — the content BLOB plus the future FTS5/`sqlite-vec` path is the designed route; nothing indexes embeddings today.
- **Multi-device replication** — content addressing makes cross-device transfer naturally deduplicating and verifiable: a receiving device verifies each object by hashing it and rebuilds its own references locally; no mutable canonical state needs to travel (ADR-011 §Decision 6). Snapshots already carry per-unit digests equal to object hashes; replication itself is not implemented.
- **Atrium (unified project runtime)** — the cross-repository union with `source_repo` provenance is the exact data shape a project-wide consumer needs: one project, many repositories, complete Engineering Knowledge in one indexed store. Atrium is the consumer this architecture was inverted for; it is not built.

## 13. Known limitations (v0.2)

Documented, accepted, not mitigated in this milestone (ADR-009/ADR-010 §Consequences):

| Limitation | Detail |
|---|---|
| **`prev_hash` lineage gap (first-writer-wins)** | re-referencing a superseded payload (already in the store as unreferenced history) points the ref at a payload whose chain does not continue from the previous reference state; rare with EKA forward-only transitions, recorded not hidden — integrity counts unreferenced payloads, and the refs row always names the current payload (ADR-011 §Decision 3) |
| **Tampered payloads stay flagged permanently** | a corrupted payload row is detected forever (the archive is immutable, there is no repair command); remediation is manual — remove the tampered row, re-pull from a clean snapshot or re-seed with `--from-docs` (ADR-011 §Decision 4). The companion behavior: `eka sync` reports "snapshot updated" when a push rewrites the snapshot with a different digest (e.g. a tampered store), instead of claiming "unchanged" |
| **Retained history grows the database** | unreferenced payloads are never deleted — no GC in v0.2; a future compaction/GC pass could break `prev_hash` chains (history would be lost or truncated), so the decision is documented before it is made |
| **Derived data requires payload parsing** | relationships and the change-log live inside `unit.json`; there are no SQL-side graph or history indexes in v0.2 — queries that need them pay a decode cost at read time |
| **Store → exchange one-way dependency** | persistence computes hashes over the exchange layer's canonical serialization (`unit.json` bytes); the dependency direction is one-way (exchange never depends on the store) and must not be reversed |
| **Migration loses payload-level history** | the v1 → v2 migration reconstructs units and recomputes hashes; `prev_hash` starts empty — v1 carried no payload history, so every migrated payload is a root (nothing that existed is lost) |
| **Single-writer assumption** | the store assumes one writer (the CLI is single-process); WAL + `busy_timeout(5000)` make concurrent processes safe without file locking — they do not remove the constraint |
| **No deletion propagation** | pull never applies deletions; repositories can accumulate units no longer in the workspace; a tombstone-based deletion protocol is reserved for a future version |
| **Last-wins conflicts** | duplicate identity across repositories in one project resolves by deterministic last-wins overwrite — a silent-overwrite risk, mitigated by the sync report record |
| **No cloud sync** | local-first by design; determinism and offline capability are preserved, server-mediated sync is deferred |
| **No hooks / no automation** | users must run `eka sync`; hooks and wrappers are explicitly rejected for v0.2 (see §9) |
| **Store-backed projection staleness** | a projection reflects the **last sync**, not live authoring edits — an edit is visible only after `eka sync`; `--from-docs` / `eka sync` are the reconciliation tools (ADR-013 §Consequences) |
| **Unregistered repositories cannot be projected** | `eka view` refuses (exit `1`) a repository not registered in the workspace; registration is one command — `eka sync` auto-registers (ADR-013 §Decision 3) |
| **Snapshot swap non-atomic window** | push stages to `.snapshots-tmp`, moves the current snapshot to `.snapshots-old`, renames the staged copy into place, then drops the old copy — a brief non-atomic window between move and rename; acceptable for v0.2, tightened later |
| **Attachments carry no unit references** | attachments are attributed to their provenance pair like objects (a push never packages another repository's attachments), but v1 carries no unit→attachment references — attachment-to-unit linking is a future concern |
| **Push namespace resolution rules** | the package label namespace is resolved as existing-snapshot header → most common namespace → error; a multi-namespace repository without a snapshot may not package the way a human would expect |
| **Workspace backup is operational** | the canonical database is machine-local; knowledge transport depends on `eka sync` being run — workspace backup is now an operational requirement |
| **Dual representation drift** | `docs/` tree and snapshot can drift if sync discipline lapses; `--from-docs` and `eka status` are the reconciliation tooling |
| **Binary size** | the embedded SQLite driver adds roughly +30 MB to the binary and grows the dependency tree — accepted, documented |

## 14. References

- [ADR-009 — Knowledge Runtime Architecture](decisions/adr-009-knowledge-runtime-architecture.md) — workspace, driver evaluation, canonical store DDL
- [ADR-010 — Synchronization Model](decisions/adr-010-synchronization-model.md) — snapshot transport, protocol, Git integration decision, lifecycle extension
- [ADR-011 — Immutable Engineering Knowledge Model](decisions/adr-011-immutable-engineering-knowledge-model.md) — schema v2, content addressing, integrity, v1 → v2 migration; supersedes ADR-009 §3 storage model
- [ADR-012 — Canonical Knowledge Object Runtime](decisions/adr-012-canonical-knowledge-object-runtime.md) — the CKO consumption pivot: Knowledge Compiler, CKO-only runtime, two validators
- [ADR-013 — Store-Backed Projections](decisions/adr-013-store-backed-projections.md) — projections read the canonical store (`store.UnitsByProject`); sync-first precondition; project-union scope; watch polls the store
- [`cko-specification.md`](cko-specification.md) — the Canonical Knowledge Object schema (field reference, derived values, integrity, Runtime Contract)
- [`reference-architecture.md`](reference-architecture.md) — the serialization this runtime stores and transports
- [`cli.md`](cli.md) — command reference, exit codes, output model
- [`migration-report-runtime-v0.2.0.md`](migration-report-runtime-v0.2.0.md) — what changed for existing repositories and how to migrate
- [`../skeleton/docs/exchange/transfer.md`](../skeleton/docs/exchange/transfer.md) — exchange conventions the snapshot realizes
- [`../skeleton/docs/lifecycle.md`](../skeleton/docs/lifecycle.md) — the lifecycle this milestone extends
- [`eka-reference-serialization-format-v1.1.md`](eka-reference-serialization-format-v1.1.md) — the format a Knowledge Snapshot serializes
</content>