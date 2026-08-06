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
