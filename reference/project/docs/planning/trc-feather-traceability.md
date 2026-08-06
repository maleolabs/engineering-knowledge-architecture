---
namespace: feather
type: trc
id: feather-traceability
instance-version: 1
revision: 1
content-state: approved
existence-state: active
dimension: planning
author: Maya Patel
created: 2026-07-15
updated: 2026-07-16
supersedes: []
derives-from:
  - scp:mvp-v1:1
  - plan:roadmap-v1:1
depends-on:
  - req:publishing-core
amends: []
validates: []
change-log:
  - date: 2026-07-15
    domain: existence-state
    from: "-"
    to: active
    by: Maya Patel
  - date: 2026-07-15
    domain: content-state
    from: "-"
    to: draft
    by: Maya Patel
  - date: 2026-07-16
    domain: content-state
    from: draft
    to: review
    by: Maya Patel
  - date: 2026-07-16
    domain: content-state
    from: review
    to: approved
    by: Maya Patel
---

# Traceability — Feather Knowledge Matrix

## Objective

The traceability map of the Feather knowledge base: how requirements, decisions, epics, and work items hang together. This artifact carries relationships only; it does not replace the relationships written on the referring artifacts.

## Scope

**Requirements ↔ decisions:**

| Requirement | Derived decisions |
|---|---|
| `req:publishing-core` | `adr:content-storage`, `adr:search-sqlite-fts`, `adr:plugin-model-deferred` |
| `req:comments-phase2` | (none yet — draft) |

**Scope ↔ epics ↔ work items:**

| Scope | Epic | Work items |
|---|---|---|
| `scp:mvp-v1` | `epc:authoring-experience` | `sto:publish-post`, `sto:draft-autosave`, `ts:markdown-renderer`, `spk:markdown-syntax-extension`, `bug:empty-title-crash` |
| `scp:mvp-v1` | `epc:distribution` (draft) | (none yet) |
| `scp:mvp-v1` | — | `td:reduce-query-count`, `ch:update-go-version` |

**Plan ↔ containers ↔ work items:**

| Plan | Container | Work items |
|---|---|---|
| `plan:roadmap-v1` | `ctr:wave-6` (completed) | foundation deliverables (prose record) |
| `plan:roadmap-v1` | `ctr:wave-7` (active) | `sto:publish-post`, `sto:draft-autosave`, `ts:markdown-renderer`, `bug:empty-title-crash`, `td:reduce-query-count`, `ch:update-go-version`, `spk:markdown-syntax-extension` |

**Upward anchor:** every chain in this matrix resolves upward to `req:publishing-core` (stratum 1), keeping the whole knowledge base stratification-traceable (R10).

## Out of Scope

- Cross-namespace traceability (Feather knowledge is single-namespace).
- Tracing content edits or revisions — the change-log of each artifact covers that.
