# Transfer — Import/Export Conventions

> Anchor EKA: Exchange Layer — transfer (import/export). Convention document, not an artifact.
> Standard: EKA v1.1, dated 2026-08-05.

This document governs the exchange of EKA artifacts between repositories (e.g. parent/component/archive repositories). All conventions apply in both directions (import = export reversed).

## 1. Round-Trip Requirements

### 1.1 Lossless (no loss)

Transfer must preserve **all** identity and state information:

- [ ] Full identity: `namespace`, `type`, `id`, `instance-version`, `revision`.
- [ ] Full state: the entire owned state vector **plus** the `change-log` history.
- [ ] Well-formed content (required sections intact).
- [ ] Relationships **by identity** (references preserved as references, not converted to text).
- [ ] Classification: `dimension`, `dimensions-secondary`.
- [ ] Preservation status (`existence-state`) is not changed by transfer.

### 1.2 Idempotent

- [ ] Re-importing the same package = **no-op** (no duplication), or an explicit *clean replace* declaration.
- [ ] Re-importing never creates duplicate artifacts.

### 1.3 Referential Integrity

- [ ] No dangling references: referenced artifacts must travel along, already exist in the target, or be allowed as a warning because the target is `draft`.

### 1.4 Identity Conflict Policy

When the identity `(namespace, type, id, instance-version)` already exists in the target:

| Option | Condition |
|---|---|
| **Reject** (default) | conflict reported; transfer cancelled |
| **Explicit re-namespace** | the entire artifact identity is declaratively moved to a new namespace (and all its references updated consistently) |

- [ ] **Never** a *silent merge* — silently merging two artifacts of different identities is forbidden.

### 1.5 Validation Before Commit

- [ ] The package to be transferred passes the entire [validation.md](validation.md) checklist.
- [ ] After import, the target is revalidated before commit.

### 1.6 Contract Version

- [ ] Every transfer package declares **two versions** in the Contract Header (Exchange Specification §9.2.1): the serialization contract version (e.g. `eka-exchange-format: 1.0`) and the EKA specification version it conforms to (e.g. `eka-spec: 1.0`).
- [ ] Import rejects packages with unsupported contract versions.

## 2. What Is Preserved

| Aspect | Description |
|---|---|
| Full identity | `(namespace, type, id, instance-version)` unchanged during transfer |
| State + change-log | the full transition history of owned domains |
| Well-formed content | required sections per type family remain intact |
| Relationships by identity | references keep following identity (not file paths) |
| Classification | `dimension`/`dimensions-secondary` preserved |
| Preservation status | `existence-state` not changed by transfer mechanics |

## 3. EX Limits on Transfer

- EX **does not judge content correctness** — only conformance and integrity.
- EX **does not change state** — state transitions remain the sole right of the state owner; transfer only copies values.
- Projections (tickets/tables) are not transferred as sources of truth; after import, projections are refreshed from the owner state in the target.

## 4. Knowledge Snapshots and Synchronization

The Knowledge Runtime (v0.2) adds a second transport over the same Exchange Package Object Model: **Knowledge Snapshots**, synchronized between the repository and the local EKA Workspace canonical store. The snapshot is the same RSF package contract as import/export — one format, two transports.

### 4.1 Snapshot = RSF package at `exchange/snapshots/`

A Knowledge Snapshot is an RSF package written in **directory layout** (not the single-file `.ekapkg`) at `exchange/snapshots/`:

| Entry | Contents |
|---|---|
| `header.json` | package header: serialization version 1.1, exchange format 1, specification 1.0, exporter `eka`, label `rsf-repo-<namespace>-1.1`, scope `repo`, namespace |
| `manifest.json` | ordered unit list (canonical identity form) with per-unit and package digests |
| `declarations.json` | closure + external reference declarations |
| `integrity.json` | SHA-256 digests: package-level, per-unit, per-attachment |
| `units/` | one directory per unit: `unit.json` (identity, state vector, change log, relationships, classification) + `content` payload, byte-exact |
| `attachments/` | non-`.md` payloads from `docs/`, byte-exact |

The same verification rules as import apply (RSF §9.4/§9.5): entry structure, strict JSON (unknown fields rejected), digest verification, manifest self-consistency. A corrupt snapshot is **refused, never silently skipped**.

### 4.2 `eka sync` workflow

- `eka sync [path]` — **pull, then push** (the default cycle).
- `eka sync pull [path]` — verify the snapshot, then upsert its units and attachments into the workspace store, attributed to the repository. **Idempotent:** an unchanged snapshot digest skips the work.
- `eka sync pull --from-docs` — re-seed from the `docs/` tree (conformance gate first) — the reconciliation tool when docs and snapshot drift.
- **Migration mode:** a repository with a `docs/` tree but no snapshot is seeded from the docs tree through the conformance gate on first pull.
- `eka sync push [path]` — assemble the repository's stored objects into the deterministic snapshot (atomic temp-dir swap; failed pushes leave the previous snapshot untouched).

### 4.3 Deletions and conflicts (v0.2)

- **Deletions are never applied** — pull is additive; units missing from a new snapshot stay in the store. A deletion protocol is reserved for a future version.
- **Duplicate identity** across repositories in one project resolves by deterministic **last-wins** overwrite, recorded in the sync report.

### 4.4 Explicit Git commit

The snapshot is ordinary repository content. After `eka sync`, commit it with the normal workflow:

- [ ] `git add exchange/snapshots` and commit after every sync that changed the snapshot.
- [ ] Review snapshot changes in pull requests like any other content (the directory layout diffs cleanly).

There are **no hooks**: synchronization is explicit by design (the workspace database is never committed; the snapshot is the transport). `eka project register <path> --name <project>` groups repositories into projects for multi-repository synchronization.
