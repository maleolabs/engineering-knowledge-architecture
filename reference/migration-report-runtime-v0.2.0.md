# Migration Report — Knowledge Runtime Architecture v0.2.0

> Milestone: EKA v0.2.0 — Knowledge Runtime Architecture (workspace + embedded canonical store + snapshot transport + synchronization protocol).
> Authoritative decisions: [ADR-009](decisions/adr-009-knowledge-runtime-architecture.md), [ADR-010](decisions/adr-010-synchronization-model.md), [ADR-011](decisions/adr-011-immutable-engineering-knowledge-model.md) (schema v2 — see the [Addendum §8](#8-addendum-immutable-engineering-knowledge-model-schema-v2)), [ADR-012](decisions/adr-012-canonical-knowledge-object-runtime.md) (the Canonical Knowledge Object consumption pivot — see the [Addendum §9](#9-addendum-canonical-knowledge-object-runtime)), [ADR-013](decisions/adr-013-store-backed-projections.md) (store-backed projections — see the [Addendum §10](#10-addendum-store-backed-projections)). Convention document, not an artifact.
> Version note: **CLI/artifact version 0.2.0**; the **EKA standard version remains 1.1** — this milestone changes the runtime of the reference implementation, not the standard.

## 1. Purpose

This report records what the Knowledge Runtime milestone changes for existing EKA repositories and their users, classifies every change, and gives the migration path. It is the migration deliverable of EKA v0.2.0.

Target reading: repository owners holding a conformant EKA repository (ADR-001…008 serialization, Git + Markdown `docs/` tree), and tooling maintainers building on the `eka` CLI.

The verdict up front: **no repository content changes are required**. Every existing command keeps its behavior; every existing repository keeps validating; the RSF format is untouched. What changes is where canonical knowledge lives and how it moves.

## 2. What changed — before and after

| Aspect | Before v0.2.0 | After v0.2.0 |
|---|---|---|
| **Canonical storage** | the repository itself: Identity, State, Content, Relationships, Classification live in `docs/` files; every read path walks them | the **local EKA Workspace** (`~/.eka/` or `$EKA_HOME`): `workspace.json` metadata + `eka.db` SQLite canonical store (schema v2, the Immutable Engineering Knowledge Model — see [Addendum §8](#8-addendum-immutable-engineering-knowledge-model-schema-v2)); the repository becomes a transport medium |
| **Repository role** | primary store | **synchronization endpoint**: `docs/` tree (kept, human-readable) + **Knowledge Snapshot** at `exchange/snapshots/` (RSF package in directory layout, deterministic, digest-verified) |
| **Query path** | full scan of `.md` files on every run (`conformance.Scan`) | indexed point lookups and metadata queries over the store (schema: `object_payloads` — immutable content-addressed payloads — and `object_refs` — mutable references with derived index columns — plus `attachments`, `sync_log`; DDL in ADR-009 as amended by ADR-011) |
| **Cross-repository knowledge** | none — each repository is an island | one workspace database, many projects, many repositories; union partitioned by `source_repo` provenance |
| **Knowledge History** | Git history = file history only | `prev_hash` lineages + forward-only instance versions + retained payloads as the derived history of the immutable store, independent of Git (no dedicated history table — see Addendum §8) |
| **Knowledge transport** | `eka export` / `eka import` on demand (`.ekapkg`) | export/import **plus** `eka sync`: deterministic snapshot pull/push between store and repository |
| **Git integration** | repository content committed normally | unchanged — the snapshot is ordinary repository content; **explicit sync** (`eka sync`), no hooks |
| **Lifecycle** | Produce → Organize → Validate → Project → Exchange → Consume | extended: Draft → Validate → Publish → **Synchronize** → Project → Consume (terminology not finalized) |

The inversion is deliberate and documented in [ADR-009](decisions/adr-009-knowledge-runtime-architecture.md): a new storage layer **below** the serialization model. Identity, the State Vector, the RSF, the Knowledge Dimensions, and the Engineering Domains are untouched.

## 3. What changes for users

### 3.1 Mental model

- **The workspace is the store; the repository is the transport.** Canonical Engineering Knowledge lives in the workspace database. The repository carries a snapshot of it for Git and for human reading.
- **Sync is explicit.** Run `eka sync`; commit the snapshot normally. Nothing happens automatically — by design (ADR-010 §Decision 4).
- **The workspace is machine-local.** Never committed, never exchanged, never part of serialized output. A workspace backup is now an operational responsibility.
- Terminology (workspace, snapshot, sync, runtime) is **not finalized** — v0.2.0 is experimental.

### 3.2 CLI

New commands (deterministic; exit codes `0` ok / `1` validation or integrity failure / `2` usage or internal error):

| Command | What it does |
|---|---|
| `eka sync [path]` | pull then push; **auto-registers** the repository (project name = directory basename) |
| `eka sync pull [path] [--from-docs]` | snapshot mode: verify digest, upsert into the store (idempotent); docs mode: conformance-gated seed from the `docs/` tree (migration path, or forced re-seed with `--from-docs`) |
| `eka sync push [path]` | assemble the repository's store objects into a deterministic snapshot at `exchange/snapshots/` (atomic temp-dir swap); no-op on an empty store |
| `eka project register [path] [--name NAME]` | register a repository under a project; same `--name` = same project (multi-repository projects, e.g. Atrium `api`/`web`/`mobile`) |
| `eka project list` | deterministic listing of projects and repositories |
| `eka status` | workspace overview: path, schema version, workspace id, store totals (Objects = references, Payloads = immutable objects, Attachments), per-repository last sync; read-only — never creates the workspace |
| `eka integrity check` | read-only store integrity scan: recompute every payload hash, strict-decode every payload, verify every reference (target + derived index columns), recompute attachment digests, check the repository registry; unreferenced payloads counted as history, never violations; exits `0` clean / `1` violations / `2` internal |

### 3.3 What stays compatible

| Area | Guarantee |
|---|---|
| All legacy commands | `eka init`, `eka validate`, `eka export`, `eka import`, `eka view`, `eka watch`, `eka version`, `eka completion` — unchanged behavior (ADR-010 §Decision 7) |
| Reference project | `reference/project/` ("Feather") — unchanged, still validates with 0 errors / 0 warnings; now also carries the Knowledge Snapshot (`exchange/snapshots/`, label `rsf-repo-feather-1.1`, 37 units) |
| RSF format | untouched — a Knowledge Snapshot is an RSF package in directory layout; the `.ekapkg` single-file form remains available through `eka export` |
| Repository conformance | `eka validate` semantics unchanged; the snapshot directory adds no `.md` files, so artifact counts and verdicts are unaffected |
| Determinism, identity, relationships, validation, projection philosophy | preserved throughout — the runtime never enters serialized output |

## 4. Migration path for existing repositories

An existing repository has a `docs/` tree and no snapshot. Migration is one command, and it is a no-op on content:

```sh
# from the repository root — first sync: conformance gate → docs-mode seed → snapshot written
eka sync

# verify the workspace
eka status

# the snapshot is ordinary content — commit it normally
git add exchange/snapshots
git commit -m "knowledge snapshot"
```

What happens on that first `eka sync` (pull side, docs mode):

1. **Conformance gate** — `conformance.Validate` (R0–R12) runs first. Blocking violations **refuse the pull** (exit `1`, full report printed) — no knowledge is seeded from a non-conformant repository.
2. **Seed** — the package is assembled exactly as `eka export` would assemble it (`exchange.RepositoryPackage`), so the migration-mode digest agrees with a normal export of the same state, and its units and attachments are upserted into the store, attributed to the repository.
3. **Push** — the repository's store objects are assembled into the snapshot at `exchange/snapshots/` (deterministic, atomic write).

After the first sync, later runs use **snapshot mode**: verify + idempotent digest check. The explicit re-seed `eka sync pull --from-docs` is available whenever the docs tree should override the snapshot (the reconciliation tool when the two drift).

Multi-repository projects register each repository under one project name before syncing:

```sh
eka project register ./api    --name atrium
eka project register ./web    --name atrium
eka project register ./mobile --name atrium
eka sync ./api && eka sync ./web && eka sync ./mobile
```

## 5. What is NOT migrated

| Area | Status | Reason / path |
|---|---|---|
| **Deletions** | not propagated | pull is additive by design in v0.2 (ADR-010 §Decision 2); units missing from a new snapshot stay in the store; a tombstone-based deletion protocol is reserved for a future version |
| **Cloud sync** | not provided | local-first by design; determinism and offline capability are preserved (ADR-010 §Alternatives) |
| **Automation** | not provided | no hooks, no auto-sync; explicit `eka sync` only (ADR-010 §Decision 4) |
| **Workspace data** | never enters the repository | `eka.db` is machine-local; backup is an operational task |
| **Multi-writer concurrency** | not supported | single-writer assumption, WAL + `busy_timeout` mitigate concurrent access only |
| **Cross-repo conflict resolution** | last-wins | duplicate identity across repositories in one project: deterministic last-wins overwrite, recorded in the sync report |

## 6. Verification evidence

- **Reference project snapshot, digest-stable** — `reference/project/exchange/snapshots/` carries the full Feather knowledge as an RSF directory package: label `rsf-repo-feather-1.1`, **37 units**, `header.json` / `manifest.json` / `declarations.json` / `integrity.json` / `units/` (with empty `attachments/`). Determinism contract: identical state → byte-identical snapshot.
- **Idempotent re-sync** — a second `eka sync` on unchanged state reports `unchanged` (`no changes: snapshot already up to date`), pulls 0 units, re-pushes **byte-identical** snapshot files, and leaves the store untouched (verified by the sync engine tests, `sync/sync_test.go`).
- **Corruption refusal** — a tampered snapshot (any byte changed) is refused with `snapshot package refused`, exit `1` — integrity failure is never silently skipped (verified by `sync/sync_test.go` and `cmd/sync_test.go`).
- **Docs gate** — a non-conformant repository (no snapshot yet) is refused by the docs-mode validation gate with the full report, exit `1`; no knowledge is seeded (verified by `cmd/sync_test.go`).
- **Multi-repository union** — two repositories under one project sync into the union in the canonical store (8 objects = 4 + 4), each attributed to its `source_repo`, each snapshot carrying exactly its own units (verified by `sync/sync_test.go`).
- **Migration digests agree with export** — docs-mode assembly reuses the export builder (`exchange.RepositoryPackage`), so migration-mode package digests equal normal export digests for the same state.

## 7. Open risks

| Risk | Mitigation / status |
|---|---|
| **Dual representation drift** — `docs/` tree and snapshot diverge when sync discipline lapses | `eka sync pull --from-docs` re-seed and `eka status` per-repository last-sync reporting; the sync report is deterministic and machine-readable |
| **Silent last-wins overwrite** across repositories | recorded in the sync report; deletion and conflict protocols are future work |
| **Snapshot swap non-atomic window** | push writes to `.snapshots-tmp` then removes + renames; a small window; tightened in a future version |
| **Single-writer assumption** | documented; WAL + `busy_timeout(5000)`; a multi-writer story is future work |
| **Machine-local store** — knowledge transport depends on `eka sync` being run | workspace backup becomes an operational requirement; no cloud sync in v0.2 |
| **Terminology instability** — snapshot / runtime / sync names not finalized | explicitly documented; a future rename is a documentation change, not a data change |
| **Binary size** — embedded SQLite driver adds ~+30 MB | accepted and documented in ADR-009 |

## 8. Addendum — Immutable Engineering Knowledge Model (schema v2)

*This addendum records the architectural clarification applied to the runtime before release: the canonical store moved from schema v1 to schema v2. The authoritative decision is [ADR-011](decisions/adr-011-immutable-engineering-knowledge-model.md); the implementation details live in [`runtime-architecture.md`](runtime-architecture.md) §4–§6. This section summarizes what it means for users.*

### 8.1 The model change

- **Immutable, content-addressed payloads.** Engineering Knowledge Objects live in `object_payloads`, keyed by `object_hash` = `SHA-256(unit.json ‖ content)` — byte-identical to the RSF per-unit digest, so object hashes agree with snapshot digests by construction. Payload rows are **insert-only**: never updated, never deleted; the store API has no update path for them.
- **Mutable references only.** `object_refs` maps each RSF Canonical Identity Form to its *current* immutable payload, with provenance and derived index columns (`namespace`, `type`, `id`, `instance_version`, `revision`, `dimension`, `domain`, `phase`).
- **The v1 tables are gone.** `objects`, `relationships`, and `change_log` are removed; relationships and the change-log now live inside the immutable `unit.json` payload (part of the RSF serialization) and are parsed at read time. No dedicated history table exists — history is derived from forward-only instance versions, `prev_hash` lineages, and retained (unreferenced) payloads.
- **SQLite is persistence only.** Immutability belongs to the model (content addressing), not to database triggers; a future storage engine swap preserves the guarantee.

### 8.2 What it means for users

- **New command: `eka integrity check`** — a read-only scan that recomputes every content-derived value: payload hashes, payload decode, reference targets + derived index columns, attachment digests, and the repository registry. Manual modification of `eka.db` is **detected, not prevented**: violations exit `1`; a clean store exits `0`; usage/internal errors exit `2`. Unreferenced payloads are reported as `History payloads` and are **not** violations.
- **`eka status` counts changed.** The overview now reports **Objects** (the number of references — current objects), **Payloads** (the number of immutable objects, including retained history), and **Attachments**; the old `Relationships` line is gone.
- **Immutability guarantees are now construction-level.** "Modify an object" is not expressible through the store API; the guarantee no longer depends on application discipline.

### 8.3 Migration behavior

- **In-place, on first open.** The v1 → v2 migration runs inside the store-open transaction on any database still at schema v1. It creates the v2 tables, reconstructs every v1 `objects` row into a canonical unit (folding `relationships` and `change_log` rows into the unit serialization), **recomputes** the content hash (the v1 digest column is never trusted), stores payload + reference pairs, then drops the v1 tables. A failure mid-way leaves the database at v1 and the migration restarts cleanly on the next open.
- **All v1 data is preserved.** Identity, State Vector, classification, provenance, content, relationships, and the full change-log survive inside the migrated payloads. Deterministic: identical v1 stores migrate to identical v2 stores (rows processed in form order).
- **History lineage starts empty.** v1 carried no payload history, so every migrated payload is a root (`prev_hash = ""`); the change-log transitions are preserved inside `unit.json`. Payload-level lineage accumulates from the first v2 write onward.
- **No user action required.** Existing repositories, snapshots, and sync flows are unaffected; snapshot digests equal object hashes, so pull/push behavior is unchanged.

### 8.4 Trade-offs (accepted, documented in ADR-011 §Consequences)

- `prev_hash` lineage is first-writer-wins: re-referencing a superseded payload produces a lineage gap at the ref (rare with forward-only transitions; recorded, not hidden).
- Retained history grows the database — no GC in v0.2; a future compaction pass could break `prev_hash` chains.
- Derived data (relationships, change-log) requires payload parsing at read time — no SQL-side graph or history indexes in v0.2.
- The store depends on the exchange layer's canonical serialization for hashing (one-way; exchange never depends on the store).
- The v1 → v2 migration loses payload-level history — v1 had none, so nothing that existed is lost.

## 9. Addendum — Canonical Knowledge Object Runtime

*This addendum records the second architectural clarification applied to the runtime before release: the consumption side was inverted from Markdown to the Canonical Knowledge Object (CKO). The authoritative decision is [ADR-012](decisions/adr-012-canonical-knowledge-object-runtime.md); the CKO schema lives in [`cko-specification.md`](cko-specification.md); the implementation details live in [`runtime-architecture.md`](runtime-architecture.md) §2.1. This section summarizes what it means for users.*

### 9.1 What changed

- **The Knowledge Compiler is the canonical gateway.** New `compile/` package — `compile.Compile(root)`: read authoring → parse → validate syntax (authoring conformance R0–R12) → normalize → generate CKO → integrity (package digest) → hand off (persist: sync pull; consume: projections). Markdown is **one authoring adapter** (the conformance package's `Scan`/`analyzeFile`); future authoring interfaces (JSON, forms, visual editors, AI) enter through the same compiler pipeline.
- **The runtime consumes only CKO.** `eka view` / `eka watch` compile authoring into Canonical Knowledge Objects and project over `exchange.Unit` — **zero Markdown parsing in the runtime consumption path** (the `view` package imports no `conformance.Scan`/`Artifact`). The ontology helpers (`ParseReference`, `DomainValues`, `OwnedDomains`, `DomainForToken`, `Stratum`) remain shared and representation-independent.
- **Two validators, two roles.** `eka validate` validates the **authoring representation** (R0–R12, adapter layer); `eka integrity check` validates **CKO integrity** (payload hash, decode, references, attachments, workspace — unchanged from Addendum §8). The runtime never re-validates Markdown.

### 9.2 What did NOT change

| Area | Status |
|---|---|
| **Authoring UX** | unchanged at the CKO pivot — Markdown files in `docs/`; the compiler runs automatically inside `eka sync` (docs mode). **Later revised by the store-backed milestone (Addendum §10):** projections read the store, so the authoring UX becomes write Markdown → `eka sync` → `eka view` |
| **Commands** | unchanged — `eka validate` / `eka view` / `eka watch` / `eka sync` behave identically, same output |
| **RSF bytes** | unchanged — `unit.json` serialization and per-unit digests are the same bytes; the CKO *is* the RSF unit entry |
| **Store schema** | unchanged — the store was already CKO-based (`object_payloads`: `unit.json` + representation-tagged content; `object_refs`) since the immutable-model milestone (Addendum §8); no schema change, no migration |

### 9.3 Migration impact

**None for users.** Same commands, same output, same repository content, same store. The pivot changes *how* the runtime reads authoring (compiled once, canonically), not *what* users write or run.

### 9.4 Trade-offs (accepted, documented in ADR-012 §Consequences)

- **`view` compiles on demand** — at the CKO pivot, authoring validation repeated per invocation; acceptable because it is deterministic and repositories are small. **Resolved by ADR-013 (Addendum §10):** projections are store-backed — per-invocation compilation and re-validation are gone, replaced by the sync-first precondition (trade-offs in §10.4).
- **Relationship rendering nuance** — projection relationship targets are presented in the authoring line convention (`ctr:wave-1`) for same-namespace references (instance version dropped), and in **full canonical form** (`<namespace>/<type>:<id>:<instance-version>`) for cross-namespace references (`view`'s `referenceForm`). Same-namespace output matches the pre-pivot authoring convention; cross-namespace targets now render canonically.
- **Authoring validation and runtime validation are distinct commands** with distinct scopes — documented, but the split could confuse users who expect one validator.

## 10. Addendum — Store-Backed Projections

*This addendum records the third architectural clarification applied to the runtime before release: projections moved off the docs tree and onto the workspace canonical store. The authoritative decision is [ADR-013](decisions/adr-013-store-backed-projections.md); the implementation details live in [`runtime-architecture.md`](runtime-architecture.md) §2.1, §8, and §11. This section summarizes what it means for users.*

### 10.1 What changed

- **Projections read the store.** `eka view` / `eka watch` now read Canonical Knowledge Objects from the EKA workspace canonical store — `store.UnitsByProject(projectID)` resolves every reference of the project to its immutable payload, `exchange.DecodeUnit` strict-decodes each (canonical-form order, digest-tagged) — instead of compiling the docs tree on demand. **Zero Markdown in the projection path.** `compile` remains the **authoring gateway** used by `eka sync` (docs mode / `--from-docs`), but is **no longer imported** by the view/watch commands.
- **Projection scope = the project (union).** The projections cover the complete Engineering Knowledge of the project — the union of every registered repository's units (e.g. Atrium `api`/`web`/`mobile` project as one knowledge set), partitioned by `source_repo` provenance. Multi-repository projects project as **one knowledge set**, with no per-repo projection logic.
- **Synchronization is the precondition.** A repository must be **registered and synced** before it can be projected. The authoring UX becomes **write Markdown → `eka sync` → `eka view`**; authoring validation (R0–R12) still runs inside sync (the compile gate), never at projection time. Failure modes, deterministic:
  - **Unregistered repository** — refused, exit `1`, deterministic message + hint: `eka: view refused: repository <abs> is not registered in the EKA workspace; run 'eka sync' (auto-registers) or 'eka project register' first`.
  - **Registered project without synced knowledge** — empty projection + informational note (`no synced knowledge for project <id>; run 'eka sync' after editing docs`), exit `0` — consistent with the existing empty-projection behavior.
- **`eka watch` polls the store.** Per tick (`--interval` unchanged) the project's units are re-read from the store; an `eka sync` run in another terminal is picked up without a restart. The **unregistered-repository refusal frame** replaces the compile-failure frame; the TTY contract and byte-comparison redraw logic are unchanged; the refusal is a rendered state, never an exit.

### 10.2 What did NOT change

| Area | Status |
|---|---|
| **Commands** | unchanged surfaces — `eka validate` / `eka view` / `eka watch` / `eka sync` keep their names, arguments, determinism, and the 0/1/2 exit-code contract; only the view/watch *input source and precondition* changed |
| **Authoring** | unchanged — Markdown files in `docs/`; `eka validate` still validates the authoring representation (R0–R12, read-only); the same gate runs inside `eka sync` before seeding |
| **RSF bytes** | unchanged — `unit.json` serialization, per-unit digests, and the snapshot format are untouched; `eka integrity check` recomputes the same object hashes the projections now read |
| **Store schema** | unchanged — no schema change, no migration; `store.UnitsByProject` is a read over the existing `object_refs` + `object_payloads` (Addendum §8) |

### 10.3 Migration impact

- **Existing users with a synced workspace: no action.** The first `eka view` after upgrading works immediately — the store already holds the project's units (any prior `eka sync` seeded them). The change is in *how* view reads, not in *what* is stored.
- **Fresh repositories: one extra step.** `eka sync` before `eka view` — the previously seamless compile-on-demand is replaced by an explicit step. The reference project's snapshot is committed, so `eka sync reference/project` then `eka view` works out of the box.
- **Unregistered clones: one command to fix.** `eka view` no longer works on a repository outside the workspace; `eka sync` (auto-registers) or `eka project register` first — the refusal message says exactly this.

### 10.4 Trade-offs (accepted, documented in ADR-013 §Consequences)

- **Staleness vs live edits** — a projection reflects the **last sync**, not live authoring edits; an edit is visible only after re-sync. Deterministic and documented, never a silent hybrid; `eka sync` / `--from-docs` are the reconciliation tools.
- **Unregistered refusal** — a repository outside the workspace cannot be projected until registered; registration is one command (`eka sync` auto-registers).
- **Explicit sync step** — ADR-012's "authoring experience unchanged (`view` compiles on demand)" is deliberately revised; the UX cost is one command before the first view.

---

*End of Migration Report — Knowledge Runtime Architecture v0.2.0.*
