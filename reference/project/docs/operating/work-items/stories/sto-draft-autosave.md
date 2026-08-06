---
namespace: feather
type: sto
id: draft-autosave
instance-version: 1
revision: 1
execution-state: in-progress
existence-state: active
dimension: requirements
author: Jonas Berg
created: 2026-07-29
updated: 2026-08-05
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-07-29
    domain: existence-state
    from: "-"
    to: active
    by: Jonas Berg
  - date: 2026-07-29
    domain: execution-state
    from: "-"
    to: planned
    by: Jonas Berg
  - date: 2026-07-31
    domain: execution-state
    from: planned
    to: todo
    by: Jonas Berg
  - date: 2026-08-03
    domain: execution-state
    from: todo
    to: in-progress
    by: Jonas Berg
---

# Story — Autosave Drafts

## Description

As an author, I never lose more than the last keystroke: the editor saves the draft buffer in the background while I type, with no explicit save action. Autosave writes the Markdown file (source of truth per `adr:content-storage`) and bumps the revision counter; it must survive reload, tab close, and concurrent edits.

## Acceptance Criteria

- The editor debounces saves (≤ 1 s after typing stops) and flushes on `beforeunload`.
- Reloading the editor restores the last saved buffer; a crash loses at most the debounce window.
- Concurrent edits: `PATCH` with a stale `If-Match` revision returns 409; the autosave loop reconciles by reloading the file, never by clobbering it.
- The autosaved draft is never publicly visible (index status stays `draft`).
- File writes are atomic (temp file + rename); `feather reindex` rebuilds the index from files without loss.
