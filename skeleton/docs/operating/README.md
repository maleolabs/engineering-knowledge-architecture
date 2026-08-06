# docs/operating/ — Operating Layer (OS)

> Anchor EKA: Operating Layer — state machine + protocol. OS holds state; OS never edits content.

## Summary

This folder serializes the project's **state** and **execution**: who owns state, how transitions happen, and how state is read without being edited.

- **State ownership (P6):** every State Domain is owned by exactly one artifact type (see the State Vector in each subfolder's README). Absence of a field = N/A for that type.
- **Single-writer:** only the state owner writes its state fields; every transition is recorded in `change-log`. No field has two writers.
- **Projections never write state:** work item tables in containers and `tkt-` are mere shadows; refreshed on read, never edited.
- **Two change channels:** content governance and state protocol are two separate mechanisms that must not be mixed.
- **Content lives in the Knowledge Layer** (knowledge dimensions in `docs/`); OS only manages the state of those artifacts — OS does not change their content.

## Documents

| File | Role |
|---|---|
| [protocol.md](protocol.md) | Operating Manual — convention, not an artifact (no `type`/`id`) |

## Subfolders

| Folder | Content | Anchored State Domain |
|---|---|---|
| [work-items/](work-items/) | 6 work item subtypes: `sto-`, `ts-`, `bug-`, `td-`, `ch-`, `spk-` | Execution State |
| [containers/](containers/) | `ctr-` Execution Container — exactly one active | Container State |
| [sessions/](sessions/) | `ses-` Session — execution notes | Existence State |
| [projections/](projections/) | `tkt-` Ticket — empty state vector (pure projection) | — (owns no state) |

## Touch Points with Other Layers

- Knowledge Layer: OS reads the content state (`content-state`) of knowledge artifacts when assessing readiness, but content is managed by the dimension owner.
- Exchange Layer: validation reads state from here; see [../exchange/](../exchange/).
