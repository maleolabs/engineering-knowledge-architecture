---
namespace: feather
type: sto
id: publish-post
instance-version: 1
revision: 1
execution-state: done
existence-state: active
dimension: requirements
author: Jonas Berg
created: 2026-07-29
updated: 2026-08-03
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
  - date: 2026-07-30
    domain: execution-state
    from: planned
    to: todo
    by: Jonas Berg
  - date: 2026-07-30
    domain: execution-state
    from: todo
    to: in-progress
    by: Jonas Berg
  - date: 2026-08-01
    domain: execution-state
    from: in-progress
    to: in-review
    by: Jonas Berg
  - date: 2026-08-03
    domain: execution-state
    from: in-review
    to: done
    by: Jonas Berg
---

# Story — Publish a Post

## Description

As an author, I can publish a draft with one action, and the post appears on the public site immediately. This is the heart of the publishing-core requirement: the draft → publish flip defined in `adr:content-storage` (file stays, index status flips, FTS index updates).

## Acceptance Criteria

- `POST /api/posts/{id}/publish` flips the post's index status from `draft` to `published` and stamps `published-at` (see `spec:publishing-api`).
- The public page renders the post after publish; the FTS index contains the post (search finds it).
- Publishing a non-existent or already-published post returns a clear error (404 / 409) and leaves state unchanged.
- Autosaved-but-unpublished content is never exposed publicly.
- Verified against the DoD (`std:definition-of-done`): tests added, `eka validate` clean, review recorded in `rvw:publishing-core-review`.
