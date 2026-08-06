---
namespace: feather
type: bug
id: empty-title-crash
instance-version: 1
revision: 1
execution-state: done
existence-state: active
dimension: requirements
author: Jonas Berg
created: 2026-08-01
updated: 2026-08-02
supersedes: []
derives-from: []
depends-on:
  - ctr:wave-7
  - plan:roadmap-v1:1
amends: []
validates: []
change-log:
  - date: 2026-08-01
    domain: existence-state
    from: "-"
    to: active
    by: Jonas Berg
  - date: 2026-08-01
    domain: execution-state
    from: "-"
    to: planned
    by: Jonas Berg
  - date: 2026-08-01
    domain: execution-state
    from: planned
    to: todo
    by: Jonas Berg
  - date: 2026-08-01
    domain: execution-state
    from: todo
    to: in-progress
    by: Jonas Berg
  - date: 2026-08-02
    domain: execution-state
    from: in-progress
    to: in-review
    by: Jonas Berg
  - date: 2026-08-02
    domain: execution-state
    from: in-review
    to: done
    by: Jonas Berg
---

# Bug — Empty Title Crashes Publish

## Description

**Symptom.** Publishing a draft whose title is empty (or whitespace-only) panics the server: `POST /api/posts/{id}/publish` dereferences the empty slug path, and the request 500s — after which the process dies (a recover() handler was missing on the publish route).

**Reproduction.**

1. Create a draft via the editor.
2. Clear the title field entirely; save (autosave writes the file).
3. Click publish.
4. Server panics; the process exits; all requests fail until restart.

**Expected behavior.** The API returns `400` with a validation message; the draft stays a draft; the process keeps serving (per `spec:publishing-api`).

## Impact

- **Severity:** high — a trivial input sequence takes down the whole site (single-binary deployment has no second instance to absorb it).
- **Users:** any author with an empty title; all readers while the process is down.
- **Root cause (identified during fix):** publish route validated the body after the slug path was already built from the title; empty title → empty path → index panic. Fix: validate title before any path construction + install a panic-recover boundary on all `/api/*` routes (severity baseline referenced by `std:definition-of-done`).
